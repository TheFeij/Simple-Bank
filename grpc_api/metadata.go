package grpc_api

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	userAgent string
	clientIP  string
}

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

func (server *GrpcServer) extractMetaData(context context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(context); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.userAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.userAgent = userAgents[0]
		}
		fmt.Printf("%+v\n", md)
		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.clientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(context); ok {
		mtdt.clientIP = p.Addr.String()
	}

	return mtdt
}
