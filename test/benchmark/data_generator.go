package benchmark

import (
	"math/rand"
	"sync"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
)

// DataGenerator is responsible for generating key-value pairs to be inserted into the database, for the sake of
// benchmarking.
type DataGenerator struct {
	// Pool of random number generators
	randPool *sync.Pool

	// A pool of randomness. Used to generate values.
	dataPool []byte

	// The seed that determines the key/value pairs generated.
	seed int64
}

// NewDataGenerator builds a data generator instance.
func NewDataGenerator(seed int64, poolSize uint64) *DataGenerator {

	randPool := &sync.Pool{
		New: func() interface{} {
			return random.NewTestRandomNoPrint()
		},
	}

	dataPool := make([]byte, poolSize)
	rng := randPool.Get().(*rand.Rand)
	rng.Read(dataPool)
	randPool.Put(randPool)

	return &DataGenerator{
		randPool: randPool,
		dataPool: dataPool,
	}
}

// GenerateKVPair generates a new key value pair. The key is always 32 bytes long and generated randomly.
// The value is sourced from reusable entropy (since it's expensive to generate huge quantities of random data).
// The resulting value is deterministic given the same index + length.
func (g *DataGenerator) GenerateKVPair(index uint64, valueLength uint64) (key []byte, value []byte) {
	rng := g.randPool.Get().(*rand.Rand)
	rng.Seed(g.seed + int64(index))

	key = make([]byte, 32)
	rng.Read(key)
	defer g.randPool.Put(rng)

	if valueLength > uint64(len(g.dataPool)) {
		// Special case: we don't have enough data in the pool to satisfy the request.
		// For the sake of completeness, just generate the data if this happens.
		// This shouldn't be encountered for sane configurations (i.e. with a pool size much larger than value sizes).
		value = make([]byte, valueLength)
		rng.Read(value)
	} else {
		startIndex := rng.Intn(len(g.dataPool) - int(valueLength))
		value = g.dataPool[startIndex : startIndex+int(valueLength)]
	}

	return key, value
}
