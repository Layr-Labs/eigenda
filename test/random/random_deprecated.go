package random

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

// InitializeRandom initializes the random number generator. If no arguments are provided, then the seed is randomly
// generated. If a single argument is provided, then the seed is fixed to that value.
// Deprecated: use TestRandom instead
func InitializeRandom(fixedSeed ...uint64) {

	var seed uint64
	if len(fixedSeed) == 0 {
		rand.Seed(uint64(time.Now().UnixNano()))
		seed = rand.Uint64()
	} else if len(fixedSeed) == 1 {
		seed = fixedSeed[0]
	} else {
		panic("too many arguments, expected exactly one seed")
	}

	fmt.Printf("Random seed: %d\n", seed)
	rand.Seed(seed)
}

// RandomBytes generates a random byte slice of a given length.
// Deprecated: use TestRandom.Bytes instead
func RandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

// RandomTime generates a random time.
// Deprecated: use TestRandom.Time instead
func RandomTime() time.Time {
	return time.Unix(int64(rand.Int31()), 0)
}

// RandomString generates a random string out of printable ASCII characters.
// Deprecated: use TestRandom.String instead
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
