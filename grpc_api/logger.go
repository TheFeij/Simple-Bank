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

// GrpcLogger logs information about gRPC requests and responses.
// It logs details such as protocol, method, duration, status code, and status text.
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

// ResponseRecorder implements http.ResponseWriter and records information about HTTP responses.
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

// WriteHeader records the status code of the HTTP response.
func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// Write records the body of the HTTP response into the ResponseRecorder.Body.
func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

// HttpLogger logs information about HTTP requests and responses.
// It logs details such as protocol, method, path, status code, status text, duration, and response body (in case of non-200 status codes).
func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.
			Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}
