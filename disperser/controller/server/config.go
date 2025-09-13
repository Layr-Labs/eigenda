package server

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
)

// Contains configuration for the controller gRPC server
type Config struct {
	common.GRPCServerConfig

	// If true, use the new payment authentication system running on the controller.
	// If false, payment authentication is disabled and request validation will always fail
	//
	// Note: This flag requires EnableServer to be true in order to function.
	EnablePaymentAuthentication bool
}

// Creates a new server config with validation
func NewConfig(
	grpcServerConfig common.GRPCServerConfig,
	enablePaymentAuthentication bool,
) (Config, error) {

	if enablePaymentAuthentication && !grpcServerConfig.EnableServer {
		return Config{}, fmt.Errorf("payment authentication requires gRPC server to be enabled")
	}

	return Config{
		GRPCServerConfig:            grpcServerConfig,
		EnablePaymentAuthentication: enablePaymentAuthentication,
	}, nil
}
