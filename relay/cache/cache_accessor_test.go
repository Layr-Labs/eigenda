package cache

import (
	"context"
	"errors"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRandomOperationsSingleThread(t *testing.T) {
	tu.InitializeRandom()

	dataSize := 1024

	baseData := make(map[int]string)
	for i := 0; i < dataSize; i++ {
		baseData[i] = tu.RandomString(10)
	}

	accessor := func(key int) (*string, error) {
		// Return an error if the key is a multiple of 17
		if key%17 == 0 {
			return nil, errors.New("intentional error")
		}

		str := baseData[key]
		return &str, nil
	}
	cacheSize := rand.Intn(dataSize) + 1

	cache := NewFIFOCache[int, *string](uint64(cacheSize), nil)
	ca, err := NewCacheAccessor[int, *string](cache, 0, accessor, nil)
	require.NoError(t, err)

	for i := 0; i < dataSize; i++ {
		value, err := ca.Get(context.Background(), i)

		if i%17 == 0 {
			require.Error(t, err)
			require.Nil(t, value)
		} else {
			require.NoError(t, err)
			require.Equal(t, baseData[i], *value)
		}
	}

	for k, v := range baseData {
		value, err := ca.Get(context.Background(), k)

		if k%17 == 0 {
			require.Error(t, err)
			require.Nil(t, value)
		} else {
			require.NoError(t, err)
			require.Equal(t, v, *value)
		}
	}
}

func TestCacheMisses(t *testing.T) {
	tu.InitializeRandom()

	cacheSize := rand.Intn(10) + 10
	keyCount := cacheSize + 1

	baseData := make(map[int]string)
	for i := 0; i < keyCount; i++ {
		baseData[i] = tu.RandomString(10)
	}

	cacheMissCount := atomic.Uint64{}

	accessor := func(key int) (*string, error) {
		cacheMissCount.Add(1)
		str := baseData[key]
		return &str, nil
	}

	cache := NewFIFOCache[int, *string](uint64(cacheSize), nil)
	ca, err := NewCacheAccessor[int, *string](cache, 0, accessor, nil)
	require.NoError(t, err)

	// Get the first cacheSize keys. This should fill the cache.
	expectedCacheMissCount := uint64(0)
	for i := 0; i < cacheSize; i++ {
		expectedCacheMissCount++
		value, err := ca.Get(context.Background(), i)
		require.NoError(t, err)
		require.Equal(t, baseData[i], *value)
		require.Equal(t, expectedCacheMissCount, cacheMissCount.Load())
	}

	// Get the first cacheSize keys again. This should not increase the cache miss count.
	for i := 0; i < cacheSize; i++ {
		value, err := ca.Get(context.Background(), i)
		require.NoError(t, err)
		require.Equal(t, baseData[i], *value)
		require.Equal(t, expectedCacheMissCount, cacheMissCount.Load())
	}

	// Read the last key. This should cause the first key to be evicted.
	expectedCacheMissCount++
	value, err := ca.Get(context.Background(), cacheSize)
	require.NoError(t, err)
	require.Equal(t, baseData[cacheSize], *value)

	// Read the keys in order. Due to the order of evictions, each read should result in a cache miss.
	for i := 0; i < cacheSize; i++ {
		expectedCacheMissCount++
		value, err := ca.Get(context.Background(), i)
		require.NoError(t, err)
		require.Equal(t, baseData[i], *value)
		require.Equal(t, expectedCacheMissCount, cacheMissCount.Load())
	}
}

func ParallelAccessTest(t *testing.T, sleepEnabled bool) {
	tu.InitializeRandom()

	dataSize := 1024

	baseData := make(map[int]string)
	for i := 0; i < dataSize; i++ {
		baseData[i] = tu.RandomString(10)
	}

	accessorLock := sync.RWMutex{}
	cacheMissCount := atomic.Uint64{}
	accessor := func(key int) (*string, error) {

		// Intentionally block if accessorLock is held by the outside scope.
		// Used to provoke specific race conditions.
		accessorLock.Lock()
		defer accessorLock.Unlock()

		cacheMissCount.Add(1)

		str := baseData[key]
		return &str, nil
	}
	cacheSize := rand.Intn(dataSize) + 1

	cache := NewFIFOCache[int, *string](uint64(cacheSize), nil)
	ca, err := NewCacheAccessor[int, *string](cache, 0, accessor, nil)
	require.NoError(t, err)

	// Lock the accessor. This will cause all cache misses to block.
	accessorLock.Lock()

	// Start several goroutines that will attempt to access the same key.
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			value, err := ca.Get(context.Background(), 0)
			require.NoError(t, err)
			require.Equal(t, baseData[0], *value)
		}()
	}

	if sleepEnabled {
		// Wait for the goroutines to start. We want to give the goroutines a chance to do naughty things if they want.
		// Eliminating this sleep will not cause the test to fail, but it may cause the test not to exercise the
		// desired race condition.
		time.Sleep(100 * time.Millisecond)
	}

	// Unlock the accessor. This will allow the goroutines to proceed.
	accessorLock.Unlock()

	// Wait for the goroutines to finish.
	wg.Wait()

	// Only one of the goroutines should have called into the accessor.
	require.Equal(t, uint64(1), cacheMissCount.Load())

	// Fetching the key again should not result in a cache miss.
	value, err := ca.Get(context.Background(), 0)
	require.NoError(t, err)
	require.Equal(t, baseData[0], *value)
	require.Equal(t, uint64(1), cacheMissCount.Load())

	// The internal lookupsInProgress map should no longer contain the key.
	require.Equal(t, 0, len(ca.(*cacheAccessor[int, *string]).lookupsInProgress))
}

