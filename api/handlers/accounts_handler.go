package handlers

import (
	"Simple-Bank/requests"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (handler *Handler) CreatAccount(context *gin.Context) {
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
