package grpc_api

import (
	"Simple-Bank/db/models"
	"Simple-Bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func convert(user models.User) *pb.User {
	return &pb.User{
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  user.FullName,
		CreatedAt: timestamppb.New(user.CreatedAt.Local().Truncate(time.Second)),
		UpdatedAt: timestamppb.New(user.UpdatedAt.Local().Truncate(time.Second)),
	}
}
