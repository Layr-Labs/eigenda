package sizing

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
)

const minBlobSize = uint64(128 * units.KiB)
const maxBlobSize = uint64(16 * units.MiB)

var defaultPayloadForm = clients.GetDefaultPayloadClientConfig().PayloadPolynomialForm

// Derive the real size of a blob for a given payload by creating a payload and converting it to a blob.
func deriveRealBlobSize(t *testing.T, payloadSize uint64) uint64 {

	rawBytes := make([]byte, payloadSize)
	payload := coretypes.NewPayload(rawBytes)
	blob, err := payload.ToBlob(defaultPayloadForm)
	require.NoError(t, err)

	return uint64(blob.BlobLengthBytes())
}

// This function generates a table containing optimum blob sizes. It is intended to be run manually.
func TestGenerateOptimumSizeTable(t *testing.T) {

	// Uncomment this to generate an optimum size table.
	t.Skip() // Do not merge with this line uncommented.

	minBlobSize := uint64(128 * units.KiB)
	maxBlobSize := uint64(16 * units.MiB)

	blobSizes, err := FindLegalBlobSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	minPayloadSizes, err := FindMinPayloadSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	maxPayloadSizes, err := FindMaxPayloadSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	columns := []string{
		"Minimum Payload Size",
		"Maximum Payload Size",
		"Blob Size      ",
		"Blob Size",
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
		minSize := 0
		if i > 0 {
			minSize = int(minPayloadSizes[i])
		}
		maxSize := maxPayloadSizes[i]
		blobSize := blobSizes[i]

		niceUnit := "kb"
		niceQuantity := int(float64(blobSize) / float64(units.KiB))
		if niceQuantity >= 1024 {
			niceUnit = "mb"
			niceQuantity = int(float64(blobSize) / float64(units.MiB))
		}

		str := fmt.Sprintf("%d bytes", minSize)
		str = fmt.Sprintf("| %-*s ", len(columns[0]), str) // Pad to column width
		sb.WriteString(str)

		str = fmt.Sprintf("%d bytes", maxSize)
		str = fmt.Sprintf("| %-*s ", len(columns[1]), str) // Pad to column width
		sb.WriteString(str)

		str = fmt.Sprintf("%d bytes", blobSize)
		str = fmt.Sprintf("| %-*s ", len(columns[2]), str) // Pad to column width
		sb.WriteString(str)

		str = fmt.Sprintf("%d %s", niceQuantity, niceUnit)
		str = fmt.Sprintf("| %-*s ", len(columns[3]), str) // Pad to column width
		sb.WriteString(str)

		sb.WriteString("|\n")
	}

	fmt.Printf(sb.String())
}

func TestMinPayloadSizes(t *testing.T) {
	legalBlobSizes, err := FindLegalBlobSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	minPayloadSizes, err := FindMinPayloadSizes(minBlobSize, maxBlobSize)
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
	legalBlobSizes, err := FindLegalBlobSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	maxPayloadSizes, err := FindMaxPayloadSizes(minBlobSize, maxBlobSize)
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
	legalBlobSizes, err := FindLegalBlobSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	minPayloadSizes, err := FindMinPayloadSizes(minBlobSize, maxBlobSize)
	require.NoError(t, err)

	maxPayloadSizes, err := FindMaxPayloadSizes(minBlobSize, maxBlobSize)
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
