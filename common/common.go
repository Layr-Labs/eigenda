package common

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
)

var (
	ErrInvalidDomainType = fmt.Errorf("invalid domain type")
)

// G1Point struct to represent G1Point in Solidity
type G1Point struct {
	X *big.Int
	Y *big.Int
}

// QuorumBlobParam struct to represent QuorumBlobParam in Solidity
type QuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// BlobHeader struct to represent BlobHeader in Solidity
type BlobHeader struct {
	Commitment       G1Point
	DataLength       uint32
	QuorumBlobParams []QuorumBlobParam
}

type Certificate disperser.BlobInfo

func (c *Certificate) BlobIndex() uint32 {
	return c.BlobVerificationProof.BlobIndex
}

func (c *Certificate) BatchHeaderRoot() []byte {
	return c.BlobVerificationProof.BatchMetadata.BatchHeader.BatchRoot
}

func (c *Certificate) ReadBlobHeader() BlobHeader {
	// parse quorum params

	qps := make([]QuorumBlobParam, len(c.BlobHeader.BlobQuorumParams))
	for i, qp := range c.BlobHeader.BlobQuorumParams {
		qps[i] = QuorumBlobParam{
			QuorumNumber:                    uint8(qp.QuorumNumber),
			AdversaryThresholdPercentage:    uint8(qp.AdversaryThresholdPercentage),
			ConfirmationThresholdPercentage: uint8(qp.ConfirmationThresholdPercentage),
			ChunkLength:                     qp.ChunkLength,
		}
	}

	return BlobHeader{
		Commitment: G1Point{
			X: new(big.Int).SetBytes(c.BlobHeader.Commitment.X),
			Y: new(big.Int).SetBytes(c.BlobHeader.Commitment.Y),
		},
		DataLength:       c.BlobHeader.DataLength,
		QuorumBlobParams: qps,
	}
}

func (c *Certificate) Proof() *disperser.BlobVerificationProof {
	return c.BlobVerificationProof
}

// DomainType is a enumeration type for the different data domains for which a
// blob can exist between
type DomainType uint8

const (
	BinaryDomain DomainType = iota
	PolyDomain
	UnknownDomain
)

func (dt DomainType) String() string {
	switch dt {
	case BinaryDomain:
		return "binary"

	case PolyDomain:
		return "polynomial"

	default:
		return "unknown"
	}
}

func StrToDomainType(s string) DomainType {
	switch s {
	case "binary":
		return BinaryDomain
	case "polynomial":
		return PolyDomain
	default:
		return UnknownDomain
	}
}

// Helper utility functions //

func EqualSlices[P comparable](s1, s2 []P) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func ParseBytesAmount(s string) (uint64, error) {
	s = strings.TrimSpace(s)

	// Extract numeric part and unit
	numStr := s
	unit := ""
	for i, r := range s {
		if !('0' <= r && r <= '9' || r == '.') {
			numStr = s[:i]
			unit = s[i:]
			break
		}
	}

	// Convert numeric part to float64
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value: %v", err)
	}

	unit = strings.ToLower(strings.TrimSpace(unit))

	// Convert to uint64 based on the unit (case-insensitive)
	switch unit {
	case "b", "":
		return uint64(num), nil
	case "kib":
		return uint64(num * 1024), nil
	case "kb":
		return uint64(num * 1000), nil // Decimal kilobyte
	case "mib":
		return uint64(num * 1024 * 1024), nil
	case "mb":
		return uint64(num * 1000 * 1000), nil // Decimal megabyte
	case "gib":
		return uint64(num * 1024 * 1024 * 1024), nil
	case "gb":
		return uint64(num * 1000 * 1000 * 1000), nil // Decimal gigabyte
	case "tib":
		return uint64(num * 1024 * 1024 * 1024 * 1024), nil
	case "tb":
		return uint64(num * 1000 * 1000 * 1000 * 1000), nil // Decimal terabyte
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unit)
	}
}

type Stats struct {
	Entries int
	Reads   int
}
