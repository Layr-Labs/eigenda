package store

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/eigenda"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/rlp"
)

// EigenDAStore does storage interactions and verifications for blobs with DA.
type EigenDAStore struct {
	client   *clients.EigenDAClient
	verifier *verify.Verifier
}

func NewEigenDAStore(ctx context.Context, client *clients.EigenDAClient, v *verify.Verifier) (*EigenDAStore, error) {
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
	blob, err := e.client.GetBlob(ctx, cert.BatchHeaderHash, cert.BlobIndex)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to retrieve blob: %w", err)
	}

	// TODO: Blob codec should be read from the blob's header, not the client's
	// write codec, which may differ from the blob's codec.  This will work for
	// now since these are the same.
	encodedBlob, err := e.client.PutCodec.EncodeBlob(blob)
	if err != nil {
		return nil, err
	}
	err = e.verifier.Verify(cert, encodedBlob)
	if err != nil {
		return nil, err
	}

	return blob, nil
}

// Put disperses a blob for some pre-image and returns the associated RLP encoded certificate commit.
func (e EigenDAStore) Put(ctx context.Context, value []byte) (comm []byte, err error) {
	cert, err := e.client.PutBlob(ctx, value)
	if err != nil {
		return nil, err
	}

	bytes, err := rlp.EncodeToBytes(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to encode DA cert to RLP format: %w", err)
	}

	return bytes, nil
}
