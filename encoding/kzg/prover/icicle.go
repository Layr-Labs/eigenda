//go:build icicle

package prover

import (
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	prover_icicle "github.com/Layr-Labs/eigenda/encoding/kzg/prover/icicle"
)

func CreateIcicleBackendProver(p *Prover, params encoding.EncodingParams, fs *fft.FFTSettings, ks *kzg.KZGSettings) (*ParametrizedProver, error) {
	_, fftPointsT, err := p.SetupFFTPoints(params)
	if err != nil {
		return nil, err
	}

	icicleDevice, err := icicle.NewIcicleDevice(icicle.IcicleDeviceConfig{
		EnableGPU:  p.Config.EnableGPU,
		NTTSize:    defaultNTTSize,
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
	multiproofBackend := &prover_icicle.KzgMultiProofIcicleBackend{
		Fs:             fs,
		FlatFFTPointsT: icicleDevice.FlatFFTPointsT,
		SRSIcicle:      icicleDevice.SRSG1Icicle,
		SFs:            sfs,
		Srs:            p.Srs,
		NttCfg:         icicleDevice.NttCfg,
		MsmCfg:         icicleDevice.MsmCfg,
		KzgConfig:      p.KzgConfig,
		Device:         icicleDevice.Device,
	}

	// Set up default commitments backend
	commitmentsBackend := &KzgCommitmentsDefaultBackend{
		Srs:        p.Srs,
		G2Trailing: p.G2Trailing,
		KzgConfig:  p.KzgConfig,
	}

	return &ParametrizedProver{
		EncodingParams:        params,
		Encoder:               p.Encoder,
		KzgConfig:             p.KzgConfig,
		Ks:                    ks,
		KzgMultiProofBackend:  multiproofBackend,
		KzgCommitmentsBackend: commitmentsBackend,
	}, nil

}
