package common_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/require"
)

func TestReadOnlyMap(t *testing.T) {
	data := map[uint8]string{
		1: "one",
		2: "two",
		3: "three",
	}
	m := common.NewReadOnlyMap(data)
	res, ok := m.Get(1)
	require.True(t, ok)
	require.Equal(t, "one", res)
	res, ok = m.Get(2)
	require.True(t, ok)
	require.Equal(t, "two", res)
	res, ok = m.Get(3)
	require.True(t, ok)
	require.Equal(t, "three", res)
	res, ok = m.Get(4)
	require.False(t, ok)
	require.Equal(t, "", res)
	require.Equal(t, 3, m.Len())
	require.ElementsMatch(t, []uint8{1, 2, 3}, m.Keys())
}
