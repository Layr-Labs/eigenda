package benchmark

// DataGenerator is responsible for generating key-value pairs to be inserted into the database, for the sake of
// benchmarking.
type DataGenerator struct {
	// A pool of randomness. Used to generate values.
	dataPool []byte

	randPool nil

	// The seed that determines the key/value pairs generated.
	seed int64
}

// NewDataGenerator builds a data generator instance.
func NewDataGenerator(seed int64, poolSize uint64) *DataGenerator {
	return &DataGenerator{}
}

// GenerateKVPair generates a new key value pair. The key is always 32 bytes long and generated randomly.
// The value is sourced from reusable entropy (since it's expensive to generate huge quantities of random data).
// The resulting value is deterministic given the same index + length.
func (g *DataGenerator) GenerateKVPair(index uint64, valueLength uint64) (key []byte, value []byte) {
	return nil, nil
}
