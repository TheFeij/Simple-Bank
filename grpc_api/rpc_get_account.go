package grpc_api

import (
	"Simple-Bank/db/services"
	"Simple-Bank/pb"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// GetAccount is a rpc to get information about an account
func (server *GrpcServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}

	violations := validateId(req.GetId())
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	account, err := server.dbServices.GetAccount(req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAccountNotFound):
			return nil, status.Errorf(codes.NotFound, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	if payload.Username != account.Owner {
		return nil, status.Errorf(codes.Unauthenticated, "account does not belong to this user")
	}

	response := &pb.GetAccountResponse{Account: &pb.Account{
		AccountId: account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		CreatedAt: timestamppb.New(account.CreatedAt.Truncate(time.Second).Local()),
		UpdatedAt: timestamppb.New(account.UpdatedAt.Truncate(time.Second).Local()),
	}}

	if account.DeletedAt.Time.IsZero() {
		response.Account.DeletedAt = timestamppb.New(account.DeletedAt.Time.Truncate(time.Second))
	} else {
		response.Account.DeletedAt = timestamppb.New(account.DeletedAt.Time.Local().Truncate(time.Second))
	}

	return response, nil
}
