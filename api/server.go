package api

import (
	"Simple-Bank/api/handlers"
	"Simple-Bank/db/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Server struct {
	router   *gin.Engine
	handlers *handlers.Handler
}

func NewServer(db *gorm.DB) Server {
	server := Server{
		router:   gin.Default(),
		handlers: handlers.New(services.New(db)),
	}

	server.router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Welcome to our bank"})
	})
	server.router.POST("/accounts", server.handlers.CreateAccount)
	server.router.GET("/accounts/:id", server.handlers.GetAccount)
	server.router.GET("/accounts", server.handlers.GetAccountsList)

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
