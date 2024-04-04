package api

import (
	mockdb "Simple-Bank/db/mock"
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"Simple-Bank/token"
	"Simple-Bank/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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
		{
			name: "InternalServerError",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, request)
			},
			buildStubs: func(services *mockdb.MockServices) {
				services.EXPECT().
					CreateAccount(gomock.Eq(randomUser.Username)).
					Times(1).
					Return(account, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

func TestGetAccount(t *testing.T) {
	randomUser, _ := randomUser(t)
	account := createAccount(randomUser.Username)

	testCases := []struct {
		name          string
		req           requests.GetAccountRequest
		setupAuth     func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker)
		buildStubs    func(services *mockdb.MockServices, req requests.GetAccountRequest)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			req: requests.GetAccountRequest{
				ID: account.ID,
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, httpReq)
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountRequest) {
				services.EXPECT().GetAccount(gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "UnAuthorized",
			req: requests.GetAccountRequest{
				ID: account.ID,
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountRequest) {
				services.EXPECT().GetAccount(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			req: requests.GetAccountRequest{
				ID: -account.ID,
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, httpReq)
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountRequest) {
				services.EXPECT().GetAccount(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			req: requests.GetAccountRequest{
				ID: account.ID,
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, httpReq)
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountRequest) {
				services.EXPECT().GetAccount(gomock.Eq(account.ID)).Times(1).Return(models.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			req: requests.GetAccountRequest{
				ID: account.ID,
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, httpReq)
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountRequest) {
				services.EXPECT().GetAccount(gomock.Eq(account.ID)).Times(1).Return(models.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			services := mockdb.NewMockServices(controller)
			testCase.buildStubs(services, testCase.req)

			server := NewTestServer(t, services)
			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", testCase.req.ID), nil)
			require.NoError(t, err)

			testCase.setupAuth(t, httpReq, server.handlers.tokenMaker)

			server.RouterServeHTTP(recorder, httpReq)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestGetAccountsList(t *testing.T) {
	randomUser, _ := randomUser(t)
	numberOfAccounts := 20
	accounts := make([]models.Account, numberOfAccounts)
	for i := 0; i < numberOfAccounts; i++ {
		accounts[i] = createAccount(randomUser.Username)
	}

	RandomPageSize := util.RandomInt(5, 10)
	RandomPageID := util.RandomInt(1, int64(math.Ceil(float64(numberOfAccounts)/float64(RandomPageSize))))

	startIndex := RandomPageSize * (RandomPageID - 1)
	endIndex := RandomPageSize*(RandomPageID-1) + RandomPageSize

	if endIndex > int64(len(accounts)) {
		endIndex = int64(len(accounts))
	}

	testCases := []struct {
		name          string
		req           requests.GetAccountsListRequest
		setupAuth     func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker)
		buildStubs    func(services *mockdb.MockServices, req requests.GetAccountsListRequest)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			req: requests.GetAccountsListRequest{
				PageID:   RandomPageID,
				PageSize: int8(RandomPageSize),
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, httpReq)
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountsListRequest) {
				services.EXPECT().
					ListAccounts(gomock.Eq(randomUser.Username), gomock.Eq(req.PageID), gomock.Eq(req.PageSize)).
					Times(1).Return(accounts[startIndex:endIndex], nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchAccountsList(t, recorder.Body, accounts[startIndex:endIndex])
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UnAuthorized",
			req: requests.GetAccountsListRequest{
				PageID:   RandomPageID,
				PageSize: int8(RandomPageSize),
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountsListRequest) {
				services.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			req: requests.GetAccountsListRequest{
				PageID:   -RandomPageID,
				PageSize: -int8(RandomPageSize),
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, httpReq)
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountsListRequest) {
				services.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			req: requests.GetAccountsListRequest{
				PageID:   RandomPageID,
				PageSize: int8(RandomPageSize),
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, httpReq)
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountsListRequest) {
				services.EXPECT().
					ListAccounts(gomock.Eq(randomUser.Username), gomock.Eq(req.PageID), gomock.Eq(req.PageSize)).
					Times(1).Return([]models.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			req: requests.GetAccountsListRequest{
				PageID:   RandomPageID,
				PageSize: int8(RandomPageSize),
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute, httpReq)
			},
			buildStubs: func(services *mockdb.MockServices, req requests.GetAccountsListRequest) {
				services.EXPECT().
					ListAccounts(gomock.Eq(randomUser.Username), gomock.Eq(req.PageID), gomock.Eq(req.PageSize)).
					Times(1).Return([]models.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			services := mockdb.NewMockServices(controller)
			testCase.buildStubs(services, testCase.req)

			server := NewTestServer(t, services)

			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", testCase.req.PageID, testCase.req.PageSize)
			httpReq, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			testCase.setupAuth(t, httpReq, server.handlers.tokenMaker)
			server.RouterServeHTTP(recorder, httpReq)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestTransfer(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	account1 := createAccount(user1.Username)
	account2 := createAccount(user2.Username)

	amount := int32(util.RandomInt(1, math.MaxInt32))

	transfer := models.Transfer{
		ID:              util.RandomID(),
		FromAccountID:   account1.ID,
		ToAccountID:     account2.ID,
		Amount:          amount,
		OutgoingEntryID: util.RandomID(),
		IncomingEntryID: util.RandomID(),
		CreatedAt:       time.Now().UTC(),
	}

	testCases := []struct {
		name          string
		req           requests.TransferRequest
		setupAuth     func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker)
		buildStubs    func(services *mockdb.MockServices, req requests.TransferRequest)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			req: requests.TransferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute, httpReq)
			},
			buildStubs: func(services *mockdb.MockServices, req requests.TransferRequest) {
				services.EXPECT().Transfer(gomock.Eq(user1.Username), gomock.Eq(req)).Return(transfer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				fmt.Println(recorder.Code)
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransfer(t, recorder.Body, transfer)
			},
		},
		{
			name: "UnAuthorized",
			req: requests.TransferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			},
			setupAuth: func(t *testing.T, httpReq *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(services *mockdb.MockServices, req requests.TransferRequest) {
				services.EXPECT().Transfer(gomock.Any(), gomock.Any()).Times(0)
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
			testCase.buildStubs(services, testCase.req)

			server := NewTestServer(t, services)

			jsonReq, err := json.Marshal(&testCase.req)
			require.NoError(t, err)

			httpReq, err := http.NewRequest(http.MethodPost, "/accounts/transfer", bytes.NewBuffer(jsonReq))
			require.NoError(t, err)

			testCase.setupAuth(t, httpReq, server.handlers.tokenMaker)

			recorder := httptest.NewRecorder()
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

func requireBodyMatchAccountsList(t *testing.T, body *bytes.Buffer, accounts []models.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var responseList responses.ListAccountsResponse
	err = json.Unmarshal(data, &responseList)
	require.NoError(t, err)

	for i, account := range accounts {
		response := responseList.Accounts[i]
		require.Equal(t, account.ID, response.AccountID)
		require.Equal(t, account.Owner, response.Owner)
		require.Equal(t, account.Balance, response.Balance)
		require.Equal(t, account.CreatedAt.Local().Truncate(time.Second), response.CreatedAt)
	}
}

func requireBodyMatchTransfer(t *testing.T, body *bytes.Buffer, transfer models.Transfer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response responses.TransferResponse
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)
	require.Equal(t, transfer.ID, response.TransferID)
	require.Equal(t, transfer.FromAccountID, response.SrcAccountID)
	require.Equal(t, transfer.ToAccountID, response.DstAccountID)
	require.Equal(t, transfer.Amount, response.Amount)
	require.Equal(t, transfer.IncomingEntryID, response.IncomingEntryID)
	require.Equal(t, transfer.OutgoingEntryID, response.OutgoingEntryID)
	require.Equal(t, transfer.CreatedAt.Local().Truncate(time.Second), response.CreatedAt.Truncate(time.Second))
}
