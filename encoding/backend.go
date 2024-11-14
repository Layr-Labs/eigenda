package encoding

import "fmt"

type BackendType string

const (
	BackendDefault BackendType = "default"
	BackendIcicle  BackendType = "icicle"
)

type Config struct {
	NumWorker   uint64
	BackendType BackendType
	EnableGPU   bool
	Verbose     bool
}

// ParseBackendType converts a string to BackendType and validates it
func ParseBackendType(backend string) (BackendType, error) {
	switch BackendType(backend) {
	case BackendDefault:
		return BackendDefault, nil
	case BackendIcicle:
		return BackendIcicle, nil
	default:
		return "", fmt.Errorf("unsupported backend type: %s. Must be one of: default, icicle", backend)
	}
}
