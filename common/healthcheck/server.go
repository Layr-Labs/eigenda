package healthcheck

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct{}

// Watch implements grpc_health_v1.HealthServer.
func (*HealthServer) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	panic("unimplemented")
}

func (s *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	// If the server is healthy, return a response with status "SERVING".
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// RegisterHealthServer registers the HealthServer with the provided gRPC server.
func RegisterHealthServer(server *grpc.Server) {
	healthServer := &HealthServer{} // Initialize your health server implementation
	grpc_health_v1.RegisterHealthServer(server, healthServer)
}
