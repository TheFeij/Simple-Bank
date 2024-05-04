package main

import (
	"Simple-Bank/api"
	"Simple-Bank/config"
	"Simple-Bank/db"
	"Simple-Bank/db/services"
	"Simple-Bank/grpc_api"
	"Simple-Bank/pb"
	"Simple-Bank/token"
	"context"
	"database/sql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
	"net"
	"net/http"
	"os"
)

func main() {
	// pretty logger for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	configs, err := config.LoadConfig("./config", "config")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load configs")
	}

	db.Init(configs.DatabaseSource)
	db := db.GetDB()
	DB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get *sql.DB object")
	}
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("cannot close db connection")
		}
	}(DB)

	tokenMaker, err := token.NewPasetoMaker(configs.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create token maker")
	}

	//runGinServer(configs, tokenMaker, db)
	go runGrpcGatewayServer(configs, tokenMaker, db)
	runGrpcServer(configs, tokenMaker, db)
}

func runGinServer(config config.Config, tokenMaker token.Maker, db *gorm.DB) {
	server, err := api.NewServer(&config, services.NewSQLServices(db), tokenMaker)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	if err = server.Start(config.HTTPServerHost + ":" + config.HTTPServerPort); err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}

func runGrpcServer(config config.Config, tokenMaker token.Maker, db *gorm.DB) {
	server := grpc_api.NewServer(&config, services.NewSQLServices(db), tokenMaker)

	grpcLogger := grpc.UnaryInterceptor(grpc_api.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	serverAddress := config.GrpcServerHost + ":" + config.GrpcServerPort
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msg("grpc server started at " + listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create grpc server")
	}
}

func runGrpcGatewayServer(config config.Config, tokenMaker token.Maker, db *gorm.DB) {
	server := grpc_api.NewServer(&config, services.NewSQLServices(db), tokenMaker)

	serveMuxOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(serveMuxOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	serverAddress := config.HTTPServerHost + ":" + config.HTTPServerPort
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msg("http gateway server started at " + listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}
