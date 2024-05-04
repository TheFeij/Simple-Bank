package grpc_api

import (
	"Simple-Bank/config"
	"Simple-Bank/db/services"
	"Simple-Bank/pb"
	"Simple-Bank/token"
)

// GrpcServer serves grpc requests for the banking service.
type GrpcServer struct {
	pb.UnimplementedSimpleBankServer
	dbServices services.Services
	tokenMaker token.Maker
	config     *config.Config
}

// NewServer creates a new grpc server.
func NewServer(config *config.Config, services services.Services, tokenMaker token.Maker) *GrpcServer {
	return &GrpcServer{
		tokenMaker: tokenMaker,
		config:     config,
		dbServices: services,
	}
}
