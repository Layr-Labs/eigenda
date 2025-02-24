package memstore

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	eigenda_common "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	DefaultPruneInterval = 500 * time.Millisecond
	BytesPerFieldElement = 32
)

/*
MemStore is a simple in-memory store for blobs which uses an expiration
time to evict blobs to best emulate the ephemeral nature of blobs dispersed to
EigenDA operators.
*/
type MemStore struct {
	// We use a SafeConfig because it is shared with the MemStore api
	// which can change its values concurrently.
	config *memconfig.SafeConfig
	log    logging.Logger

	// mu protects keyStarts and store (and verifier??)
	mu        sync.RWMutex
	keyStarts map[string]time.Time
	store     map[string][]byte
	// We only use the verifier for kzgCommitment verification.
	// MemStore generates random certs which can't be verified.
	// TODO: we should probably refactor the Verifier to be able to only take in a BlobVerifier here.
	verifier *verify.Verifier
	codec    codecs.BlobCodec

	reads int
}

var _ common.GeneratedKeyStore = (*MemStore)(nil)

// New ... constructor
func New(
	ctx context.Context, verifier *verify.Verifier, log logging.Logger, config *memconfig.SafeConfig,
) (*MemStore, error) {
	store := &MemStore{
		log:       log,
		config:    config,
		keyStarts: make(map[string]time.Time),
		store:     make(map[string][]byte),
		verifier:  verifier,
		codec:     codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}

	if store.config.BlobExpiration() != 0 {
		log.Info("memstore expiration enabled", "time", store.config.BlobExpiration)
		go store.pruningLoop(ctx)
	}

	return store, nil
}

// pruningLoop ... runs a background goroutine to prune expired blobs from the store on a regular interval.
func (e *MemStore) pruningLoop(ctx context.Context) {
	timer := time.NewTicker(DefaultPruneInterval)

	for {
		select {
		case <-ctx.Done():
			return

		case <-timer.C:
			e.pruneExpired()
		}
	}
}

// pruneExpired ... removes expired blobs from the store based on the expiration time.
func (e *MemStore) pruneExpired() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for commit, dur := range e.keyStarts {
		if time.Since(dur) >= e.config.BlobExpiration() {
			delete(e.keyStarts, commit)
			delete(e.store, commit)

			e.log.Debug("blob pruned", "commit", commit)
		}
	}
}

// Get fetches a value from the store.
func (e *MemStore) Get(_ context.Context, commit []byte) ([]byte, error) {
	time.Sleep(e.config.LatencyGETRoute())
	e.reads++
	e.mu.RLock()
	defer e.mu.RUnlock()

	var encodedBlob []byte
	var exists bool
	if encodedBlob, exists = e.store[crypto.Keccak256Hash(commit).String()]; !exists {
		return nil, fmt.Errorf("commitment key not found")
	}

	return e.codec.DecodeBlob(encodedBlob)
}

// Put inserts a value into the store.
func (e *MemStore) Put(_ context.Context, value []byte) ([]byte, error) {
	time.Sleep(e.config.LatencyPUTRoute())
	if e.config.PutReturnsFailoverError() {
		return nil, api.NewErrorFailover(errors.New("memstore in failover simulation mode"))
	}
	encodedVal, err := e.codec.EncodeBlob(value)
	if err != nil {
		return nil, err
	}

	if uint64(len(encodedVal)) > e.config.MaxBlobSizeBytes() {
		return nil, fmt.Errorf("%w: blob length %d, max blob size %d", common.ErrProxyOversizedBlob, len(value), e.config.MaxBlobSizeBytes())
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	commitment, err := e.verifier.Commit(encodedVal)
	if err != nil {
		return nil, err
	}

	// generate batch header hash
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
			DataLength: uint32((len(encodedVal) + BytesPerFieldElement - 1) / BytesPerFieldElement), // #nosec G115
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
				Fee:                     []byte{},
				ConfirmationBlockNumber: num,
				BatchHeaderHash:         []byte{},
			},
			BatchId:        69,
			BlobIndex:      420,
			InclusionProof: entropy,
			QuorumIndexes:  []byte{0x1, 0x0},
		},
	}

	certBytes, err := rlp.EncodeToBytes(cert)
	if err != nil {
		return nil, err
	}

	certKey := crypto.Keccak256Hash(certBytes).String()

	// construct key
	if _, exists := e.store[certKey]; exists {
		return nil, fmt.Errorf("commitment key already exists")
	}

	e.store[certKey] = encodedVal
	// add expiration
	e.keyStarts[certKey] = time.Now()

	return certBytes, nil
}

func (e *MemStore) Verify(_ context.Context, _, _ []byte) error {
	return nil
}

func (e *MemStore) BackendType() common.BackendType {
	return common.MemoryBackendType
}
