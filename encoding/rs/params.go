package rs

import (
	"errors"

	"github.com/Layr-Labs/eigenda/encoding"
)

var (
	ErrInvalidParams = errors.New("invalid encoding params")
)

type EncodingParams struct {
	NumChunks uint64 // number of total chunks that are padded to power of 2
	ChunkLen  uint64 // number of Fr symbol stored inside a chunk
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
	dataLen := RoundUpDivision(dataSize, encoding.BYTES_PER_SYMBOL)
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
	dataLen := RoundUpDivision(dataSize, encoding.BYTES_PER_SYMBOL)
	chunkLen := RoundUpDivision(dataLen, numSys)

	return ParamsFromMins(numNodes, chunkLen)
}
