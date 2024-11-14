//go:build icicle

package icicle

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	icicle_runtime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

// IcicleDevice wraps the core device setup and configurations
type IcicleDevice struct {
	Device         icicle_runtime.Device
	NttCfg         core.NTTConfig[[icicle_bn254.SCALAR_LIMBS]uint32]
	MsmCfg         core.MSMConfig
	FlatFFTPointsT []icicle_bn254.Affine
	SRSG1Icicle    []icicle_bn254.Affine
}

// IcicleDeviceConfig holds configuration options for device setup
type IcicleDeviceConfig struct {
	EnableGPU bool
	NTTSize   uint8
	// MSM setup parameters (optional)
	FFTPointsT [][]bn254.G1Affine
	SRSG1      []bn254.G1Affine
}

// NewIcicleDevice creates and initializes a new IcicleDevice
func NewIcicleDevice(config IcicleDeviceConfig) (*IcicleDevice, error) {
	icicle_runtime.LoadBackendFromEnvOrDefault()

	device := setupDevice(config.EnableGPU)

	var wg sync.WaitGroup
	wg.Add(1)
	var (
		nttCfg         core.NTTConfig[[icicle_bn254.SCALAR_LIMBS]uint32]
		msmCfg         core.MSMConfig
		flatFftPointsT []icicle_bn254.Affine
		srsG1Icicle    []icicle_bn254.Affine
		setupErr       error
		icicleErr      icicle_runtime.EIcicleError
	)

	// Setup NTT and optionally MSM on device
	icicle_runtime.RunOnDevice(&device, func(args ...any) {
		defer wg.Done()

		// Setup NTT
		nttCfg, icicleErr = SetupNTT(config.NTTSize)
		if icicleErr != icicle_runtime.Success {
			setupErr = fmt.Errorf("could not setup NTT")
			return
		}

		// Setup MSM if parameters are provided
		if config.FFTPointsT != nil && config.SRSG1 != nil {
			flatFftPointsT, srsG1Icicle, msmCfg, _, icicleErr = SetupMsm(
				config.FFTPointsT,
				config.SRSG1,
			)
			if icicleErr != icicle_runtime.Success {
				setupErr = fmt.Errorf("could not setup MSM")
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
func setupDevice(enableGPU bool) icicle_runtime.Device {
	if enableGPU {
		return setupGPUDevice()
	}

	return setupCPUDevice()
}

// setupGPUDevice attempts to initialize a CUDA device, falling back to CPU if unavailable
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

// setupCPUDevice initializes a CPU device
func setupCPUDevice() icicle_runtime.Device {
	device := icicle_runtime.CreateDevice("CPU", 0)
	if icicle_runtime.IsDeviceAvailable(&device) {
		slog.Info("CPU device available, setting device")
	}
	icicle_runtime.SetDevice(&device)

	return device
}
