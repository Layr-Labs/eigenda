package random

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

// Tests that random seeding produces random results, and that consistent seeding produces consistent results
func TestSetup(t *testing.T) {
	testRandom1 := NewTestRandom(t)
	x := testRandom1.Int()

	testRandom2 := NewTestRandom(t)
	y := testRandom2.Int()

	assert.NotEqual(t, x, y)

	seed := rand.Int63()
	testRandom3 := NewTestRandom(t, seed)
	a := testRandom3.Int()

	testRandom4 := NewTestRandom(t, seed)
	b := testRandom4.Int()

	assert.Equal(t, a, b)
}
