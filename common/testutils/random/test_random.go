package random

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"fmt"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

// charset is the set of characters that can be used to generate random strings
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
//
// The testing.T object is optional but highly recommended. If nil and an error occurs, this utility panics.
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

// VariableBytes generates a random byte slice of a length between min (inclusive) and max (exclusive).
func (r *TestRandom) VariableBytes(min int, max int) []byte {
	length := r.Intn(max-min) + min
	return r.Bytes(length)
}

// Time generates a random time.
func (r *TestRandom) Time() time.Time {
	return time.Unix(r.Int63(), r.Int63())
}

// String generates a random string out of printable ASCII characters.
func (r *TestRandom) String(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// VariableString generates a random string out of printable ASCII characters of a length between
// min (inclusive) and max (exclusive).
func (r *TestRandom) VariableString(min int, max int) string {
	length := r.Intn(max-min) + min
	return r.String(length)
}

// Uint32n generates a random uint32 less than n.
func (r *TestRandom) Uint32n(n uint32) uint32 {
	return r.Uint32() % n
}

// Uint64n generates a random uint64 less than n.
func (r *TestRandom) Uint64n(n uint64) uint64 {
	return r.Uint64() % n
}

// Gaussian generates a random float64 from a Gaussian distribution with the given mean and standard deviation.
func (r *TestRandom) Gaussian(mean float64, stddev float64) float64 {
	return r.NormFloat64()*stddev + mean
}

// BoundedGaussian generates a random float64 from a Gaussian distribution with the given mean and standard deviation,
// but bounded by the given min and max values. If a generated value exceeds the bounds, the bound is returned instead.
func (r *TestRandom) BoundedGaussian(mean float64, stddev float64, min float64, max float64) float64 {
	val := r.Gaussian(mean, stddev)
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

var _ io.Reader = &randIOReader{}

// randIOReader is an io.Reader that reads from a random number generator.
type randIOReader struct {
	rand *TestRandom
}

// Read reads random bytes into the provided buffer, returning the number of bytes read.
func (i *randIOReader) Read(p []byte) (n int, err error) {
	return i.rand.Read(p)
}

// IOReader creates an io.Reader that reads from a random number generator.
func (r *TestRandom) IOReader() io.Reader {
	return &randIOReader{r}
}

// ECDSA generates a random ECDSA key. Note that the returned keys are not deterministic due to limitations
// **intentionally** imposed by the Go standard libraries. (╯°□°)╯︵ ┻━┻
//
// NOT CRYPTOGRAPHICALLY SECURE!!! FOR TESTING PURPOSES ONLY. DO NOT USE THESE KEYS FOR SECURITY PURPOSES.
func (r *TestRandom) ECDSA() (*ecdsa.PublicKey, *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(crypto.S256(), crand.Reader)
	r.requireNoError(err)
	return &key.PublicKey, key
}

// BLS generates a random BLS key pair.
//
// NOT CRYPTOGRAPHICALLY SECURE!!! FOR TESTING PURPOSES ONLY. DO NOT USE THESE KEYS FOR SECURITY PURPOSES.
func (r *TestRandom) BLS() *core.KeyPair {
	//Max random value is order of the curve
	maxValue := new(big.Int)
	maxValue.SetString(fr.Modulus().String(), 10)

	//Generate cryptographically strong pseudo-random between 0 - max
	n, err := crand.Int(r.IOReader(), maxValue)
	r.requireNoError(err)

	sk := new(core.PrivateKey).SetBigInt(n)
	return core.MakeKeyPair(sk)
}

func (r *TestRandom) requireNoError(err error) {
	if err != nil && r.t == nil {
		panic(err)
	}
	r.requireNoError(err)
}
