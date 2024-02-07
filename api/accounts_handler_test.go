package api

import (
	mockdb "Simple-Bank/db/mock"
	"Simple-Bank/responses"
	"Simple-Bank/util"
	"bytes"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mockdb.NewMockServices(ctrl)
	// build stubs
	services.EXPECT().
		GetAccount(gomock.Eq(account.AccountID)).
		Times(1).
		Return(account, responseError)

	// testing http server
	server := NewServer(services)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.AccountID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.RouterServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func randomAccount() responses.GetAccountResponse {
	return responses.GetAccountResponse{
		AccountID: uint64(uint(util.RandomInt(1, 1000))),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
		Owner:     util.RandomString(int(util.RandomInt(1, 50))),
		Balance:   uint64(util.RandomInt(0, 9999)),
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
