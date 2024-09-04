package chunkgroup

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestUint64ToBytes(t *testing.T) {
	bytes := make([]byte, 8)

	// Test the conversion of 0 to bytes.
	uint64ToBytes(0, bytes)
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, bytes)

	// Test the conversion of 1 to bytes.
	uint64ToBytes(1, bytes)
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 1}, bytes)

	// Test the conversion of 256 to bytes.
	uint64ToBytes(256, bytes)
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 1, 0}, bytes)

	// Test the conversion of 0xAABBCCDDEEFF9988 to bytes.
	uint64ToBytes(0xAABBCCDDEEFF9988, bytes)
	assert.Equal(t, []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x99, 0x88}, bytes)

	// Test the conversion of 0xFFFFFFFFFFFFFFFF to bytes.
	uint64ToBytes(0xFFFFFFFFFFFFFFFF, bytes)
	assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, bytes)

	bytes = make([]byte, 32)

	// Test the conversion of 0xAABBCCDDEEFF9988, 0x1122334455667788, 0x99AABBCCDDEEFF00, 0x0011223344556677 to bytes.
	uint64ToBytes(0xAABBCCDDEEFF9988, bytes[:8])
	uint64ToBytes(0x1122334455667788, bytes[8:16])
	uint64ToBytes(0x99AABBCCDDEEFF00, bytes[16:24])
	uint64ToBytes(0x0011223344556677, bytes[24:32])
	assert.Equal(t, []byte{
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x99, 0x88,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00,
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
	}, bytes)
}

func TestBytesToUint64(t *testing.T) {
	// Test the conversion of 0 from bytes.
	assert.Equal(t, uint64(0), bytesToUint64([]byte{0, 0, 0, 0, 0, 0, 0, 0}))

	// Test the conversion of 1 from bytes.
	assert.Equal(t, uint64(1), bytesToUint64([]byte{0, 0, 0, 0, 0, 0, 0, 1}))

	// Test the conversion of 256 from bytes.
	assert.Equal(t, uint64(256), bytesToUint64([]byte{0, 0, 0, 0, 0, 0, 1, 0}))

	// Test the conversion of 0xAABBCCDDEEFF9988 from bytes.
	assert.Equal(t, uint64(0xAABBCCDDEEFF9988), bytesToUint64([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x99, 0x88}))

	// Test the conversion of 0xFFFFFFFFFFFFFFFF from bytes.
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), bytesToUint64([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}))
}

func TestByteRoundTrip(t *testing.T) {
	tu.InitializeRandom()
	bytes := make([]byte, 8)
	for i := 0; i < 1000; i++ {
		value := uint64(rand.Int63())
		uint64ToBytes(value, bytes)
		assert.Equal(t, value, bytesToUint64(bytes))
	}

	bytes = make([]byte, 32)
	for i := 0; i < 1000; i++ {
		value1 := uint64(rand.Int63())
		value2 := uint64(rand.Int63())
		value3 := uint64(rand.Int63())
		value4 := uint64(rand.Int63())
		uint64ToBytes(value1, bytes[:8])
		uint64ToBytes(value2, bytes[8:16])
		uint64ToBytes(value3, bytes[16:24])
		uint64ToBytes(value4, bytes[24:32])
		assert.Equal(t, value1, bytesToUint64(bytes[:8]))
		assert.Equal(t, value2, bytesToUint64(bytes[8:16]))
		assert.Equal(t, value3, bytesToUint64(bytes[16:24]))
		assert.Equal(t, value4, bytesToUint64(bytes[24:32]))
	}
}

func TestComputeShuffleOffset(t *testing.T) {
	tu.InitializeRandom()

	shufflePeriod := time.Second * time.Duration(rand.Intn(1000)+1)

	uniqueOffsets := make(map[time.Duration]bool)

	// Ensure that shuffle offset is always between 0 and shufflePeriod.
	iterations := 1000
	for i := 0; i < iterations; i++ {
		seed := rand.Uint64()
		index := uint64(rand.Intn(64))
		offset := ComputeShuffleOffset(seed, index, shufflePeriod)
		assert.True(t, offset >= 0)
		assert.True(t, offset < shufflePeriod)
		uniqueOffsets[offset] = true

		// Computing the offset a second time should yield the same result.
		offset2 := ComputeShuffleOffset(seed, index, shufflePeriod)
		assert.Equal(t, offset, offset2)

		// Changing the seed should yield a different result.
		index2 := (index + 1) % 64
		offset3 := ComputeShuffleOffset(seed, index2, shufflePeriod)
		assert.NotEqual(t, offset, offset3)
	}

	// Although it's possible we end up with two offsets that are the same, the probability
	// is not high. Although it's unlikely for there to be even one collision, to avoid a flaky test
	//only assert that we have at least 99% unique offsets.
	assert.True(t, len(uniqueOffsets) >= int(float64(iterations)*0.99))
}

