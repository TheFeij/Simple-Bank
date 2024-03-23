package api

import (
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"Simple-Bank/token"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (handler *Handler) CreateAccount(context *gin.Context) {
	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)

	newAccount, err := handler.services.CreateAccount(authPayload.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := responses.CreateAccountResponse{
		AccountID: newAccount.ID,
		Owner:     newAccount.Owner,
		Balance:   newAccount.Balance,
		CreatedAt: newAccount.CreatedAt.Truncate(time.Second).Local(),
	}
	context.JSON(http.StatusOK, res)
}

func (handler *Handler) GetAccount(context *gin.Context) {
	var req requests.GetAccountRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := handler.services.GetAccount(req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	auhPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	if auhPayload.Username != account.Owner {
		err := fmt.Errorf("users cannot create account for other users")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res := responses.GetAccountResponse{
		AccountID: account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt.Truncate(time.Second).Local(),
		UpdatedAt: account.UpdatedAt.Truncate(time.Second).Local(),
	}

	if account.DeletedAt.Time.IsZero() {
		res.DeletedAt = account.DeletedAt.Time.Truncate(time.Second)
	} else {
		res.DeletedAt = account.DeletedAt.Time.Local().Truncate(time.Second)
	}

	context.JSON(http.StatusOK, res)
}

func (handler *Handler) GetAccountsList(context *gin.Context) {
	var req requests.GetAccountsListRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accounts, err := handler.services.ListAccounts(req.PageID, req.PageSize)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := responses.ListAccountsResponse{}
	for i := range accounts {
		account := responses.GetAccountResponse{
			AccountID: accounts[i].ID,
			Owner:     accounts[i].Owner,
			Balance:   accounts[i].Balance,
			CreatedAt: accounts[i].CreatedAt.Truncate(time.Second).Local(),
			UpdatedAt: accounts[i].UpdatedAt.Truncate(time.Second).Local(),
		}

		res.Accounts = append(res.Accounts, account)

		if accounts[i].DeletedAt.Time.IsZero() {
			res.Accounts[i].DeletedAt = accounts[i].DeletedAt.Time.Truncate(time.Second)
		} else {
			res.Accounts[i].DeletedAt = accounts[i].DeletedAt.Time.Local().Truncate(time.Second)
		}
	}
	context.JSON(http.StatusOK, res)
}
