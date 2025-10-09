package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	proverv2 "github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	verifierv2 "github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
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
	TraceFile     string
	EnableVerify  bool
}

func parseFlags() Config {
	config := Config{}
	flag.StringVar(&config.OutputFile, "output", "benchmark_results.json", "Output file for results")
	flag.Uint64Var(&config.MinBlobLength, "min-blob-length", 524288, "Minimum blob length (power of 2)")
	flag.Uint64Var(&config.MaxBlobLength, "max-blob-length", 524288, "Maximum blob length (power of 2)")
	flag.Uint64Var(&config.NumChunks, "num-chunks", 8192, "Minimum number of chunks (power of 2)")
	flag.StringVar(&config.CPUProfile, "cpuprofile", "", "Write CPU profile to file")
	flag.StringVar(&config.MemProfile, "memprofile", "", "Write memory profile to file")
	flag.StringVar(&config.TraceFile, "trace", "trace.out", "Write execution trace to file")
	flag.BoolVar(&config.EnableVerify, "enable-verify", true, "Verify blobs after encoding")
	flag.Parse()
	return config
}

var proverKzgConfig *proverv2.KzgConfig
var verifierKzgConfig *verifierv2.KzgConfig

func main() {
	config := parseFlags()

	fmt.Println("Config output", config.OutputFile)

	// Setup phase
	proverKzgConfig = &proverv2.KzgConfig{
		G1Path:          "/home/ubuntu/eigenda/resources/srs/g1.point",
		G2Path:          "/home/ubuntu/eigenda/resources/srs/g2.point",
		G2TrailingPath:  "/home/ubuntu/eigenda/resources/srs/g2.trailing.point",
		CacheDir:        "/home/ubuntu/eigenda/resources/srs/SRSTables",
		SRSNumberToLoad: 524288,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
		PreloadEncoder:  true,
	}

	verifierKzgConfig = &verifierv2.KzgConfig{
		G1Path:          "/home/ubuntu/eigenda/resources/srs/g1.point",
		SRSNumberToLoad: 524288,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	fmt.Printf("* Task Starts\n")

	cfg := &encoding.Config{
		BackendType: encoding.GnarkBackend,
		GPUEnable:   false,
		NumWorker:   uint64(runtime.GOMAXPROCS(0)),
	}
	p, err := proverv2.NewProver(proverKzgConfig, cfg)

	if err != nil {
		log.Fatalf("Failed to create prover: %v", err)
	}

	if config.CPUProfile != "" {
		f, err := os.Create(config.CPUProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer core.CloseLogOnError(f, f.Name(), nil)
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if config.TraceFile != "" {
		f, err := os.Create(config.TraceFile)
		if err != nil {
			log.Fatal("could not create trace file: ", err)
		}
		defer core.CloseLogOnError(f, f.Name(), nil)
		if err := trace.Start(f); err != nil {
			log.Fatal("could not start trace: ", err)
		}
		defer trace.Stop()
	}

	results := runBenchmark(p, &config)
	if config.MemProfile != "" {
		f, err := os.Create(config.MemProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer core.CloseLogOnError(f, f.Name(), nil)
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

func runBenchmark(p *proverv2.Prover, config *Config) []BenchmarkResult {
	var results []BenchmarkResult

	// Fixed coding ratio of 8
	codingRatio := uint64(8)

	for blobLength := config.MinBlobLength; blobLength <= config.MaxBlobLength; blobLength *= 2 {
		chunkLen := (blobLength * codingRatio) / config.NumChunks
		if chunkLen < 1 {
			continue // Skip invalid configurations
		}
		result := benchmarkEncodeAndVerify(p, blobLength, config.NumChunks, chunkLen, config.EnableVerify)
		results = append(results, result)
	}
	return results
}

func benchmarkEncodeAndVerify(
	p *proverv2.Prover,
	blobLength uint64,
	numChunks uint64,
	chunkLen uint64,
	verifyResults bool,
) BenchmarkResult {
	params := encoding.EncodingParams{
		NumChunks:   numChunks,
		ChunkLength: chunkLen,
	}

	fmt.Printf("Running benchmark: numChunks=%d, chunkLen=%d, blobLength=%d\n", params.NumChunks, params.ChunkLength, blobLength)

	enc, err := p.GetKzgEncoder(params)
	if err != nil {
		log.Fatalf("Failed to get KZG encoder: %v", err)
	}

	// Create polynomial
	inputSize := blobLength
	inputFr := make([]fr.Element, inputSize)
	for i := uint64(0); i < inputSize; i++ {
		inputFr[i].SetInt64(int64(i + 1))
	}

	commit, _, _, err := enc.GetCommitments(inputFr)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	frames, _, err := enc.GetFrames(inputFr)
	if err != nil {
		log.Fatal(err)
	}
	duration := time.Since(start)

	verifyResult := true
	verifyStart := time.Now()

	if verifyResults {
		v, err := verifierv2.NewVerifier(verifierKzgConfig, nil)
		if err != nil {
			log.Fatalf("Failed to create verifier: %v", err)
		}

		samples := []encoding.Sample{}
		for i, frame := range frames {
			samples = append(samples, encoding.Sample{
				Commitment:      (*encoding.G1Commitment)(commit),
				Chunk:           &frame,
				AssignmentIndex: uint(i),
				BlobIndex:       0,
			})
		}

		err = v.UniversalVerifySubBatch(params, samples, 1)
		if err != nil {
			log.Fatal("Wtf", err)
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
