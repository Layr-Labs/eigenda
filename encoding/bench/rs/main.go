package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	rs_cpu "github.com/Layr-Labs/eigenda/encoding/rs/cpu"
	rs_gpu "github.com/Layr-Labs/eigenda/encoding/rs/gpu"
	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	icicle_runtime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

type BenchmarkResult struct {
	NumChunks    uint64        `json:"num_chunks"`
	ChunkLength  uint64        `json:"chunk_length"`
	BlobLength   uint64        `json:"blob_length"`
	EncodeTime   time.Duration `json:"encode_time"`
	VerifyTime   time.Duration `json:"verify_time"`
	VerifyResult bool          `json:"verify_result"`
}

type Config struct {
	MinBlobLength uint64 `json:"min_blob_length"`
	MaxBlobLength uint64 `json:"max_blob_length"`
	OutputFile    string
	BlobLength    uint64
	NumChunks     uint64
	NumRuns       uint64
	CPUProfile    string
	MemProfile    string
	EnableVerify  bool
}

func parseFlags() Config {
	config := Config{}
	flag.StringVar(&config.OutputFile, "output", "benchmark_results.json", "Output file for results")
	flag.Uint64Var(&config.MinBlobLength, "min-blob-length", 1024, "Minimum blob length (power of 2)")
	flag.Uint64Var(&config.MaxBlobLength, "max-blob-length", 1048576, "Maximum blob length (power of 2)")
	flag.Uint64Var(&config.NumChunks, "num-chunks", 8192, "Minimum number of chunks (power of 2)")
	flag.StringVar(&config.CPUProfile, "cpuprofile", "", "Write CPU profile to file")
	flag.StringVar(&config.MemProfile, "memprofile", "", "Write memory profile to file")
	flag.BoolVar(&config.EnableVerify, "enable-verify", true, "Verify blobs after encoding")
	flag.Parse()
	return config
}

