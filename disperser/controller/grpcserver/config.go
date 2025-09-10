package grpcserver

import (
	"fmt"
	"time"
)

// Contains configuration for the controller gRPC server
type Config struct {
	// Enables the gRPC server
	//
	// TODO(litt3): the option to disable the server will be removed once any feature has been fully rolled out which
	// requires the server.
	Enable bool

	// Port that the gRPC server listens on
	GrpcPort string

	// Maximum size of a gRPC message that the server will accept (in bytes)
	MaxGRPCMessageSize int

	// Maximum time a connection can be idle before it is closed.
	MaxIdleConnectionAge time.Duration
}

// Creates a new server config with validation
func NewConfig(
	enable bool,
	grpcPort string,
	maxGRPCMessageSize int,
	maxIdleConnectionAge time.Duration,
) (Config, error) {

	if enable {
		if grpcPort == "" {
			return Config{}, fmt.Errorf("grpc port is required")
		}
		if maxGRPCMessageSize < 0 {
			return Config{}, fmt.Errorf("max grpc message size must be >= 0, got %d", maxGRPCMessageSize)
		}
		if maxIdleConnectionAge < 0 {
			return Config{}, fmt.Errorf("max idle connection age must be >= 0, got %v", maxIdleConnectionAge)
		}
	}

	return Config{
		Enable:               enable,
		GrpcPort:             grpcPort,
		MaxGRPCMessageSize:   maxGRPCMessageSize,
		MaxIdleConnectionAge: maxIdleConnectionAge,
	}, nil
}
