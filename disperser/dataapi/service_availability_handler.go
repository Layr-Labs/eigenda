package dataapi

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	mutex sync.Mutex
)

type EigenDAServiceAvailabilityCheck struct {
	disperserHostName string
	churnerHostName   string
}

func NewEigenDAServiceHealthCheck(disperserHostName, churnerHostName string) EigenDAServiceChecker {
	return &EigenDAServiceAvailabilityCheck{
		disperserHostName: disperserHostName,
		churnerHostName:   churnerHostName,
	}
}

func (s *server) getServiceAvailability(ctx context.Context, services []string) ([]*ServiceAvailability, error) {
	if services == nil {
		return nil, fmt.Errorf("services cannot be nil")
	}

	availabilityStatuses := make([]*ServiceAvailability, len(services))

	for i, serviceName := range services {
		var availabilityStatus *ServiceAvailability
		s.logger.Info("checking service health", "service", serviceName)
		response, err := s.eigenDAServiceChecker.CheckHealth(ctx, serviceName)
		if err != nil {
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

// CheckServiceHealth matches the HealthCheckService interface
func (sac *EigenDAServiceAvailabilityCheck) CheckHealth(ctx context.Context, serviceName string) (*grpc_health_v1.HealthCheckResponse, error) {
	serviceName = strings.ToLower(serviceName) // Normalize service name to lower case.
	var serverAddr string

	switch serviceName {
	case "disperser":
		serverAddr = sac.disperserHostName
	case "churner":
		serverAddr = sac.churnerHostName
	default:
		return nil, fmt.Errorf("unsupported service: %s", serviceName)
	}

	//Create connection to the server
	conn, err := grpc.DialContext(ctx, serverAddr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	return client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
}
