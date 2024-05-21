package store

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	"github.com/Layr-Labs/op-plasma-eigenda/verify"
	"github.com/ethereum/go-ethereum/rlp"
)

// EigenDAStore does storage interactions and verifications for blobs with DA.
type EigenDAStore struct {
	client   *eigenda.EigenDAClient
	verifier *verify.Verifier
}

func NewEigenDAStore(ctx context.Context, client *eigenda.EigenDAClient, v *verify.Verifier) (*EigenDAStore, error) {
	return &EigenDAStore{
		client:   client,
		verifier: v,
	}, nil
}

// Get fetches a blob from DA using certificate fields and verifies blob
// against commitment to ensure data is valid and non-tampered.
func (e EigenDAStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	var cert eigenda.Cert
	err := rlp.DecodeBytes(key, &cert)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DA cert to RLP format: %w", err)
	}
	blob, err := e.client.RetrieveBlob(ctx, cert.BatchHeaderHash, cert.BlobIndex)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to retrieve blob: %w", err)
	}

	err = e.verifier.Verify(cert, eigenda.EncodeToBlob(blob))
	if err != nil {
		return nil, err
	}

	return blob, nil
}

// Put disperses a blob for some pre-image and returns the associated RLP encoded certificate commit.
func (e EigenDAStore) Put(ctx context.Context, value []byte) (comm []byte, err error) {
	cert, err := e.client.DisperseBlob(ctx, value)
	if err != nil {
		return nil, err
	}

	bytes, err := rlp.EncodeToBytes(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to encode DA cert to RLP format: %w", err)
	}

	return bytes, nil
}
