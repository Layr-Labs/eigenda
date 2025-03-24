package memstore

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/ephemeraldb"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	eigenda_common "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	BytesPerFieldElement = 32
)

/*
MemStore is a simple in-memory store for blobs which uses an expiration
time to evict blobs to best emulate the ephemeral nature of blobs dispersed to
EigenDA V1 operators.
*/
type MemStore struct {
	*ephemeraldb.DB
	log logging.Logger

	// We only use the verifier for kzgCommitment verification.
	// MemStore generates random certs which can't be verified.
	// TODO: we should probably refactor the Verifier to be able to only take in a BlobVerifier here.
	verifier *verify.Verifier
	codec    codecs.BlobCodec
}

var _ common.GeneratedKeyStore = (*MemStore)(nil)

// New ... constructor
func New(
	ctx context.Context, verifier *verify.Verifier, log logging.Logger, config *memconfig.SafeConfig,
) (*MemStore, error) {
	return &MemStore{
		ephemeraldb.New(ctx, config, log),
		log,
		verifier,
		codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}, nil
}

// generateRandomCert ... generates random EigenDA V1 certificate
func (e *MemStore) generateRandomCert(blobValue []byte) (*verify.Certificate, error) {
	commitment, err := e.verifier.Commit(blobValue)
	if err != nil {
		return nil, err
	}

	// generate batch root hash
	entropy := make([]byte, 10)
	_, err = rand.Read(entropy)
	if err != nil {
		return nil, err
	}
	mockBatchRoot := crypto.Keccak256Hash(entropy)
	blockNum, _ := rand.Int(rand.Reader, big.NewInt(1000))

	num := uint32(blockNum.Uint64()) // #nosec G115

	cert := &verify.Certificate{
		BlobHeader: &disperser.BlobHeader{
			Commitment: &eigenda_common.G1Commitment{
				X: commitment.X.Marshal(),
				Y: commitment.Y.Marshal(),
			},
			DataLength: uint32((len(blobValue) + BytesPerFieldElement - 1) / BytesPerFieldElement), // #nosec G115
			BlobQuorumParams: []*disperser.BlobQuorumParam{
				{
					QuorumNumber:                    1,
					AdversaryThresholdPercentage:    29,
					ConfirmationThresholdPercentage: 30,
					ChunkLength:                     300,
				},
			},
		},
		BlobVerificationProof: &disperser.BlobVerificationProof{
			BatchMetadata: &disperser.BatchMetadata{
				BatchHeader: &disperser.BatchHeader{
					BatchRoot:               mockBatchRoot[:],
					QuorumNumbers:           []byte{0x1, 0x0},
					QuorumSignedPercentages: []byte{0x60, 0x90},
					ReferenceBlockNumber:    num,
				},
				SignatoryRecordHash:     mockBatchRoot[:],
				Fee:                     []byte{0x00},
				ConfirmationBlockNumber: num,
				BatchHeaderHash:         []byte{},
			},
			BatchId:        69,
			BlobIndex:      420,
			InclusionProof: entropy,
			QuorumIndexes:  []byte{0x1, 0x0},
		},
	}

	// compute batch header hash
	// this is necessary since EigenDA x Arbitrum reconstructs
	// the batch header hash before querying proxy since hash field
	// isn't persisted via the onchain cert posted to the inbox
	bh := cert.BlobVerificationProof.BatchMetadata.BatchHeader

	reducedHeader := core.BatchHeader{
		BatchRoot:            [32]byte(bh.GetBatchRoot()),
		ReferenceBlockNumber: uint(bh.GetReferenceBlockNumber()),
	}

	headerHash, err := reducedHeader.GetBatchHeaderHash()
	if err != nil {
		return nil, fmt.Errorf("generating batch header hash: %w", err)
	}

	cert.BlobVerificationProof.BatchMetadata.BatchHeaderHash = headerHash[:]

	return cert, nil
}

// Get fetches a value from the store.
func (e *MemStore) Get(_ context.Context, commit []byte) ([]byte, error) {
	encodedBlob, err := e.FetchEntry(crypto.Keccak256Hash(commit).Bytes())
	if err != nil {
		return nil, fmt.Errorf("fetching entry via v1 memstore: %w", err)
	}

	return e.codec.DecodeBlob(encodedBlob)
}

// Put inserts a value into the store.
func (e *MemStore) Put(_ context.Context, value []byte) ([]byte, error) {
	encodedVal, err := e.codec.EncodeBlob(value)
	if err != nil {
		return nil, err
	}

	cert, err := e.generateRandomCert(encodedVal)
	if err != nil {
		return nil, err
	}

	certBytes, err := rlp.EncodeToBytes(cert)
	if err != nil {
		return nil, err
	}

	certKey := crypto.Keccak256Hash(certBytes).Bytes()

	err = e.InsertEntry(certKey, encodedVal)
	if err != nil {
		return nil, err
	}

	return certBytes, nil
}

func (e *MemStore) Verify(_ context.Context, _, _ []byte) error {
	return nil
}

func (e *MemStore) BackendType() common.BackendType {
	return common.MemstoreV1BackendType
}
