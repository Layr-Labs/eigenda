//go:build icicle

package rs

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	icicle "github.com/Layr-Labs/eigenda/encoding/rs/icicle"
	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	icicle_runtime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

func CreateIcicleBackendEncoder(e *Encoder, params encoding.EncodingParams, fs *fft.FFTSettings) (*ParametrizedEncoder, error) {
	icicle_runtime.LoadBackendFromEnvOrDefault()

	device := setupIcicleDevice(e.Config.EnableGPU)

	var wg sync.WaitGroup
	wg.Add(1)

	var (
		nttCfg     core.NTTConfig[[icicle_bn254.SCALAR_LIMBS]uint32]
		setupErr   error
		icicle_err icicle_runtime.EIcicleError
	)

	// Setup NTT on device
	icicle_runtime.RunOnDevice(&device, func(args ...any) {
		defer wg.Done()

		// Setup NTT
		nttCfg, icicle_err = gpu_utils.SetupNTT(defaultNTTSize)
		if icicle_err != icicle_runtime.Success {
			setupErr = fmt.Errorf("could not setup NTT")
			return
		}
	})

	wg.Wait()

	if setupErr != nil {
		return nil, setupErr
	}

	return &ParametrizedEncoder{
		Config:         e.Config,
		EncodingParams: params,
		Fs:             fs,
		RSEncoderComputer: &icicle.RsIcicleComputeDevice{
			NttCfg: nttCfg,
			Device: device,
		},
	}, nil
}

func setupIcicleDevice(enableGPU bool) icicle_runtime.Device {
	fmt.Println(enableGPU)
	if enableGPU {
		return setupGPUDevice()
	}
	return setupCPUDevice()
}

func setupGPUDevice() icicle_runtime.Device {
	deviceCuda := icicle_runtime.CreateDevice("CUDA", 0)
	if icicle_runtime.IsDeviceAvailable(&deviceCuda) {
		device := icicle_runtime.CreateDevice("CUDA", 0)
		slog.Info("CUDA device available, setting device")
		icicle_runtime.SetDevice(&device)
		return device
	}

	slog.Info("CUDA device not available, falling back to CPU")
	return setupCPUDevice()
}

func setupCPUDevice() icicle_runtime.Device {
	device := icicle_runtime.CreateDevice("CPU", 0)
	if icicle_runtime.IsDeviceAvailable(&device) {
		slog.Info("CPU device available, setting device")
	}
	icicle_runtime.SetDevice(&device)
	return device
}
