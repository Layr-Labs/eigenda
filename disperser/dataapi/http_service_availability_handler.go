package dataapi

import (
	"context"
	"net/http"
)

// Simple struct with a Service Name and its HealthEndPt.
type HttpServiceAvailabilityCheck struct {
	ServiceName string
	HealthEndPt string
}

type HttpServiceAvailability struct{}

func (s *server) getServiceHealth(ctx context.Context, services []HttpServiceAvailabilityCheck) ([]*ServiceAvailability, error) {

	availabilityStatuses := make([]*ServiceAvailability, len(services))
	for i, service := range services {
		var availabilityStatus *ServiceAvailability
		s.logger.Info("checking service health", "service", service.ServiceName)

		resp, err := s.eigenDAHttpServiceChecker.CheckHealth(service.HealthEndPt)
		if err != nil {
			s.logger.Error("Error querying service health:", "err", err)
		}

		availabilityStatus = &ServiceAvailability{
			ServiceName:   service.ServiceName,
			ServiceStatus: resp,
		}
		availabilityStatuses[i] = availabilityStatus
	}
	return availabilityStatuses, nil
}

// ServiceAvailability represents the status of a service.
func (sa *HttpServiceAvailability) CheckHealth(endpt string) (string, error) {
	resp, err := http.Get(endpt)
	if err != nil {
		return "UNKNOWN", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "SERVING", nil
	}

	return "NOT_SERVING", nil
}
