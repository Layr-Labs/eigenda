package node_test

import (
	"encoding/binary"
	"testing"

	nodegrpc "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeGobData constructs valid GOB-format encoded data from the given chunks.
// In GOB format (format=0), each chunk is preceded by its uint64 length.
// The first 8 bytes double as the header (format=0, chunkLen=first chunk size).
func makeGobData(chunks [][]byte) []byte {
	data := make([]byte, 0)
	for _, chunk := range chunks {
		length := make([]byte, 8)
		binary.LittleEndian.PutUint64(length, uint64(len(chunk)))
		data = append(data, length...)
		data = append(data, chunk...)
	}
	return data
}

// makeGnarkData constructs valid Gnark-format encoded data from the given chunkLen and chunk count.
// In Gnark format (format=1), the header has format=1 in the top byte and chunkLen in the lower bytes.
// Each chunk is exactly SizeOfG1AffineCompressed + BYTES_PER_SYMBOL*chunkLen bytes.
func makeGnarkData(chunkLen uint64, numChunks int) []byte {
	header := make([]byte, 8)
	val := (uint64(1) << 56) | chunkLen
	binary.LittleEndian.PutUint64(header, val)

	chunkSize := bn254.SizeOfG1AffineCompressed + encoding.BYTES_PER_SYMBOL*int(chunkLen)
	data := make([]byte, 8+chunkSize*numChunks)
	copy(data, header)
	// Fill chunk data with non-zero values for verifiability.
	for i := 8; i < len(data); i++ {
		data[i] = byte(i % 251)
	}
	return data
}

// --- EncodeChunks ---

func TestEncodeChunks(t *testing.T) {
	chunks := [][]byte{
		{1, 2, 3},
		{4, 5},
		{6, 7, 8, 9},
	}
	encoded, err := node.EncodeChunks(chunks)
	require.NoError(t, err)

	// 3 length prefixes (3*8=24) + data (3+2+4=9) = 33 bytes total
	assert.Len(t, encoded, 33)

	off := 0
	for _, chunk := range chunks {
		size := binary.LittleEndian.Uint64(encoded[off : off+8])
		assert.Equal(t, uint64(len(chunk)), size)
		assert.Equal(t, chunk, encoded[off+8:off+8+len(chunk)])
		off += 8 + len(chunk)
	}
}

func TestEncodeChunksEmpty(t *testing.T) {
	encoded, err := node.EncodeChunks([][]byte{})
	require.NoError(t, err)
	assert.Empty(t, encoded)
}

func TestEncodeChunksSingleEmpty(t *testing.T) {
	encoded, err := node.EncodeChunks([][]byte{{}})
	require.NoError(t, err)
	// 8 bytes for the length prefix (value 0), no chunk data.
	assert.Len(t, encoded, 8)
	assert.Equal(t, uint64(0), binary.LittleEndian.Uint64(encoded))
}

// --- DecodeChunks ---

func TestDecodeChunksEmpty(t *testing.T) {
	chunks, format, err := node.DecodeChunks([]byte{})
	require.NoError(t, err)
	assert.Equal(t, nodegrpc.ChunkEncodingFormat_UNKNOWN, format)
	assert.Empty(t, chunks)
}

func TestDecodeChunksTooShort(t *testing.T) {
	_, _, err := node.DecodeChunks([]byte{1, 2, 3})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "less 8 bytes")
}

func TestDecodeChunksGobFormat(t *testing.T) {
	chunkData := []byte{0xAA, 0xBB, 0xCC}
	data := makeGobData([][]byte{chunkData})

	chunks, format, err := node.DecodeChunks(data)
	require.NoError(t, err)
	assert.Equal(t, nodegrpc.ChunkEncodingFormat_GOB, format)
	require.Len(t, chunks, 1)
	assert.Equal(t, chunkData, chunks[0])
}

func TestDecodeChunksGnarkFormat(t *testing.T) {
	data := makeGnarkData(1, 1)

	chunks, format, err := node.DecodeChunks(data)
	require.NoError(t, err)
	assert.Equal(t, nodegrpc.ChunkEncodingFormat_GNARK, format)
	require.Len(t, chunks, 1)
}

func TestDecodeChunksInvalidFormat(t *testing.T) {
	header := make([]byte, 8)
	val := (uint64(2) << 56) | 1 // format=2 is invalid
	binary.LittleEndian.PutUint64(header, val)

	_, _, err := node.DecodeChunks(header)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid data encoding format")
}

// --- DecodeGobChunks ---

