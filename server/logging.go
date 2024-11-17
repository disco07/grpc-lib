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

// ANSI color codes for background
const (
	Reset    = "\033[0m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgRed    = "\033[41m"
	FgWhite  = "\033[97m"
)

// getStatusColor returns the background color for a given status code
func getStatusColor(code codes.Code) string {
	switch {
	case code == codes.OK:
		return BgGreen
	case code == codes.Canceled || code == codes.InvalidArgument || code == codes.NotFound || code == codes.AlreadyExists || code == codes.PermissionDenied:
		return BgYellow
	case code == codes.Internal || code == codes.Unavailable || code == codes.DataLoss || code == codes.Unauthenticated:
		return BgRed
	default:
		return BgRed
	}
}

// padStatus adds a space to the left and right of the status string
func padStatus(status string) string {
	return fmt.Sprintf(" %s ", status)
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

		// Get background color based on the status code
		color := getStatusColor(code)

		// Format the log message
		statusMessage := fmt.Sprintf("%s%s%s%s", color, FgWhite, padStatus(code.String()), Reset)
		logMessage := fmt.Sprintf(
			"[%s] Method: %s | Duration: %v | Status: %s | Error: %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			info.FullMethod,
			duration,
			statusMessage,
			err,
		)

		log.Println(logMessage)
		return resp, err
	}
}
