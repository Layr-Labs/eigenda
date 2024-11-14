//go:build icicle

package prover

import (
	"fmt"
	"log/slog"
	"math"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/cpu"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/gpu"
	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	icicle_runtime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

func CreateIcicleBackendProver(p *Prover, params encoding.EncodingParams, fs *fft.FFTSettings, ks *kzg.KZGSettings) (*ParametrizedProver, error) {
	icicle_runtime.LoadBackendFromEnvOrDefault()

	// Setup device (GPU or CPU)
	device := setupIcicleDevice(p.Config.EnableGPU)

	_, fftPointsT, err := p.SetupFFTPoints(params)
	if err != nil {
		return nil, err
	}

	// Setup NTT
	// Run on the device
	var wg sync.WaitGroup
	wg.Add(1)

	var (
		nttCfg         core.NTTConfig[[icicle_bn254.SCALAR_LIMBS]uint32]
		flatFftPointsT []icicle_bn254.Affine
		srsG1Icicle    []icicle_bn254.Affine
		msmCfg         core.MSMConfig
		setupErr       error
		icicle_err     icicle_runtime.EIcicleError
	)

	icicle_runtime.RunOnDevice(&device, func(args ...any) {
		defer wg.Done()

		// Setup NTT
		nttCfg, icicle_err = gpu_utils.SetupNTT(defaultNTTSize)
		if icicle_err != icicle_runtime.Success {
			setupErr = fmt.Errorf("could not setup NTT")
			return
		}

		// Setup MSM
		flatFftPointsT, srsG1Icicle, msmCfg, _, icicle_err = gpu_utils.SetupMsm(
			fftPointsT,
			p.Srs.G1[:p.KzgConfig.SRSNumberToLoad],
		)
		if icicle_err != icicle_runtime.Success {
			setupErr = fmt.Errorf("could not setup MSM")
			return
		}
	})

	wg.Wait()

	if setupErr != nil {
		return nil, setupErr
	}

	// Create subgroup FFT settings
	t := uint8(math.Log2(float64(2 * params.NumChunks)))
	sfs := fft.NewFFTSettings(t)

	// Set up GPU proof computer
	proofComputer := &gpu.KzgGpuProofDevice{
		Fs:             fs,
		FlatFFTPointsT: flatFftPointsT,
		SRSIcicle:      srsG1Icicle,
		SFs:            sfs,
		Srs:            p.Srs,
		NttCfg:         nttCfg,
		MsmCfg:         msmCfg,
		KzgConfig:      p.KzgConfig,
		Device:         device,
	}

	// Set up CPU commitments computer (same as default backend)
	commitmentsComputer := &cpu.KzgCPUCommitmentsDevice{
		Srs:        p.Srs,
		G2Trailing: p.G2Trailing,
		KzgConfig:  p.KzgConfig,
	}

	return &ParametrizedProver{
		EncodingParams:      params,
		Encoder:             p.Encoder,
		KzgConfig:           p.KzgConfig,
		Ks:                  ks,
		ProofComputer:       proofComputer,
		CommitmentsComputer: commitmentsComputer,
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
