package dataapi

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

func ConvertHexadecimalToBytes(byteHash []byte) ([32]byte, error) {
	hexString := strings.TrimPrefix(string(byteHash), "0x")

	// Now decode the hex string to bytes
	decodedBytes, err := hex.DecodeString(hexString)
	if err != nil {
		return [32]byte{}, err
	}

	// We expect the resulting byte slice to have a length of 32 bytes.
	if len(decodedBytes) != 32 {
		return [32]byte{}, errors.New("error decoding hash")
	}

	// Convert the byte slice to a [32]byte array
	var byteArray [32]byte
	copy(byteArray[:], decodedBytes[:32])

	return byteArray, nil
}

func ConvertNanosecondToSecond(timestamp uint64) uint64 {
	return timestamp / uint64(time.Second)
}
