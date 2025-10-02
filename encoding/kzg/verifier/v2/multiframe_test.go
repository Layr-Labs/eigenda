package verifier_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniversalVerify(t *testing.T) {
	harness := getTestHarness()

	group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.Nil(t, err)

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.Nil(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.Nil(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))
	enc, err := group.GetKzgEncoder(params)
	require.Nil(t, err)

	numBlob := 5
	samples := make([]verifier.Sample, 0)
	for z := 0; z < numBlob; z++ {
		inputFr, err := rs.ToFrArray(harness.paddedGettysburgAddressBytes)
		require.Nil(t, err)

		commit, _, _, err := committer.GetCommitments(inputFr)
		require.Nil(t, err)
		frames, fIndices, err := enc.GetFrames(inputFr)
		require.Nil(t, err)

		// create samples
		for i := 0; i < len(frames); i++ {
			f := frames[i]
			j := fIndices[i]

			q, err := rs.GetLeadingCosetIndex(uint64(i), harness.numSys+harness.numPar)
			require.Nil(t, err)

			assert.Equal(t, j, q, "leading coset inconsistency")

			sample := verifier.Sample{
				Commitment: *commit,
				Proof:      f.Proof,
				RowIndex:   z,
				Coeffs:     f.Coeffs,
				X:          uint(q),
			}
			samples = append(samples, sample)
		}
	}

	assert.True(t, v.UniversalVerify(params, samples, numBlob) == nil, "universal batch verification failed\n")
}

func TestUniversalVerifyWithPowerOf2G2(t *testing.T) {
	harness := getTestHarness()
	group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.Nil(t, err)

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.Nil(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	assert.NoError(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))
	enc, err := group.GetKzgEncoder(params)
	assert.NoError(t, err)

	numBlob := 5
	samples := make([]verifier.Sample, 0)
	for z := 0; z < numBlob; z++ {
		inputFr, err := rs.ToFrArray(harness.paddedGettysburgAddressBytes)
		require.Nil(t, err)

		commit, _, _, err := committer.GetCommitments(inputFr)
		require.Nil(t, err)
		frames, fIndices, err := enc.GetFrames(inputFr)
		require.Nil(t, err)

		// create samples
		for i := 0; i < len(frames); i++ {
			f := frames[i]
			j := fIndices[i]

			q, err := rs.GetLeadingCosetIndex(uint64(i), harness.numSys+harness.numPar)
			require.Nil(t, err)

			assert.Equal(t, j, q, "leading coset inconsistency")

			sample := verifier.Sample{
				Commitment: *commit,
				Proof:      f.Proof,
				RowIndex:   z,
				Coeffs:     f.Coeffs,
				X:          uint(q),
			}
			samples = append(samples, sample)
		}
	}

	assert.True(t, v.UniversalVerify(params, samples, numBlob) == nil, "universal batch verification failed\n")
}