func main() {
	config := parseFlags()

	fmt.Println("Config output", config.OutputFile)

	fmt.Printf("* Task Starts\n")

	// GPU Setup
	icicle_runtime.LoadBackendFromEnvOrDefault()

	// trying to choose CUDA if available, or fallback to CPU otherwise (default device)
	var device icicle_runtime.Device
	deviceCuda := icicle_runtime.CreateDevice("CUDA", 0) // GPU-0
	if icicle_runtime.IsDeviceAvailable(&deviceCuda) {
		device = icicle_runtime.CreateDevice("CUDA", 0) // GPU-0
		slog.Debug("CUDA device available, setting device")
		icicle_runtime.SetDevice(&device)
	} else {
		slog.Debug("CUDA device not available, falling back to CPU")
		device = icicle_runtime.CreateDevice("CPU", 0)
		icicle_runtime.SetDevice(&device)
	}

	gpuLock := sync.Mutex{}

	// Setup FFT Settings
	// time this part
	fs := fft.NewFFTSettings(uint8(math.Log2(float64(8192 * 1024))))

	for i := 1; i <= 1024; i *= 2 {
		start := time.Now()
		_ = fft.NewFFTSettings(uint8(math.Log2(float64(8192 * i))))
		fmt.Printf("FFT Settings took (%v, %v) : %v\n", 8192, i, time.Since(start))
	}

	// Setup NTT
	start := time.Now()
	fmt.Println("wtf")

	var nttCfg core.NTTConfig[[icicle_bn254.SCALAR_LIMBS]uint32]
	var icicle_err icicle_runtime.EIcicleError
	var wg sync.WaitGroup
	wg.Add(1)

	icicle_runtime.RunOnDevice(&device, func(args ...any) {
		defer wg.Done()
		maxSettings := uint8(math.Log2(float64(8192 * 1024)))
		nttCfg, icicle_err = gpu_utils.SetupNTT(maxSettings)
		if icicle_err != icicle_runtime.Success {
			log.Fatal("could not setup NTT")
		}
		fmt.Printf("NTT Setup took: %v\n", time.Since(start))

		for i := 1; i <= 1024; i *= 2 {
			nttCfg, icicle_err = gpu_utils.SetupNTT(uint8(math.Log2(float64(8192 * i))))
			if icicle_err != icicle_runtime.Success {
				log.Fatal("could not setup NTT")
			}
			fmt.Printf("NTT Setup took (%v, %v): %v\n", 8192, i, time.Since(start))
		}
	})

	wg.Wait()

	fmt.Println("Wwww")

	RsComputeDevice := &rs_gpu.RsGpuComputeDevice{
		NttCfg:  nttCfg,
		GpuLock: &gpuLock,
		Device:  device,
	}

	if config.CPUProfile != "" {
		f, err := os.Create(config.CPUProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	results := runBenchmark(&config, RsComputeDevice, fs)
	if config.MemProfile != "" {
		f, err := os.Create(config.MemProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	// Output results as JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(config.OutputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write results to file: %v", err)
	}

	fmt.Printf("Benchmark results written to %s\n", config.OutputFile)
}

func runBenchmark(config *Config, compute rs.RsComputeDevice, fs *fft.FFTSettings) []BenchmarkResult {
	var results []BenchmarkResult

	// Fixed coding ratio of 8
	codingRatio := uint64(8)

	for blobLength := config.MinBlobLength; blobLength <= config.MaxBlobLength; blobLength *= 2 {
		chunkLen := (blobLength * codingRatio) / config.NumChunks
		if chunkLen < 1 {
			continue // Skip invalid configurations
		}
		result := benchmarkRSEncode(blobLength, config.NumChunks, chunkLen, config.EnableVerify, compute, fs)
		results = append(results, result)
	}
	return results
}

func benchmarkRSEncode(blobLength uint64, numChunks uint64, chunkLen uint64, verifyResults bool, compute rs.RsComputeDevice, fs *fft.FFTSettings) BenchmarkResult {
	params := encoding.EncodingParams{
		NumChunks:   numChunks,
		ChunkLength: chunkLen,
	}

	fmt.Printf("Running benchmark: numChunks=%d, chunkLen=%d, blobLength=%d\n", params.NumChunks, params.ChunkLength, blobLength)
	encoder, err := rs.NewEncoderFFT(params, fs)
	if err != nil {
		log.Fatal(err)
	}

	// Set RS CPU computer
	RsComputeDevice := &rs_cpu.RsCpuComputeDevice{
		Fs: encoder.Fs,
	}
	encoder.Computer = RsComputeDevice

	// Set RS GPU computer
	encoder.Computer = compute

	// Create array of bytes
	inputSize := blobLength * 32
	inputData := make([]byte, inputSize)
	for i := 0; i < len(inputData); i++ {
		var element fr.Element
		element.SetRandom()
		bytes := element.Bytes()
		copy(inputData[i:min(i+encoding.BYTES_PER_SYMBOL, len(inputData))], bytes[:])
	}

	start := time.Now()
	frames, indices, err := encoder.EncodeBytes(inputData)
	if err != nil {
		log.Fatal(err)
	}
	duration := time.Since(start)

	verifyResult := true
	verifyStart := time.Now()

	if verifyResults {

		// Convert indices to uint32
		res := make([]uint64, len(indices))
		for i, d := range indices {
			res[i] = uint64(d)
		}

		samples, indices := sampleFrames(frames, uint64((len(frames))))
		data, err := encoder.Decode(samples, indices, uint64(len(inputData)))
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(inputData)
		// fmt.Println(data)

		verifyResult = len(data) == int(inputSize)
		log.Println("Decoded data length: ", len(data), "Expected length: ", inputSize)
		// Check if the decoded data is correct
		for i := 0; i < len(data); i++ {
			if data[i] != inputData[i] {
				verifyResult = false
				break
			}
		}
	}

	verifyTime := time.Since(verifyStart)

	return BenchmarkResult{
		NumChunks:    numChunks,
		ChunkLength:  chunkLen,
		BlobLength:   blobLength,
		EncodeTime:   duration,
		VerifyTime:   verifyTime,
		VerifyResult: verifyResult,
	}
}

func sampleFrames(frames []rs.Frame, num uint64) ([]rs.Frame, []uint64) {
	samples := make([]rs.Frame, num)
	indices := rand.Perm(len(frames))
	indices = indices[:num]

	frameIndices := make([]uint64, num)
	for i, j := range indices {
		samples[i] = frames[j]
		frameIndices[i] = uint64(j)
	}
	return samples, frameIndices
}
