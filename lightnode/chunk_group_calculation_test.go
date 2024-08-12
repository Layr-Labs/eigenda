package lightnode

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestComputeShuffleOffset(t *testing.T) {
	tu.InitializeRandom()

	shufflePeriod := time.Second * time.Duration(rand.Intn(1000)+1)

	uniqueOffsets := make(map[time.Duration]bool)

	// Ensure that shuffle offset is always between 0 and shufflePeriod.
	iterations := 1000
	for i := 0; i < iterations; i++ {
		seed := rand.Uint64()
		offset := ComputeShuffleOffset(seed, shufflePeriod)
		assert.True(t, offset >= 0)
		assert.True(t, offset < shufflePeriod)
		uniqueOffsets[offset] = true

		// Computing the offset a second time should yield the same result.
		offset2 := ComputeShuffleOffset(seed, shufflePeriod)
		assert.Equal(t, offset, offset2)
	}

	// Although it's possible we end up with two offsets that are the same, the probability
	// is not high. Although it's unlikely for there to be even one collision, to avoid a flaky test
	//only assert that we have at least 99% unique offsets.
	assert.True(t, len(uniqueOffsets) >= int(float64(iterations)*0.99))
}

func TestComputeShuffleOffsetNegativeShufflePeriod(t *testing.T) {
	assert.Panics(t, func() {
		_ = ComputeShuffleOffset(0, -time.Second)
	})
}

func TestComputeShuffleEpochHandCraftedScenario(t *testing.T) {
	genesis := time.Unix(0, 0)
	shufflePeriod := time.Second * 100
	shuffleOffset := time.Second * 50

	// Requesting the epoch before genesis should return an error.
	assert.Panics(t, func() {
		_ = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(-time.Second))
	})

	// Requesting the epoch at genesis should return 0.
	epoch := ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis)
	assert.Equal(t, uint64(0), epoch)

	// Requesting the epoch one nanosecond after genesis should return 0.
	epoch = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(1))
	assert.Equal(t, uint64(0), epoch)

	// Requesting the epoch in-between genesis and the first shuffle should return 0.
	epoch = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(shuffleOffset/2))
	assert.Equal(t, uint64(0), epoch)

	// Requesting the epoch one nanosecond before the first shuffle should return 0.
	epoch = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(shuffleOffset-1))
	assert.Equal(t, uint64(0), epoch)

	// Requesting the epoch at the exact first shuffle time should return 1.
	epoch = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(shuffleOffset))
	assert.Equal(t, uint64(1), epoch)

	// Requesting the epoch one nanosecond after the first shuffle should return 1.
	epoch = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(shuffleOffset+1))
	assert.Equal(t, uint64(1), epoch)

	// Requesting the epoch at the mid-point between the first and second shuffle should return 1.
	epoch = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(shuffleOffset+shufflePeriod/2))
	assert.Equal(t, uint64(1), epoch)

	// Requesting the epoch one nanosecond before the second shuffle should return 1.
	epoch = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(shuffleOffset+shufflePeriod-1))
	assert.Equal(t, uint64(1), epoch)

	// Requesting the epoch at the exact second shuffle time should return 2.
	epoch = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(shuffleOffset+shufflePeriod))
	assert.Equal(t, uint64(2), epoch)

	// Moving forward 1000 shuffle periods should put us in epoch 1002.
	epoch = ComputeShuffleEpoch(genesis, shufflePeriod, shuffleOffset, genesis.Add(shuffleOffset+shufflePeriod*1001))
	assert.Equal(t, uint64(1002), epoch)
}

// Simulate a scenario where we move down a timeline. Verify that the epoch increases as expected.
func TestComputeShuffleEpochTimeWalk(t *testing.T) {
	tu.InitializeRandom()

	genesis := time.Unix(int64(rand.Intn(1_000_000)), 0)
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	seed := rand.Uint64()
	offset := ComputeShuffleOffset(seed, shufflePeriod)

	// If we move forward in time 1000 shuffle periods, we should see the epoch increase 1000 times.
	endTime := genesis.Add(shufflePeriod * 1000)
	epochCount := 0

	previousEpoch := ComputeShuffleEpoch(genesis, shufflePeriod, offset, genesis)

	for now := genesis; !now.After(endTime); now = now.Add(time.Second) {

		currentEpoch := ComputeShuffleEpoch(genesis, shufflePeriod, offset, now)
		if currentEpoch != previousEpoch {
			epochCount++
			previousEpoch = currentEpoch
		}
	}

	assert.Equal(t, 1000, epochCount)
}

