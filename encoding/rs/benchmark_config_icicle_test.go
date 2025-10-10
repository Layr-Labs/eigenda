//go:build icicle

package rs_test

import (
	"github.com/Layr-Labs/eigenda/encoding"
)

// getBenchmarkConfig returns the default config using Gnark backend (CPU-only).
// This file is only compiled when the icicle build tag is NOT present.
func getBenchmarkConfig() *encoding.Config {
	cfg := encoding.DefaultConfig()
	cfg.BackendType = encoding.IcicleBackend
	cfg.GPUEnable = true
	return cfg
}
