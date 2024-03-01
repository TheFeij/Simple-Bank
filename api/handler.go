package api

import (
	"Simple-Bank/config"
	"Simple-Bank/db/services"
	"Simple-Bank/token"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services   services.Services
	tokenMaker token.Maker
	config     *config.Config
}

func New(services services.Services, tokenMaker token.Maker, config *config.Config) *Handler {
	return &Handler{
		services:   services,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