func TestComputeShuffleEpochInvalidShufflePeriod(t *testing.T) {
	assert.Panics(t, func() {
		_ = ComputeShuffleEpoch(time.Now(), -time.Second, time.Second, time.Now())
	})

	assert.Panics(t, func() {
		_ = ComputeShuffleEpoch(time.Now(), 0, time.Second, time.Now())
	})
}

func TestComputeShuffleEpochInvalidShuffleOffset(t *testing.T) {
	assert.Panics(t, func() {
		_ = ComputeShuffleEpoch(time.Now(), time.Second, -time.Second, time.Now())
	})
}

func TestComputeChunkGroup(t *testing.T) {
	tu.InitializeRandom()

	chunkGroupCount := uint64(rand.Intn(1_000) + 10_000)
	seed := rand.Uint64()
	firstEpoch := uint64(rand.Intn(1000))

	uniqueChunkGroups := make(map[uint64]bool)

	// Ensure that the chunk group is always between 0 and chunkGroupCount.
	iterations := 1000
	for i := 0; i < iterations; i++ {
		epoch := firstEpoch + uint64(i)
		chunkGroup := ComputeChunkGroup(seed, epoch, chunkGroupCount)
		assert.True(t, chunkGroup < chunkGroupCount)
		uniqueChunkGroups[chunkGroup] = true
	}

	// We will have some collisions, but we should have at least 50% unique chunk groups with high probability.
	assert.True(t, len(uniqueChunkGroups) >= int(float64(iterations)*0.5))
}

func TestComputeEndOfShuffleEpochErrors(t *testing.T) {
	assert.Panics(t, func() {
		_ = ComputeEndOfShuffleEpoch(time.Now(), -time.Second, time.Second, 0)
	})

	assert.Panics(t, func() {
		_ = ComputeEndOfShuffleEpoch(time.Now(), time.Second, -time.Second, 0)
	})
}

func TestComputeEndOfShuffleEpochHandCraftedScenario(t *testing.T) {

	genesis := time.Unix(0, 0)
	shufflePeriod := time.Second * 100
	shuffleOffset := time.Second * 50

	assert.Equal(t, time.Unix(50, 0), ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 0))
	assert.Equal(t, time.Unix(150, 0), ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 1))
	assert.Equal(t, time.Unix(250, 0), ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 2))
	assert.Equal(t, time.Unix(350, 0), ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 3))
	assert.Equal(t, time.Unix(450, 0), ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 4))
}

func TestComputeEndOfShuffleEpoch(t *testing.T) {
	tu.InitializeRandom()

	genesis := time.Unix(int64(rand.Intn(1_000_000)), 0)
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	shuffleOffset := time.Second * time.Duration(rand.Intn(10)+1)

	assert.Equal(t, genesis.Add(shuffleOffset),
		ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 0))
	assert.Equal(t, genesis.Add(shuffleOffset).Add(shufflePeriod),
		ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 1))
	assert.Equal(t, genesis.Add(shuffleOffset).Add(shufflePeriod*2),
		ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 2))
	assert.Equal(t, genesis.Add(shuffleOffset).Add(shufflePeriod*3),
		ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 3))
	assert.Equal(t, genesis.Add(shuffleOffset).Add(shufflePeriod*4),
		ComputeEndOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 4))
}

func TestComputeStartOfShuffleEpochHandCraftedScenario(t *testing.T) {

	genesis := time.Unix(0, 0)
	shufflePeriod := time.Second * 100
	shuffleOffset := time.Second * 50

	assert.Equal(t, genesis, ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 0))
	assert.Equal(t, time.Unix(50, 0), ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 1))
	assert.Equal(t, time.Unix(150, 0), ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 2))
	assert.Equal(t, time.Unix(250, 0), ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 3))
	assert.Equal(t, time.Unix(350, 0), ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 4))
}

func TestComputeStartOfShuffleEpoch(t *testing.T) {
	tu.InitializeRandom()

	genesis := time.Unix(int64(rand.Intn(1_000_000)), 0)
	shufflePeriod := time.Second * time.Duration(rand.Intn(10)+1)
	shuffleOffset := time.Second * time.Duration(rand.Intn(10)+1)

	assert.Equal(t, genesis,
		ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 0))
	assert.Equal(t, genesis.Add(shuffleOffset),
		ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 1))
	assert.Equal(t, genesis.Add(shuffleOffset).Add(shufflePeriod),
		ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 2))
	assert.Equal(t, genesis.Add(shuffleOffset).Add(shufflePeriod*2),
		ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 3))
	assert.Equal(t, genesis.Add(shuffleOffset).Add(shufflePeriod*3),
		ComputeStartOfShuffleEpoch(genesis, shufflePeriod, shuffleOffset, 4))
}
