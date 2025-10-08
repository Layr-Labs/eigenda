package main

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/config"
)

var _ config.VerifiableConfig = (*TestConfig)(nil)

type TestConfig struct {
	A string
	B int
	C bool
	D time.Duration
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
	cfg, err := config.Bootstrap(DefaultTestConfig, "TEST")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Test configuration: %+v\n", cfg)
}
