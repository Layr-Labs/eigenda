package segment

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddress(t *testing.T) {
	rand := random.NewTestRandom(t)

	index := rand.Uint32()
	offset := rand.Uint32()
	address := NewAddress(index, offset)

	require.Equal(t, index, address.Index())
	require.Equal(t, offset, address.Offset())
}
