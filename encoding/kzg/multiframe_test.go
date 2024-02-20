package kzgEncoder_test

import (
	"testing"

	kzgRs "github.com/Layr-Labs/eigenda/encoding/kzg"
	rs "github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniversalVerify(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgRs.NewKzgEncoderGroup(kzgConfig, true)
	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, err := group.NewKzgEncoder(params)
	require.Nil(t, err)

	numBlob := 5
	samples := make([]kzgRs.Sample, 0)
	for z := 0; z < numBlob; z++ {
		inputFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)

		commit, _, _, frames, fIndices, err := enc.Encode(inputFr)
		require.Nil(t, err)

		// create samples
		for i := 0; i < len(frames); i++ {
			f := frames[i]
			j := fIndices[i]

			q, err := rs.GetLeadingCosetIndex(uint64(i), numSys+numPar)
			require.Nil(t, err)

			assert.Equal(t, j, q, "leading coset inconsistency")

			sample := kzgRs.Sample{
				Commitment: *commit,
				Proof:      f.Proof,
				RowIndex:   z,
				Coeffs:     f.Coeffs,
				X:          uint(q),
			}
			samples = append(samples, sample)
		}
	}

	assert.True(t, group.UniversalVerify(params, samples, numBlob) == nil, "universal batch verification failed\n")
}

func TestUniversalVerifyWithPowerOf2G2(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgRs.NewKzgEncoderGroup(kzgConfig, true)
	group.KzgConfig.G2Path = ""
	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, err := group.NewKzgEncoder(params)
	require.Nil(t, err)

	numBlob := 5
	samples := make([]kzgRs.Sample, 0)
	for z := 0; z < numBlob; z++ {
		inputFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)

		commit, _, _, frames, fIndices, err := enc.Encode(inputFr)
		require.Nil(t, err)

		// create samples
		for i := 0; i < len(frames); i++ {
			f := frames[i]
			j := fIndices[i]

			q, err := rs.GetLeadingCosetIndex(uint64(i), numSys+numPar)
			require.Nil(t, err)

			assert.Equal(t, j, q, "leading coset inconsistency")

			sample := kzgRs.Sample{
				Commitment: *commit,
				Proof:      f.Proof,
				RowIndex:   z,
				Coeffs:     f.Coeffs,
				X:          uint(q),
			}
			samples = append(samples, sample)
		}
	}

	assert.True(t, group.UniversalVerify(params, samples, numBlob) == nil, "universal batch verification failed\n")
}
