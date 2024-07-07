package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"runtime"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	cr "github.com/ingonyama-zk/icicle/v2/wrappers/golang/cuda_runtime"
)

func main() {
	TestKzgRs()
	//err := kzg.WriteGeneratorPoints(600000)
	//if err != nil {
	//	log.Println("WriteGeneratorPoints failed:", err)
	//}
	//readpoints()
}

/*
func readpoints() {
	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../../inabox/resources/kzg/g1.point",
		G2Path:          "../../inabox/resources/kzg/g2.point",
		CacheDir:        "SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	// create encoding object
	kzgGroup, _ := prover.NewProver(kzgConfig, true)
	fmt.Println("there are ", len(kzgGroup.Srs.G1), "points")
	for i := 0; i < len(kzgGroup.Srs.G1); i++ {

		fmt.Printf("%v %v\n", i, string(kzgGroup.Srs.G1[i].String()))
	}
	if kzgGroup.Srs.G1[0].X == kzg.GenG1.X && kzgGroup.Srs.G1[0].Y == kzg.GenG1.Y {
		fmt.Println("start with gen")
	}
}
*/

func TestKzgRs() {
	isSmallTest := false

	numSymbols := 4096 / 8
	// encode parameters
	numNode := uint64(4096) // 200
	numSys := uint64(512)   // 180
	numPar := numNode - numSys

	numDevices, _ := cr.GetDeviceCount()
	fmt.Println("num device ", numDevices)

	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../../inabox/resources/kzg/g1.point.600000",
		G2Path:          "../../inabox/resources/kzg/g2.point.600000",
		CacheDir:        "SRSTables",
		SRSOrder:        300000,
		SRSNumberToLoad: 300000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		Verbose:         true,
	}

	if isSmallTest {
		numSymbols = 2
		numNode = 4
		numSys = 2
		numPar = numNode - numSys
		kzgConfig = &kzg.KzgConfig{
			G1Path:          "../../inabox/resources/kzg/g1.point.300000",
			G2Path:          "../../inabox/resources/kzg/g2.point.300000",
			CacheDir:        "SRSTables",
			SRSOrder:        3000,
			SRSNumberToLoad: 3000,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)),
			Verbose:         true,
		}
	}

	// Prepare data
	fmt.Printf("* Task Starts\n")
	fmt.Printf("    Num Sys: %v\n", numSys)
	fmt.Printf("    Num Par: %v\n", numPar)
	//fmt.Printf("    Data size(byte): %v\n", len(inputBytes))

	// create encoding object
	p, _ := prover.NewProver(kzgConfig, true)

	p.UseGpu = false

	params := encoding.EncodingParams{NumChunks: numNode, ChunkLength: uint64(numSymbols) / uint64(numSys)}
	enc, _ := p.GetKzgEncoder(params)

	//inputFr := kzg.ToFrArray(inputBytes)
	inputSize := uint64(numSymbols)
	inputFr := make([]fr.Element, inputSize)
	for i := uint64(0); i < inputSize; i++ {
		inputFr[i].SetInt64(int64(i + 1))
	}

	fmt.Printf("Input \n")
	//printFr(inputFr)

	//inputSize := uint64(len(inputFr))
	commit, lengthCommit, lengthProof, frames, fIndices, err := enc.Encode(inputFr)
	_ = lengthProof
	_ = lengthCommit
	_ = commit
	_ = frames
	_ = fIndices
	if err != nil {
		log.Fatal(err)
	}
	// Optionally verify

	startVerify := time.Now()

	//os.Exit(0)
	for i := 0; i < len(frames); i++ {
		//for i, f := range frames {
		f := frames[i]
		j := fIndices[i]
		q, err := rs.GetLeadingCosetIndex(uint64(i), numSys+numPar)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if j != q {
			log.Fatal("leading coset inconsistency")
		}

		// special case when chunk Len =1, we need to artificially use root doubled
		if params.ChunkLength == uint64(1) {
			j *= 2
		}

		lc := enc.Fs.ExpandedRootsOfUnity[uint64(j)]

		g2Atn, err := kzg.ReadG2Point(uint64(len(f.Coeffs)), kzgConfig)
		if err != nil {
			log.Fatalf("Load g2 %v failed\n", err)
		}
		err = verifier.VerifyFrame(&f, enc.Ks, commit, &lc, &g2Atn)
		if err != nil {
			log.Fatalf("Proof %v failed\n", i)
		}
	}
	fmt.Printf("* Verify %v frames -> all correct. together using %v\n",
		len(frames), time.Since(startVerify))
	// sample some frames
	samples, indices := SampleFrames(frames, uint64(len(frames)-3))
	//samples, indices := SampleFrames(frames, numSys)
	//fmt.Printf("* Sampled %v frames\n", numSys)
	//// Decode data from samples

	dataFr, err := enc.Decode(samples, indices, inputSize)
	if err != nil {
		log.Fatal(err)
	}

	//printFr(dataFr)
	//dataFr, err := kzg.DecodeSys(samples, indices, inputSize)
	//if err != nil {
	//log.Fatalf("%v", err)
	//}
	_ = dataFr

	//fmt.Println(dataFr)
	// printFr(dataFr)
	//deData := kzg.ToByteArray(dataFr, inputByteSize)
	//fmt.Println("dataFr")
	// printFr(dataFr)
	//fmt.Println(deData)
	// Verify data is original in Fr
	//compareData(inputFr, dataFr)
	// Verify data is original in Byte
	//compareDataByte(deData, inputBytes)
	//fmt.Printf("* Compared original %v bytes with reconstructed -> PASS\n", inputByteSize)
	//_ = deData
}

