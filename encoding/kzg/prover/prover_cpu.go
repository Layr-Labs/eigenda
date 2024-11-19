//go:build !gpu
// +build !gpu

package prover

import (
	"log"
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	kzg_prover_cpu "github.com/Layr-Labs/eigenda/encoding/kzg/prover/cpu"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	rs_cpu "github.com/Layr-Labs/eigenda/encoding/rs/cpu"
	"github.com/consensys/gnark-crypto/ecc/bn254"

	_ "go.uber.org/automaxprocs"
)

func (g *Prover) newProver(params encoding.EncodingParams) (*ParametrizedProver, error) {
	if err := encoding.ValidateEncodingParams(params, g.SRSOrder); err != nil {
		return nil, err
	}

	encoder, err := rs.NewEncoder(params, g.Verbose)
	if err != nil {
		log.Println("Could not create encoder: ", err)
		return nil, err
	}

	subTable, err := NewSRSTable(g.CacheDir, g.Srs.G1, g.NumWorker)
	if err != nil {
		log.Println("Could not create srs table:", err)
		return nil, err
	}

	fftPoints, err := subTable.GetSubTables(encoder.NumChunks, encoder.ChunkLength)
	if err != nil {
		log.Println("could not get sub tables", err)
		return nil, err
	}

	fftPointsT := make([][]bn254.G1Affine, len(fftPoints[0]))
	for i := range fftPointsT {
		fftPointsT[i] = make([]bn254.G1Affine, len(fftPoints))
		for j := uint64(0); j < encoder.ChunkLength; j++ {
			fftPointsT[i][j] = fftPoints[j][i]
		}
	}
	_ = fftPoints
	n := uint8(math.Log2(float64(encoder.NumEvaluations())))
	if encoder.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * encoder.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	ks, err := kzg.NewKZGSettings(fs, g.Srs)
	if err != nil {
		return nil, err
	}

	t := uint8(math.Log2(float64(2 * encoder.NumChunks)))
	sfs := fft.NewFFTSettings(t)

	// Set KZG Prover CPU computer
	computer := &kzg_prover_cpu.KzgCpuProofDevice{
		Fs:         fs,
		FFTPointsT: fftPointsT,
		SFs:        sfs,
		Srs:        g.Srs,
		G2Trailing: g.G2Trailing,
		KzgConfig:  g.KzgConfig,
	}

	// Set RS CPU computer
	RsComputeDevice := &rs_cpu.RsCpuComputeDevice{
		Fs:             fs,
		EncodingParams: params,
	}
	encoder.Computer = RsComputeDevice

	return &ParametrizedProver{
		Encoder:   encoder,
		KzgConfig: g.KzgConfig,
		Ks:        ks,
		Computer:  computer,
	}, nil
}
