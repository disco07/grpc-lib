package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

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
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				reflection.Register(server)

				if err = server.Serve(listener); err != nil {
					logger.Warn(err.Error())
				}
			}()
			logger.With(fmt.Sprintf("%s://%s", listener.Addr().Network(), listener.Addr().String())).Info("Listening")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.GracefulStop()

			return nil
		},
	})

	return server
}
