package rs

import (
	"errors"

	"github.com/Layr-Labs/eigenda/encoding"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

var (
	ErrInvalidParams = errors.New("invalid encoding params")
)

type EncodingParams encoding.EncodingParams

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

func GetNumSys(dataSize uint64, chunkLen uint64) uint64 {
	dataLen := RoundUpDivision(dataSize, bls.BYTES_PER_COEFFICIENT)
	numSys := dataLen / chunkLen
	return numSys
}

func ParamsFromMins(numChunks, chunkLen uint64) EncodingParams {

	chunkLen = NextPowerOf2(chunkLen)
	numChunks = NextPowerOf2(numChunks)

	return EncodingParams{
		NumChunks:   numChunks,
		ChunkLength: chunkLen,
	}

}

func ParamsFromSysPar(numSys, numPar, dataSize uint64) EncodingParams {

	numNodes := numSys + numPar
	dataLen := RoundUpDivision(dataSize, bls.BYTES_PER_COEFFICIENT)
	chunkLen := RoundUpDivision(dataLen, numSys)
	return ParamsFromMins(numNodes, chunkLen)

}
