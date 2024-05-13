package api

import (
	mockdb "Simple-Bank/db/mock"
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"Simple-Bank/token"
	mocktokenmaker "Simple-Bank/token/mock"
	"Simple-Bank/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
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
	configs := getTestConfig()
	createdUser, password := randomUser(t)

	var accessTokenPayload *token.Payload
	var accessToken string
	var refreshTokenPayload *token.Payload
	var refreshToken string

	var session models.Session

	testCases := []struct {
		name     string
		subCases []struct {
			name string
			req  requests.CreateUserRequest
		}
		buildStubs func(
			services *mockdb.MockServices,
			tokenMaker *mocktokenmaker.MockMaker,
			req requests.CreateUserRequest,
		)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker)
	}{
		{
			name: "OK",
			subCases: []struct {
				name string
				req  requests.CreateUserRequest
			}{
				{
					name: "OK",
					req: requests.CreateUserRequest{
						Username: createdUser.Username,
						FullName: createdUser.FullName,
						Email:    createdUser.Email,
						Password: password,
					},
				},
			},
			buildStubs: func(services *mockdb.MockServices, tokenMaker *mocktokenmaker.MockMaker, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Eq(req)).
					Times(1).
					Return(createdUser, nil)
				tokenMaker.EXPECT().
					CreateToken(req.Username, configs.TokenAccessTokenDuration).
					Times(1).
					Return(accessToken, accessTokenPayload, nil)
				tokenMaker.EXPECT().
					CreateToken(req.Username, configs.TokenRefreshTokenDuration).
					Times(1).
					Return(refreshToken, refreshTokenPayload, nil)
				services.EXPECT().
					CreateSession(newSessionMatcher(session)).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchLogin(
					t,
					recorder.Body,
					createdUser,
					configs.TokenSymmetricKey,
					refreshToken,
					accessToken,
				)
			},
		},
		{
			name: "InternalError",
			subCases: []struct {
				name string
				req  requests.CreateUserRequest
			}{
				{
					name: "ErrConnDone",
					req: requests.CreateUserRequest{
						Username: createdUser.Username,
						FullName: createdUser.FullName,
						Email:    createdUser.Email,
						Password: password,
					},
				},
			},
			buildStubs: func(services *mockdb.MockServices, tokenMaker *mocktokenmaker.MockMaker, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Eq(req)).
					Times(1).
					Return(models.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Forbidden",
			subCases: []struct {
				name string
				req  requests.CreateUserRequest
			}{
				{
					name: "DuplicateEmail",
					req: requests.CreateUserRequest{
						Username: createdUser.Username,
						FullName: createdUser.FullName,
						Email:    createdUser.Email,
						Password: password,
					},
				},
			},
			buildStubs: func(services *mockdb.MockServices, tokenMaker *mocktokenmaker.MockMaker, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Eq(req)).
					Times(1).
					Return(models.User{}, &pgconn.PgError{ConstraintName: "users_email_key"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Forbidden",
			subCases: []struct {
				name string
				req  requests.CreateUserRequest
			}{
				{
					name: "DuplicateUsername",
					req: requests.CreateUserRequest{
						Username: createdUser.Username,
						FullName: createdUser.FullName,
						Email:    createdUser.Email,
						Password: password,
					},
				},
			},
			buildStubs: func(services *mockdb.MockServices, tokenMaker *mocktokenmaker.MockMaker, req requests.CreateUserRequest) {
				services.EXPECT().
					CreateUser(gomock.Eq(req)).
					Times(1).
					Return(models.User{}, &pgconn.PgError{ConstraintName: "users_pkey"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			subCases: []struct {
				name string
				req  requests.CreateUserRequest
			}{
				{
					name: "InvalidPassword",
					req: requests.CreateUserRequest{
						Username: createdUser.Username,
						FullName: createdUser.FullName,
						Email:    createdUser.Email,
						Password: "123456",
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
				},
				{
					name: "InvalidUsername",
					req: requests.CreateUserRequest{
						Username: util.RandomEmail(),
						FullName: createdUser.FullName,
						Email:    createdUser.Email,
						Password: password,
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
				},
				{
					name: "InvalidUsername",
					req: requests.CreateUserRequest{
						Username: util.RandomEmail(),
						FullName: createdUser.FullName,
						Email:    createdUser.Email,
						Password: password,
					},
				},
			},
			buildStubs: func(services *mockdb.MockServices, tokenMaker *mocktokenmaker.MockMaker, req requests.CreateUserRequest) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, subCase := range testCase.subCases {
				t.Run(subCase.name, func(t *testing.T) {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()

					services := mockdb.NewMockServices(ctrl)
					tokenMaker := mocktokenmaker.NewMockMaker(ctrl)

					accessToken, accessTokenPayload = createToken(t, createdUser.Username, configs.TokenAccessTokenDuration, configs.TokenSymmetricKey)
					refreshToken, refreshTokenPayload = createToken(t, createdUser.Username, configs.TokenRefreshTokenDuration, configs.TokenSymmetricKey)

					session = createSession(refreshTokenPayload.ID, createdUser.Username, refreshToken)

					// build stubs
					testCase.buildStubs(services, tokenMaker, subCase.req)

					// testing http server
					server := NewTestServer(t, services, tokenMaker)
					recorder := httptest.NewRecorder()

					jsonReq, err := json.Marshal(subCase.req)
					require.NoError(t, err)

					httpReq, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonReq))
					require.NoError(t, err)

					server.RouterServeHTTP(recorder, httpReq)
					testCase.checkResponse(t, recorder, server.handlers.tokenMaker)
				})
			}
		})
	}

}

func TestGetUserAPI(t *testing.T) {
	configs := getTestConfig()
	createdUser, _ := randomUser(t)

	testCases := []struct {
		name          string
		username      string
		buildStubs    func(services *mockdb.MockServices, username string)
		setupAuth     func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker)
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
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, createdUser.Username, time.Minute, httpReq)
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
					Return(models.User{}, gorm.ErrRecordNotFound)
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, createdUser.Username, time.Minute, httpReq)
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
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, createdUser.Username, time.Minute, httpReq)
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
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, createdUser.Username, time.Minute, httpReq)
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

			tokenMaker, err := token.NewPasetoMaker(configs.TokenSymmetricKey)

			// testing http server
			server := NewTestServer(t, services, tokenMaker)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", testCase.username)

			httpReq, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, httpReq, server.handlers.tokenMaker)

			server.RouterServeHTTP(recorder, httpReq)
			testCase.checkResponse(t, recorder)
		})

	}

}

