package dataapi

import (
	"context"
	"fmt"
	"net/http"
)

// Service represents a simple struct with a service name and its URL.
type HttpServiceAvailabilityCheck struct {
	ServiceName string
	URL         string
}

func (s *server) getServiceHealth(ctx context.Context, services []HttpServiceAvailabilityCheck) ([]*ServiceAvailability, error) {

	availabilityStatuses := make([]*ServiceAvailability, len(services))
	for i, service := range services {
		var availabilityStatus *ServiceAvailability
		s.logger.Info("checking service health", "service", service.ServiceName)
		resp, err := http.Get(service.URL)
		if err != nil {
			s.logger.Error("Error querying service health:", "err", err)
			availabilityStatus := &ServiceAvailability{
				ServiceName:   service.ServiceName,
				ServiceStatus: "UNKNOWN",
			}
			availabilityStatuses[i] = availabilityStatus
			continue
		}
		defer resp.Body.Close()

		// Check if the HTTP status code is 200 OK, which typically indicates healthiness.
		// Adjust the logic if the service uses different conventions.
		if resp.StatusCode == http.StatusOK {
			availabilityStatus = &ServiceAvailability{
				ServiceName:   service.ServiceName,
				ServiceStatus: "SERVING",
			}
			s.logger.Info("Service healthy", "service", service.ServiceName)
			availabilityStatuses[i] = availabilityStatus
		} else {
			fmt.Printf("Service may not be healthy. Received status code: %d\n", resp.StatusCode)
			availabilityStatus = &ServiceAvailability{
				ServiceName:   service.ServiceName,
				ServiceStatus: "NOT_SERVING",
			}
			s.logger.Info("Service unhealthy", "service", service.ServiceName)
			availabilityStatuses[i] = availabilityStatus
		}
	}
	return availabilityStatuses, nil
}
