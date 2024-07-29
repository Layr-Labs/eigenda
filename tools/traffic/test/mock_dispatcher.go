package test

import (
	"context"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type mockDisperserClient struct {
	t *testing.T
	// if true, DisperseBlobAuthenticated is expected to be used, otherwise DisperseBlob is expected to be used
	authenticated bool

	// The next status, key, and error to return from DisperseBlob or DisperseBlobAuthenticated
	StatusToReturn disperser.BlobStatus
	KeyToReturn    []byte
	ErrorToReturn  error

	// The previous values passed to DisperseBlob or DisperseBlobAuthenticated
	ProvidedData   []byte
	ProvidedQuorum []uint8

	// Incremented each time DisperseBlob or DisperseBlobAuthenticated is called.
	Count uint

	lock *sync.Mutex
}

func newMockDisperserClient(t *testing.T, lock *sync.Mutex, authenticated bool) *mockDisperserClient {
	return &mockDisperserClient{
		t:             t,
		lock:          lock,
		authenticated: authenticated,
	}
}

func (m *mockDisperserClient) DisperseBlob(
	ctx context.Context,
	data []byte,
	customQuorums []uint8) (*disperser.BlobStatus, []byte, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	assert.False(m.t, m.authenticated, "writer configured to use non-authenticated disperser method")
	m.ProvidedData = data
	m.ProvidedQuorum = customQuorums
	m.Count++
	return &m.StatusToReturn, m.KeyToReturn, m.ErrorToReturn
}

func (m *mockDisperserClient) DisperseBlobAuthenticated(
	ctx context.Context,
	data []byte,
	customQuorums []uint8) (*disperser.BlobStatus, []byte, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	assert.True(m.t, m.authenticated, "writer configured to use authenticated disperser method")
	m.ProvidedData = data
	m.ProvidedQuorum = customQuorums
	m.Count++
	return &m.StatusToReturn, m.KeyToReturn, m.ErrorToReturn
}

func (m *mockDisperserClient) GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error) {
	panic("this method should not be called in this test")
}

func (m *mockDisperserClient) RetrieveBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error) {
	panic("this method should not be called in this test")
}
