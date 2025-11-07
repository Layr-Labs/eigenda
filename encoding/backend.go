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
	// Increase this value to allow more concurrent GPU frame (chunk+proof) tasks.
	// Only used by V2 when Backend=icicle and GPUEnable=true.
	// Note Chunk generation (encoding/v2/rs) and multiproofs generation (encoding/v2/kzg/prover)
	// each have their own separate semaphore which is weighted using this same value.
	//
	// This protects against out-of-memory errors on the GPU, not GPU time usage.
	// WARNING: setting this value too high may lead to out-of-memory errors on the GPU.
	// If this ever happens, the GPU device needs to be rebooted as it can be left in a bad state.
	//
	// For now we use this very coarse-grained approach, instead of using a RAM-usage based semaphore,
	// because that would feel brittle and require approximations of RAM usage per MSM/NTT operation.
	// This would need to take into account RAM used by:
	// - msm/ntt initialization (srs points and ntt roots that are kept on GPU)
	// - msm as a fct of input size (see https://dev.ingonyama.com/api/cpp/msm#memory-usage-estimation)
	// - ntt as a fct of input size (afaiu ntt uses input+output=2*input size in RAM)
	// If we want to enable more optimal GPU usage, we should revisit this.
	GPUConcurrentFrameGenerationDangerous int64
}

// DefaultConfig returns a Config struct with default values.
// If icicle is available (binary built with icicle tag), it sets the backend to icicle and enables GPU.
// Make sure to set GPUEnable to false if you want to run on CPU only.
// If icicle is not available (build without icicle tag), it sets the backend to gnark.
func DefaultConfig() *Config {
	if icicle.IsAvailable {
		return &Config{
			NumWorker:                             uint64(runtime.GOMAXPROCS(0)),
			BackendType:                           IcicleBackend,
			GPUEnable:                             true,
			GPUConcurrentFrameGenerationDangerous: 1,
		}
	}
	return &Config{
		NumWorker:                             uint64(runtime.GOMAXPROCS(0)),
		BackendType:                           GnarkBackend,
		GPUEnable:                             false,
		GPUConcurrentFrameGenerationDangerous: 0, // Not used
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
