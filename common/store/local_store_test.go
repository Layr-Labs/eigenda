package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/stretchr/testify/assert"
)

var (
	inmemBucketStoreSize = 1000
)

func TestLocalStore(t *testing.T) {

	localStore, err := store.NewLocalParamStore[common.RateBucketParams](inmemBucketStoreSize, "disperser_lock")
	assert.NoError(t, err)

	ctx := context.Background()

	p := &common.RateBucketParams{
		BucketLevels:    []time.Duration{time.Second, time.Minute},
		LastRequestTime: time.Now(),
	}

	p2, err := localStore.GetItem(ctx, "testRetriever")
	assert.Error(t, err)
	assert.Nil(t, p2)

	err = localStore.UpdateItem(ctx, "testRetriever", p)
	assert.NoError(t, err)

	p2, err = localStore.GetItem(ctx, "testRetriever")

	assert.NoError(t, err)
	assert.Equal(t, p, p2)

}

func TestLockingMechanism(t *testing.T) {
	localStore, err := store.NewLocalParamStore[common.RateBucketParams](inmemBucketStoreSize, "disperser_lock")
	assert.NoError(t, err)
	key := "lockKey"
	p := &common.RateBucketParams{
		BucketLevels:    []time.Duration{time.Second, time.Minute},
		LastRequestTime: time.Now(),
	}

	ctx := context.Background()

	if localStore.AcquireLock(key, 0) {

		p2, err := localStore.GetItem(ctx, key)
		assert.Error(t, err)
		assert.Nil(t, p2)

		err = localStore.UpdateItem(ctx, key, p)
		assert.NoError(t, err)

		p2, err = localStore.GetItem(ctx, key)

		assert.NoError(t, err)
		assert.Equal(t, p, p2)

		err = localStore.ReleaseLock(key)
		assert.NoError(t, err)
	} else {
		t.Error("Failed to acquire lock")
	}
}
