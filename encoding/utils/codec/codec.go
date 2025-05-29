package codec

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
)

// BlobHeaderSize is the number of bytes needed for a blob header.
const BlobHeaderSize = uint32(32)

// ConvertByPaddingEmptyByte takes bytes and insert an empty byte at the front of every 31 byte.
// The empty byte is padded at the low address, because we use big endian to interpret a field element.
// This ensures every 32 bytes is within the valid range of a field element for bn254 curve.
// If the input data is not a multiple of 31, the remainder is added to the output by
// inserting a 0 and the remainder. The output is thus not necessarily a multiple of 32.
//
// TODO (litt3): usage of this function should be migrated to use PadPayload instead. I've left it unchanged for now,
//
//	since v1 logic and tests rely on the specific assumptions of this implementation.
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
//
//	for now, since v1 logic and tests rely on the specific assumptions of this implementation.
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

// PayloadSizeToBlobSize takes a payload size in bytes and returns the corresponding blob size in bytes.
// The blob size is the size used for determining payments and throttling by EigenDA. Two payloads of
// differing length that have the same blob size cost the same and use the same amount of bandwidth.
func PayloadSizeToBlobSize(payloadSize uint32) uint32 {
	return encoding.NextPowerOf2(GetPaddedDataLength(payloadSize) + BlobHeaderSize)
}

// FindLegalBlobSizes finds a list of blob sizes that are legal for EigenDA. A legal blob size is
// a blob size that is a power of 2 and is between the minimum and maximum blob sizes (inclusive).
func FindLegalBlobSizes(minBlobSize uint32, maxBlobSize uint32) ([]uint32, error) {
	if minBlobSize > maxBlobSize {
		return nil, fmt.Errorf("min blob size %d is greater than max blob size %d", minBlobSize, maxBlobSize)
	}
	if !encoding.IsPowerOfTwo(minBlobSize) {
		return nil, fmt.Errorf("min blob size %d is not a power of 2", minBlobSize)
	}
	if !encoding.IsPowerOfTwo(maxBlobSize) {
		return nil, fmt.Errorf("max blob size %d is not a power of 2", maxBlobSize)
	}

	sizes := make([]uint32, 0)

	for i := minBlobSize; i <= maxBlobSize; i *= 2 {
		sizes = append(sizes, i)
	}

	return sizes, nil
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

// GetUnpaddedDataLength accepts the length of an array that has been padded with PadPayload
//
// It returns what the length of the output array would be, if you called RemoveInternalPadding on it.
func GetUnpaddedDataLength(inputLen uint32) (uint32, error) {
	if inputLen%encoding.BYTES_PER_SYMBOL != 0 {
		return 0, fmt.Errorf(
			"%d isn't a multiple of encoding.BYTES_PER_SYMBOL (%d)",
			inputLen, encoding.BYTES_PER_SYMBOL)
	}

	chunkCount := inputLen / encoding.BYTES_PER_SYMBOL
	bytesPerChunk := uint32(encoding.BYTES_PER_SYMBOL - 1)

	unpaddedLength := chunkCount * bytesPerChunk

	return unpaddedLength, nil
}

// BlobSymbolsToMaxPayloadSize accepts a blob length in symbols and returns the size in bytes of the largest payload
// that could fit inside the blob.
func BlobSymbolsToMaxPayloadSize(blobLengthSymbols uint32) (uint32, error) {
	if blobLengthSymbols == 0 {
		return 0, fmt.Errorf("input blobLengthSymbols is zero")
	}

	if !encoding.IsPowerOfTwo(uint64(blobLengthSymbols)) {
		return 0, fmt.Errorf("blobLengthSymbols %d is not a power of two", blobLengthSymbols)
	}

	maxPayloadLength, err := GetUnpaddedDataLength(blobLengthSymbols*encoding.BYTES_PER_SYMBOL - BlobHeaderSize)
	if err != nil {
		return 0, fmt.Errorf("get unpadded data length: %w", err)
	}

	return maxPayloadLength, nil
}

// BlobSizeToMaxPayloadSize accepts a blob length in bytes and returns the size in bytes of the largest payload
// that could fit inside the blob.
func BlobSizeToMaxPayloadSize(blobLengthBytes uint32) (uint32, error) {
	return BlobSymbolsToMaxPayloadSize(blobLengthBytes / encoding.BYTES_PER_SYMBOL)
}

// FindMaxPayloadSizes finds a list of payload sizes that are as large as possible for a given blob size.
// Increasing the size of a maximum payload by a single byte will result in a blob that is the next tier larger.
func FindMaxPayloadSizes(minBlobSize uint32, maxBlobSize uint32) ([]uint32, error) {
	legalBlobSizes, err := FindLegalBlobSizes(minBlobSize, maxBlobSize)
	if err != nil {
		return nil, fmt.Errorf("failed to find legal blob sizes: %w", err)
	}

	sizes := make([]uint32, 0, len(legalBlobSizes))

	for _, blobSize := range legalBlobSizes {
		maxPayloadSize, err := BlobSizeToMaxPayloadSize(blobSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get maximum payload size for blob size %d: %w", blobSize, err)
		}
		sizes = append(sizes, maxPayloadSize)
	}

	return sizes, nil
}

// BlobSizeToMinPayloadSize takes a given a blob size and determines the minimum payload size
// that yields that blob size.
func BlobSizeToMinPayloadSize(blobSize uint32) (uint32, error) {
	if !encoding.IsPowerOfTwo(blobSize) {
		return 0, fmt.Errorf("blob size %d is not a power of 2", blobSize)
	}

	paddedLength := blobSize/2 - BlobHeaderSize + 1

	payloadSizeAdjustment := uint32(0)
	if paddedLength%encoding.BYTES_PER_SYMBOL != 0 {
		// If the padded length is not a multiple of BYTES_PER_SYMBOL, this means that there is a "partial" symbol.
		// That is to say, we don't need all the bytes in the last symbol to represent the data. Subtract away
		// this partial symbol before converting to unpadded size, then add 1 byte to the final answer to determine the
		// minimum size required to result in the partial symbol that we subtract in this step.
		payloadSizeAdjustment = 1
		paddedLength -= paddedLength % encoding.BYTES_PER_SYMBOL
	}

	size, err := GetUnpaddedDataLength(paddedLength)
	if err != nil {
		return 0, fmt.Errorf("get unpadded data length: %w", err)
	}

	size += payloadSizeAdjustment

	return size, nil
}

// FindMinPayloadSizes finds a list of payload sizes that are the minimum possible payload size for a given blob size.
// Decreasing the size of a minimum payload by a single byte will result in a blob that is the next tier smaller.
func FindMinPayloadSizes(minBlobSize uint32, maxBlobSize uint32) ([]uint32, error) {
	legalBlobSizes, err := FindLegalBlobSizes(minBlobSize, maxBlobSize)
	if err != nil {
		return nil, fmt.Errorf("failed to find legal blob sizes: %w", err)
	}

	sizes := make([]uint32, 0, len(legalBlobSizes))

	for _, blobSize := range legalBlobSizes {
		minPayloadSize, err := BlobSizeToMinPayloadSize(blobSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get minimum payload size for blob size %d: %w", blobSize, err)
		}
		sizes = append(sizes, minPayloadSize)
	}

	return sizes, nil
}
