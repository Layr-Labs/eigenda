package srs_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/resources/srs"
	"github.com/stretchr/testify/require"
)

func TestG2PowerOf2SRSContains28Points(t *testing.T) {
	require.Equal(t, 28, len(srs.G2PowerOf2SRS))
	t.Log(srs.G2PowerOf2SRS[0])
}
