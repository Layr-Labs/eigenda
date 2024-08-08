package test

import (
	"context"
	"github.com/Layr-Labs/eigenda/api/clients"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var _ clients.DisperserClient = (*mockDisperserClient)(nil)

type mockDisperserClient struct {
	t *testing.T
	// if true, DisperseBlobAuthenticated is expected to be used, otherwise DisperseBlob is expected to be used
	authenticated bool

	// The next status, key, and error to return from DisperseBlob or DisperseBlobAuthenticated
	StatusToReturn        disperser.BlobStatus
	KeyToReturn           []byte
	DispenseErrorToReturn error

	// The previous values passed to DisperseBlob or DisperseBlobAuthenticated
	ProvidedData   []byte
	ProvidedQuorum []uint8

	// Incremented each time DisperseBlob or DisperseBlobAuthenticated is called.
	DisperseCount uint

	// A map from key (in string form) to the status to return from GetBlobStatus. If nil, then an error is returned.
	StatusMap map[string]disperser_rpc.BlobStatus

	// Incremented each time GetBlobStatus is called.
	GetStatusCount uint

	lock *sync.Mutex
}

func newMockDisperserClient(t *testing.T, lock *sync.Mutex, authenticated bool) *mockDisperserClient {
	return &mockDisperserClient{
		t:             t,
		lock:          lock,
		authenticated: authenticated,
		StatusMap:     make(map[string]disperser_rpc.BlobStatus),
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
	m.DisperseCount++
	return &m.StatusToReturn, m.KeyToReturn, m.DispenseErrorToReturn
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
	m.DisperseCount++
	return &m.StatusToReturn, m.KeyToReturn, m.DispenseErrorToReturn
}

func (m *mockDisperserClient) GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	status := m.StatusMap[string(key)]
	m.GetStatusCount++

	return &disperser_rpc.BlobStatusReply{
		Status: status,
		Info: &disperser_rpc.BlobInfo{
			BlobVerificationProof: &disperser_rpc.BlobVerificationProof{
				BatchMetadata: &disperser_rpc.BatchMetadata{},
			},
		},
	}, nil
}

func (m *mockDisperserClient) RetrieveBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error) {
	panic("this method should not be called by the generator utility")
}
