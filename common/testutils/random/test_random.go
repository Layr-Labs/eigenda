package random

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"testing"
	"time"
)

// TestRandom provides all the functionality of math/rand.Rand, plus additional randomness functionality useful for testing
type TestRandom struct {
	// The source of randomness
	*rand.Rand

	// The testing object
	t *testing.T

	// The seed used to initialize the random number generator
	seed int64
}

// NewTestRandom creates a new instance of TestRandom
// This method may either be seeded, or not seeded. If no seed is provided, then current unix nano time is used.
func NewTestRandom(t *testing.T, fixedSeed ...int64) *TestRandom {
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
		Rand: rand.New(rand.NewSource(seed)),
		t:    t,
		seed: seed,
	}
}

// Reset resets the random number generator to the state it was in when it was first created.
// This method is not thread safe with respect to other methods in this struct.
func (r *TestRandom) Reset() {
	r.Seed(r.seed)
}

// Bytes generates a random byte slice of a given length.
func (r *TestRandom) Bytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := r.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

// Time generates a random time.
func (r *TestRandom) Time() time.Time {
	return time.Unix(r.Int63(), r.Int63())
}

// String generates a random string out of printable ASCII characters.
func (r *TestRandom) String(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

var _ io.Reader = &randIOReader{}

// randIOReader is an io.Reader that reads from a random number generator.
type randIOReader struct {
	rand *TestRandom
}

func (i *randIOReader) Read(p []byte) (n int, err error) {
	return i.rand.Read(p)
}

// ECDSA generates a random ECDSA key. FOR TESTING PURPOSES ONLY. DO NOT USE THESE KEYS FOR SECURITY PURPOSES.
func (r *TestRandom) ECDSA() (*ecdsa.PublicKey, *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(crypto.S256(), &randIOReader{r})
	require.NoError(r.t, err)
	return &key.PublicKey, key
}
