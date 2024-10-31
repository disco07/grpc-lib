package healthcheck

import (
	"context"
	"github.com/disco07/grpc-lib/protogen/go/health"
	"github.com/golang/protobuf/ptypes/empty"
)

type healthCheck struct {
	health.HealthServiceServer
}

func newHealthCheck() health.HealthServiceServer {
	return &healthCheck{}
}

func (h *healthCheck) Check(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
