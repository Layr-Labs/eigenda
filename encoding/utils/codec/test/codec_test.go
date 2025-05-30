package test

import (
	"crypto/rand"
	"fmt"
	"strings"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
)

// This test is defined in its own package to avoid import cycles between the codec package and the coretypes package.
// Unit tests in this file call into both to ensure that codec packages calculations agree with the results of the
// actual operations in the coretypes package.

const minBlobSize = uint32(128 * units.KiB)
const maxBlobSize = uint32(16 * units.MiB)

var defaultPayloadForm = clients.GetDefaultPayloadClientConfig().PayloadPolynomialForm

// Derive the real size of a blob for a given payload by creating a payload and converting it to a blob.
func deriveRealBlobSize(t *testing.T, payloadSize uint32) uint32 {

	rawBytes := make([]byte, payloadSize)
	payload := coretypes.NewPayload(rawBytes)
	blob, err := payload.ToBlob(defaultPayloadForm)
	require.NoError(t, err)

	// We should get the same answer when we use the equation to calculate the blob size.
	calculatedBlobSize := codec.PayloadSizeToBlobSize(payloadSize)
	require.Equal(t, blob.BlobLengthBytes(), calculatedBlobSize)

	return blob.BlobLengthBytes()
}

// This function generates a table containing optimum blob sizes. It is intended to be run manually.
func TestGenerateOptimumSizeTable(t *testing.T) {

	// Comment this to generate an optimum size table.
	t.Skip() // Do not merge with this test enabled

	blobSizes, err := codec.FindLegalBlobSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	maxPayloadSizes, err := codec.FindMaxPayloadSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	columns := []string{
		"Maximum Payload Size",
		"Blob Size              ",
	}

	sb := strings.Builder{}

	// Write header
	for _, col := range columns {
		sb.WriteString(fmt.Sprintf("| %s ", col))
	}
	sb.WriteString("|\n")

	// Write separator
	for _, col := range columns {
		sb.WriteString("|:")
		sb.WriteString(strings.Repeat("-", len(col)+1))
	}
	sb.WriteString("|\n")

	for i := 0; i < len(blobSizes); i++ {
		maxSize := maxPayloadSizes[i]
		blobSize := blobSizes[i]

		niceUnit := "KiB"
		niceQuantity := int(float64(blobSize) / float64(units.KiB))
		if niceQuantity >= 1024 {
			niceUnit = "MiB"
			niceQuantity = int(float64(blobSize) / float64(units.MiB))
		}

		str := fmt.Sprintf("%d bytes", maxSize)
		str = fmt.Sprintf("| %-*s ", len(columns[0]), str) // Pad to column width
		sb.WriteString(str)

		str = fmt.Sprintf("%d bytes (%d %s)", blobSize, niceQuantity, niceUnit)
		str = fmt.Sprintf("| %-*s ", len(columns[1]), str) // Pad to column width
		sb.WriteString(str)

		sb.WriteString("|\n")
	}

	fmt.Print(sb.String())
}

func TestMinPayloadSizes(t *testing.T) {
	legalBlobSizes, err := codec.FindLegalBlobSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	minPayloadSizes, err := codec.FindMinPayloadSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	require.Equal(t, len(legalBlobSizes), len(minPayloadSizes))

	for i := 0; i < len(legalBlobSizes); i++ {
		blobSize := legalBlobSizes[i]
		minPayloadSize := minPayloadSizes[i]

		realBlobSize := deriveRealBlobSize(t, minPayloadSize)
		require.Equal(t, blobSize, realBlobSize)

		// Subtracting 1 byte from the minimum payload size should result in a blob that is the next tier smaller.
		if i > 0 {
			previousTier := legalBlobSizes[i-1]
			realBlobSize = deriveRealBlobSize(t, minPayloadSize-1)
			require.Equal(t, previousTier, realBlobSize)
		}
	}
}

func TestMaxPayloadSizes(t *testing.T) {
	legalBlobSizes, err := codec.FindLegalBlobSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	maxPayloadSizes, err := codec.FindMaxPayloadSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	require.Equal(t, len(legalBlobSizes), len(maxPayloadSizes))

	for i := 0; i < len(legalBlobSizes); i++ {
		blobSize := legalBlobSizes[i]
		maxPayloadSize := maxPayloadSizes[i]

		realBlobSize := deriveRealBlobSize(t, maxPayloadSize)
		require.Equal(t, blobSize, realBlobSize)

		// Adding 1 byte to the maximum payload size should result in a blob that is the next tier larger.
		if i < len(legalBlobSizes)-1 {
			nextTier := legalBlobSizes[i+1]
			realBlobSize = deriveRealBlobSize(t, maxPayloadSize+1)
			require.Equal(t, nextTier, realBlobSize)
		}
	}
}

