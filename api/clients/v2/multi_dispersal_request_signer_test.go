package clients

import (
	"context"
	"testing"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSigner is a mock DispersalRequestSigner for testing.
type mockSigner struct {
	signFunc func(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error)
}

func (m *mockSigner) SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
	return m.signFunc(ctx, request)
}

func TestNewMultiDispersalRequestSigner(t *testing.T) {
	signer1 := &mockSigner{
		signFunc: func(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
			return []byte("signature_1"), nil
		},
	}
	signer0 := &mockSigner{
		signFunc: func(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
			return []byte("signature_0"), nil
		},
	}

	t.Run("valid configuration", func(t *testing.T) {
		config := MultiDispersalRequestSignerConfig{
			Signers: map[uint32]DispersalRequestSigner{
				1: signer1,
				0: signer0,
			},
			DisperserIDs: []uint32{1, 0},
		}

		signer, err := NewMultiDispersalRequestSigner(config)
		require.NoError(t, err)
		require.NotNil(t, signer)
		assert.Equal(t, []uint32{1, 0}, signer.GetDisperserIDs())
		assert.True(t, signer.HasDisperserID(1))
		assert.True(t, signer.HasDisperserID(0))
		assert.False(t, signer.HasDisperserID(2))
	})

	t.Run("no signers", func(t *testing.T) {
		config := MultiDispersalRequestSignerConfig{
			Signers:      map[uint32]DispersalRequestSigner{},
			DisperserIDs: []uint32{1},
		}

		signer, err := NewMultiDispersalRequestSigner(config)
		require.Error(t, err)
		assert.Nil(t, signer)
		assert.Contains(t, err.Error(), "at least one signer is required")
	})

	t.Run("no disperser IDs", func(t *testing.T) {
		config := MultiDispersalRequestSignerConfig{
			Signers: map[uint32]DispersalRequestSigner{
				1: signer1,
			},
			DisperserIDs: []uint32{},
		}

		signer, err := NewMultiDispersalRequestSigner(config)
		require.Error(t, err)
		assert.Nil(t, signer)
		assert.Contains(t, err.Error(), "at least one disperser ID is required")
	})

	t.Run("disperser ID without signer", func(t *testing.T) {
		config := MultiDispersalRequestSignerConfig{
			Signers: map[uint32]DispersalRequestSigner{
				1: signer1,
			},
			DisperserIDs: []uint32{1, 0}, // 0 has no signer
		}

		signer, err := NewMultiDispersalRequestSigner(config)
		require.Error(t, err)
		assert.Nil(t, signer)
		assert.Contains(t, err.Error(), "no signer configured for disperser ID 0")
	})
}

func TestMultiDispersalRequestSigner_SignStoreChunksRequest(t *testing.T) {
	signer1 := &mockSigner{
		signFunc: func(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
			return []byte("signature_1"), nil
		},
	}
	signer0 := &mockSigner{
		signFunc: func(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
			return []byte("signature_0"), nil
		},
	}

	config := MultiDispersalRequestSignerConfig{
		Signers: map[uint32]DispersalRequestSigner{
			1: signer1,
			0: signer0,
		},
		DisperserIDs: []uint32{1, 0},
	}

	multiSigner, err := NewMultiDispersalRequestSigner(config)
	require.NoError(t, err)

	t.Run("sign with ID 1", func(t *testing.T) {
		request := &grpc.StoreChunksRequest{}
		signature, err := multiSigner.SignStoreChunksRequest(context.Background(), request, 1)

		require.NoError(t, err)
		assert.Equal(t, []byte("signature_1"), signature)
		assert.Equal(t, uint32(1), request.GetDisperserID())
	})

	t.Run("sign with ID 0", func(t *testing.T) {
		request := &grpc.StoreChunksRequest{}
		signature, err := multiSigner.SignStoreChunksRequest(context.Background(), request, 0)

		require.NoError(t, err)
		assert.Equal(t, []byte("signature_0"), signature)
		assert.Equal(t, uint32(0), request.GetDisperserID())
	})

	t.Run("sign with unknown ID", func(t *testing.T) {
		request := &grpc.StoreChunksRequest{}
		signature, err := multiSigner.SignStoreChunksRequest(context.Background(), request, 99)

		require.Error(t, err)
		assert.Nil(t, signature)
		assert.Contains(t, err.Error(), "no signer configured for disperser ID 99")
	})
}
