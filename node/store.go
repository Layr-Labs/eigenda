package node

import (
	"encoding/binary"
	"errors"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// parseHeader parses the header and returns the encoding format and the chunk length.
func parseHeader(data []byte) (core.BundleEncodingFormat, uint64, error) {
	if len(data) < 8 {
		return 0, 0, errors.New("no header found, the data size is less 8 bytes")
	}
	meta := binary.LittleEndian.Uint64(data)
	format := binary.LittleEndian.Uint64(data) >> (core.NumBundleHeaderBits - core.NumBundleEncodingFormatBits)
	chunkLen := (meta << core.NumBundleEncodingFormatBits) >> core.NumBundleEncodingFormatBits
	return uint8(format), chunkLen, nil
}

// EncodeChunks flattens an array of byte arrays (chunks) into a single byte array.
// EncodeChunks(chunks) = (len(chunks[0]), chunks[0], len(chunks[1]), chunks[1], ...)
func EncodeChunks(chunks [][]byte) ([]byte, error) {
	totalSize := 0
	for _, chunk := range chunks {
		totalSize += len(chunk) + 8 // Add size of uint64 for length
	}
	result := make([]byte, totalSize)
	buf := result
	for _, chunk := range chunks {
		binary.LittleEndian.PutUint64(buf, uint64(len(chunk)))
		buf = buf[8:]
		copy(buf, chunk)
		buf = buf[len(chunk):]
	}
	return result, nil
}

// DecodeChunks converts a flattened array of chunks into an array of its constituent chunks,
// throwing an error in case the chunks were not serialized correctly.
func DecodeChunks(data []byte) ([][]byte, node.ChunkEncodingFormat, error) {
	// Empty chunk is valid, but there is nothing to decode.
	if len(data) == 0 {
		return [][]byte{}, node.ChunkEncodingFormat_UNKNOWN, nil
	}
	format, _, err := parseHeader(data)
	if err != nil {
		return nil, node.ChunkEncodingFormat_UNKNOWN, err
	}

	// Note: the encoding format IDs may not be the same as the field ID in protobuf.
	// For example, GobBundleEncodingFormat is 1 but node.ChunkEncodingFormat_GOB has proto
	// field ID 2.
	switch format {
	case 0:
		chunks, err := DecodeGobChunks(data)
		return chunks, node.ChunkEncodingFormat_GOB, err
	case 1:
		chunks, err := DecodeGnarkChunks(data)
		return chunks, node.ChunkEncodingFormat_GNARK, err
	default:
		return nil, node.ChunkEncodingFormat_UNKNOWN, errors.New("invalid data encoding format")
	}
}

// DecodeGobChunks decodes chunks in GOB format.
// DecodeGobChunks((len(chunks[0]), chunks[0], len(chunks[1]), chunks[1], ...)) = chunks
func DecodeGobChunks(data []byte) ([][]byte, error) {
	format, chunkLen, err := parseHeader(data)
	if err != nil {
		return nil, err
	}
	if format != core.GobBundleEncodingFormat {
		return nil, errors.New("invalid bundle data encoding format")
	}
	if chunkLen == 0 {
		return nil, errors.New("chunk length must be greater than zero")
	}
	chunks := make([][]byte, 0)
	buf := data
	for len(buf) > 0 {
		if len(buf) < 8 {
			return nil, errors.New("invalid data to decode")
		}
		chunkSize := binary.LittleEndian.Uint64(buf)
		buf = buf[8:]

		if len(buf) < int(chunkSize) {
			return nil, errors.New("invalid data to decode")
		}
		chunks = append(chunks, buf[:chunkSize])
		buf = buf[chunkSize:]
	}
	return chunks, nil
}

// DecodeGnarkChunks decodes chunks in Gnark format.
func DecodeGnarkChunks(data []byte) ([][]byte, error) {
	format, chunkLen, err := parseHeader(data)
	if err != nil {
		return nil, err
	}
	if format != core.GnarkBundleEncodingFormat {
		return nil, errors.New("invalid bundle data encoding format")
	}
	if chunkLen == 0 {
		return nil, errors.New("chunk length must be greater than zero")
	}
	chunkSize := bn254.SizeOfG1AffineCompressed + encoding.BYTES_PER_SYMBOL*int(chunkLen)
	chunks := make([][]byte, 0)
	buf := data[8:]
	for len(buf) > 0 {
		if len(buf) < chunkSize {
			return nil, errors.New("invalid data to decode")
		}
		chunks = append(chunks, buf[:chunkSize])
		buf = buf[chunkSize:]
	}
	return chunks, nil
}
