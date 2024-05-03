package api

import (
	"Simple-Bank/config"
	"Simple-Bank/db/services"
	"Simple-Bank/token"
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

func NewServer(config *config.Config, services services.Services, tokenMaker token.Maker) (*Server, error) {
	server := &Server{
		router:   gin.Default(),
		handlers: New(services, tokenMaker, config),
	}

	registerCustomValidators()

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	server.router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Welcome to our bank"})
	})

	authRoutes := server.router.Group("/").Use(authMiddleWare(server.handlers.tokenMaker))
	authRoutes.POST("/accounts", server.handlers.CreateAccount)
	authRoutes.GET("/accounts/:id", server.handlers.GetAccount)
	authRoutes.GET("/accounts", server.handlers.GetAccountsList)
	authRoutes.POST("/accounts/transfer", server.handlers.Transfer)
	server.router.POST("/users", server.handlers.CreateUser)
	authRoutes.GET("/users/:username", server.handlers.GetUser)
	server.router.POST("/users/login", server.handlers.Login)
	server.router.POST("/tokens/renew_access_token", server.handlers.RenewAccessToken)
}

func registerCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("validUsername", ValidUsername); err != nil {
			log.Fatal("could not register validUsername validator")
		}
		if err := v.RegisterValidation("validPassword", ValidPassword); err != nil {
			log.Fatal("could not register validPassword validator")
		}
		if err := v.RegisterValidation("validFullname", ValidFullname); err != nil {
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
