package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"
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
	OutputFile   string
	BlobLength   uint64
	NumChunks    uint64
	NumRuns      uint64
	CPUProfile   string
	MemProfile   string
	EnableVerify bool
}

func parseFlags() Config {
	config := Config{}
	flag.StringVar(&config.OutputFile, "output", "benchmark_results.json", "Output file for results")
	flag.Uint64Var(&config.BlobLength, "blob-length", 1048576, "Blob length (power of 2)")
	flag.Uint64Var(&config.NumChunks, "num-chunks", 8192, "Minimum number of chunks (power of 2)")
	flag.Uint64Var(&config.NumRuns, "num-runs", 10, "Number of times to run the benchmark")
	flag.StringVar(&config.CPUProfile, "cpuprofile", "", "Write CPU profile to file")
	flag.StringVar(&config.MemProfile, "memprofile", "", "Write memory profile to file")
	flag.BoolVar(&config.EnableVerify, "enable-verify", false, "Verify blobs after encoding")
	flag.Parse()
	return config
}

func main() {
	config := parseFlags()

	fmt.Println("Config output", config.OutputFile)

	// Setup phase
	kzgConfig := &kzg.KzgConfig{
		G1Path:          "/home/ec2-user/resources/kzg/g1.point",
		G2Path:          "/home/ec2-user/resources/kzg/g2.point",
		CacheDir:        "/home/ec2-user/resources/kzg/SRSTables",
		SRSOrder:        268435456,
		SRSNumberToLoad: 2097152,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		Verbose:         true,
	}

	fmt.Printf("* Task Starts\n")

	// create encoding object
	p, _ := prover.NewProver(kzgConfig, true)

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

	results := runBenchmark(p, &config)
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

func runBenchmark(p *prover.Prover, config *Config) []BenchmarkResult {
	var results []BenchmarkResult

	// Fixed coding ratio of 8
	codingRatio := uint64(8)
	for i := uint64(0); i < config.NumRuns; i++ {
		chunkLen := (config.BlobLength * codingRatio) / config.NumChunks
		if chunkLen < 1 {
			continue // Skip invalid configurations
		}
		result := benchmarkEncodeAndVerify(p, config.BlobLength, config.NumChunks, chunkLen, config.EnableVerify)
		results = append(results, result)
	}
	return results
}

func benchmarkEncodeAndVerify(p *prover.Prover, blobLength uint64, numChunks uint64, chunkLen uint64, verifyResults bool) BenchmarkResult {
	params := encoding.EncodingParams{
		NumChunks:   numChunks,
		ChunkLength: chunkLen,
	}

	fmt.Printf("Running benchmark: numChunks=%d, chunkLen=%d, blobLength=%d\n", params.NumChunks, params.ChunkLength, blobLength)

	enc, _ := p.GetKzgEncoder(params)

	// Create polynomial
	inputSize := blobLength
	inputFr := make([]fr.Element, inputSize)
	for i := uint64(0); i < inputSize; i++ {
		inputFr[i].SetInt64(int64(i + 1))
	}

	start := time.Now()
	commit, _, _, frames, fIndices, err := enc.Encode(inputFr)
	if err != nil {
		log.Fatal(err)
	}
	duration := time.Since(start)

	verifyResult := true
	verifyStart := time.Now()

	if verifyResults {
		for i := 0; i < len(frames); i++ {
			f := frames[i]
			j := fIndices[i]
			q, err := rs.GetLeadingCosetIndex(uint64(i), numChunks)
			if err != nil {
				log.Fatalf("%v", err)
			}

			if j != q {
				log.Fatal("leading coset inconsistency")
			}

			lc := enc.Fs.ExpandedRootsOfUnity[uint64(j)]

			g2Atn, err := kzg.ReadG2Point(uint64(len(f.Coeffs)), p.KzgConfig)
			if err != nil {
				log.Fatalf("Load g2 %v failed\n", err)
			}

			err = verifier.VerifyFrame(&f, enc.Ks, commit, &lc, &g2Atn)
			if err != nil {
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
