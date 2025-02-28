package dataapi

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCConn interface {
	Dial(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
}

type GRPCDialerSkipTLS struct{}

type EigenDAServiceAvailabilityCheck struct {
	disperserConn *grpc.ClientConn
	churnerConn   *grpc.ClientConn
}

func (s *server) getServiceAvailability(ctx context.Context, services []string) ([]*ServiceAvailability, error) {
	if services == nil {
		return nil, fmt.Errorf("services cannot be nil")
	}

	availabilityStatuses := make([]*ServiceAvailability, len(services))

	for i, serviceName := range services {
		var availabilityStatus *ServiceAvailability
		s.logger.Info("checking service health", "service", serviceName)

		response, err := s.eigenDAGRPCServiceChecker.CheckHealth(ctx, serviceName)
		if err != nil {

			if err.Error() == "disperser connection is nil" {
				s.logger.Error("disperser connection is nil")
				availabilityStatus = &ServiceAvailability{
					ServiceName:   serviceName,
					ServiceStatus: grpc_health_v1.HealthCheckResponse_UNKNOWN.String(),
				}
				availabilityStatuses[i] = availabilityStatus
				continue
			}

			if err.Error() == "churner connection is nil" {
				s.logger.Error("churner connection is nil")
				availabilityStatus = &ServiceAvailability{
					ServiceName:   serviceName,
					ServiceStatus: grpc_health_v1.HealthCheckResponse_UNKNOWN.String(),
				}
				availabilityStatuses[i] = availabilityStatus
				continue
			}

			s.logger.Error("failed to check service health", "service", serviceName, "err", err)
			availabilityStatus = &ServiceAvailability{
				ServiceName:   serviceName,
				ServiceStatus: grpc_health_v1.HealthCheckResponse_NOT_SERVING.String(),
			}
			availabilityStatuses[i] = availabilityStatus
		} else {
			s.logger.Info("service status", "service", serviceName, "status", response.Status.String())
			availabilityStatus = &ServiceAvailability{
				ServiceName:   serviceName,
				ServiceStatus: response.Status.String(),
			}
			availabilityStatuses[i] = availabilityStatus
		}
	}
	return availabilityStatuses, nil
}

func NewEigenDAServiceHealthCheck(grpcConnection GRPCConn, disperserHostName, churnerHostName string) EigenDAGRPCServiceChecker {

	// Create Pre-configured connections to the services
	// Saves from having to create new connection on each request

	disperserConn, err := grpcConnection.Dial(disperserHostName, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))

	if err != nil {
		return nil
	}

	churnerConn, err := grpcConnection.Dial(churnerHostName, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))

	if err != nil {
		return nil
	}

	return &EigenDAServiceAvailabilityCheck{
		disperserConn: disperserConn,
		churnerConn:   churnerConn,
	}
}

// Create Connection to the service
func (rc *GRPCDialerSkipTLS) Dial(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// Create client options with timeout
	opts = append(opts, grpc.WithConnectParams(grpc.ConnectParams{
		MinConnectTimeout: 10 * time.Second,
	}))

	return grpc.NewClient(serviceName, opts...)
}

// CheckServiceHealth matches the HealthCheckService interface
func (sac *EigenDAServiceAvailabilityCheck) CheckHealth(ctx context.Context, serviceName string) (*grpc_health_v1.HealthCheckResponse, error) {
	serviceName = strings.ToLower(serviceName) // Normalize service name to lower case.

	var client grpc_health_v1.HealthClient

	switch serviceName {
	case "disperser":

		if sac.disperserConn == nil {
			return nil, fmt.Errorf("disperser connection is nil")
		}
		client = grpc_health_v1.NewHealthClient(sac.disperserConn)
	case "churner":

		if sac.churnerConn == nil {
			return nil, fmt.Errorf("churner connection is nil")
		}
		client = grpc_health_v1.NewHealthClient(sac.churnerConn)
	default:
		return nil, fmt.Errorf("unsupported service: %s", serviceName)
	}

	return client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
}

// Close Open connections
func (sac *EigenDAServiceAvailabilityCheck) CloseConnections() error {
	if sac.disperserConn != nil {
		sac.disperserConn.Close()
	}
	if sac.churnerConn != nil {
		sac.churnerConn.Close()
	}

	return nil
}
