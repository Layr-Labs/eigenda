package healthcheck

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// RegisterHealthServer registers the generic gRPC health check server implementation
// with the given gRPC server.
func RegisterHealthServer(name string, server *grpc.Server) {
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	healthServer.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
}
