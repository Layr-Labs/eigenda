package dispersal

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/reputation"
)

var _ config.VerifiableConfig = (*DisperserClientMultiplexerConfig)(nil)

// Configuration for the [DisperserClientMultiplexer]
type DisperserClientMultiplexerConfig struct {
	// Dispersers to use beyond the default set from the DisperserRegistry contract, which specifies the default
	// dispersers for network participants to interact with.
	AdditionalDispersers []uint32
	// Dispersers to never interact with.
	//
	// This field may be used to avoid interacting with dispersers in the default set.
	DisperserBlacklist []uint32
	// Configuration for the reputation system used to select dispersers
	ReputationConfig reputation.ReputationConfig
	// Whether to use secure gRPC connections (TLS) when connecting to dispersers
	UseSecureGrpcFlag bool
}

func DefaultDisperserClientMultiplexerConfig() *DisperserClientMultiplexerConfig {
	return &DisperserClientMultiplexerConfig{
		AdditionalDispersers: nil,
		DisperserBlacklist:   nil,
		ReputationConfig:     reputation.DefaultConfig(),
		UseSecureGrpcFlag:    true,
	}
}

// Verify implements [config.VerifiableConfig].
func (c *DisperserClientMultiplexerConfig) Verify() error {
	err := c.ReputationConfig.Verify()
	if err != nil {
		return fmt.Errorf("verify reputation config: %w", err)
	}

	return nil
}
