package encoding

import (
	"errors"
	"fmt"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

var (
	ErrInvalidParams = errors.New("invalid encoding params")
)

type EncodingParams struct {
	ChunkLen  uint64 // ChunkSize is the length of the chunk in symbols
	NumChunks uint64
}

func (p EncodingParams) ChunkDegree() uint64 {
	return p.ChunkLen - 1
}

func (p EncodingParams) NumEvaluations() uint64 {
	return p.NumChunks * p.ChunkLen
}

func (p EncodingParams) Validate() error {

	if NextPowerOf2(p.NumChunks) != p.NumChunks {
		return ErrInvalidParams
	}

	if NextPowerOf2(p.ChunkLen) != p.ChunkLen {
		return ErrInvalidParams
	}

	return nil
}

func GetNumSys(dataSize uint64, chunkLen uint64) uint64 {
	dataLen := roundUpDivide(dataSize, bls.BYTES_PER_COEFFICIENT)
	numSys := dataLen / chunkLen
	return numSys
}

func ParamsFromMins(numChunks, chunkLen uint64) EncodingParams {

	chunkLen = NextPowerOf2(chunkLen)
	numChunks = NextPowerOf2(numChunks)

	return EncodingParams{
		NumChunks: numChunks,
		ChunkLen:  chunkLen,
	}

}

func GetEncodingParams(numSys, numPar, dataSize uint64) EncodingParams {

	numNodes := numSys + numPar
	dataLen := roundUpDivide(dataSize, bls.BYTES_PER_COEFFICIENT)
	chunkLen := roundUpDivide(dataLen, numSys)
	return ParamsFromMins(numNodes, chunkLen)

}

// ValidateEncodingParams takes in the encoding parameters and returns an error if they are invalid.
func ValidateEncodingParams(params EncodingParams, blobLength, SRSOrder int) error {

	if int(params.ChunkLen*params.NumChunks) >= SRSOrder {
		return fmt.Errorf("the supplied encoding parameters are not valid with respect to the SRS. ChunkLength: %d, NumChunks: %d, SRSOrder: %d", params.ChunkLen, params.NumChunks, SRSOrder)
	}

	if int(params.ChunkLen*params.NumChunks) < blobLength {
		return fmt.Errorf("the supplied encoding parameters are not sufficient for the size of the data input")
	}

	return nil

}
