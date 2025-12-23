package common

import (
	"fmt"
	"time"
)

// Contains configuration for a gRPC server
type GRPCServerConfig struct {
	// Port that the gRPC server listens on
	GrpcPort uint16

	// Maximum size of a gRPC message that the server will accept (in bytes)
	MaxGRPCMessageSize int

	// Maximum time a connection can be idle before it is closed.
	MaxIdleConnectionAge time.Duration

	// Maximum age of a request in the past that the server will accept.
	// Requests older than this will be rejected to prevent replay attacks.
	RequestMaxPastAge time.Duration

	// Maximum age of a request in the future that the server will accept.
	// Requests with timestamps too far in the future will be rejected.
	RequestMaxFutureAge time.Duration
}

// DefaultGRPCServerConfig returns the default gRPC server configuration.
func DefaultGRPCServerConfig() GRPCServerConfig {
	return GRPCServerConfig{
		GrpcPort:             32010,
		MaxGRPCMessageSize:   1024 * 1024, // 1 MB
		MaxIdleConnectionAge: 5 * time.Minute,
		RequestMaxPastAge:    5 * time.Minute,
		RequestMaxFutureAge:  3 * time.Minute,
	}
}

func (c *GRPCServerConfig) Verify() error {
	if c.MaxGRPCMessageSize <= 0 {
		return fmt.Errorf("max gRPC message size must be positive, got %d", c.MaxGRPCMessageSize)
	}
	if c.MaxIdleConnectionAge <= 0 {
		return fmt.Errorf("max idle connection age must be positive, got %v", c.MaxIdleConnectionAge)
	}
	if c.RequestMaxPastAge <= 0 {
		return fmt.Errorf("request max past age must be positive, got %v", c.RequestMaxPastAge)
	}
	if c.RequestMaxFutureAge <= 0 {
		return fmt.Errorf("request max future age must be positive, got %v", c.RequestMaxFutureAge)
	}
	return nil
}

// NewGRPCServerConfig creates a new gRPC server config with validation
func NewGRPCServerConfig(
	grpcPort uint16,
	maxGRPCMessageSize int,
	maxIdleConnectionAge time.Duration,
	requestMaxPastAge time.Duration,
	requestMaxFutureAge time.Duration,
) (GRPCServerConfig, error) {

	if maxGRPCMessageSize < 0 {
		return GRPCServerConfig{}, fmt.Errorf("max grpc message size must be >= 0, got %d", maxGRPCMessageSize)
	}
	if maxIdleConnectionAge < 0 {
		return GRPCServerConfig{}, fmt.Errorf("max idle connection age must be >= 0, got %v", maxIdleConnectionAge)
	}
	if requestMaxPastAge < 0 {
		return GRPCServerConfig{}, fmt.Errorf("request max past age must be >= 0, got %v", requestMaxPastAge)
	}
	if requestMaxFutureAge < 0 {
		return GRPCServerConfig{}, fmt.Errorf("request max future age must be >= 0, got %v", requestMaxFutureAge)
	}

	return GRPCServerConfig{
		GrpcPort:             grpcPort,
		MaxGRPCMessageSize:   maxGRPCMessageSize,
		MaxIdleConnectionAge: maxIdleConnectionAge,
		RequestMaxPastAge:    requestMaxPastAge,
		RequestMaxFutureAge:  requestMaxFutureAge,
	}, nil
}
