package random

import (
	"fmt"
	"math/rand"
	"time"
)

// TestRandom provides all the functionality of math/rand.Rand, plus additional randomness functionality useful for testing
type TestRandom struct {
	*rand.Rand
}

// NewTestRandom creates a new instance of TestRandom
// This method may either be seeded, or not seeded. If no seed is provided, then current unix nano time is used.
func NewTestRandom(fixedSeed ...int64) *TestRandom {
	var seed int64
	if len(fixedSeed) == 0 {
		seed = time.Now().UnixNano()
	} else if len(fixedSeed) == 1 {
		seed = fixedSeed[0]
	} else {
		panic("too many arguments, expected exactly one seed")
	}

	fmt.Printf("Random seed: %d\n", seed)
	return &TestRandom{
		rand.New(rand.NewSource(seed)),
	}
}

// GetRand returns the underlying random instance
func (r *TestRandom) GetRand() *rand.Rand {
	return r.Rand
}

// RandomBytes generates a random byte slice of a given length.
func (r *TestRandom) RandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := r.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

// RandomTime generates a random time.
func (r *TestRandom) RandomTime() time.Time {
	return time.Unix(r.Int63(), r.Int63())
}

// RandomString generates a random string out of printable ASCII characters.
func (r *TestRandom) RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
