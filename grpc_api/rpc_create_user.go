package grpc_api

import (
	"Simple-Bank/pb"
	"Simple-Bank/requests"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *GrpcServer) CreateUser(context context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := ValidateCreateUserRequest(req)
	if violations != nil {
		err := invalidArgumentError(violations)
		return nil, err
	}

	newUser, err := server.dbServices.CreateUser(requests.CreateUserRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		FullName: req.GetFullname(),
		Email:    req.GetEmail(),
	})
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.ConstraintName {
			case "users_pkey":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists")
			case "users_email_key":
				return nil, status.Errorf(codes.AlreadyExists, "email already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	response := &pb.CreateUserResponse{User: convert(newUser)}

	return response, nil
}
