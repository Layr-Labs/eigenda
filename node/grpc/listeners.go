package grpc

import (
	"fmt"
	"net"
)

// V1Listeners holds the network listeners for v1 gRPC servers.
type V1Listeners struct {
	Dispersal net.Listener
	Retrieval net.Listener
}

// Close closes both listeners.
func (l *V1Listeners) Close() {
	if l.Dispersal != nil {
		_ = l.Dispersal.Close()
	}
	if l.Retrieval != nil {
		_ = l.Retrieval.Close()
	}
}

// CreateV1Listeners creates network listeners for v1 gRPC servers.
// On error, any successfully created listeners are closed before returning.
func CreateV1Listeners(dispersalAddr, retrievalAddr string) (V1Listeners, error) {
	var listeners V1Listeners

	var err error
	listeners.Dispersal, err = net.Listen("tcp", dispersalAddr)
	if err != nil {
		return listeners, fmt.Errorf("failed to create v1 dispersal listener: %w", err)
	}

	listeners.Retrieval, err = net.Listen("tcp", retrievalAddr)
	if err != nil {
		listeners.Close()
		return listeners, fmt.Errorf("failed to create v1 retrieval listener: %w", err)
	}

	return listeners, nil
}

// V2Listeners holds the network listeners for v2 gRPC servers.
type V2Listeners struct {
	Dispersal net.Listener
	Retrieval net.Listener
}

// Close closes both listeners.
func (l *V2Listeners) Close() {
	if l.Dispersal != nil {
		_ = l.Dispersal.Close()
	}
	if l.Retrieval != nil {
		_ = l.Retrieval.Close()
	}
}

// CreateV2Listeners creates network listeners for v2 gRPC servers.
// On error, any successfully created listeners are closed before returning.
func CreateV2Listeners(dispersalAddr, retrievalAddr string) (V2Listeners, error) {
	var listeners V2Listeners

	var err error
	listeners.Dispersal, err = net.Listen("tcp", dispersalAddr)
	if err != nil {
		return listeners, fmt.Errorf("failed to create v2 dispersal listener: %w", err)
	}

	listeners.Retrieval, err = net.Listen("tcp", retrievalAddr)
	if err != nil {
		listeners.Close()
		return listeners, fmt.Errorf("failed to create v2 retrieval listener: %w", err)
	}

	return listeners, nil
}
