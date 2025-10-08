package main

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/config"
)

var _ config.DocumentedConfig = (*TestConfig)(nil)

type TestConfig struct {
	A string
	B int
	C bool
	D time.Duration
}

// GetEnvVarPrefix implements config.DocumentedConfig.
func (t *TestConfig) GetEnvVarPrefix() string {
	return "TEST"
}

// GetName implements config.DocumentedConfig.
func (t *TestConfig) GetName() string {
	return "TestConfig"
}

// GetPackagePaths implements config.DocumentedConfig.
func (t *TestConfig) GetPackagePaths() []string {
	return []string{"github.com/Layr-Labs/eigenda/common/config/bootstrap_test"}
}

// Verify implements config.VerifiableConfig.
func (t *TestConfig) Verify() error {
	if t.B < 0 {
		return fmt.Errorf("variable B must be non-negative, got %d", t.B)
	}
	return nil
}

func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		A: "defaultA",
		B: 42,
		C: false,
		D: 5 * time.Second,
	}
}

func main() {
	cfg, err := config.Bootstrap(DefaultTestConfig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Test configuration: %+v\n", cfg)
}