type sessionMatcher struct {
	session models.Session
}

func (s sessionMatcher) Matches(x any) bool {
	inputSession, ok := x.(models.Session)
	if !ok {
		return false
	}

	if s.session.ID == inputSession.ID ||
		s.session.Username == inputSession.Username ||
		s.session.UserAgent == inputSession.UserAgent ||
		s.session.ClientIP == inputSession.ClientIP ||
		s.session.IsBlocked == inputSession.IsBlocked ||
		s.session.RefreshToken == inputSession.RefreshToken ||
		s.session.DeletedAt == inputSession.DeletedAt ||
		s.session.CreatedAt.Truncate(time.Second) == inputSession.CreatedAt.Truncate(time.Second) ||
		s.session.ExpiresAt.Truncate(time.Second) == inputSession.ExpiresAt.Truncate(time.Second) {
		return true
	}

	return false
}

func (s sessionMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", s.session, s.session)
}

func newSessionMatcher(session models.Session) gomock.Matcher {
	return sessionMatcher{
		session: session,
	}
}

func TestLogin(t *testing.T) {
	randomUser, password := randomUser(t)
	configs := getTestConfig()

	var accessTokenPayload *token.Payload
	var accessToken string
	var refreshTokenPayload *token.Payload
	var refreshToken string

	var session models.Session

	testCases := []struct {
		name       string
		req        requests.LoginRequest
		buildStubs func(
			services *mockdb.MockServices,
			tokenMaker *mocktokenmaker.MockMaker,
			req requests.LoginRequest,
		)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker)
	}{
		{
			name: "OK",
			req: requests.LoginRequest{
				Username: randomUser.Username,
				Password: password,
			},
			buildStubs: func(
				services *mockdb.MockServices,
				tokenMaker *mocktokenmaker.MockMaker,
				req requests.LoginRequest,
			) {
				services.EXPECT().
					GetUser(gomock.Eq(req.Username)).
					Times(1).
					Return(randomUser, nil)
				tokenMaker.EXPECT().
					CreateToken(req.Username, configs.TokenAccessTokenDuration).
					Times(1).
					Return(accessToken, accessTokenPayload, nil)
				tokenMaker.EXPECT().
					CreateToken(req.Username, configs.TokenRefreshTokenDuration).
					Times(1).
					Return(refreshToken, refreshTokenPayload, nil)
				services.EXPECT().CreateSession(newSessionMatcher(session)).Return(session, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchLogin(t, recorder.Body, randomUser, configs.TokenSymmetricKey, refreshToken, accessToken)
			},
		},
		{
			name: "BadRequest",
			req: requests.LoginRequest{
				Username: util.RandomEmail(), // emails have @ and are invalid Usernames
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, tokenMaker *mocktokenmaker.MockMaker, req requests.LoginRequest) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			req: requests.LoginRequest{
				Username: randomUser.Username,
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, tokenMaker *mocktokenmaker.MockMaker, req requests.LoginRequest) {
				services.EXPECT().
					GetUser(gomock.Eq(req.Username)).
					Times(1).
					Return(models.User{}, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "DBInternalServerError",
			req: requests.LoginRequest{
				Username: randomUser.Username,
				Password: password,
			},
			buildStubs: func(services *mockdb.MockServices, tokenMaker *mocktokenmaker.MockMaker, req requests.LoginRequest) {
				services.EXPECT().
					GetUser(gomock.Eq(req.Username)).
					Times(1).
					Return(models.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "WrongPasswordUnAuthorized",
			req: requests.LoginRequest{
				Username: randomUser.Username,
				Password: "Wrong_Password_1",
			},
			buildStubs: func(
				services *mockdb.MockServices,
				tokenMaker *mocktokenmaker.MockMaker,
				req requests.LoginRequest,
			) {
				services.EXPECT().
					GetUser(gomock.Eq(req.Username)).
					Times(1).
					Return(randomUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "CreateAccessTokenInternalServerError",
			req: requests.LoginRequest{
				Username: randomUser.Username,
				Password: password,
			},
			buildStubs: func(
				services *mockdb.MockServices,
				tokenMaker *mocktokenmaker.MockMaker,
				req requests.LoginRequest,
			) {
				services.EXPECT().
					GetUser(gomock.Eq(req.Username)).
					Times(1).
					Return(randomUser, nil)
				tokenMaker.EXPECT().
					CreateToken(req.Username, configs.TokenAccessTokenDuration).
					Times(1).
					Return("", &token.Payload{}, errors.New("failed to encode payload to []byte"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "CreateRefreshTokenInternalServerError",
			req: requests.LoginRequest{
				Username: randomUser.Username,
				Password: password,
			},
			buildStubs: func(
				services *mockdb.MockServices,
				tokenMaker *mocktokenmaker.MockMaker,
				req requests.LoginRequest,
			) {
				services.EXPECT().
					GetUser(gomock.Eq(req.Username)).
					Times(1).
					Return(randomUser, nil)
				tokenMaker.EXPECT().
					CreateToken(req.Username, configs.TokenAccessTokenDuration).
					Times(1).
					Return(accessToken, accessTokenPayload, nil)
				tokenMaker.EXPECT().
					CreateToken(req.Username, configs.TokenRefreshTokenDuration).
					Times(1).
					Return("", &token.Payload{}, errors.New("failed to encode payload to []byte"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "CreateSessionInternalServerError",
			req: requests.LoginRequest{
				Username: randomUser.Username,
				Password: password,
			},
			buildStubs: func(
				services *mockdb.MockServices,
				tokenMaker *mocktokenmaker.MockMaker,
				req requests.LoginRequest,
			) {
				services.EXPECT().
					GetUser(gomock.Eq(req.Username)).
					Times(1).
					Return(randomUser, nil)
				tokenMaker.EXPECT().
					CreateToken(req.Username, configs.TokenAccessTokenDuration).
					Times(1).
					Return(accessToken, accessTokenPayload, nil)
				tokenMaker.EXPECT().
					CreateToken(req.Username, configs.TokenRefreshTokenDuration).
					Times(1).
					Return(refreshToken, refreshTokenPayload, nil)
				services.EXPECT().CreateSession(newSessionMatcher(session)).Return(models.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, tokenMaker token.Maker) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			services := mockdb.NewMockServices(ctrl)
			tokenMaker := mocktokenmaker.NewMockMaker(ctrl)

			accessToken, accessTokenPayload = createToken(t, randomUser.Username, configs.TokenAccessTokenDuration, configs.TokenSymmetricKey)
			refreshToken, refreshTokenPayload = createToken(t, randomUser.Username, configs.TokenRefreshTokenDuration, configs.TokenSymmetricKey)

			session = createSession(refreshTokenPayload.ID, randomUser.Username, refreshToken)

			testCase.buildStubs(services, tokenMaker, testCase.req)

			server := NewTestServer(t, services, tokenMaker)

			jsonReq, err := json.Marshal(&testCase.req)
			require.NoError(t, err)

			httpReq, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(jsonReq))
			httpReq.Header.Set("User-Agent", "test agent")
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.RouterServeHTTP(recorder, httpReq)

			testCase.checkResponse(t, recorder, server.handlers.tokenMaker)
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

func createToken(t *testing.T, username string, duration time.Duration, tokenSymmetricKey string) (string, *token.Payload) {
	tokenMaker, err := token.NewPasetoMaker(tokenSymmetricKey)
	require.NoError(t, err)

	accessToken, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, payload)

	require.NoError(t, payload.Valid())
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiredAt, time.Second)
	require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)

	return accessToken, payload
}

func createSession(sessionID uuid.UUID, username string, refreshToken string) models.Session {
	return models.Session{
		ID:           sessionID,
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    "test agent",
		ClientIP:     "",
		IsBlocked:    false,
		CreatedAt:    time.Now().UTC().Truncate(time.Second),
		ExpiresAt:    time.Now().UTC().Truncate(time.Second),
		DeletedAt:    gorm.DeletedAt{},
	}
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
	require.Equal(t, user.CreatedAt.Local().Truncate(time.Second), response.CreatedAt.Local())
	require.Equal(t, user.UpdatedAt.Local().Truncate(time.Second), response.UpdatedAt.Local())
	require.Equal(t, user.DeletedAt.Time.Truncate(time.Second), response.DeletedAt)
}

func requireBodyMatchLogin(
	t *testing.T,
	body *bytes.Buffer,
	user models.User,
	tokenSymmetricKey string,
	originalRefreshToken string,
	originalAccessToken string,
) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var loginResponse responses.LoginResponse
	err = json.Unmarshal(data, &loginResponse)

	response := loginResponse.UserInformation
	require.NoError(t, err)
	require.Equal(t, user.Username, response.Username)
	require.Equal(t, user.FullName, response.FullName)
	require.Equal(t, user.Email, response.Email)
	require.Equal(t, user.CreatedAt.Local().Truncate(time.Second), response.CreatedAt.Local())
	require.Equal(t, user.UpdatedAt.Local().Truncate(time.Second), response.UpdatedAt.Local())
	require.Equal(t, user.DeletedAt.Time.Truncate(time.Second), response.DeletedAt)

	accessToken := loginResponse.AccessToken
	require.NotEmpty(t, accessToken)
	require.Equal(t, originalAccessToken, accessToken)

	tokenMaker, err := token.NewPasetoMaker(tokenSymmetricKey)
	require.NoError(t, err)

	accessTokenPayload, err := tokenMaker.VerifyToken(accessToken)
	require.NoError(t, err)
	require.NotEmpty(t, accessTokenPayload)

	refreshToken := loginResponse.RefreshToken
	require.NotEmpty(t, refreshToken)
	require.Equal(t, originalRefreshToken, refreshToken)

	refreshTokenPayload, err := tokenMaker.VerifyToken(refreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, refreshTokenPayload)

	require.NotEmpty(t, loginResponse.SessionID)
	require.Equal(t, refreshTokenPayload.ID, loginResponse.SessionID)
}
