package testutils

import (
	"fmt"
	"math/rand"
	"time"
)

// TestRandom provides all the functionality of math/rand.Rand, plus additional randomness functionality useful for testing
// This struct wraps an instance of rand.Rand, and directly delegates all standard random methods to this internal instance
type TestRandom struct {
	internalRand *rand.Rand
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
		internalRand: rand.New(rand.NewSource(seed)),
	}
}

// GetRand returns the underlying random instance
func (r *TestRandom) GetRand() *rand.Rand {
	return r.internalRand
}

// RandomBytes generates a random byte slice of a given length.
func (r *TestRandom) RandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := r.internalRand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

// RandomTime generates a random time.
func (r *TestRandom) RandomTime() time.Time {
	return time.Unix(int64(r.internalRand.Int31()), 0)
}

// RandomString generates a random string out of printable ASCII characters.
func (r *TestRandom) RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.internalRand.Intn(len(charset))]
	}
	return string(b)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (r *TestRandom) Int63() int64 { return r.internalRand.Int63() }

// Uint32 returns a pseudo-random 32-bit value as a uint32.
func (r *TestRandom) Uint32() uint32 { return r.internalRand.Uint32() }

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func (r *TestRandom) Uint64() uint64 { return r.internalRand.Uint64() }

// Int31 returns a non-negative pseudo-random 31-bit integer as an int32.
func (r *TestRandom) Int31() int32 { return r.internalRand.Int31() }

// Int returns a non-negative pseudo-random int.
func (r *TestRandom) Int() int { return r.internalRand.Int() }

// Int63n returns, as an int64, a non-negative pseudo-random number in the half-open interval [0,n).
// It panics if n <= 0.
func (r *TestRandom) Int63n(n int64) int64 { return r.internalRand.Int63n(n) }

// Int31n returns, as an int32, a non-negative pseudo-random number in the half-open interval [0,n).
// It panics if n <= 0.
func (r *TestRandom) Int31n(n int32) int32 { return r.internalRand.Int31n(n) }

// Intn returns, as an int, a non-negative pseudo-random number in the half-open interval [0,n).
// It panics if n <= 0.
func (r *TestRandom) Intn(n int) int { return r.internalRand.Intn(n) }

// Float64 returns, as a float64, a pseudo-random number in the half-open interval [0.0,1.0).
func (r *TestRandom) Float64() float64 { return r.internalRand.Float64() }

// Float32 returns, as a float32, a pseudo-random number in the half-open interval [0.0,1.0).
func (r *TestRandom) Float32() float32 { return r.internalRand.Float32() }

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers
// in the half-open interval [0,n).
func (r *TestRandom) Perm(n int) []int { return r.internalRand.Perm(n) }

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func (r *TestRandom) Shuffle(n int, swap func(i, j int)) { r.internalRand.Shuffle(n, swap) }

// Read generates len(p) random bytes and writes them into p. It
// always returns len(p) and a nil error.
// Read should not be called concurrently with any other Rand method.
func (r *TestRandom) Read(p []byte) (n int, err error) { return r.internalRand.Read(p) }
