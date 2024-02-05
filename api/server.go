package api

import (
	"Simple-Bank/db/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Server struct {
	router   *gin.Engine
	services *services.Services
}

func NewServer(db *gorm.DB) Server {
	server := Server{
		router:   gin.Default(),
		services: services.New(db),
	}

	server.router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Welcome to our bank"})
	})

	return server
}
