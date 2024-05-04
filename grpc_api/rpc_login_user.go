package grpc_api

import (
	"Simple-Bank/db/models"
	"Simple-Bank/pb"
	"Simple-Bank/util"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

func (server *GrpcServer) LoginUser(context context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUseRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.dbServices.GetUser(req.GetUsername())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(
		req.GetUsername(),
		server.config.TokenAccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		req.Username,
		server.config.TokenRefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	metadata := server.extractMetaData(context)
	session := models.Session{
		ID:           refreshTokenPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.userAgent,
		ClientIP:     metadata.clientIP,
		IsBlocked:    false,
		CreatedAt:    time.Now().UTC(),
		ExpiresAt:    time.Now().UTC(),
		DeletedAt:    gorm.DeletedAt{},
	}
	session, err = server.dbServices.CreateSession(session)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	response := &pb.LoginUserResponse{
		User:                  convert(user),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		SessionId:             session.ID.String(),
	}

	return response, nil
}
