package codecs

import (
	"encoding/binary"
	"fmt"
)

func EncodeCodecBlobHeader(version byte, length uint32) []byte {
	codecBlobHeader := make([]byte, 32)
	// first byte is always 0 to ensure the codecBlobHeader is a valid bn254 element
	// encode version byte
	codecBlobHeader[1] = version

	// encode length as uint32
	binary.BigEndian.PutUint32(codecBlobHeader[2:6], length) // uint32 should be more than enough to store the length (approx 4gb)
	return codecBlobHeader
}

func DecodeCodecBlobHeader(codecBlobHeader []byte) (byte, uint32, error) {
	// make sure the codecBlobHeader is 32 bytes long
	if len(codecBlobHeader) != 32 {
		err := fmt.Errorf("codecBlobHeader must be exactly 32 bytes long, but got %d bytes", len(codecBlobHeader))
		return 0, 0, err
	}

	version := codecBlobHeader[1]
	length := binary.BigEndian.Uint32(codecBlobHeader[2:6])

	return version, length, nil
}
