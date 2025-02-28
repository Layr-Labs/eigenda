//go:build icicle

package prover

import (
	"math"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	gnarkprover "github.com/Layr-Labs/eigenda/encoding/kzg/prover/gnark"
	icicleprover "github.com/Layr-Labs/eigenda/encoding/kzg/prover/icicle"
)

const (
	// MAX_NTT_SIZE is the maximum NTT domain size needed to compute FFTs for the
	// largest supported blobs. Assuming a coding ratio of 1/8 and symbol size of 32 bytes:
	// - Encoded size: 2^{MAX_NTT_SIZE} * 32 bytes ≈ 1 GB
	// - Original blob size: 2^{MAX_NTT_SIZE} * 32 / 8 = 2^{MAX_NTT_SIZE + 2} ≈ 128 MB
	MAX_NTT_SIZE = 25
)

func CreateIcicleBackendProver(p *Prover, params encoding.EncodingParams, fs *fft.FFTSettings, ks *kzg.KZGSettings) (*ParametrizedProver, error) {
	_, fftPointsT, err := p.SetupFFTPoints(params)
	if err != nil {
		return nil, err
	}
	icicleDevice, err := icicle.NewIcicleDevice(icicle.IcicleDeviceConfig{
		GPUEnable:  p.Config.GPUEnable,
		NTTSize:    MAX_NTT_SIZE,
		FFTPointsT: fftPointsT,
		SRSG1:      p.Srs.G1[:p.KzgConfig.SRSNumberToLoad],
	})
	if err != nil {
		return nil, err
	}

	// Create subgroup FFT settings
	t := uint8(math.Log2(float64(2 * params.NumChunks)))
	sfs := fft.NewFFTSettings(t)

	// Set up icicle multiproof backend
	multiproofBackend := &icicleprover.KzgMultiProofIcicleBackend{
		Fs:             fs,
		FlatFFTPointsT: icicleDevice.FlatFFTPointsT,
		SRSIcicle:      icicleDevice.SRSG1Icicle,
		SFs:            sfs,
		Srs:            p.Srs,
		NttCfg:         icicleDevice.NttCfg,
		MsmCfg:         icicleDevice.MsmCfg,
		KzgConfig:      p.KzgConfig,
		Device:         icicleDevice.Device,
		GpuLock:        sync.Mutex{},
	}

	// Set up gnark commitments backend
	commitmentsBackend := &gnarkprover.KzgCommitmentsGnarkBackend{
		Srs:        p.Srs,
		G2Trailing: p.G2Trailing,
		KzgConfig:  p.KzgConfig,
	}

	return &ParametrizedProver{
		EncodingParams:        params,
		Encoder:               p.encoder,
		KzgConfig:             p.KzgConfig,
		Ks:                    ks,
		KzgMultiProofBackend:  multiproofBackend,
		KzgCommitmentsBackend: commitmentsBackend,
	}, nil
}
