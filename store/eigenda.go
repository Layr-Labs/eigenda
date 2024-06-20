package store

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/rlp"
)

// EigenDAStore does storage interactions and verifications for blobs with DA.
type EigenDAStore struct {
	client           *clients.EigenDAClient
	verifier         *verify.Verifier
	maxBlobSizeBytes uint64
}

var _ Store = (*EigenDAStore)(nil)

func NewEigenDAStore(ctx context.Context, client *clients.EigenDAClient, v *verify.Verifier, maxBlobSizeBytes uint64) (*EigenDAStore, error) {
	return &EigenDAStore{
		client:           client,
		verifier:         v,
		maxBlobSizeBytes: maxBlobSizeBytes,
	}, nil
}

// Get fetches a blob from DA using certificate fields and verifies blob
// against commitment to ensure data is valid and non-tampered.
func (e EigenDAStore) Get(ctx context.Context, key []byte, domain common.DomainType) ([]byte, error) {
	var cert common.Certificate
	err := rlp.DecodeBytes(key, &cert)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DA cert to RLP format: %w", err)
	}

	decodedBlob, err := e.client.GetBlob(ctx, cert.BlobVerificationProof.BatchMetadata.BatchHeaderHash, cert.BlobVerificationProof.BlobIndex)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to retrieve decoded blob: %w", err)
	}
	// reencode blob for verification
	encodedBlob, err := e.client.GetCodec().EncodeBlob(decodedBlob)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to reencode blob: %w", err)
	}

	err = e.verifier.Verify(cert.BlobHeader.Commitment, encodedBlob)
	if err != nil {
		return nil, err
	}

	switch domain {
	case common.BinaryDomain:
		return decodedBlob, nil
	case common.PolyDomain:
		return encodedBlob, nil
	default:
		return nil, fmt.Errorf("unexpected domain type: %d", domain)
	}
}

// Put disperses a blob for some pre-image and returns the associated RLP encoded certificate commit.
func (e EigenDAStore) Put(ctx context.Context, value []byte) (comm []byte, err error) {
	if uint64(len(value)) > e.maxBlobSizeBytes {
		return nil, fmt.Errorf("blob is larger than max blob size: blob length %d, max blob size %d", len(value), e.maxBlobSizeBytes)
	}
	cert, err := e.client.PutBlob(ctx, value)
	if err != nil {
		return nil, err
	}

	encodedBlob, err := e.client.GetCodec().EncodeBlob(value)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to reencode blob: %w", err)
	}
	err = e.verifier.Verify(cert.BlobHeader.Commitment, encodedBlob)
	if err != nil {
		return nil, err
	}

	bytes, err := rlp.EncodeToBytes(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to encode DA cert to RLP format: %w", err)
	}

	return bytes, nil
}

// Entries are a no-op for EigenDA Store
func (e EigenDAStore) Stats() *common.Stats {
	return nil
}
