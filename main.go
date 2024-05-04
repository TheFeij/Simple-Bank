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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
)

func main() {
	configs, err := config.LoadConfig("./config", "config")
	if err != nil {
		log.Fatalln(err)
	}

	db.Init(configs.DatabaseSource)
	db := db.GetDB()
	DB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(DB)

	tokenMaker, err := token.NewPasetoMaker(configs.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("cannot create token maker: %v", err)
	}

	//runGinServer(configs, tokenMaker, db)
	go runGrpcGatewayServer(configs, tokenMaker, db)
	runGrpcServer(configs, tokenMaker, db)
}

func runGinServer(config config.Config, tokenMaker token.Maker, db *gorm.DB) {
	server, err := api.NewServer(&config, services.NewSQLServices(db), tokenMaker)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	if err = server.Start(config.HTTPServerHost + ":" + config.HTTPServerPort); err != nil {
		log.Fatalln("cannot start server: ", err)
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
		log.Fatalln("cannot create listener: ", err)
	}

	log.Println("grpc server started at " + listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln("cannot create grpc server ", err)
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
		log.Fatalln("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	serverAddress := config.HTTPServerHost + ":" + config.HTTPServerPort
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalln("cannot create listener: ", err)
	}

	log.Println("http gateway server started at " + listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalln("cannot start HTTP gateway server ", err)
	}
}
