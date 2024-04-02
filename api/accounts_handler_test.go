package api

import (
	mockdb "Simple-Bank/db/mock"
	"Simple-Bank/db/models"
	"Simple-Bank/responses"
	"Simple-Bank/token"
	"Simple-Bank/util"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateAccount(t *testing.T) {
	randomUser, _ := randomUser(t)
	account := createAccount(randomUser.Username)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(services *mockdb.MockServices)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, request)
			},
			buildStubs: func(services *mockdb.MockServices) {
				services.EXPECT().
					CreateAccount(gomock.Eq(randomUser.Username)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "UnAuthorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(services *mockdb.MockServices) {
				services.EXPECT().
					CreateAccount(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			services := mockdb.NewMockServices(controller)
			testCase.buildStubs(services)

			server := NewTestServer(t, services)
			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(http.MethodPost, "/accounts", nil)
			require.NoError(t, err)

			testCase.setupAuth(t, httpReq, server.handlers.tokenMaker)

			server.RouterServeHTTP(recorder, httpReq)
			testCase.checkResponse(t, recorder)
		})
	}

}

func createAccount(owner string) models.Account {
	return models.Account{
		ID:        util.RandomInt(1, math.MaxInt64),
		Owner:     owner,
		Balance:   util.RandomBalance(),
		CreatedAt: time.Now().Truncate(time.Second).UTC(),
		UpdatedAt: time.Now().Truncate(time.Second).UTC(),
		DeletedAt: gorm.DeletedAt{},
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account models.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response responses.CreateAccountResponse
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)
	require.Equal(t, account.ID, response.AccountID)
	require.Equal(t, account.Owner, response.Owner)
	require.Equal(t, account.Balance, response.Balance)
	require.Equal(t, account.CreatedAt.Local().Truncate(time.Second), response.CreatedAt)
}
