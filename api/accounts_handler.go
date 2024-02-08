package api

import (
	"Simple-Bank/requests"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (handler *Handler) CreateAccount(context *gin.Context) {
	var req requests.CreateAccountRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := handler.services.CreateAccount(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, res)
}

func (handler *Handler) GetAccount(context *gin.Context) {
	var req requests.GetAccountRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := handler.services.GetAccount(req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, res)
}

func (handler *Handler) GetAccountsList(context *gin.Context) {
	var req requests.GetAccountsListRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := handler.services.ListAccounts(req.PageID, req.PageSize)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, res)
}

func (handler *Handler) CreateUser(context *gin.Context) {
	var req requests.CreateUserRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := handler.services.CreateUser(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, res)
}

func (handler *Handler) GetUser(context *gin.Context) {
	var req requests.GetUserRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := handler.services.GetUser(req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, res)
}
