package codecs

import (
	"encoding/binary"
	"fmt"
)

func EncodeCodecBlobHeader(version byte, length uint32) []byte {
	codecBlobHeader := make([]byte, 5)
	// the first byte is always 0 so we are always smaller than the field modulo

	// encode version byte
	codecBlobHeader[0] = version

	// encode length as uint32
	binary.BigEndian.PutUint32(codecBlobHeader[1:5], length) // uint32 should be more than enough to store the length (approx 4gb)
	return codecBlobHeader
}

func DecodeCodecBlobHeader(codecBlobHeader []byte) (byte, uint32, error) {
	// make sure the codecBlobHeader is 5 bytes long
	if len(codecBlobHeader) != 5 {
		err := fmt.Errorf("codecBlobHeader must be exactly 5 bytes long, but got %d bytes", len(codecBlobHeader))
		return 0, 0, err
	}

	version := codecBlobHeader[0]
	length := binary.BigEndian.Uint32(codecBlobHeader[1:5])

	return version, length, nil
}
