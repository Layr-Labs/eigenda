package prover_test

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

func TestProveAllCosetThreads(t *testing.T) {
	harness := getTestHarness()

	group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	c, err := committer.NewFromConfig(*harness.committerConfig)
	require.NoError(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))

	commitments, err := c.GetCommitmentsForPaddedLength(harness.paddedGettysburgAddressBytes)
	require.Nil(t, err)
	frames, err := group.GetFrames(harness.paddedGettysburgAddressBytes, params)
	require.Nil(t, err)

	verifier, err := verifier.NewVerifier(harness.verifierV2KzgConfig)
	require.Nil(t, err)

	indices := []encoding.ChunkNumber{}
	for i := range len(frames) {
		indices = append(indices, encoding.ChunkNumber(i))
	}
	err = verifier.VerifyFrames(frames, indices, commitments, params)
	require.Nil(t, err)
}

func TestEncodeDecodeFrame_AreInverses(t *testing.T) {
	harness := getTestHarness()

	group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))

	blobLength := uint64(encoding.GetBlobLengthPowerOf2(uint32(len(harness.paddedGettysburgAddressBytes))))
	p, err := group.GetKzgProver(params, blobLength)

	require.Nil(t, err)
	require.NotNil(t, p)

	// Convert to inputFr
	inputFr, err := rs.ToFrArray(harness.paddedGettysburgAddressBytes)
	require.Nil(t, err)

	frames, _, err := p.GetFrames(inputFr)
	require.Nil(t, err)
	require.NotNil(t, frames, err)

	b, err := frames[0].SerializeGob()
	require.Nil(t, err)
	require.NotNil(t, b)

	frame, err := new(encoding.Frame).DeserializeGob(b)
	require.Nil(t, err)
	require.NotNil(t, frame)

	assert.Equal(t, *frame, frames[0])
}
