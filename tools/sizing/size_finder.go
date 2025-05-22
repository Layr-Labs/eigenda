package sizing

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
)

// next_power_of_2(len(payload)*32/31)

// BlobHeaderSize is the size of the header data on a blob, in bytes.
const BlobHeaderSize = uint64(31)

// PayloadSizeToBlobSize takes a payload size in bytes and returns the corresponding blob size in bytes.
// The blob size is the size used for determining payments and throttling by EigenDA. Two payloads of
// differing length that have the same blob size cost the same and use the same amount of bandwidth.
func PayloadSizeToBlobSize(payloadSize uint64) uint64 {
	return encoding.NextPowerOf2(payloadSize*32/31 + BlobHeaderSize)
}

// BlobSizeToMaxPayloadSize takes a given a blob size and determines the maximum payload size
// that yields that blob size.
func BlobSizeToMaxPayloadSize(blobSize uint64) (uint64, error) {
	if !encoding.IsPowerOf2(blobSize) {
		return 0, fmt.Errorf("blob size %d is not a power of 2", blobSize)
	}

	return (blobSize - BlobHeaderSize) * 31 / 32, nil
}

// BlobSizeToMinPayloadSize takes a given a blob size and determines the minimum payload size
// that yields that blob size.
func BlobSizeToMinPayloadSize(blobSize uint64) (uint64, error) {
	if !encoding.IsPowerOf2(blobSize) {
		return 0, fmt.Errorf("blob size %d is not a power of 2", blobSize)
	}

	return (blobSize/2 - BlobHeaderSize + 1) * 31 / 32, nil
}

// FindLegalBlobSizes finds a list of blob sizes that are legal for EigenDA. A legal blob size is
// a blob size that is a power of 2 and is between the minimum and maximum blob sizes (inclusive).
func FindLegalBlobSizes(minBlobSize uint64, maxBlobSize uint64) ([]uint64, error) {
	if minBlobSize > maxBlobSize {
		return nil, fmt.Errorf("min blob size %d is greater than max blob size %d", minBlobSize, maxBlobSize)
	}
	if !encoding.IsPowerOf2(minBlobSize) {
		return nil, fmt.Errorf("min blob size %d is not a power of 2", minBlobSize)
	}
	if !encoding.IsPowerOf2(maxBlobSize) {
		return nil, fmt.Errorf("max blob size %d is not a power of 2", maxBlobSize)
	}

	sizes := make([]uint64, 0)

	for i := minBlobSize; i <= maxBlobSize; i *= 2 {
		sizes = append(sizes, i)
	}

	return sizes, nil
}

// FindMaxPayloadSizes finds a list of payload sizes that are as large as possible for a given blob size.
// Increasing the size of a maximum payload by a single byte will result in a blob that is the next tier larger.
func FindMaxPayloadSizes(minBlobSize uint64, maxBlobSize uint64) ([]uint64, error) {
	legalBlobSizes, err := FindLegalBlobSizes(minBlobSize, maxBlobSize)
	if err != nil {
		return nil, fmt.Errorf("failed to find legal blob sizes: %w", err)
	}

	sizes := make([]uint64, 0, len(legalBlobSizes))

	for _, blobSize := range legalBlobSizes {
		maxPayloadSize, err := BlobSizeToMaxPayloadSize(blobSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get maximum payload size for blob size %d: %w", blobSize, err)
		}
		sizes = append(sizes, maxPayloadSize)
	}

	return sizes, nil
}

// FindMinPayloadSizes finds a list of payload sizes that are the minimum possible payload size for a given blob size.
// Decreasing the size of a minimum payload by a single byte will result in a blob that is the next tier smaller.
func FindMinPayloadSizes(minBlobSize uint64, maxBlobSize uint64) ([]uint64, error) {
	legalBlobSizes, err := FindLegalBlobSizes(minBlobSize, maxBlobSize)
	if err != nil {
		return nil, fmt.Errorf("failed to find legal blob sizes: %w", err)
	}

	sizes := make([]uint64, 0, len(legalBlobSizes))

	for _, blobSize := range legalBlobSizes {
		minPayloadSize, err := BlobSizeToMinPayloadSize(blobSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get minimum payload size for blob size %d: %w", blobSize, err)
		}
		sizes = append(sizes, minPayloadSize)
	}

	return sizes, nil
}
