package api

import (
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"time"
)

func (handler *Handler) CreateUser(context *gin.Context) {
	var req requests.CreateUserRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := handler.services.CreateUser(req)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			fmt.Println(pgError.ConstraintName)
			switch pgError.ConstraintName {
			case "users_pkey", "users_email_key":
				context.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
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

	user, err := handler.services.GetUser(req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := responses.UserInformationResponse{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Local().Truncate(time.Second),
		UpdatedAt: user.CreatedAt.Local().Truncate(time.Second),
	}

	if user.DeletedAt.Time.IsZero() {
		res.DeletedAt = user.DeletedAt.Time.Truncate(time.Second)
	} else {
		res.DeletedAt = user.DeletedAt.Time.Local().Truncate(time.Second)
	}

	context.JSON(http.StatusOK, res)
}
