package healthcheck

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// RegisterHealthServer registers the default gRPC health check server implementation
// with the given gRPC server.
func RegisterHealthServer(name string, server *grpc.Server) {
	healthServer := health.NewServer()
	healthServer.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)
}
