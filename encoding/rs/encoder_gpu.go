//go:build gpu
// +build gpu

package rs

import (
	"fmt"
	"log/slog"
	"math"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	gpu "github.com/Layr-Labs/eigenda/encoding/rs/gpu"
	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"

	_ "go.uber.org/automaxprocs"
)

func (g *Encoder) newEncoder(params encoding.EncodingParams) (*ParametrizedEncoder, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	n := uint8(math.Log2(float64(params.NumEvaluations())))
	if params.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * params.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	// GPU Setup
	runtime.LoadBackendFromEnvOrDefault()

	// trying to choose CUDA if available, or fallback to CPU otherwise (default device)
	var device runtime.Device
	deviceCuda := runtime.CreateDevice("CUDA", 0) // GPU-0
	if runtime.IsDeviceAvailable(&deviceCuda) {
		device = runtime.CreateDevice("CUDA", 0) // GPU-0
		slog.Info("CUDA device available, setting device")
		runtime.SetDevice(&device)
	} else {
		slog.Info("CUDA device not available, falling back to CPU")
		device = runtime.CreateDevice("CPU", 0)
	}

	gpuLock := sync.Mutex{}

	// Setup NTT
	nttCfg, icicle_err := gpu_utils.SetupNTT(25)
	if icicle_err != runtime.Success {
		return nil, fmt.Errorf("could not setup NTT")
	}

	// Set RS CPU computer
	RsComputeDevice := &gpu.RsGpuComputeDevice{
		NttCfg:  nttCfg,
		GpuLock: &gpuLock,
		Device:  device,
	}

	return &ParametrizedEncoder{
		Config:            g.Config,
		EncodingParams:    params,
		Fs:                fs,
		RSEncoderComputer: RsComputeDevice,
	}, nil
}