func TestDecodeGobChunksSingle(t *testing.T) {
	chunkData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	data := makeGobData([][]byte{chunkData})

	chunks, err := node.DecodeGobChunks(data)
	require.NoError(t, err)
	require.Len(t, chunks, 1)
	assert.Equal(t, chunkData, chunks[0])
}

func TestDecodeGobChunksMultiple(t *testing.T) {
	chunk1 := []byte{1, 2, 3}
	chunk2 := []byte{4, 5}
	chunk3 := []byte{6, 7, 8, 9, 10, 11}
	data := makeGobData([][]byte{chunk1, chunk2, chunk3})

	chunks, err := node.DecodeGobChunks(data)
	require.NoError(t, err)
	require.Len(t, chunks, 3)
	assert.Equal(t, chunk1, chunks[0])
	assert.Equal(t, chunk2, chunks[1])
	assert.Equal(t, chunk3, chunks[2])
}

func TestDecodeGobChunksWrongFormat(t *testing.T) {
	// Use Gnark header (format=1) — should fail format check.
	data := makeGnarkData(1, 1)
	_, err := node.DecodeGobChunks(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid bundle data encoding format")
}

func TestDecodeGobChunksZeroChunkLen(t *testing.T) {
	// Header with format=0 and chunkLen=0.
	header := make([]byte, 8)
	binary.LittleEndian.PutUint64(header, 0)

	_, err := node.DecodeGobChunks(header)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chunk length must be greater than zero")
}

func TestDecodeGobChunksTruncatedChunkData(t *testing.T) {
	// Header says first chunk is 100 bytes but only 5 bytes follow.
	header := make([]byte, 8)
	binary.LittleEndian.PutUint64(header, 100)
	data := append(header, make([]byte, 5)...)

	_, err := node.DecodeGobChunks(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid data to decode")
}

func TestDecodeGobChunksPartialSecondHeader(t *testing.T) {
	// Valid first chunk followed by 3 trailing bytes (not enough for a length prefix).
	chunkData := []byte{0x01, 0x02}
	data := makeGobData([][]byte{chunkData})
	data = append(data, []byte{0xFF, 0xFF, 0xFF}...)

	_, err := node.DecodeGobChunks(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid data to decode")
}

func TestDecodeGobChunksTooShort(t *testing.T) {
	_, err := node.DecodeGobChunks([]byte{1, 2})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "less 8 bytes")
}

// --- DecodeGnarkChunks ---

func TestDecodeGnarkChunksSingle(t *testing.T) {
	data := makeGnarkData(1, 1)

	chunks, err := node.DecodeGnarkChunks(data)
	require.NoError(t, err)
	require.Len(t, chunks, 1)

	expectedChunkSize := bn254.SizeOfG1AffineCompressed + encoding.BYTES_PER_SYMBOL
	assert.Len(t, chunks[0], expectedChunkSize)
}

func TestDecodeGnarkChunksMultiple(t *testing.T) {
	data := makeGnarkData(2, 3)

	chunks, err := node.DecodeGnarkChunks(data)
	require.NoError(t, err)
	require.Len(t, chunks, 3)

	expectedChunkSize := bn254.SizeOfG1AffineCompressed + encoding.BYTES_PER_SYMBOL*2
	for i, chunk := range chunks {
		assert.Len(t, chunk, expectedChunkSize, "chunk %d has wrong size", i)
	}
}

func TestDecodeGnarkChunksWrongFormat(t *testing.T) {
	// Use GOB header (format=0) — should fail format check.
	chunkData := []byte{0x01, 0x02, 0x03}
	data := makeGobData([][]byte{chunkData})

	_, err := node.DecodeGnarkChunks(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid bundle data encoding format")
}

func TestDecodeGnarkChunksZeroChunkLen(t *testing.T) {
	header := make([]byte, 8)
	val := uint64(1) << 56 // format=1, chunkLen=0
	binary.LittleEndian.PutUint64(header, val)

	_, err := node.DecodeGnarkChunks(header)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chunk length must be greater than zero")
}

func TestDecodeGnarkChunksTruncated(t *testing.T) {
	// chunkLen=1 means each chunk should be 64 bytes, but only provide 10.
	header := make([]byte, 8)
	val := (uint64(1) << 56) | 1
	binary.LittleEndian.PutUint64(header, val)
	data := append(header, make([]byte, 10)...)

	_, err := node.DecodeGnarkChunks(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid data to decode")
}

func TestDecodeGnarkChunksTooShort(t *testing.T) {
	_, err := node.DecodeGnarkChunks([]byte{1, 2})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "less 8 bytes")
}
