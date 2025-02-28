package prover_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeFrame_AreInverses(t *testing.T) {
	group, err := prover.NewProver(kzgConfig, nil)
	require.NoError(t, err)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(gettysburgAddressBytes)))

	p, err := group.GetKzgEncoder(params)

	require.Nil(t, err)
	require.NotNil(t, p)

	// Convert to inputFr
	inputFr, err := rs.ToFrArray(gettysburgAddressBytes)
	require.Nil(t, err)

	frames, _, err := p.GetFrames(inputFr)
	require.Nil(t, err)
	require.NotNil(t, frames, err)

	b, err := frames[0].Encode()
	require.Nil(t, err)
	require.NotNil(t, b)

	frame, err := encoding.Decode(b)
	require.Nil(t, err)
	require.NotNil(t, frame)

	assert.Equal(t, frame, frames[0])
}
