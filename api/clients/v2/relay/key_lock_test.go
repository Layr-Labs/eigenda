package relay

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

func TestKeyLock(t *testing.T) {
	// test in a field of 100 unique keys
	keyCount := 100

	// keep an atomic count, and a non-atomic count for each key
	// the atomic count can be used at the end of the test, to make sure that the non-atomic count was handled correctly
	atomicKeyAccessCounts := make([]atomic.Uint32, keyCount)
	nonAtomicKeyAccessCounts := make([]uint32, keyCount)
	for i := 0; i < keyCount; i++ {
		atomicKeyAccessCounts = append(atomicKeyAccessCounts, atomic.Uint32{})
		nonAtomicKeyAccessCounts = append(nonAtomicKeyAccessCounts, uint32(0))
	}

	keyLock := NewKeyLock[uint32]()

	var waitGroup sync.WaitGroup

	targetValue := uint32(1000)
	worker := func() {
		workerRandom := random.NewTestRandom()

		for {
			// randomly select a key to access
			keyToAccess := uint32(workerRandom.Intn(keyCount))
			newValue := atomicKeyAccessCounts[keyToAccess].Add(1)

			unlock := keyLock.AcquireKeyLock(keyToAccess)
			// increment the non-atomic count after acquiring access
			// if the access controls are working correctly, this is a safe operation
			nonAtomicKeyAccessCounts[keyToAccess] = nonAtomicKeyAccessCounts[keyToAccess] + 1
			unlock()

			// each worker stops looping after it sees a counter that has increased to targetValue
			if newValue >= targetValue {
				break
			}
		}

		waitGroup.Done()
	}

	// start up 100 concurrent workers
	for i := 0; i < 100; i++ {
		waitGroup.Add(1)
		go worker()
	}
	waitGroup.Wait()

	for i := 0; i < keyCount; i++ {
		require.Equal(t, atomicKeyAccessCounts[i].Load(), nonAtomicKeyAccessCounts[i])
	}
}
