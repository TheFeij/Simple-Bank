package grpc_api

import (
	"Simple-Bank/token"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	// authorizationHeader is the name of the header containing authorization information, including access token.
	authorizationHeader = "authorization"
	// authorizationTypeBearer is the type of authorization bearer token.
	authorizationTypeBearer = "bearer"
)

// authorizeUser authorizes the user based on the access token provided in the context.
// It extracts the access token from the authorization header and verifies it using the token maker.
// It returns the token payload if the access token is valid, otherwise it returns an error.
func (server *GrpcServer) authorizeUser(ctx context.Context) (*token.Payload, error) {
	mtdt, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := mtdt.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid aithorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationTypeBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", accessToken)
	}

	return payload, nil
}
