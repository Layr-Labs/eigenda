package grpc

import (
	"fmt"
	"net"
)

// Listeners holds the network listeners for gRPC servers.
type Listeners struct {
	Dispersal net.Listener
	Retrieval net.Listener
}

// Close closes both listeners.
func (l *Listeners) Close() {
	if l.Dispersal != nil {
		_ = l.Dispersal.Close()
	}
	if l.Retrieval != nil {
		_ = l.Retrieval.Close()
	}
}

// CreateListeners creates network listeners for gRPC servers.
// Ports should be specified as strings (e.g., "32003" or "0" for auto-assignment).
// On error, any successfully created listeners are closed before returning.
func CreateListeners(dispersalPort, retrievalPort string) (Listeners, error) {
	var listeners Listeners

	dispersalAddr := fmt.Sprintf("0.0.0.0:%s", dispersalPort)
	retrievalAddr := fmt.Sprintf("0.0.0.0:%s", retrievalPort)

	var err error
	listeners.Dispersal, err = net.Listen("tcp", dispersalAddr)
	if err != nil {
		return listeners, fmt.Errorf("failed to create dispersal listener: %w", err)
	}

	listeners.Retrieval, err = net.Listen("tcp", retrievalAddr)
	if err != nil {
		listeners.Close()
		return listeners, fmt.Errorf("failed to create retrieval listener: %w", err)
	}

	return listeners, nil
}
