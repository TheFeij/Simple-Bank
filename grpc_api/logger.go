package grpc_api

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func GrpcLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {

	start := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(start)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.
		Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Dur("duration", duration).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Msg("received a gRPC request")

	return resp, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	Body       []byte
}

func (recorder *ResponseRecorder) WriteHeader(statusCode int) {
	recorder.statusCode = statusCode
	recorder.ResponseWriter.WriteHeader(statusCode)
}

func (recorder *ResponseRecorder) Write(body []byte) (int, error) {
	recorder.Body = body
	return recorder.ResponseWriter.Write(body)
}

func HttpLogger(
	handler http.Handler,
) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, req *http.Request) {
		recorder := &ResponseRecorder{
			ResponseWriter: response,
			statusCode:     http.StatusOK,
		}

		start := time.Now()
		handler.ServeHTTP(recorder, req)
		duration := time.Since(start)

		logger := log.Info()
		if recorder.statusCode != http.StatusOK {
			logger = log.Error().Bytes("body", recorder.Body)
		}

		logger.
			Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Dur("duration", duration).
			Int("status_code", recorder.statusCode).
			Str("status_text", http.StatusText(recorder.statusCode)).
			Msg("received an HTTP request")
	})
}