func TestComputeShuffleOffsetNegativeShufflePeriod(t *testing.T) {
	assert.Panics(t, func() {
		_ = ComputeShuffleOffset(0, 0, -time.Second)
	})
}

func TestComputeShuffleOffsetTooLargeAssignmentIndex(t *testing.T) {
	assert.Panics(t, func() {
		_ = ComputeShuffleOffset(0, 64, time.Second)
	})
	assert.Panics(t, func() {
		_ = ComputeShuffleOffset(0, 12345, time.Second)
	})
}

func TestComputeShuffleEpochHandCraftedScenario(t *testing.T) {
	shufflePeriod := time.Second * 100
	shuffleOffset := time.Second * 50

	unixEpoch := time.Unix(0, 0)

	// Requesting the epoch before genesis should return an error.
	assert.Panics(t, func() {
		_ = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(-time.Second))
	})

	// Requesting the epoch at genesis should return 0.
	epoch := ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch)
	assert.Equal(t, uint64(0), epoch)

	// Requesting the epoch one nanosecond after genesis should return 0.
	epoch = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(1))
	assert.Equal(t, uint64(0), epoch)

	// Requesting the epoch in-between genesis and the first shuffle should return 0.
	epoch = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(shuffleOffset/2))
	assert.Equal(t, uint64(0), epoch)

	// Requesting the epoch one nanosecond before the first shuffle should return 0.
	epoch = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(shuffleOffset-1))
	assert.Equal(t, uint64(0), epoch)

	// Requesting the epoch at the exact first shuffle time should return 1.
	epoch = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(shuffleOffset))
	assert.Equal(t, uint64(1), epoch)

	// Requesting the epoch one nanosecond after the first shuffle should return 1.
	epoch = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(shuffleOffset+1))
	assert.Equal(t, uint64(1), epoch)

	// Requesting the epoch at the mid-point between the first and second shuffle should return 1.
	epoch = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(shuffleOffset+shufflePeriod/2))
	assert.Equal(t, uint64(1), epoch)

	// Requesting the epoch one nanosecond before the second shuffle should return 1.
	epoch = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(shuffleOffset+shufflePeriod-1))
	assert.Equal(t, uint64(1), epoch)

	// Requesting the epoch at the exact second shuffle time should return 2.
	epoch = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(shuffleOffset+shufflePeriod))
	assert.Equal(t, uint64(2), epoch)

	// Moving forward 1000 shuffle periods should put us in epoch 1002.
	epoch = ComputeShuffleEpoch(shufflePeriod, shuffleOffset, unixEpoch.Add(shuffleOffset+shufflePeriod*1001))
	assert.Equal(t, uint64(1002), epoch)
}

// Simulate a scenario where we move down a timeline. Verify that the epoch increases as expected.
func TestComputeShuffleEpochTimeWalk(t *testing.T) {
	tu.InitializeRandom()

	unixEpoch := time.Unix(0, 0)

	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	seed := rand.Uint64()
	index := uint64(rand.Intn(64))
	offset := ComputeShuffleOffset(seed, index, shufflePeriod)

	// If we move forward in time 1000 shuffle periods, we should see the epoch increase 1000 times.
	endTime := unixEpoch.Add(shufflePeriod * 1000)
	epochCount := 0

	previousEpoch := ComputeShuffleEpoch(shufflePeriod, offset, unixEpoch)

	for now := unixEpoch; !now.After(endTime); now = now.Add(time.Second) {

		currentEpoch := ComputeShuffleEpoch(shufflePeriod, offset, now)
		if currentEpoch != previousEpoch {
			epochCount++
			previousEpoch = currentEpoch
		}
	}

	assert.Equal(t, 1000, epochCount)
}

func TestComputeShuffleEpochInvalidShufflePeriod(t *testing.T) {
	assert.Panics(t, func() {
		_ = ComputeShuffleEpoch(-time.Second, time.Second, time.Now())
	})

	assert.Panics(t, func() {
		_ = ComputeShuffleEpoch(0, time.Second, time.Now())
	})
}

