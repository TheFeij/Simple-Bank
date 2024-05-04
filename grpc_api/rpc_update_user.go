package grpc_api

import (
	"Simple-Bank/db/services"
	"Simple-Bank/pb"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *GrpcServer) UpdateUser(context context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		err := invalidArgumentError(violations)
		return nil, err
	}

	updatedUser, err := server.dbServices.UpdateUser(services.UpdateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Fullname: req.Fullname,
		Email:    req.Email,
	})
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.ConstraintName {
			case "users_email_key":
				return nil, status.Errorf(codes.AlreadyExists, "email already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to update user")
	}

	response := &pb.UpdateUserResponse{User: convert(updatedUser)}

	return response, nil
}