// func getData(inputSize uint64) []fr.Element {
// 	inputFr := make([]fr.Element, inputSize)
// 	for i := uint64(0); i < inputSize; i++ {
// 		bls.AsFr(&inputFr[i], i+1)
// 	}
// 	return inputFr
// }
//
// func compareData(inputFr, dataFr []fr.Element) {
// 	if len(inputFr) != len(dataFr) {
// 		log.Fatalf("Error. Diff length. input %v, data %v\n", len(inputFr), len(dataFr))
// 	}
//
// 	for i := 0; i < len(inputFr); i++ {
// 		if !bls.EqualFr(&inputFr[i], &dataFr[i]) {
// 			log.Fatalf("Error. Diff value at %v. input %v, data %v\n",
// 				i, inputFr[i].String(), dataFr[i].String())
// 		}
// 	}
// }
//
// func compareDataByte(inputFr, dataFr []byte) {
// 	if len(inputFr) != len(dataFr) {
// 		log.Fatalf("Error. Diff length. input %v, data %v\n", len(inputFr), len(dataFr))
// 	}
//
// 	for i := 0; i < len(inputFr); i++ {
// 		if inputFr[i] != dataFr[i] {
// 			log.Fatalf("Error. Diff Data byte value at %v. input %v, data %v\n",
// 				i, inputFr[i:], dataFr[i:])
// 		}
// 	}
// }
//
// func initPoly(size int) ([]fr.Element, []fr.Element) {
// 	v := make([]uint64, size)
// 	for i := 0; i < size; i++ {
// 		v[i] = uint64(i + 1)
// 	}
// 	polyFr := makeFr(v)
// 	fs := kzg.NewFFTSettings(3)
// 	dataFr, _ := fs.FFT(polyFr, false)
// 	return polyFr, dataFr
// }
//
// func initData(size uint64) ([]fr.Element, []fr.Element) {
// 	v := make([]uint64, size)
// 	for i := uint64(0); i < size; i++ {
// 		v[i] = uint64(i + 1)
// 	}
// 	dataFr := makeFr(v)
// 	order := kzg.CeilIntPowerOf2Num(size)
// 	fs := kzg.NewFFTSettings(uint8(order))
// 	polyFr, err := fs.FFT(dataFr, true)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return polyFr, dataFr
// }
//
// func makeFr(input []uint64) []fr.Element {
// 	inputFr := make([]fr.Element, len(input))
// 	for i := 0; i < len(input); i++ {
// 		bls.AsFr(&inputFr[i], input[i])
// 	}
// 	return inputFr
// }

func printFr(d []fr.Element) {
	for _, e := range d {
		fmt.Printf("%v ", e.String())
	}
	fmt.Printf("\n")
}

// func printG1(d []bn254.G1Affine) {
// 	for i, e := range d {
// 		fmt.Printf("%v: %v \n", i, e.String())
// 	}
// 	fmt.Printf("\n")
// }

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
