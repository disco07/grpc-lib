package client

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		newGRPCClientConn,
		newServeMux,
	),
	fx.Invoke(
		startHTTPClient,
	),
)