func TestParallelAccess(t *testing.T) {
	// To show that the sleep is not necessary, we run the test twice: once with the sleep enabled and once without.
	// The purpose of the sleep is to make a certain type of race condition more likely to occur.

	ParallelAccessTest(t, false)
	ParallelAccessTest(t, true)
}

func TestParallelAccessWithError(t *testing.T) {
	tu.InitializeRandom()

	accessorLock := sync.RWMutex{}
	cacheMissCount := atomic.Uint64{}
	accessor := func(key int) (*string, error) {
		// Intentionally block if accessorLock is held by the outside scope.
		// Used to provoke specific race conditions.
		accessorLock.Lock()
		defer accessorLock.Unlock()

		cacheMissCount.Add(1)

		return nil, errors.New("intentional error")
	}
	cacheSize := 100

	cache := NewFIFOCache[int, *string](uint64(cacheSize), nil)
	ca, err := NewCacheAccessor[int, *string](cache, 0, accessor, nil)
	require.NoError(t, err)

	// Lock the accessor. This will cause all cache misses to block.
	accessorLock.Lock()

	// Start several goroutines that will attempt to access the same key.
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			value, err := ca.Get(context.Background(), 0)
			require.Nil(t, value)
			require.Equal(t, errors.New("intentional error"), err)
		}()
	}

	// Wait for the goroutines to start. We want to give the goroutines a chance to do naughty things if they want.
	// Eliminating this sleep will not cause the test to fail, but it may cause the test not to exercise the
	// desired race condition.
	time.Sleep(100 * time.Millisecond)

	// Unlock the accessor. This will allow the goroutines to proceed.
	accessorLock.Unlock()

	// Wait for the goroutines to finish.
	wg.Wait()

	// At least one of the goroutines should have called into the accessor. In theory all of them could have,
	// but most likely it will be exactly one.
	count := cacheMissCount.Load()
	require.True(t, count >= 1)

	// Fetching the key again should result in another cache miss since the previous fetch failed.
	value, err := ca.Get(context.Background(), 0)
	require.Nil(t, value)
	require.Equal(t, errors.New("intentional error"), err)
	require.Equal(t, count+1, cacheMissCount.Load())

	// The internal lookupsInProgress map should no longer contain the key.
	require.Equal(t, 0, len(ca.(*cacheAccessor[int, *string]).lookupsInProgress))
}

func TestConcurrencyLimiter(t *testing.T) {
	tu.InitializeRandom()

	dataSize := 1024

	baseData := make(map[int]string)
	for i := 0; i < dataSize; i++ {
		baseData[i] = tu.RandomString(10)
	}

	maxConcurrency := 10 + rand.Intn(10)

	accessorLock := sync.RWMutex{}
	accessorLock.Lock()
	activeAccessors := atomic.Int64{}
	accessor := func(key int) (*string, error) {
		activeAccessors.Add(1)
		accessorLock.Lock()
		defer func() {
			activeAccessors.Add(-1)
		}()
		accessorLock.Unlock()

		value := baseData[key]
		return &value, nil
	}

	cacheSize := 100

	cache := NewFIFOCache[int, *string](uint64(cacheSize), nil)
	ca, err := NewCacheAccessor[int, *string](cache, maxConcurrency, accessor, nil)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(dataSize)
	for i := 0; i < dataSize; i++ {
		boundI := i
		go func() {
			value, err := ca.Get(context.Background(), boundI)
			require.NoError(t, err)
			require.Equal(t, baseData[boundI], *value)
			wg.Done()
		}()
	}

	// Wait for the goroutines to start. We want to give the goroutines a chance to do naughty things if they want.
	// Eliminating this sleep will not cause the test to fail, but it may cause the test not to exercise the
	// desired race condition.
	time.Sleep(100 * time.Millisecond)

	// The number of active accessors should be less than or equal to the maximum concurrency.
	require.True(t, activeAccessors.Load() <= int64(maxConcurrency))

	// Unlock the accessor. This will allow the goroutines to proceed.
	accessorLock.Unlock()
	wg.Wait()
}

