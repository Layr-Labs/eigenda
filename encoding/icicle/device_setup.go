//go:build icicle

package icicle

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	runtime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

// IcicleDevice wraps the core device setup and configurations
type IcicleDevice struct {
	Device         runtime.Device
	NttCfg         core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
	MsmCfg         core.MSMConfig
	FlatFFTPointsT []iciclebn254.Affine
	SRSG1Icicle    []iciclebn254.Affine
}

// IcicleDeviceConfig holds configuration options for a single device.
//   - The GPUEnable parameter is used to enable GPU acceleration.
//   - The NTTSize parameter is used to set the maximum domain size for NTT configuration.
//   - The FFTPointsT and SRSG1 parameters are used to set up the MSM configuration.
//   - MSM setup is optional and can be skipped by not providing these parameters.
//     The reason for this is that not all applications require an MSM setup. For example
//     in the case of reed-solomon, it only requires the NTT setup.
type IcicleDeviceConfig struct {
	GPUEnable bool
	NTTSize   uint8

	// MSM setup parameters (optional)
	FFTPointsT [][]bn254.G1Affine
	SRSG1      []bn254.G1Affine
}

// NewIcicleDevice creates and initializes a new IcicleDevice
func NewIcicleDevice(config IcicleDeviceConfig) (*IcicleDevice, error) {
	runtime.LoadBackendFromEnvOrDefault()

	device, err := setupDevice(config.GPUEnable)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	var (
		nttCfg         core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
		msmCfg         core.MSMConfig
		flatFftPointsT []iciclebn254.Affine
		srsG1Icicle    []iciclebn254.Affine
		setupErr       error
		icicleErr      runtime.EIcicleError
	)

	// Setup NTT and optionally MSM on device
	runtime.RunOnDevice(&device, func(args ...any) {
		defer wg.Done()

		// Setup NTT
		nttCfg, icicleErr = SetupNTT(config.NTTSize)
		if icicleErr != runtime.Success {
			setupErr = fmt.Errorf("could not setup NTT: %v", icicleErr.AsString())
			return
		}

		// Setup MSM if parameters are provided
		if config.FFTPointsT != nil && config.SRSG1 != nil {
			flatFftPointsT, srsG1Icicle, msmCfg, icicleErr = SetupMsmG1(
				config.FFTPointsT,
				config.SRSG1,
			)
			if icicleErr != runtime.Success {
				setupErr = fmt.Errorf("could not setup MSM: %v", icicleErr.AsString())
				return
			}
		}
	})

	wg.Wait()

	if setupErr != nil {
		return nil, setupErr
	}

	return &IcicleDevice{
		Device:         device,
		NttCfg:         nttCfg,
		MsmCfg:         msmCfg,
		FlatFFTPointsT: flatFftPointsT,
		SRSG1Icicle:    srsG1Icicle,
	}, nil
}

// setupDevice initializes either a GPU or CPU device
func setupDevice(gpuEnable bool) (runtime.Device, error) {
	if gpuEnable {
		return setupGPUDevice()
	}

	return setupCPUDevice()
}

// setupGPUDevice attempts to initialize a CUDA device, falling back to CPU if unavailable
func setupGPUDevice() (runtime.Device, error) {
	deviceCuda := runtime.CreateDevice("CUDA", 0)
	if runtime.IsDeviceAvailable(&deviceCuda) {
		device := runtime.CreateDevice("CUDA", 0)
		slog.Info("CUDA device available, setting device")
		runtime.SetDevice(&device)

		return device, nil
	}

	slog.Info("CUDA device not available, falling back to CPU")
	return setupCPUDevice()
}

// setupCPUDevice initializes a CPU device
func setupCPUDevice() (runtime.Device, error) {
	device := runtime.CreateDevice("CPU", 0)
	if !runtime.IsDeviceAvailable(&device) {
		slog.Error("CPU device is not available")
		return device, errors.New("cpu device is not available")
	}

	runtime.SetDevice(&device)
	return device, nil
}
