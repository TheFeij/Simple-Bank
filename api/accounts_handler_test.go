package api

import (
	mockdb "Simple-Bank/db/mock"
	"Simple-Bank/responses"
	"Simple-Bank/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()
	var responseError error = nil

	testCases := []struct {
		name          string
		AccountID     int64
		buildStubs    func(services *mockdb.MockServices)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			AccountID: int64(account.AccountID),
			buildStubs: func(services *mockdb.MockServices) {
				services.EXPECT().
					GetAccount(gomock.Eq(account.AccountID)).
					Times(1).
					Return(account, responseError)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			AccountID: int64(account.AccountID),
			buildStubs: func(services *mockdb.MockServices) {
				services.EXPECT().
					GetAccount(gomock.Eq(account.AccountID)).
					Times(1).
					Return(responses.GetAccountResponse{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			AccountID: int64(account.AccountID),
			buildStubs: func(services *mockdb.MockServices) {
				services.EXPECT().
					GetAccount(gomock.Eq(account.AccountID)).
					Times(1).
					Return(responses.GetAccountResponse{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			AccountID: util.RandomInt(-1000, 0),
			buildStubs: func(services *mockdb.MockServices) {
				services.EXPECT().
					GetAccount(gomock.Any()).
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
			testCase.buildStubs(services)

			// testing http server
			server := NewServer(services)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", testCase.AccountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.RouterServeHTTP(recorder, req)
			testCase.checkResponse(t, recorder)
		})

	}

}

func randomAccount() responses.GetAccountResponse {
	return responses.GetAccountResponse{
		AccountID: util.RandomID(),
		CreatedAt: time.Now().Truncate(time.Nanosecond).UTC(),
		UpdatedAt: time.Now().Truncate(time.Nanosecond).UTC(),
		DeletedAt: gorm.DeletedAt{},
		Owner:     util.RandomUsername(),
		Balance:   util.RandomBalance(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account responses.GetAccountResponse) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var response responses.GetAccountResponse
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)
	require.Equal(t, account, response)
}
