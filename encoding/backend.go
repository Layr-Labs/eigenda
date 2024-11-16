package encoding

import "fmt"

type BackendType string

const (
	GnarkBackend  BackendType = "gnark"
	IcicleBackend BackendType = "icicle"
)

type Config struct {
	NumWorker   uint64
	BackendType BackendType
	GPUEnable   bool
	Verbose     bool
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
