package main

import (
	"fmt"
	"log"
	"math/rand"

	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// TODO(samlaf): is anyone using this...? Can we delete it?
func main() {
	TestKzgRs()
}

func TestKzgRs() {
	numSymbols := 1024
	// encode parameters
	numNode := uint64(4) // 200
	numSys := uint64(2)  // 180
	numPar := numNode - numSys
	// Prepare data
	fmt.Printf("* Task Starts\n")
	fmt.Printf("    Num Sys: %v\n", numSys)
	fmt.Printf("    Num Par: %v\n", numPar)
	//fmt.Printf("    Data size(byte): %v\n", len(inputBytes))

	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../../../resources/srs/g1.point",
		G2Path:          "../../../resources/srs/g2.point",
		CacheDir:        "SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	// create encoding object
	p, err := prover.NewProver(kzgConfig, nil)
	if err != nil {
		log.Fatalf("Failed to create prover: %v", err)
	}
	v, err := verifier.NewVerifier(kzgConfig, nil)
	if err != nil {
		log.Fatalf("Failed to create verifier: %v", err)
	}

	params := encoding.EncodingParams{NumChunks: numNode, ChunkLength: uint64(numSymbols) / numSys}
	enc, err := p.GetKzgEncoder(params)
	if err != nil {
		log.Fatalf("Failed to create encoder: %v", err)
	}
	verifier, err := v.GetKzgVerifier(params)
	if err != nil {
		log.Fatalf("Failed to create verifier: %v", err)
	}

	//inputFr := kzg.ToFrArray(inputBytes)
	inputSize := uint64(numSymbols)
	inputFr := make([]fr.Element, inputSize)
	for i := uint64(0); i < inputSize; i++ {
		inputFr[i].SetInt64(int64(i + 1))
	}

	fmt.Printf("Input \n")
	printFr(inputFr)

	//inputSize := uint64(len(inputFr))
	commit, lengthCommit, lengthProof, frames, fIndices, err := enc.Encode(inputFr)
	_ = lengthProof
	_ = lengthCommit
	if err != nil {
		log.Fatal(err)
	}
	// Optionally verify
	startVerify := time.Now()

	for i := 0; i < len(frames); i++ {
		err = verifier.VerifyFrame(&frames[i], uint64(fIndices[i]), commit, params.NumChunks)
		if err != nil {
			log.Fatalf("Failed to verify frame %d: %v", i, err)
		}
	}
	fmt.Printf("* Verify %v frames -> all correct. together using %v\n",
		len(frames), time.Since(startVerify))
	// sample some frames
	samples, indices := SampleFrames(frames, uint64(len(frames)-3))

	dataFr, err := enc.Decode(samples, indices, inputSize)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x\n", dataFr)
}

func printFr(d []fr.Element) {
	for _, e := range d {
		fmt.Printf("%v ", e.String())
	}
	fmt.Printf("\n")
}

func SampleFrames(frames []encoding.Frame, num uint64) ([]encoding.Frame, []uint64) {
	samples := make([]encoding.Frame, num)
	indices := rand.Perm(len(frames))
	indices = indices[:num]

	frameIndices := make([]uint64, num)
	for i, j := range indices {
		samples[i] = frames[j]
		frameIndices[i] = uint64(j)
	}
	return samples, frameIndices
}
