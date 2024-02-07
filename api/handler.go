package api

import (
	"Simple-Bank/db/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services services.Services
}

func New(services services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
