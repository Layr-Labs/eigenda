// The V1 kzg/prover does both KZG commitment generation and multiproof generation.
// For V2, we split off the committer functionality into this package,
// and kzg/prover/v2 only does multiproof generation.
package committer

import (
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/hashicorp/go-multierror"
)

// Committer is responsible for computing [encoding.BlobCommitments],
// which are needed by clients to create BlobHeaders and disperse blobs.
type Committer struct {
	// G1 SRS points are used for computing Blob commitments.
	g1SRS []bn254.G1Affine
	// G2 SRS points are used for computing Blob length commitments+proofs.
	g2SRS         []bn254.G2Affine
	g2TrailingSRS []bn254.G2Affine
}

func New(g1SRS []bn254.G1Affine, g2SRS []bn254.G2Affine, g2TrailingSRS []bn254.G2Affine) (*Committer, error) {
	if len(g1SRS) == 0 {
		return nil, fmt.Errorf("g1SRS is empty")
	}
	if len(g2SRS) == 0 {
		return nil, fmt.Errorf("g2SRS is empty")
	}
	if len(g2TrailingSRS) == 0 {
		return nil, fmt.Errorf("g2TrailingSRS is empty")
	}
	if len(g1SRS) != len(g2SRS) {
		return nil, fmt.Errorf("g1SRS and g2SRS must be the same length")
	}
	if len(g2SRS) != len(g2TrailingSRS) {
		return nil, fmt.Errorf("g2SRS and g2TrailingSRS must be the same length")
	}

	return &Committer{
		g1SRS:         g1SRS,
		g2SRS:         g2SRS,
		g2TrailingSRS: g2TrailingSRS,
	}, nil
}

type Config struct {
	// Number of SRS points to load from all 3 SRS files. Must be a power of 2.
	// Committer will only be able to compute commitments for blobs of size up to this number of field elements.
	// e.g. if SRSNumberToLoad=2^19, then the committer can compute commitments for blobs of size up to
	// 2^19 field elements = 2^19 * 32 bytes = 16 MiB.
	SRSNumberToLoad   uint64
	G1SRSPath         string
	G2SRSPath         string
	G2TrailingSRSPath string
}

func NewFromConfig(config Config) (*Committer, error) {
	if config.G1SRSPath == "" {
		return nil, fmt.Errorf("G1SRSPath is empty")
	}
	if config.G2SRSPath == "" {
		return nil, fmt.Errorf("G2SRSPath is empty")
	}
	if config.G2TrailingSRSPath == "" {
		return nil, fmt.Errorf("G2TrailingSRSPath is empty")
	}

	// ReadG1/G2Points is CPU bound, the actual reading is very fast, but the parsing is slow.
	// We just spin up as many goroutines as we have CPUs.
	numWorkers := uint64(runtime.GOMAXPROCS(0))
	g1SRS, err := kzg.ReadG1Points(config.G1SRSPath, config.SRSNumberToLoad, numWorkers)
	if err != nil {
		return nil, fmt.Errorf("read G1 points from %s: %w", config.G1SRSPath, err)
	}
	g2SRS, err := kzg.ReadG2Points(config.G2SRSPath, config.SRSNumberToLoad, numWorkers)
	if err != nil {
		return nil, fmt.Errorf("read G2 points from %s: %w", config.G2SRSPath, err)
	}
	// TODO(samlaf): we should have a function ReadG2TrailingPoints that reads from the end of the file.
	numG2point, err := kzg.NumberOfPointsInSRSFile(config.G2TrailingSRSPath, kzg.G2PointBytes)
	if err != nil {
		return nil, fmt.Errorf("number of points in srs file %v: %w", config.G2TrailingSRSPath, err)
	}
	if numG2point < config.SRSNumberToLoad {
		return nil, fmt.Errorf("kzgConfig.G2TrailingPath=%v contains %v G2 Points, "+
			"which is < kzgConfig.SRSNumberToLoad=%v",
			config.G2TrailingSRSPath, numG2point, config.SRSNumberToLoad)
	}
	g2TrailingSRS, err := kzg.ReadG2PointSection(
		config.G2TrailingSRSPath, numG2point-config.SRSNumberToLoad, numG2point, numWorkers)
	if err != nil {
		return nil, fmt.Errorf("read G2 trailing points from %s: %w", config.G2TrailingSRSPath, err)
	}

	return New(g1SRS, g2SRS, g2TrailingSRS)
}

// GetCommitmentsForPaddedLength takes in a byte slice representing a list of bn254
// field elements (32 bytes each, except potentially the last element),
// pads the (potentially incomplete) last element with zeroes, and returns the commitments for the padded list.
func (c *Committer) GetCommitmentsForPaddedLength(data []byte) (encoding.BlobCommitments, error) {
	symbols, err := rs.ToFrArray(data)
	if err != nil {
		return encoding.BlobCommitments{}, fmt.Errorf("ToFrArray: %w", err)
	}

	commit, lengthCommit, lengthProof, err := c.GetCommitments(symbols)
	if err != nil {
		return encoding.BlobCommitments{}, fmt.Errorf("get commitments: %w", err)
	}

	commitments := encoding.BlobCommitments{
		Commitment:       (*encoding.G1Commitment)(commit),
		LengthCommitment: (*encoding.G2Commitment)(lengthCommit),
		LengthProof:      (*encoding.G2Commitment)(lengthProof),
		Length:           math.NextPowOf2u32(uint32(len(symbols))),
	}

	return commitments, nil
}

