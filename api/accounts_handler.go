package api

import (
	"Simple-Bank/db/services"
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"Simple-Bank/token"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	auhPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)

	accounts, err := handler.services.ListAccounts(services.ListAccountsRequest{
		Owner:      auhPayload.Username,
		PageSize:   int(req.PageSize),
		PageNumber: int(req.PageID),
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

func (handler *Handler) Transfer(context *gin.Context) {
	var req requests.TransferRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)

	transfer, err := handler.services.Transfer(services.TransferRequest{
		Owner:         authPayload.Username,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, responses.TransferResponse{
		TransferID:      transfer.ID,
		SrcAccountID:    transfer.FromAccountID,
		DstAccountID:    transfer.ToAccountID,
		Amount:          transfer.Amount,
		CreatedAt:       transfer.CreatedAt.Local(),
		IncomingEntryID: transfer.IncomingEntryID,
		OutgoingEntryID: transfer.OutgoingEntryID,
	})
}
