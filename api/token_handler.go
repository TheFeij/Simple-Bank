package api

import (
	"Simple-Bank/db/models"
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

func (handler *Handler) RenewAccessToken(context *gin.Context) {
	var req requests.RenewAccessTokenRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshTokenPayload, err := handler.tokenMaker.VerifyToken(req.RefreshToken)
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := handler.services.GetSession(refreshTokenPayload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := checkSession(session, req.RefreshToken, refreshTokenPayload); err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := handler.tokenMaker.CreateToken(
		session.Username,
		handler.config.TokenAccessTokenDuration,
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := responses.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	}
	context.JSON(http.StatusOK, response)
}

func checkSession(session models.Session, refreshToken string, refreshTokenPayload *token.Payload) error {
	if session.IsBlocked {
		return fmt.Errorf("session is blocked")
	}
	if session.Username != refreshTokenPayload.Username {
		return fmt.Errorf("incorrect session user")
	}
	if session.RefreshToken != refreshToken {
		return fmt.Errorf("mismatch session token")
	}
	if time.Now().After(session.ExpiresAt) {
		return fmt.Errorf("expired session")
	}

	return nil
}
