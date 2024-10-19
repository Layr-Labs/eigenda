package codec

import (
	"github.com/Layr-Labs/eigenda/encoding"
)

// ConvertByPaddingEmptyByte takes bytes and insert an empty byte at the front of every 31 byte.
// The empty byte is padded at the low address, because we use big endian to interpret a field element.
// This ensures every 32 bytes is within the valid range of a field element for bn254 curve.
// If the input data is not a multiple of 31, the remainder is added to the output by
// inserting a 0 and the remainder. The output is thus not necessarily a multiple of 32.
func ConvertByPaddingEmptyByte(data []byte) []byte {
	dataSize := len(data)
	parseSize := encoding.BYTES_PER_SYMBOL - 1
	putSize := encoding.BYTES_PER_SYMBOL

	dataLen := (dataSize + parseSize - 1) / parseSize

	validData := make([]byte, dataLen*putSize)
	validEnd := len(validData)

	for i := 0; i < dataLen; i++ {
		start := i * parseSize
		end := (i + 1) * parseSize
		if end > len(data) {
			end = len(data)
			// 1 is the empty byte
			validEnd = end - start + 1 + i*putSize
		}

		// with big endian, set first byte is always 0 to ensure data is within valid range of
		validData[i*encoding.BYTES_PER_SYMBOL] = 0x00
		copy(validData[i*encoding.BYTES_PER_SYMBOL+1:(i+1)*encoding.BYTES_PER_SYMBOL], data[start:end])

	}
	return validData[:validEnd]
}

// RemoveEmptyByteFromPaddedBytes takes bytes and remove the first byte from every 32 bytes.
// This reverses the change made by the function ConvertByPaddingEmptyByte.
// The function does not assume the input is a multiple of BYTES_PER_SYMBOL(32 bytes).
// For the reminder of the input, the first byte is taken out, and the rest is appended to
// the output.
func RemoveEmptyByteFromPaddedBytes(data []byte) []byte {
	dataSize := len(data)
	parseSize := encoding.BYTES_PER_SYMBOL
	dataLen := (dataSize + parseSize - 1) / parseSize

	putSize := encoding.BYTES_PER_SYMBOL - 1

	validData := make([]byte, dataLen*putSize)
	validLen := len(validData)

	for i := 0; i < dataLen; i++ {
		// add 1 to leave the first empty byte untouched
		start := i*parseSize + 1
		end := (i + 1) * parseSize

		if end > len(data) {
			end = len(data)
			validLen = end - start + i*putSize
		}

		copy(validData[i*putSize:(i+1)*putSize], data[start:end])
	}
	return validData[:validLen]
}
