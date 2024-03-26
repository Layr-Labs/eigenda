package verifier

import (
	"fmt"
	"math"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

type Config struct {
	G1Path    string
	G2Path    string
	NumPoint  uint64
	NumBatch  uint64
	NumWorker int
}

const numUpdate = 20

func VerifySRS(config Config) {
	numPoint := config.NumPoint
	numBatch := config.NumBatch

	batchSize := uint64(math.Ceil(float64(numPoint) / float64(numBatch)))

	processStart := time.Now()

	updateSize := int64(numBatch / numUpdate)

	fmt.Printf("In total, we will verify %v batches. Each batch contains %v points.\n", numBatch, batchSize)
	fmt.Printf("For the first 3 batches, we show the time taken to verify each batch, then estimate the total verification hours.\n")
	fmt.Printf("After the first 3 batches, we will update every %v batches\n", updateSize)

	flag := false
	var g1Gen bn254.G1Affine
	var g2Gen bn254.G2Affine
	var g2Tau bn254.G2Affine

	for i := int64(0); i < int64(numBatch); i++ {
		begin := time.Now()
		from := i*int64(batchSize) - 1 // -1 for covering previous loop
		to := (i + 1) * int64(batchSize)
		if from < 0 {
			from = 0
		}
		if uint64(to) > numPoint {
			to = int64(numPoint)
		}

		// read in sections to avoid memory overflow
		g1points, err := ReadG1PointSection(config.G1Path, uint64(from), uint64(to), 8)
		if err != nil {
			fmt.Println("err", err)
			return
		}

		g2points, err := ReadG2PointSection(config.G2Path, uint64(from), uint64(to), 8)
		if err != nil {
			fmt.Println("err", err)
			return
		}

		// get generator and initial points
		if !flag {
			g1Gen = g1points[0]
			g2Gen = g2points[0]
			g2Tau = g2points[1]
			flag = true
		}

		verifyBegin := time.Now()
		err = G1Check(g1points, g2points, &g2Gen, &g2Tau, config.NumWorker)
		if err != nil {
			fmt.Println("Verify SRS G1 Check error", err)
			return
		}

		err = G2Check(g1points, g2points, &g1Gen, &g2Gen, config.NumWorker)
		if err != nil {
			fmt.Println("Verify SRS G2 Check error", err)
			return
		}

		if i < 3 {
			elapsed := time.Since(begin)
			expectedFinishDuration := uint64(elapsed.Seconds()) * numBatch
			fmt.Printf("Verify 1 batch takes %v. Verify takes %v\n", elapsed, time.Since(verifyBegin))
			fmt.Printf("verify %v batches will take %v Hours\n", numBatch, expectedFinishDuration/3600.0)
		} else if i%updateSize == 0 {
			fmt.Printf("Verified %v-th batches. Time spent so far is %v\n", i, time.Since(processStart))
		}
	}

	fmt.Println("Done. Everything is correct")
}

// https://github.com/ethereum/kzg-ceremony-specs/blob/master/docs/sequencer/sequencer.md#pairing-checks
func G1Check(g1points []bn254.G1Affine, g2points []bn254.G2Affine, g2Gen *bn254.G2Affine, g2Tau *bn254.G2Affine, numWorker int) error {
	n := uint64(len(g1points))
	if len(g1points) != len(g2points) {
		panic("not equal length")
	}

	workerLoad := uint64(math.Ceil(float64(n) / float64(numWorker)))

	results := make(chan error, numWorker)

	for w := uint64(0); w < uint64(numWorker); w++ {
		start := w * workerLoad
		end := (w + 1) * workerLoad
		if end >= n {
			end = n - 1
		}

		go G1CheckWorker(g1points, g2points, g2Gen, g2Tau, start, end, results)
	}

	for i := 0; i < numWorker; i++ {
		err := <-results
		if err != nil {
			fmt.Println("err", err)
			return err
		}
	}

	return nil
}

// https://github.com/ethereum/kzg-ceremony-specs/blob/master/docs/sequencer/sequencer.md#pairing-checks
func G2Check(g1points []bn254.G1Affine, g2points []bn254.G2Affine, g1Gen *bn254.G1Affine, g2Gen *bn254.G2Affine, numWorker int) error {
	n := uint64(len(g1points))
	if len(g1points) != len(g2points) {
		panic("not equal length")
	}

	workerLoad := uint64(math.Ceil(float64(n) / float64(numWorker)))

	results := make(chan error, numWorker)

	for w := uint64(0); w < uint64(numWorker); w++ {
		start := w * workerLoad
		end := (w + 1) * workerLoad
		if end > n {
			end = n
		}

		go G2CheckWorker(g1points, g2points, g1Gen, g2Gen, start, end, results)
	}

	for i := 0; i < numWorker; i++ {
		err := <-results
		if err != nil {
			fmt.Println("G2Checker err", err)
			return err
		}
	}

	return nil

}

func PairingCheck(a1 *bn254.G1Affine, a2 *bn254.G2Affine, b1 *bn254.G1Affine, b2 *bn254.G2Affine) error {
	var negB1 bn254.G1Affine
	negB1.Neg((*bn254.G1Affine)(b1))

	P := [2]bn254.G1Affine{*(*bn254.G1Affine)(a1), negB1}
	Q := [2]bn254.G2Affine{*(*bn254.G2Affine)(a2), *(*bn254.G2Affine)(b2)}

	ok, err := bn254.PairingCheck(P[:], Q[:])
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("PairingCheck pairing not ok. SRS is invalid")
	}

	return nil
}

func G1CheckWorker(
	g1points []bn254.G1Affine,
	g2points []bn254.G2Affine,
	g2Gen *bn254.G2Affine,
	g2Tau *bn254.G2Affine,
	start uint64, // in element, not in byte
	end uint64,
	results chan<- error,
) {
	for i := start; i < end; i++ {
		err := PairingCheck(&g1points[i+1], g2Gen, &g1points[i], g2Tau)
		if err != nil {
			fmt.Println("pairing check failed at ", i)
			results <- err
			return
		}
	}
	results <- nil
}

func G2CheckWorker(
	g1points []bn254.G1Affine,
	g2points []bn254.G2Affine,
	g1Gen *bn254.G1Affine,
	g2Gen *bn254.G2Affine,
	start uint64, // in element, not in byte
	end uint64,
	results chan<- error,
) {
	for i := start; i < end; i++ {
		err := PairingCheck(&g1points[i], g2Gen, g1Gen, &g2points[i])
		if err != nil {
			results <- err
			return
		}
	}
	results <- nil
}