func TestComputeShuffleEpochInvalidShuffleOffset(t *testing.T) {
	assert.Panics(t, func() {
		_ = ComputeShuffleEpoch(time.Second, -time.Second, time.Now())
	})
}

func TestComputeChunkGroup(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint64(rand.Intn(1_000) + 10_000)
	seed := rand.Uint64()
	index := uint64(rand.Intn(64))
	firstEpoch := uint64(rand.Intn(1000))

	uniqueChunkGroups := make(map[uint64]bool)

	// Ensure that the chunk group is always between 0 and chunkGroupCount.
	iterations := 1000
	for i := 0; i < iterations; i++ {
		epoch := firstEpoch + uint64(i)
		chunkGroup := ComputeChunkGroup(seed, index, epoch, chunkGroupCount)
		assert.True(t, chunkGroup < chunkGroupCount)
		uniqueChunkGroups[chunkGroup] = true
	}

	// We will have some collisions, but we should have at least 50% unique chunk groups with high probability.
	assert.True(t, len(uniqueChunkGroups) >= int(float64(iterations)*0.5))
}

func TestComputeEndOfShuffleEpochErrors(t *testing.T) {
	assert.Panics(t, func() {
		_ = ComputeEndOfShuffleEpoch(-time.Second, time.Second, 0)
	})

	assert.Panics(t, func() {
		_ = ComputeEndOfShuffleEpoch(time.Second, -time.Second, 0)
	})
}

func TestComputeEndOfShuffleEpochHandCraftedScenario(t *testing.T) {
	shufflePeriod := time.Second * 100
	shuffleOffset := time.Second * 50

	assert.Equal(t, time.Unix(50, 0), ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 0))
	assert.Equal(t, time.Unix(150, 0), ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 1))
	assert.Equal(t, time.Unix(250, 0), ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 2))
	assert.Equal(t, time.Unix(350, 0), ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 3))
	assert.Equal(t, time.Unix(450, 0), ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 4))
}

func TestComputeEndOfShuffleEpoch(t *testing.T) {
	tu.InitializeRandom()

	unixEpoch := time.Unix(0, 0)
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	shuffleOffset := time.Second * time.Duration(rand.Intn(10)+1)

	assert.Equal(t, unixEpoch.Add(shuffleOffset),
		ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 0))
	assert.Equal(t, unixEpoch.Add(shuffleOffset).Add(shufflePeriod),
		ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 1))
	assert.Equal(t, unixEpoch.Add(shuffleOffset).Add(shufflePeriod*2),
		ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 2))
	assert.Equal(t, unixEpoch.Add(shuffleOffset).Add(shufflePeriod*3),
		ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 3))
	assert.Equal(t, unixEpoch.Add(shuffleOffset).Add(shufflePeriod*4),
		ComputeEndOfShuffleEpoch(shufflePeriod, shuffleOffset, 4))
}

func TestComputeStartOfShuffleEpochHandCraftedScenario(t *testing.T) {

	shufflePeriod := time.Second * 100
	shuffleOffset := time.Second * 50

	assert.Equal(t, time.Unix(-50, 0), ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 0))
	assert.Equal(t, time.Unix(50, 0), ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 1))
	assert.Equal(t, time.Unix(150, 0), ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 2))
	assert.Equal(t, time.Unix(250, 0), ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 3))
	assert.Equal(t, time.Unix(350, 0), ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 4))
}

func TestComputeStartOfShuffleEpoch(t *testing.T) {
	tu.InitializeRandom()

	unixEpoch := time.Unix(0, 0)
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	shuffleOffset := time.Second * time.Duration(rand.Intn(10)+1)

	assert.Equal(t, unixEpoch.Add(shuffleOffset-shufflePeriod),
		ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 0))
	assert.Equal(t, unixEpoch.Add(shuffleOffset),
		ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 1))
	assert.Equal(t, unixEpoch.Add(shuffleOffset).Add(shufflePeriod),
		ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 2))
	assert.Equal(t, unixEpoch.Add(shuffleOffset).Add(shufflePeriod*2),
		ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 3))
	assert.Equal(t, unixEpoch.Add(shuffleOffset).Add(shufflePeriod*3),
		ComputeStartOfShuffleEpoch(shufflePeriod, shuffleOffset, 4))
}
