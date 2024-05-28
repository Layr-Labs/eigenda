package codecs

import (
	"encoding/binary"
	"fmt"
)

func EncodeCodecBlobHeader(version byte, length uint64) []byte {
	codecBlobHeader := make([]byte, 32)
	// the first byte is always 0 so we are always smaller than the field modulo

	// encode version byte
	codecBlobHeader[1] = version

	// encode length as uint64
	binary.BigEndian.PutUint64(codecBlobHeader[2:10], length) // uint64 should be more than enough to store the length
	return codecBlobHeader
}

func DecodeCodecBlobHeader(codecBlobHeader []byte) (byte, uint64, error) {
	// make sure the codecBlobHeader is 32 bytes long
	if len(codecBlobHeader) != 32 {
		err := fmt.Errorf("codecBlobHeader must be exactly 32 bytes long, but got %d bytes", len(codecBlobHeader))
		return 0, 0, err
	}
	// make sure first byte is 0
	if codecBlobHeader[0] != 0 {
		err := fmt.Errorf("codecBlobHeader must start with 0, but got %d", codecBlobHeader[0])
		return 0, 0, err
	}

	version := codecBlobHeader[1]
	length := binary.BigEndian.Uint64(codecBlobHeader[2:10])

	return version, length, nil
}
