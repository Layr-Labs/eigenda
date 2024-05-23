package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testPreimage = "Four score and seven years ago"
)

func TestGetSet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ms, err := NewMemStore(
		ctx,
		&MemStoreConfig{
			Enabled:        true,
			BlobExpiration: time.Hour * 1000,
		},
	)

	assert.NoError(t, err)

	expected := []byte(testPreimage)
	key, err := ms.Put(ctx, expected)
	assert.NoError(t, err)

	actual, err := ms.Get(ctx, key)
	assert.NoError(t, err)

	assert.Equal(t, actual, expected)
}

func TestExpiration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ms, err := NewMemStore(
		ctx,
		&MemStoreConfig{
			Enabled:        true,
			BlobExpiration: time.Millisecond * 10,
		},
	)

	assert.NoError(t, err)

	preimage := []byte(testPreimage)
	key, err := ms.Put(ctx, preimage)
	assert.NoError(t, err)

	// sleep 1ms and verify that older entries are removed
	time.Sleep(time.Second * 1)

	_, err = ms.Get(ctx, key)
	assert.Error(t, err)

}
