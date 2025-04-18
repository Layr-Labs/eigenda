package network_benchmark

import testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"

type reusableRandomness struct {
	data []byte
}

// newReusableRandomness provides random data from a pre-generated pool.
func newReusableRandomness(size int, seed int64) *reusableRandomness {
	rand := testrandom.NewTestRandomNoPrint(seed)
	data := rand.PrintableBytes(size)
	return &reusableRandomness{
		data: data,
	}
}

func (r *reusableRandomness) getData(size int64, seed int64) []byte {
	if size > int64(len(r.data)) {
		panic("Requested size exceeds available data size")
	}
	if seed < 0 {
		seed = -seed
	}

	maxStartingIndex := int64(len(r.data)) - size
	startingIndex := seed % maxStartingIndex
	return r.data[startingIndex : startingIndex+size]
}
