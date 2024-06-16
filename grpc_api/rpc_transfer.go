package grpc_api

import (
	"Simple-Bank/db/services"
	"Simple-Bank/pb"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Transfer is a rpc to transfer money from one account to another
func (server *GrpcServer) Transfer(context context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	// authorize user and get the authorization payload
	payload, err := server.authorizeUser(context)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}

	// validate user request params
	violations := validateTransferRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// transfer money
	transfer, err := server.dbServices.Transfer(services.TransferRequest{
		Owner:         payload.Username,
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	})
	if err != nil {
		// TODO: better error handling required, to be implemented
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	// return the successful transfer response
	return &pb.TransferResponse{
		TransferId:      transfer.ID,
		SrcAccountId:    transfer.FromAccountID,
		DstAccountId:    transfer.ToAccountID,
		IncomingEntryId: transfer.IncomingEntryID,
		OutgoingEntryId: transfer.OutgoingEntryID,
		CreateAt:        timestamppb.New(transfer.CreatedAt.UTC().Truncate(time.Second)),
		Amount:          transfer.Amount,
	}, nil
}
