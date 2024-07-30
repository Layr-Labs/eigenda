package testutils

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
)

func TestRandomSetup(t *testing.T) {
	InitializeRandom()
	x := rand.Int()

	InitializeRandom()
	y := rand.Int()

	assert.NotEqual(t, x, y)

	seed := uint64(rand.Int())
	InitializeRandomWithSeed(seed)
	a := rand.Int()

	InitializeRandomWithSeed(seed)
	b := rand.Int()

	assert.Equal(t, a, b)
}
