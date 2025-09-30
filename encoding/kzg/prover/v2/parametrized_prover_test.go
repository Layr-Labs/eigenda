package prover_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProveAllCosetThreads(t *testing.T) {
	harness := getTestHarness()

	group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	// TODO(samlaf): committer should have its own builder for loading SRS
	// Or we should be loading SRS points completely separately and injecting
	// them into the prover/committer/verifier.
	c, err := committer.New(group.Srs.G1, group.Srs.G2, group.G2Trailing)
	require.NoError(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))

	commitments, err := c.GetCommitmentsForPaddedLength(harness.paddedGettysburgAddressBytes)
	require.Nil(t, err)
	frames, err := group.GetFrames(harness.paddedGettysburgAddressBytes, params)
	require.Nil(t, err)

	verifierGroup, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.Nil(t, err)
	verifier, err := verifierGroup.GetKzgVerifier(params)
	require.Nil(t, err)

	for i, frame := range frames {
		err = verifier.VerifyFrame(frame, uint64(i), (*bn254.G1Affine)(commitments.Commitment), params.NumChunks)
		require.Nil(t, err)
	}
}

func TestEncodeDecodeFrame_AreInverses(t *testing.T) {
	harness := getTestHarness()

	group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))

	p, err := group.GetKzgEncoder(params)

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
