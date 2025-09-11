package service

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
	EnableServer bool

	// If true, use the new payment authentication system running on the controller.
	// If false, payment authentication is disabled and request validation will always fail
	//
	// Note: This flag requires EnableServer to be true in order to function.
	EnablePaymentAuthentication bool

	// Port that the gRPC server listens on
	GrpcPort string

	// Maximum size of a gRPC message that the server will accept (in bytes)
	MaxGRPCMessageSize int

	// Maximum time a connection can be idle before it is closed.
	MaxIdleConnectionAge time.Duration

	// Maximum age of an authorization request in the past that the server will accept.
	// Requests older than this will be rejected to prevent replay attacks.
	AuthorizationRequestMaxPastAge time.Duration

	// Maximum age of an authorization request in the future that the server will accept.
	// Requests with timestamps too far in the future will be rejected.
	AuthorizationRequestMaxFutureAge time.Duration
}

// Creates a new server config with validation
func NewConfig(
	enableServer bool,
	enablePaymentAuthentication bool,
	grpcPort string,
	maxGRPCMessageSize int,
	maxIdleConnectionAge time.Duration,
	authorizationRequestMaxPastAge time.Duration,
	authorizationRequestMaxFutureAge time.Duration,
) (Config, error) {

	if enableServer {
		if grpcPort == "" {
			return Config{}, fmt.Errorf("grpc port is required")
		}
		if maxGRPCMessageSize < 0 {
			return Config{}, fmt.Errorf("max grpc message size must be >= 0, got %d", maxGRPCMessageSize)
		}
		if maxIdleConnectionAge < 0 {
			return Config{}, fmt.Errorf("max idle connection age must be >= 0, got %v", maxIdleConnectionAge)
		}
		if authorizationRequestMaxPastAge < 0 {
			return Config{}, fmt.Errorf("authorization request max past age must be >= 0, got %v",
				authorizationRequestMaxPastAge)
		}
		if authorizationRequestMaxFutureAge < 0 {
			return Config{}, fmt.Errorf("authorization request max future age must be >= 0, got %v",
				authorizationRequestMaxFutureAge)
		}
	}

	if enablePaymentAuthentication && !enableServer {
		return Config{}, fmt.Errorf("payment authentication requires gRPC server to be enabled")
	}

	return Config{
		EnableServer:                     enableServer,
		EnablePaymentAuthentication:      enablePaymentAuthentication,
		GrpcPort:                         grpcPort,
		MaxGRPCMessageSize:               maxGRPCMessageSize,
		MaxIdleConnectionAge:             maxIdleConnectionAge,
		AuthorizationRequestMaxPastAge:   authorizationRequestMaxPastAge,
		AuthorizationRequestMaxFutureAge: authorizationRequestMaxFutureAge,
	}, nil
}
