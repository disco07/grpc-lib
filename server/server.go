package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func metadataInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		slog.WarnContext(ctx, "No metadata found in context")
		return nil, status.Errorf(codes.InvalidArgument, "No metadata found in context")
	}

	slog.DebugContext(ctx, "Metadata found in context", slog.Any("metadata", md))

	return handler(ctx, req)
}

func newGPRCServer(lifecycle fx.Lifecycle, logger *slog.Logger, config GRPCConfigServer) grpc.ServiceRegistrar {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port()))
	if err != nil {
		logger.Warn(err.Error())
	}

	keepaliveOptions := grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    time.Minute,
		Timeout: 3 * time.Second,
	})

	keepaliveEnforcementOptions := grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             30 * time.Second,
		PermitWithoutStream: true,
	})

	server := grpc.NewServer(
		keepaliveOptions,
		keepaliveEnforcementOptions,
		grpc.UnaryInterceptor(metadataInterceptor),
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				reflection.Register(server)

				if err = server.Serve(listener); err != nil {
					logger.Warn(err.Error())
				}
			}()

			logger.Info(fmt.Sprintf("%s://%s", listener.Addr().Network(), listener.Addr().String()))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.GracefulStop()

			return nil
		},
	})

	return server
}
