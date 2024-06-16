package grpc_api

import (
	"Simple-Bank/pb"
	"Simple-Bank/util"
	"errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := util.ValidateFullname(req.GetFullname()); err != nil {
		violations = append(violations, fieldViolation("fullname", err))
	}
	if err := util.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}

func validateLoginUseRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.Password != nil {
		if err := util.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if req.Fullname != nil {
		if err := util.ValidateFullname(req.GetFullname()); err != nil {
			violations = append(violations, fieldViolation("fullname", err))
		}
	}

	if req.Email != nil {
		if err := util.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}

// validateTransferRequest validates transfer request's params
func validateTransferRequest(req *pb.TransferRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.FromAccountId < 1 {
		err := errors.New("invalid id, should be equal or higher than 1")
		violations = append(violations, fieldViolation("from_account_id", err))
	}
	if req.ToAccountId < 1 {
		err := errors.New("invalid id, should be equal or higher than 1")
		violations = append(violations, fieldViolation("to_account_id", err))
	}
	if req.ToAccountId == req.FromAccountId {
		err := errors.New("cannot transfer money from an account to itself")
		violations = append(violations, fieldViolation("to_account_id", err))
	}
	if req.Amount < 1 {
		err := errors.New("invalid amount, amount should be a positive integer")
		violations = append(violations, fieldViolation("amount", err))
	}

	return violations
}
