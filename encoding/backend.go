package encoding

import (
	"fmt"
	"runtime"

	_ "go.uber.org/automaxprocs/maxprocs"
)

type BackendType string

const (
	// GnarkBackend is the default backend, using the gnark-crypto library.
	// It only supports CPU execution.
	GnarkBackend BackendType = "gnark"
	// IcicleBackend uses the icicle performanced-oriented library.
	// It is optimized for GPU (CUDA and metal) execution, but also supports CPU.
	IcicleBackend BackendType = "icicle"
)

type Config struct {
	NumWorker   uint64
	BackendType BackendType
	GPUEnable   bool
}

// DefaultConfig returns a Config struct with default values
func DefaultConfig() *Config {
	return &Config{
		NumWorker:   uint64(runtime.GOMAXPROCS(0)),
		BackendType: GnarkBackend,
		GPUEnable:   false,
	}
}

// ParseBackendType converts a string to BackendType and validates it
func ParseBackendType(backend string) (BackendType, error) {
	switch BackendType(backend) {
	case GnarkBackend:
		return GnarkBackend, nil
	case IcicleBackend:
		return IcicleBackend, nil
	default:
		return "", fmt.Errorf("unsupported backend type: %s. Must be one of: gnark, icicle", backend)
	}
}
