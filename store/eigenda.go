package store

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

type EigenDAStoreConfig struct {
	MaxBlobSizeBytes     uint64
	EthConfirmationDepth uint64

	// The total amount of time that the client will spend waiting for EigenDA to confirm a blob
	StatusQueryTimeout time.Duration
}

// EigenDAStore does storage interactions and verifications for blobs with DA.
type EigenDAStore struct {
	client   *clients.EigenDAClient
	verifier *verify.Verifier
	cfg      *EigenDAStoreConfig
	log      log.Logger
}

var _ Store = (*EigenDAStore)(nil)

func NewEigenDAStore(ctx context.Context, client *clients.EigenDAClient, v *verify.Verifier, log log.Logger, cfg *EigenDAStoreConfig) (*EigenDAStore, error) {
	return &EigenDAStore{
		client:   client,
		verifier: v,
		log:      log,
		cfg:      cfg,
	}, nil
}

// Get fetches a blob from DA using certificate fields and verifies blob
// against commitment to ensure data is valid and non-tampered.
func (e EigenDAStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	var cert verify.Certificate
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
		return nil, fmt.Errorf("EigenDA client failed to re-encode blob: %w", err)
	}

	err = e.verifier.VerifyCommitment(cert.BlobHeader.Commitment, encodedBlob)
	if err != nil {
		return nil, err
	}

	return decodedBlob, nil
}

// Put disperses a blob for some pre-image and returns the associated RLP encoded certificate commit.
func (e EigenDAStore) Put(ctx context.Context, value []byte) (comm []byte, err error) {
	dispersalStart := time.Now()
	blobInfo, err := e.client.PutBlob(ctx, value)
	if err != nil {
		return nil, err
	}
	cert := (*verify.Certificate)(blobInfo)

	encodedBlob, err := e.client.GetCodec().EncodeBlob(value)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to re-encode blob: %w", err)
	}
	if uint64(len(encodedBlob)) > e.cfg.MaxBlobSizeBytes {
		return nil, fmt.Errorf("encoded blob is larger than max blob size: blob length %d, max blob size %d", len(value), e.cfg.MaxBlobSizeBytes)
	}


	err = e.verifier.VerifyCommitment(cert.BlobHeader.Commitment, encodedBlob)
	if err != nil {
		return nil, err
	}

	dispersalDuration := time.Since(dispersalStart)
	remainingTimeout := e.cfg.StatusQueryTimeout - dispersalDuration

	ticker := time.NewTicker(12 * time.Second) // avg. eth block time
	defer ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), remainingTimeout)
	defer cancel()

	done := false
	for !done {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out when trying to verify the DA certificate for a blob batch after dispersal")
		case <-ticker.C:
			err = e.verifier.VerifyCert(cert)
			if err == nil {
				done = true

			} else if !errors.Is(err, verify.ErrBatchMetadataHashNotFound) {
				return nil, err
			} else {
				e.log.Info("Blob confirmed, waiting for sufficient confirmation depth...", "targetDepth", e.cfg.EthConfirmationDepth)
			}
		}
	}

	bytes, err := rlp.EncodeToBytes(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to encode DA cert to RLP format: %w", err)
	}

	return bytes, nil
}

// Entries are a no-op for EigenDA Store
func (e EigenDAStore) Stats() *Stats {
	return nil
}

// Key is used to recover certificate fields and that verifies blob
// against commitment to ensure data is valid and non-tampered.
func (e EigenDAStore) EncodeAndVerify(ctx context.Context, key []byte, value []byte) ([]byte, error) {
	var cert verify.Certificate
	err := rlp.DecodeBytes(key, &cert)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DA cert to RLP format: %w", err)
	}

	// reencode blob for verification
	encodedBlob, err := e.client.GetCodec().EncodeBlob(value)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to re-encode blob: %w", err)
	}
	if uint64(len(encodedBlob)) > e.cfg.MaxBlobSizeBytes {
		return nil, fmt.Errorf("encoded blob is larger than max blob size: blob length %d, max blob size %d", len(value), e.cfg.MaxBlobSizeBytes)
	}

	err = e.verifier.VerifyCommitment(cert.BlobHeader.Commitment, encodedBlob)
	if err != nil {
		return nil, err
	}

	return value, nil
}
