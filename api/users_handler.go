package api

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"Simple-Bank/token"
	"Simple-Bank/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (handler *Handler) CreateUser(context *gin.Context) {
	var req requests.CreateUserRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	newUser, err := handler.services.CreateUser(req)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.ConstraintName {
			case "users_pkey", "users_email_key":
				context.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userInformation := responses.UserInformationResponse{
		Username:  newUser.Username,
		Email:     newUser.Email,
		FullName:  newUser.FullName,
		CreatedAt: newUser.CreatedAt.Local().Truncate(time.Second),
		UpdatedAt: newUser.UpdatedAt.Local().Truncate(time.Second),
		DeletedAt: newUser.DeletedAt.Time.Truncate(time.Second),
	}

	accessToken, accessTokenPayload, err := handler.tokenMaker.CreateToken(
		userInformation.Username,
		handler.config.TokenAccessTokenDuration,
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	refreshToken, refreshTokenPayload, err := handler.tokenMaker.CreateToken(
		req.Username,
		handler.config.TokenRefreshTokenDuration)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session := models.Session{
		ID:           refreshTokenPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    context.Request.UserAgent(),
		ClientIP:     context.ClientIP(),
		IsBlocked:    false,
		CreatedAt:    time.Now().UTC(),
		ExpiresAt:    time.Now().UTC(),
		DeletedAt:    gorm.DeletedAt{},
	}
	session, err = handler.services.CreateSession(session)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := responses.LoginResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		SessionID:             session.ID,
		UserInformation:       userInformation,
	}
	context.JSON(http.StatusOK, res)
}

func (handler *Handler) GetUser(context *gin.Context) {
	var req requests.GetUserRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != req.Username {
		err := fmt.Errorf("users cannot see other user`s information")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := handler.services.GetUser(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

func (handler *Handler) Login(context *gin.Context) {
	var req requests.LoginRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := handler.services.GetUser(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userInformation := responses.UserInformationResponse{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Local().Truncate(time.Second),
		UpdatedAt: user.CreatedAt.Local().Truncate(time.Second),
	}

	accessToken, accessTokenPayload, err := handler.tokenMaker.CreateToken(
		req.Username,
		handler.config.TokenAccessTokenDuration,
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := handler.tokenMaker.CreateToken(
		req.Username,
		handler.config.TokenRefreshTokenDuration)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session := models.Session{
		ID:           refreshTokenPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    context.Request.UserAgent(),
		ClientIP:     context.ClientIP(),
		IsBlocked:    false,
		CreatedAt:    time.Now().UTC(),
		ExpiresAt:    time.Now().UTC(),
		DeletedAt:    gorm.DeletedAt{},
	}
	session, err = handler.services.CreateSession(session)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := responses.LoginResponse{
		UserInformation:       userInformation,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		SessionID:             session.ID,
	}
	context.JSON(http.StatusOK, response)
}
