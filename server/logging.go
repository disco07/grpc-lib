package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Blue   = "\033[34m"
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
		// Skip logging for specific methods
		if info.FullMethod == "/health.HealthService/Check" {
			return handler(ctx, req)
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		code := status.Code(err)

		// Get color based on the status code
		color := getStatusColor(code)

		// Construct log message
		statusMessage := fmt.Sprintf("%s%-8s%s", color, code.String(), Reset)
		logMessage := fmt.Sprintf(
			"%s[%s] %sMethod: %s%s\n%sDuration: %v | Status: %s | Error: %v%s\n",
			Blue, time.Now().Format("2006-01-02 15:04:05"), Reset,
			info.FullMethod,
			Reset,
			" ", duration, statusMessage, err, Reset,
		)

		log.Println(logMessage)
		return resp, err
	}
}
