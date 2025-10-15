package encoding

import (
	"errors"
	"fmt"
	gomath "math"

	"github.com/Layr-Labs/eigenda/common/math"
	"golang.org/x/exp/constraints"
)

type EncodingParams struct {
	// number of Fr symbols stored inside a chunk
	ChunkLength uint64
	// number of total chunks (always a power of 2)
	NumChunks uint64
	// number of systematic chunk (always power of 2)
	blobLength uint64
}

func (p *EncodingParams) SetBlobLength(blobLength uint64) {
	p.blobLength = blobLength
}

func (p *EncodingParams) GetBlobLength() (uint64, error) {
	if p.blobLength == 0 {
		return 0, fmt.Errorf("Blob length is not set in EncodingParams")
	}
	return p.blobLength, nil
}

func (p EncodingParams) NumEvaluations() uint64 {
	return p.NumChunks * p.ChunkLength
}

func (p EncodingParams) Validate() error {
	if !math.IsPowerOfTwo(p.NumChunks) {
		return fmt.Errorf("number of chunks must be a power of 2, got %d", p.NumChunks)
	}
	if !math.IsPowerOfTwo(p.ChunkLength) {
		return fmt.Errorf("chunk length must be a power of 2, got %d", p.ChunkLength)
	}
	return nil
}

func ParamsFromMins[T constraints.Integer](minChunkLength, minNumChunks T) EncodingParams {
	return EncodingParams{
		NumChunks:   math.NextPowOf2u64(uint64(minNumChunks)),
		ChunkLength: math.NextPowOf2u64(uint64(minChunkLength)),
	}
}

// ParamsFromSysPar takes in the number of systematic and parity chunks, as well as the data size in bytes,
// and returns the corresponding encoding parameters.
func ParamsFromSysPar(numSys, numPar, dataSize uint64) EncodingParams {

	numNodes := numSys + numPar
	dataLen := math.RoundUpDivide(dataSize, BYTES_PER_SYMBOL)
	chunkLen := math.RoundUpDivide(dataLen, numSys)
	return ParamsFromMins(chunkLen, numNodes)

}

func GetNumSys(dataSize uint64, chunkLen uint64) uint64 {
	dataLen := math.RoundUpDivide(dataSize, BYTES_PER_SYMBOL)
	numSys := dataLen / chunkLen
	return numSys
}

// ValidateEncodingParams takes in the encoding parameters and returns an error if they are invalid.
func ValidateEncodingParams(params EncodingParams, SRSOrder uint64) error {
	if params.NumChunks == 0 {
		return errors.New("number of chunks must be greater than 0")
	}
	if params.ChunkLength == 0 {
		return errors.New("chunk length must be greater than 0")
	}

	if params.NumChunks > gomath.MaxUint64/params.ChunkLength {
		return fmt.Errorf("multiplication overflow: ChunkLength: %d, NumChunks: %d", params.ChunkLength, params.NumChunks)
	}

	// Check that the parameters are valid with respect to the SRS. The precomputed terms of the amortized KZG
	// prover use up to order params.ChunkLen*params.NumChunks-1 for the SRS, so we must have
	// params.ChunkLen*params.NumChunks-1 <= g.SRSOrder. The condition below could technically
	// be relaxed to params.ChunkLen*params.NumChunks > g.SRSOrder+1, but because all of the paramters are
	// powers of 2, the stricter condition is equivalent.
	if params.ChunkLength*params.NumChunks > SRSOrder {
		return fmt.Errorf("the supplied encoding parameters are not valid with respect to the SRS. ChunkLength: %d, NumChunks: %d, SRSOrder: %d", params.ChunkLength, params.NumChunks, SRSOrder)
	}

	return nil

}

// ValidateEncodingParamsAndBlobLength takes in the encoding parameters and blob length and returns an error if they are collectively invalid.
func ValidateEncodingParamsAndBlobLength(params EncodingParams, blobLength, SRSOrder uint64) error {

	if err := ValidateEncodingParams(params, SRSOrder); err != nil {
		return err
	}

	if params.ChunkLength*params.NumChunks < blobLength {
		return errors.New("the supplied encoding parameters are not sufficient for the size of the data input")
	}

	return nil

}
