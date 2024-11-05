package encoding

import (
	"errors"
	"fmt"

	"golang.org/x/exp/constraints"
)

var (
	ErrInvalidParams = errors.New("invalid encoding params")
)

type EncodingParams struct {
	ChunkLength uint64 // ChunkSize is the length of the chunk in symbols
	NumChunks   uint64
}

func (p EncodingParams) ChunkDegree() uint64 {
	return p.ChunkLength - 1
}

func (p EncodingParams) NumEvaluations() uint64 {
	return p.NumChunks * p.ChunkLength
}

func (p EncodingParams) Validate() error {

	if NextPowerOf2(p.NumChunks) != p.NumChunks {
		return ErrInvalidParams
	}

	if NextPowerOf2(p.ChunkLength) != p.ChunkLength {
		return ErrInvalidParams
	}

	return nil
}

func ParamsFromMins[T constraints.Integer](minChunkLength, minNumChunks T) EncodingParams {
	return EncodingParams{
		NumChunks:   NextPowerOf2(uint64(minNumChunks)),
		ChunkLength: NextPowerOf2(uint64(minChunkLength)),
	}
}

func ParamsFromSysPar(numSys, numPar, dataSize uint64) EncodingParams {

	numNodes := numSys + numPar
	dataLen := roundUpDivide(dataSize, BYTES_PER_SYMBOL)
	chunkLen := roundUpDivide(dataLen, numSys)
	return ParamsFromMins(chunkLen, numNodes)

}

func GetNumSys(dataSize uint64, chunkLen uint64) uint64 {
	dataLen := roundUpDivide(dataSize, BYTES_PER_SYMBOL)
	numSys := dataLen / chunkLen
	return numSys
}

// ValidateEncodingParams takes in the encoding parameters and returns an error if they are invalid.
func ValidateEncodingParams(params EncodingParams, SRSOrder uint64) error {

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
