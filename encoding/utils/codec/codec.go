package codec

import (
	"github.com/Layr-Labs/eigenda/encoding"
)

// ConvertByPaddingEmptyByte takes bytes and insert an empty at the front of every 31 byte
// This ensure every 32 bytes are within the valid range of a field element for bn254 curve
// If the input data is not a multiple of 31, the reminder is added to the output by
// inserting a 0 and the reminder. The output does not necessarily be a multipler of 32
func ConvertByPaddingEmptyByte(data []byte) []byte {
	dataSize := len(data)
	parseSize := encoding.BYTES_PER_COEFFICIENT - 1
	putSize := uint64(encoding.BYTES_PER_COEFFICIENT)

	dataLen := uint64((dataSize + parseSize - 1) / parseSize)

	validData := make([]byte, dataLen*putSize)
	validEnd := uint64(len(validData))

	for i := uint64(0); i < uint64(dataLen); i++ {
		start := i * uint64(parseSize)
		end := (i + 1) * uint64(parseSize)
		if end > uint64(len(data)) {
			end = uint64(len(data))
			// 1 is the empty byte
			validEnd = end - start + 1 + i*putSize
		}

		// with big endian, set first byte is always 0 to ensure data is within valid range of
		validData[i*encoding.BYTES_PER_COEFFICIENT] = 0x00
		copy(validData[i*encoding.BYTES_PER_COEFFICIENT+1:(i+1)*encoding.BYTES_PER_COEFFICIENT], data[start:end])

	}
	return validData[:validEnd]
}

// RemoveEmptyByteFromPaddedBytes takes bytes and remove the first byte from every 32 bytes.
// This reverses the change made by the function ConvertByPaddingEmptyByte.
// The function does not assume the input is a multiple of BYTES_PER_COEFFICIENT(32 bytes).
// For the reminder of the input, the first byte is taken out, and the rest is appended to
// the output.
func RemoveEmptyByteFromPaddedBytes(data []byte) []byte {
	dataSize := len(data)
	parseSize := encoding.BYTES_PER_COEFFICIENT
	dataLen := uint64((dataSize + parseSize - 1) / parseSize)

	putSize := uint64(encoding.BYTES_PER_COEFFICIENT - 1)

	validData := make([]byte, dataLen*putSize)
	validLen := int64(len(validData))

	for i := uint64(0); i < uint64(dataLen); i++ {
		// add 1 to leave the first empty byte untouched
		start := i*uint64(parseSize) + 1
		end := (i + 1) * uint64(parseSize)

		if end > uint64(len(data)) {
			end = uint64(len(data))
			validLen = int64(end - start + i*putSize)
		}

		copy(validData[i*putSize:(i+1)*putSize], data[start:end])
	}
	return validData[:validLen]
}
