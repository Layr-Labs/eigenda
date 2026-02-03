package clients

import (
	"context"
	"fmt"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/pkg/errors"
)

// MultiDispersalRequestSigner manages multiple disperser IDs and their corresponding signers.
// It allows signing requests with any configured disperser ID, supporting seamless migration
// between different disperser identities.
type MultiDispersalRequestSigner struct {
	// signers maps disperser ID to its corresponding signer
	signers map[uint32]DispersalRequestSigner
	// disperserIDs is the ordered list of disperser IDs (priority order)
	disperserIDs []uint32
}

// MultiDispersalRequestSignerConfig configures the multi-signer with disperser IDs and their signers.
type MultiDispersalRequestSignerConfig struct {
	// Signers maps each disperser ID to its corresponding signer
	Signers map[uint32]DispersalRequestSigner
	// DisperserIDs is the ordered list of IDs to try (in priority order)
	DisperserIDs []uint32
}

// NewMultiDispersalRequestSigner creates a new MultiDispersalRequestSigner.
func NewMultiDispersalRequestSigner(
	config MultiDispersalRequestSignerConfig,
) (*MultiDispersalRequestSigner, error) {
	if len(config.Signers) == 0 {
		return nil, errors.New("at least one signer is required")
	}
	if len(config.DisperserIDs) == 0 {
		return nil, errors.New("at least one disperser ID is required")
	}

	// Verify all disperser IDs have corresponding signers
	for _, id := range config.DisperserIDs {
		if _, ok := config.Signers[id]; !ok {
			return nil, fmt.Errorf("no signer configured for disperser ID %d", id)
		}
	}

	return &MultiDispersalRequestSigner{
		signers:      config.Signers,
		disperserIDs: config.DisperserIDs,
	}, nil
}

// SignStoreChunksRequest signs a StoreChunksRequest using the signer for the specified disperser ID.
// It sets the disperser ID in the request before signing.
func (m *MultiDispersalRequestSigner) SignStoreChunksRequest(
	ctx context.Context,
	request *grpc.StoreChunksRequest,
	disperserID uint32,
) ([]byte, error) {
	signer, ok := m.signers[disperserID]
	if !ok {
		return nil, fmt.Errorf("no signer configured for disperser ID %d", disperserID)
	}

	request.DisperserID = disperserID
	signature, err := signer.SignStoreChunksRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to sign with disperser ID %d: %w", disperserID, err)
	}
	return signature, nil
}

// GetDisperserIDs returns the ordered list of disperser IDs.
func (m *MultiDispersalRequestSigner) GetDisperserIDs() []uint32 {
	return m.disperserIDs
}

// HasDisperserID returns true if the signer is configured for the given disperser ID.
func (m *MultiDispersalRequestSigner) HasDisperserID(id uint32) bool {
	_, ok := m.signers[id]
	return ok
}
