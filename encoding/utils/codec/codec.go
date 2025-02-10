package codec

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
)

// ConvertByPaddingEmptyByte takes bytes and insert an empty byte at the front of every 31 byte.
// The empty byte is padded at the low address, because we use big endian to interpret a field element.
// This ensures every 32 bytes is within the valid range of a field element for bn254 curve.
// If the input data is not a multiple of 31, the remainder is added to the output by
// inserting a 0 and the remainder. The output is thus not necessarily a multiple of 32.
//
// TODO (litt3): usage of this function should be migrated to use PadPayload instead. I've left it unchanged for now,
//  since v1 logic and tests rely on the specific assumptions of this implementation.
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
//
// TODO (litt3): usage of this function should be migrated to use RemoveInternalPadding instead. I've left it unchanged
//  for now, since v1 logic and tests rely on the specific assumptions of this implementation.
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

// PadPayload internally pads the input data by prepending a 0x00 to each chunk of 31 bytes. This guarantees that
// the data will be a valid field element for the bn254 curve
//
// Additionally, this function will add necessary padding to align the output to 32 bytes
//
// NOTE: this method is a reimplementation of ConvertByPaddingEmptyByte, with one meaningful difference: the alignment
// of the output to encoding.BYTES_PER_SYMBOL. This alignment actually makes the padding logic simpler, and the
// code that uses this function needs an aligned output anyway.
func PadPayload(inputData []byte) []byte {
	// 31 bytes, for the bn254 curve
	bytesPerChunk := uint32(encoding.BYTES_PER_SYMBOL - 1)

	// this is the length of the output, which is aligned to 32 bytes
	outputLength := GetPaddedDataLength(uint32(len(inputData)))
	paddedOutput := make([]byte, outputLength)

	// pre-pad the input, so that it aligns to 31 bytes. This means that the internally padded result will automatically
	// align to 32 bytes. Doing this padding in advance simplifies the for loop.
	requiredPad := (bytesPerChunk - uint32(len(inputData))%bytesPerChunk) % bytesPerChunk
	prePaddedPayload := append(inputData, make([]byte, requiredPad)...)

	for element := uint32(0); element < outputLength/encoding.BYTES_PER_SYMBOL; element++ {
		// add the 0x00 internal padding to guarantee that the data is in the valid range
		zeroByteIndex := element * encoding.BYTES_PER_SYMBOL
		paddedOutput[zeroByteIndex] = 0x00

		destIndex := zeroByteIndex + 1
		srcIndex := element * bytesPerChunk

		// copy 31 bytes of data from the payload to the padded output
		copy(paddedOutput[destIndex:destIndex+bytesPerChunk], prePaddedPayload[srcIndex:srcIndex+bytesPerChunk])
	}

	return paddedOutput
}

// RemoveInternalPadding accepts an array of padded data, and removes the internal padding that was added in PadPayload
//
// This function assumes that the input aligns to 32 bytes. Since it is removing 1 byte for every 31 bytes kept, the
// output from this function is not guaranteed to align to 32 bytes.
//
// NOTE: this method is a reimplementation of RemoveEmptyByteFromPaddedBytes, with one meaningful difference: this
// function relies on the assumption that the input is aligned to encoding.BYTES_PER_SYMBOL, which makes the padding
// removal logic simpler.
func RemoveInternalPadding(paddedData []byte) ([]byte, error) {
	if len(paddedData)%encoding.BYTES_PER_SYMBOL != 0 {
		return nil, fmt.Errorf(
			"padded data (length %d) must be multiple of encoding.BYTES_PER_SYMBOL %d",
			len(paddedData),
			encoding.BYTES_PER_SYMBOL)
	}

	bytesPerChunk := encoding.BYTES_PER_SYMBOL - 1

	symbolCount := len(paddedData) / encoding.BYTES_PER_SYMBOL
	outputLength := symbolCount * bytesPerChunk

	outputData := make([]byte, outputLength)

	for i := 0; i < symbolCount; i++ {
		dstIndex := i * bytesPerChunk
		srcIndex := i*encoding.BYTES_PER_SYMBOL + 1

		copy(outputData[dstIndex:dstIndex+bytesPerChunk], paddedData[srcIndex:srcIndex+bytesPerChunk])
	}

	return outputData, nil
}

// GetPaddedDataLength accepts the length of a byte array, and returns the length that the array would be after
// adding internal byte padding
//
// The value returned from this function will always be a multiple of encoding.BYTES_PER_SYMBOL
func GetPaddedDataLength(inputLen uint32) uint32 {
	bytesPerChunk := uint32(encoding.BYTES_PER_SYMBOL - 1)
	chunkCount := inputLen / bytesPerChunk

	if inputLen%bytesPerChunk != 0 {
		chunkCount++
	}

	return chunkCount * encoding.BYTES_PER_SYMBOL
}
