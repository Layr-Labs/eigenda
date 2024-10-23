//go:build gpu
// +build gpu

package prover

import (
	"fmt"
	"log"
	"log/slog"
	"math"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	kzg_prover "github.com/Layr-Labs/eigenda/encoding/kzg/prover/gpu"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	rs_encoder "github.com/Layr-Labs/eigenda/encoding/rs/cpu"
	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"

	_ "go.uber.org/automaxprocs"
)

func (g *Prover) newProver(params encoding.EncodingParams) (*ParametrizedProver, error) {

	// Check that the parameters are valid with respect to the SRS.
	if params.ChunkLength*params.NumChunks >= g.SRSOrder {
		return nil, fmt.Errorf("the supplied encoding parameters are not valid with respect to the SRS. ChunkLength: %d, NumChunks: %d, SRSOrder: %d", params.ChunkLength, params.NumChunks, g.SRSOrder)
	}

	encoder, err := rs.NewEncoder(params)
	if err != nil {
		log.Println("Could not create encoder: ", err)
		return nil, err
	}

	subTable, err := NewSRSTable(g.CacheDir, g.Srs.G1, g.NumWorker)
	if err != nil {
		log.Println("Could not create srs table:", err)
		return nil, err
	}

	log.Println("Getting sub tables")
	fftPoints, err := subTable.GetSubTables(encoder.NumChunks, encoder.ChunkLength)
	if err != nil {
		log.Println("could not get sub tables", err)
		return nil, err
	}

	log.Println("Transposing sub tables")
	fftPointsT := make([][]bn254.G1Affine, len(fftPoints[0]))
	for i := range fftPointsT {
		fftPointsT[i] = make([]bn254.G1Affine, len(fftPoints))
		for j := uint64(0); j < encoder.ChunkLength; j++ {
			fftPointsT[i][j] = fftPoints[j][i]
		}
	}
	_ = fftPoints

	log.Println("Creating FFT settings")
	n := uint8(math.Log2(float64(encoder.NumEvaluations())))
	if encoder.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * encoder.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	log.Println("Creating KZG settings")
	ks, err := kzg.NewKZGSettings(fs, g.Srs)
	if err != nil {
		return nil, err
	}

	t := uint8(math.Log2(float64(2 * encoder.NumChunks)))
	sfs := fft.NewFFTSettings(t)

	// GPU Setup
	runtime.LoadBackendFromEnvOrDefault()

	// trying to choose CUDA if available, or fallback to CPU otherwise (default device)
	var device runtime.Device
	deviceCuda := runtime.CreateDevice("CUDA", 0) // GPU-0
	if runtime.IsDeviceAvailable(&deviceCuda) {
		device = runtime.CreateDevice("CUDA", 0) // GPU-0
		slog.Debug("CUDA device available, setting device")
		runtime.SetDevice(&device)
	} else {
		slog.Debug("CUDA device not available, falling back to CPU")
		device = runtime.CreateDevice("CPU", 0)
		runtime.SetDevice(&device)
	}

	gpuLock := sync.Mutex{}

	// Setup NTT
	nttCfg, icicle_err := gpu_utils.SetupNTT()
	if icicle_err != runtime.Success {
		return nil, fmt.Errorf("could not setup NTT")
	}

	// Setup MSM
	flatFftPointsT, srsG1Icicle, msmCfg, msmCfgG2, icicle_err := gpu_utils.SetupMsm(fftPointsT, g.Srs.G1[:g.SRSNumberToLoad])
	if icicle_err != runtime.Success {
		return nil, fmt.Errorf("could not setup MSM")
	}

	stream, icicle_err := runtime.CreateStream()
	if icicle_err != runtime.Success {
		return nil, fmt.Errorf("could not create stream")
	}

	// Set KZG Prover GPU computer
	computer := &kzg_prover.KzgGpuProofDevice{
		Fs:             fs,
		FlatFFTPointsT: flatFftPointsT,
		SRSIcicle:      srsG1Icicle,
		SFs:            sfs,
		Srs:            g.Srs,
		G2Trailing:     g.G2Trailing,
		NttCfg:         nttCfg,
		MsmCfg:         msmCfg,
		MsmCfgG2:       msmCfgG2,
		KzgConfig:      g.KzgConfig,
		GpuLock:        &gpuLock,
		Stream:         &stream,
		Device:         device,
	}

	RsComputeDevice := &rs_encoder.RsCpuComputeDevice{
		Fs: fs,
	}

	encoder.Computer = RsComputeDevice

	return &ParametrizedProver{
		Encoder:   encoder,
		KzgConfig: g.KzgConfig,
		Ks:        ks,
		Computer:  computer,
	}, nil
}