func TestOriginalRequesterTimesOut(t *testing.T) {
	tu.InitializeRandom()

	dataSize := 1024

	baseData := make(map[int]string)
	for i := 0; i < dataSize; i++ {
		baseData[i] = tu.RandomString(10)
	}

	accessorLock := sync.RWMutex{}
	cacheMissCount := atomic.Uint64{}
	accessor := func(key int) (*string, error) {

		// Intentionally block if accessorLock is held by the outside scope.
		// Used to provoke specific race conditions.
		accessorLock.Lock()
		defer accessorLock.Unlock()

		cacheMissCount.Add(1)

		str := baseData[key]
		return &str, nil
	}
	cacheSize := rand.Intn(dataSize) + 1

	cache := NewFIFOCache[int, *string](uint64(cacheSize), nil)
	ca, err := NewCacheAccessor[int, *string](cache, 0, accessor, nil)
	require.NoError(t, err)

	// Lock the accessor. This will cause all cache misses to block.
	accessorLock.Lock()

	// Start several goroutines that will attempt to access the same key.
	wg := sync.WaitGroup{}
	wg.Add(10)
	errCount := atomic.Uint64{}
	for i := 0; i < 10; i++ {

		var ctx context.Context
		if i == 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.Background(), 1*time.Millisecond)
			defer cancel()
		} else {
			ctx = context.Background()
		}

		go func() {
			defer wg.Done()
			value, err := ca.Get(ctx, 0)

			if err != nil {
				errCount.Add(1)
			} else {
				require.Equal(t, baseData[0], *value)
			}
		}()

		if i == 0 {
			// Give the thread with the small timeout a chance to start. Although this sleep statement is
			// not required for the test to pass, it makes it much more likely for this test to exercise
			// the intended code pathway.
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Unlock the accessor. This will allow the goroutines to proceed.
	accessorLock.Unlock()

	// Wait for the goroutines to finish.
	wg.Wait()

	// Only one of the goroutines should have called into the accessor.
	require.Equal(t, uint64(1), cacheMissCount.Load())

	// At most, one goroutine should have timed out.
	require.True(t, errCount.Load() <= 1)

	// Fetching the key again should not result in a cache miss.
	value, err := ca.Get(context.Background(), 0)
	require.NoError(t, err)
	require.Equal(t, baseData[0], *value)
	require.Equal(t, uint64(1), cacheMissCount.Load())

	// The internal lookupsInProgress map should no longer contain the key.
	require.Equal(t, 0, len(ca.(*cacheAccessor[int, *string]).lookupsInProgress))
}

func TestSecondaryRequesterTimesOut(t *testing.T) {
	tu.InitializeRandom()

	dataSize := 1024

	baseData := make(map[int]string)
	for i := 0; i < dataSize; i++ {
		baseData[i] = tu.RandomString(10)
	}

	accessorLock := sync.RWMutex{}
	cacheMissCount := atomic.Uint64{}
	accessor := func(key int) (*string, error) {

		// Intentionally block if accessorLock is held by the outside scope.
		// Used to provoke specific race conditions.
		accessorLock.Lock()
		defer accessorLock.Unlock()

		cacheMissCount.Add(1)

		str := baseData[key]
		return &str, nil
	}
	cacheSize := rand.Intn(dataSize) + 1

	cache := NewFIFOCache[int, *string](uint64(cacheSize), nil)
	ca, err := NewCacheAccessor[int, *string](cache, 0, accessor, nil)
	require.NoError(t, err)

	// Lock the accessor. This will cause all cache misses to block.
	accessorLock.Lock()

	// Start several goroutines that will attempt to access the same key.
	wg := sync.WaitGroup{}
	wg.Add(10)
	errCount := atomic.Uint64{}
	for i := 0; i < 10; i++ {

		var ctx context.Context
		if i == 1 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.Background(), 1*time.Millisecond)
			defer cancel()
		} else {
			ctx = context.Background()
		}

		go func() {
			defer wg.Done()
			value, err := ca.Get(ctx, 0)

			if err != nil {
				errCount.Add(1)
			} else {
				require.Equal(t, baseData[0], *value)
			}
		}()

		if i == 0 {
			// Give the thread with the context that won't time out a chance to start. Although this sleep statement is
			// not required for the test to pass, it makes it much more likely for this test to exercise
			// the intended code pathway.
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Give context a chance to time out. Although this sleep statement is not required for the test to pass, it makes
	// it much more likely for this test to exercise the intended code pathway.
	time.Sleep(100 * time.Millisecond)

	// Unlock the accessor. This will allow the goroutines to proceed.
	accessorLock.Unlock()

	// Wait for the goroutines to finish.
	wg.Wait()

	// Only one of the goroutines should have called into the accessor.
	require.Equal(t, uint64(1), cacheMissCount.Load())

	// At most, one goroutine should have timed out.
	require.True(t, errCount.Load() <= 1)

	// Fetching the key again should not result in a cache miss.
	value, err := ca.Get(context.Background(), 0)
	require.NoError(t, err)
	require.Equal(t, baseData[0], *value)
	require.Equal(t, uint64(1), cacheMissCount.Load())

	// The internal lookupsInProgress map should no longer contain the key.
	require.Equal(t, 0, len(ca.(*cacheAccessor[int, *string]).lookupsInProgress))
}
