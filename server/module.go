package server

import (
	"github.com/disco07/grpc-lib/healthcheck"
	"github.com/disco07/grpc-lib/protogen/go/health"
	"go.uber.org/fx"
)

var Module = fx.Options(
	healthcheck.Module,
	fx.Provide(
		newGPRCServer,
	),
	fx.Invoke(
		health.RegisterHealthServiceServer,
	),
)