func TestMinAgreesWithMax(t *testing.T) {
	legalBlobSizes, err := codec.FindLegalBlobSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	minPayloadSizes, err := codec.FindMinPayloadSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	maxPayloadSizes, err := codec.FindMaxPayloadSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	// Each minimum payload size should be exactly one larger than the maximum payload size of the previous tier.
	for i := 0; i < len(legalBlobSizes); i++ {
		if i > 0 {
			minPayloadSize := minPayloadSizes[i]
			maxPayloadSize := maxPayloadSizes[i-1]

			require.Equal(t, minPayloadSize, maxPayloadSize+1)
		}
	}
}

func TestSimplePaddingCodec(t *testing.T) {
	gettysburgAddressBytes := []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")

	paddedData := codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes)

	restored := codec.RemoveEmptyByteFromPaddedBytes(paddedData)

	require.Equal(t, gettysburgAddressBytes, restored[:len(gettysburgAddressBytes)])
}

func TestSimplePadding_IsValid(t *testing.T) {
	gettysburgAddressBytes := []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")

	paddedData := codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes)

	_, err := rs.ToFrArray(paddedData)
	require.Nil(t, err)
}

func TestSimplePaddingCodec_Fuzz(t *testing.T) {
	numFuzz := 100

	dataSizeList := make([]int, 0)
	for i := 32; i < 3000; i = i + 10 {
		dataSizeList = append(dataSizeList, i)
	}

	for i := 0; i < numFuzz; i++ {
		for j := 0; j < len(dataSizeList); j++ {
			data := make([]byte, dataSizeList[j])
			_, err := rand.Read(data)
			require.Nil(t, err)
			paddedData := codec.ConvertByPaddingEmptyByte(data)
			_, err = rs.ToFrArray(paddedData)
			require.Nil(t, err)
			restored := codec.RemoveEmptyByteFromPaddedBytes(paddedData)
			require.Equal(t, data, restored)
		}
	}
}

// TestGetPaddedDataLength tests that GetPaddedDataLength behaves relative to hardcoded expected results
func TestGetPaddedDataLengthAgainstKnowns(t *testing.T) {
	startLengths := []uint32{0, 30, 31, 32, 33, 68}
	expectedResults := []uint32{0, 32, 32, 64, 64, 96}

	for i := range startLengths {
		require.Equal(t, codec.GetPaddedDataLength(startLengths[i]), expectedResults[i])
	}
}

// TestGetUnpaddedDataLengthAgainstKnowns tests that GetPaddedDataLength behaves relative to hardcoded expected results
func TestGetUnpaddedDataLengthAgainstKnowns(t *testing.T) {
	startLengths := []uint32{0, 32, 64, 128}
	expectedResults := []uint32{0, 31, 62, 124}

	for i := range startLengths {
		unpaddedDataLength, err := codec.GetUnpaddedDataLength(startLengths[i])
		require.Nil(t, err)

		require.Equal(t, expectedResults[i], unpaddedDataLength)
	}

	unpaddedDataLength, err := codec.GetUnpaddedDataLength(129)
	require.Error(t, err)
	require.Equal(t, uint32(0), unpaddedDataLength)
}

// TestPadUnpad makes sure that padding and unpadding doesn't corrupt underlying data
func TestPadUnpad(t *testing.T) {
	testRandom := random.NewTestRandom()
	testIterations := 1000

	for i := 0; i < testIterations; i++ {
		originalBytes := testRandom.Bytes(testRandom.Intn(1024))

		paddedBytes := codec.PadPayload(originalBytes)
		require.Equal(t, len(paddedBytes)%32, 0)

		unpaddedBytes, err := codec.RemoveInternalPadding(paddedBytes)
		require.Nil(t, err)

		expectedUnpaddedLength, err := codec.GetUnpaddedDataLength(uint32(len(paddedBytes)))
		require.Nil(t, err)
		require.Equal(t, expectedUnpaddedLength, uint32(len(unpaddedBytes)))

		// unpadded payload may have up to 31 extra trailing zeros, since RemoveInternalPadding doesn't consider these
		require.Greater(t, len(originalBytes), len(unpaddedBytes)-32)
		require.LessOrEqual(t, len(originalBytes), len(unpaddedBytes))

		require.Equal(t, originalBytes, unpaddedBytes[:len(originalBytes)])
	}
}
