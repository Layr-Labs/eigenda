// The V1 kzg/prover does both KZG commitment generation and multiproof generation.
// For V2, we split off the committer functionality into this package,
// and kzg/prover/v2 only does multiproof generation.
package committer

import (
	"fmt"
	"runtime"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Committer is responsible for computing [encoding.BlobCommitments],
// which are needed by clients to create BlobHeaders and disperse blobs.
type Committer struct {
	// G1 SRS points are used for computing Blob commitments.
	g1SRS []bn254.G1Affine
	// G2 SRS points are used for computing Blob length commitments.
	g2SRS []bn254.G2Affine
	// G2 trailing SRS points are used for computing Blob length proofs.
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
	// Number of SRS points to load from SRS files. Must be a power of 2.
	// Committer will only be able to compute commitments for blobs of size up to this number of field elements.
	// e.g. if SRSNumberToLoad=2^19, then the committer can compute commitments for blobs of size up to
	// 2^19 field elements = 2^19 * 32 bytes = 16 MiB.
	SRSNumberToLoad uint64
	G1SRSPath       string
	// There are 2 ways to configure G2 points:
	// 1. Entire G2 SRS file (16GiB) is provided via G2SRSPath (G2TrailingSRSPath is not used).
	// 2. G2SRSPath and G2TrailingSRSPath both contain at least SRSNumberToLoad points,
	//    where G2SRSPath contains the first SRSNumberToLoad points of the full G2 SRS file,
	//    and G2TrailingSRSPath contains the last SRSNumberToLoad points of the G2 SRS file.
	//
	// TODO(samlaf): to prevent misconfigurations and simplify the code, we should probably
	// not multiplex G2SRSPath like this, and instead use a G2PrefixPath config.
	// Then EITHER G2SRSPath is used, OR both G2PrefixSRSPath and G2TrailingSRSPath are used.
	G2SRSPath         string
	G2TrailingSRSPath string
}

var _ config.VerifiableConfig = (*Config)(nil)

func (c *Config) Verify() error {
	if c.SRSNumberToLoad <= 0 {
		return fmt.Errorf("SRSNumberToLoad must be specified for disperser version 2")
	}
	if c.G1SRSPath == "" {
		return fmt.Errorf("G1Path must be specified for disperser version 2")
	}
	if c.G2SRSPath == "" {
		return fmt.Errorf("G2Path must be specified for disperser version 2")
	}
	// G2TrailingSRSPath is optional but its need depends on the content of G2SRSPath
	// so we can't check it here. It is checked inside [NewFromConfig].
	return nil
}

func NewFromConfig(config Config) (*Committer, error) {
	if err := config.Verify(); err != nil {
		return nil, fmt.Errorf("config verify: %w", err)
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

	var g2TrailingSRS []bn254.G2Affine
	hasG2TrailingFile := len(config.G2TrailingSRSPath) != 0
	if hasG2TrailingFile {
		// TODO(samlaf): this function/check should probably be done in ReadG2PointSection
		numG2point, err := kzg.NumberOfPointsInSRSFile(config.G2TrailingSRSPath, kzg.G2PointBytes)
		if err != nil {
			return nil, fmt.Errorf("number of points in srs file %v: %w", config.G2TrailingSRSPath, err)
		}
		if numG2point < config.SRSNumberToLoad {
			return nil, fmt.Errorf("config.G2TrailingPath=%v contains %v G2 Points, "+
				"which is < config.SRSNumberToLoad=%v",
				config.G2TrailingSRSPath, numG2point, config.SRSNumberToLoad)
		}

		// use g2 trailing file
		g2TrailingSRS, err = kzg.ReadG2PointSection(
			config.G2TrailingSRSPath,
			numG2point-config.SRSNumberToLoad,
			numG2point, // last exclusive
			numWorkers,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to read G2 trailing points (%v to %v) from file %v: %w",
				numG2point-config.SRSNumberToLoad, numG2point, config.G2TrailingSRSPath, err)
		}
	} else {
		// require entire G2SRSPath to contain all 2^28 points, from which we can read the trailing points
		numG2point, err := kzg.NumberOfPointsInSRSFile(config.G2SRSPath, kzg.G2PointBytes)
		if err != nil {
			return nil, fmt.Errorf("number of points in srs file: %w", err)
		}
		if numG2point < encoding.SRSOrder {
			return nil, fmt.Errorf("no config.G2TrailingPath was passed, yet the G2 SRS file %v is incomplete: contains %v < 2^28 G2 Points", config.G2SRSPath, numG2point)
		}
		g2TrailingSRS, err = kzg.ReadG2PointSection(
			config.G2SRSPath,
			encoding.SRSOrder-config.SRSNumberToLoad,
			encoding.SRSOrder, // last exclusive
			numWorkers,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to read G2 points (%v to %v) from file %v: %w",
				encoding.SRSOrder-config.SRSNumberToLoad, encoding.SRSOrder, config.G2SRSPath, err)
		}
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

	// We compute all 3 commitments sequentially, since each individual computation
	// already saturates all cores by default.
	commit, err := c.computeCommitmentV2(inputFr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("compute commitment: %w", err)
	}

	lengthCommitment, err := c.computeLengthCommitmentV2(inputFr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("compute length commitment: %w", err)
	}

	lenProof, err := c.computeLengthProofV2(inputFr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("compute length proof: %w", err)
	}

	return commit, lengthCommitment, lenProof, nil
}

func (c *Committer) computeCommitmentV2(coeffs []fr.Element) (*bn254.G1Affine, error) {
	var commitment bn254.G1Affine
	_, err := commitment.MultiExp(c.g1SRS[:len(coeffs)], coeffs, ecc.MultiExpConfig{})
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