func (c *Committer) GetCommitments(
	inputFr []fr.Element,
) (*bn254.G1Affine, *bn254.G2Affine, *bn254.G2Affine, error) {
	// We've checked in the constructor that len(g1SRS)=len(g2SRS)=len(g2TrailingSRS)
	// so we only need to check against one of them here.
	if len(inputFr) > len(c.g1SRS) {
		return nil, nil, nil, fmt.Errorf("input length %v > number SRS points %v",
			len(inputFr), len(c.g1SRS))
	}

	encodeStart := time.Now()

	lengthCommitmentChan := make(chan lengthCommitmentResult, 1)
	lengthProofChan := make(chan lengthProofResult, 1)
	commitmentChan := make(chan commitmentResult, 1)

	// compute commit for the full poly
	go func() {
		start := time.Now()
		commit, err := c.computeCommitmentV2(inputFr)
		commitmentChan <- commitmentResult{
			Commitment: commit,
			Err:        err,
			Duration:   time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		lengthCommitment, err := c.computeLengthCommitmentV2(inputFr)
		lengthCommitmentChan <- lengthCommitmentResult{
			LengthCommitment: lengthCommitment,
			Err:              err,
			Duration:         time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		lengthProof, err := c.computeLengthProofV2(inputFr)
		lengthProofChan <- lengthProofResult{
			LengthProof: lengthProof,
			Err:         err,
			Duration:    time.Since(start),
		}
	}()

	lengthProofResult := <-lengthProofChan
	lengthCommitmentResult := <-lengthCommitmentChan
	commitmentResult := <-commitmentChan

	if lengthProofResult.Err != nil || lengthCommitmentResult.Err != nil ||
		commitmentResult.Err != nil {
		return nil, nil, nil, multierror.Append(lengthProofResult.Err, lengthCommitmentResult.Err, commitmentResult.Err)
	}
	totalProcessingTime := time.Since(encodeStart)

	slog.Info("Commitment process details",
		"Input_size_bytes", len(inputFr)*encoding.BYTES_PER_SYMBOL,
		"Total_duration", totalProcessingTime,
		"Committing_duration", commitmentResult.Duration,
		"LengthCommit_duration", lengthCommitmentResult.Duration,
		"lengthProof_duration", lengthProofResult.Duration,
		"SRSOrder", encoding.SRSOrder,
		// TODO(samlaf): should we take NextPowerOf2(len(inputFr)) instead?
		"SRSOrder_shift", encoding.SRSOrder-uint64(len(inputFr)),
	)

	return commitmentResult.Commitment, lengthCommitmentResult.LengthCommitment, lengthProofResult.LengthProof, nil
}

func (c *Committer) computeCommitmentV2(coeffs []fr.Element) (*bn254.G1Affine, error) {
	// compute commit for the full poly
	config := ecc.MultiExpConfig{}
	var commitment bn254.G1Affine
	_, err := commitment.MultiExp(c.g1SRS[:len(coeffs)], coeffs, config)
	if err != nil {
		return nil, fmt.Errorf("multi exp: %w", err)
	}
	return &commitment, nil
}

func (c *Committer) computeLengthCommitmentV2(coeffs []fr.Element) (*bn254.G2Affine, error) {
	var lengthCommitment bn254.G2Affine
	_, err := lengthCommitment.MultiExp(c.g2SRS[:len(coeffs)], coeffs, ecc.MultiExpConfig{})
	if err != nil {
		return nil, fmt.Errorf("multi exp: %w", err)
	}
	return &lengthCommitment, nil
}

func (c *Committer) computeLengthProofV2(coeffs []fr.Element) (*bn254.G2Affine, error) {
	// blobLen must always be a power of 2 in V2
	// coeffs is not modified because padding with 0s doesn't change the commitment,
	// but we need to pretend like it was actually padded with 0s to get the correct length proof.
	blobLen := math.NextPowOf2u32(uint32(len(coeffs)))

	start := uint32(len(c.g2TrailingSRS)) - blobLen
	shiftedSecret := c.g2TrailingSRS[start : start+uint32(len(coeffs))]

	// The proof of low degree is commitment of the polynomial shifted to the largest srs degree
	var lengthProof bn254.G2Affine
	_, err := lengthProof.MultiExp(shiftedSecret, coeffs, ecc.MultiExpConfig{})
	if err != nil {
		return nil, fmt.Errorf("multi exp: %w", err)
	}

	return &lengthProof, nil
}

type lengthCommitmentResult struct {
	LengthCommitment *bn254.G2Affine
	Duration         time.Duration
	Err              error
}

type lengthProofResult struct {
	LengthProof *bn254.G2Affine
	Duration    time.Duration
	Err         error
}

type commitmentResult struct {
	Commitment *bn254.G1Affine
	Duration   time.Duration
	Err        error
}
