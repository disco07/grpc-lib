package server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Color constants
const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
)

// getStatusColor returns the color for a given status code
func getStatusColor(code codes.Code) string {
	switch {
	case code == codes.OK:
		return Green
	case code == codes.Canceled || code == codes.InvalidArgument || code == codes.NotFound || code == codes.AlreadyExists || code == codes.PermissionDenied:
		return Yellow
	case code == codes.Internal || code == codes.Unavailable || code == codes.DataLoss || code == codes.Unauthenticated:
		return Red
	default:
		return Red
	}
}

// LoggingInterceptor logs the details of each request and response
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		code := status.Code(err)

		// Get color based on the status code
		color := getStatusColor(code)

		// Log format
		log.Printf(
			"%sMethod: %s | Status: %v | Duration: %v | Error: %v%s",
			color,
			info.FullMethod,
			code,
			duration,
			err,
			Reset,
		)

		return resp, err
	}
}
