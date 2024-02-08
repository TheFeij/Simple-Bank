package api

import (
	"Simple-Bank/db/services"
	"Simple-Bank/requests"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"net/http/httptest"
)

type Server struct {
	router   *gin.Engine
	handlers *Handler
}

func NewServer(services services.Services) Server {
	server := Server{
		router:   gin.Default(),
		handlers: New(services),
	}

	registerCustomValidators()

	server.router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Welcome to our bank"})
	})
	server.router.POST("/accounts", server.handlers.CreateAccount)
	server.router.GET("/accounts/:id", server.handlers.GetAccount)
	server.router.GET("/accounts", server.handlers.GetAccountsList)
	server.router.POST("/users", server.handlers.CreateUser)
	server.router.GET("/users/:username", server.handlers.GetUser)

	return server
}

func registerCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("validUsername", requests.ValidUsername); err != nil {
			log.Fatal("could not register validUsername validator")
		}
		if err := v.RegisterValidation("validPassword", requests.ValidPassword); err != nil {
			log.Fatal("could not register validPassword validator")
		}
		if err := v.RegisterValidation("validFullname", requests.ValidFullname); err != nil {
			log.Fatal("could not register validFullname validator")
		}
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) RouterServeHTTP(recorder *httptest.ResponseRecorder, req *http.Request) {
	server.router.ServeHTTP(recorder, req)
}
