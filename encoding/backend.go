package encoding

import (
	"fmt"
	"runtime"

	"github.com/Layr-Labs/eigenda/encoding/icicle"
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

// DefaultConfig returns a Config struct with default values.
// If icicle is availeble (binary built with icicle tag), it sets the backend to icicle and enables GPU.
// Make sure to set GPUEnable to false if you want to run on CPU only.
// If icicle is not available (build without icicle tag), it sets the backend to gnark.
func DefaultConfig() *Config {
	if icicle.IsAvailable {
		return &Config{
			NumWorker:   uint64(runtime.GOMAXPROCS(0)),
			BackendType: IcicleBackend,
			GPUEnable:   true,
		}
	}
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
