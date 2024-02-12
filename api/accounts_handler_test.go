package api

import (
	mockdb "Simple-Bank/db/mock"
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"Simple-Bank/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateUserAPI(t *testing.T) {
	createdUser, password := randomUser(t)

	testCases := []struct {
		name          string
		req           requests.CreateUserRequest
		buildStubs    func(services *mockdb.MockServices, req requests.CreateUserRequest)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			req: requests.CreateUserRequest{
				Username: createdUser.Username,
				FullName: createdUser.FullName,
				Email:    createdUser.Email,
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Eq(req)).
					Times(1).
					Return(createdUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, createdUser)
			},
		},
		{
			name: "InternalError",
			req: requests.CreateUserRequest{
				Username: createdUser.Username,
				FullName: createdUser.FullName,
				Email:    createdUser.Email,
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Eq(req)).
					Times(1).
					Return(models.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateEmail",
			req: requests.CreateUserRequest{
				Username: createdUser.Username,
				FullName: createdUser.FullName,
				Email:    createdUser.Email,
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Eq(req)).
					Times(1).
					Return(models.User{}, &pgconn.PgError{ConstraintName: "users_email_key"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			req: requests.CreateUserRequest{
				Username: createdUser.Username,
				FullName: createdUser.FullName,
				Email:    createdUser.Email,
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Eq(req)).
					Times(1).
					Return(models.User{}, &pgconn.PgError{ConstraintName: "users_pkey"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidFullname",
			req: requests.CreateUserRequest{
				Username: createdUser.Username,
				FullName: util.RandomEmail(),
				Email:    createdUser.Email,
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			req: requests.CreateUserRequest{
				Username: util.RandomEmail(),
				FullName: createdUser.FullName,
				Email:    createdUser.Email,
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			req: requests.CreateUserRequest{
				Username: createdUser.Username,
				FullName: createdUser.FullName,
				Email:    util.RandomUsername(),
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			req: requests.CreateUserRequest{
				Username: createdUser.Username,
				FullName: createdUser.FullName,
				Email:    createdUser.Email,
				Password: "123456",
			},
			buildStubs: func(services *mockdb.MockServices, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {

		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			services := mockdb.NewMockServices(ctrl)
			// build stubs
			testCase.buildStubs(services, testCase.req)

			// testing http server
			server := NewServer(services)
			recorder := httptest.NewRecorder()

			jsonReq, err := json.Marshal(testCase.req)
			require.NoError(t, err)

			httpReq, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonReq))
			require.NoError(t, err)

			server.RouterServeHTTP(recorder, httpReq)
			testCase.checkResponse(t, recorder)
		})

	}

}

func TestGetUserAPI(t *testing.T) {
	createdUser, _ := randomUser(t)

	testCases := []struct {
		name          string
		username      string
		buildStubs    func(services *mockdb.MockServices, username string)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: createdUser.Username,
			buildStubs: func(services *mockdb.MockServices, username string) {
				services.EXPECT().
					GetUser(gomock.Eq(username)).
					Times(1).
					Return(createdUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, createdUser)
			},
		},
		{
			name:     "NotFound",
			username: createdUser.Username,
			buildStubs: func(services *mockdb.MockServices, username string) {
				services.EXPECT().
					GetUser(gomock.Eq(username)).
					Times(1).
					Return(models.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			username: createdUser.Username,
			buildStubs: func(services *mockdb.MockServices, username string) {
				services.EXPECT().
					GetUser(gomock.Eq(username)).
					Times(1).
					Return(models.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidUsername",
			username: util.RandomEmail(), // emails are invalid usernames
			buildStubs: func(services *mockdb.MockServices, username string) {
				services.EXPECT().
					GetUser(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {

		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			services := mockdb.NewMockServices(ctrl)
			// build stubs
			testCase.buildStubs(services, testCase.username)

			// testing http server
			server := NewServer(services)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", testCase.username)

			httpReq, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.RouterServeHTTP(recorder, httpReq)
			testCase.checkResponse(t, recorder)
		})

	}

}

func randomUser(t *testing.T) (models.User, string) {
	password := util.RandomPassword()
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	return models.User{
		Username:       util.RandomUsername(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomFullname(),
		CreatedAt:      time.Now().Truncate(time.Second).UTC(),
		UpdatedAt:      time.Now().Truncate(time.Second).UTC(),
		DeletedAt:      gorm.DeletedAt{},
	}, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user models.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response responses.UserInformationResponse
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)
	require.Equal(t, user.Username, response.Username)
	require.Equal(t, user.FullName, response.FullName)
	require.Equal(t, user.Email, response.Email)
	require.Equal(t, user.CreatedAt.Local().Truncate(time.Second), response.CreatedAt)
	require.Equal(t, user.UpdatedAt.Local().Truncate(time.Second), response.UpdatedAt)
	require.Equal(t, user.DeletedAt.Time.Truncate(time.Second), response.DeletedAt)
}
