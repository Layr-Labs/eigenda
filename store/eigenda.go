package store

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	"github.com/Layr-Labs/op-plasma-eigenda/verify"
	"github.com/ethereum/go-ethereum/rlp"
)

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

// Get retrieves the given key if it's present in the key-value data store.
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

// PutWithCommitment attempts to insert the given key and value into the key-value data store
// Since EigenDA only has a commitment after blob dispersal this method is unsupported
func (e EigenDAStore) PutWithComm(ctx context.Context, key []byte, value []byte) error {
	return fmt.Errorf("EigenDA plasma store does not support PutWithComm()")
}

// PutWithoutComm inserts the given value into the key-value data store and returns the corresponding commitment
func (e EigenDAStore) PutWithoutComm(ctx context.Context, value []byte) (comm []byte, err error) {
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
