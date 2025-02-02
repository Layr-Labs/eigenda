package memstore

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/verify"
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

type Config struct {
	MaxBlobSizeBytes uint64
	BlobExpiration   time.Duration
	// artificial latency added for memstore backend to mimic eigenda's latency
	PutLatency time.Duration
	GetLatency time.Duration
}

/*
MemStore is a simple in-memory store for blobs which uses an expiration
time to evict blobs to best emulate the ephemeral nature of blobs dispersed to
EigenDA operators.
*/
type MemStore struct {
	sync.RWMutex

	config    Config
	log       logging.Logger
	keyStarts map[string]time.Time
	store     map[string][]byte
	verifier  *verify.Verifier
	codec     codecs.BlobCodec

	reads int
}

var _ common.GeneratedKeyStore = (*MemStore)(nil)

// New ... constructor
func New(
	ctx context.Context, verifier *verify.Verifier, log logging.Logger, config Config,
) (*MemStore, error) {
	store := &MemStore{
		log:       log,
		config:    config,
		keyStarts: make(map[string]time.Time),
		store:     make(map[string][]byte),
		verifier:  verifier,
		codec:     codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
	}

	if store.config.BlobExpiration != 0 {
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
	e.Lock()
	defer e.Unlock()

	for commit, dur := range e.keyStarts {
		if time.Since(dur) >= e.config.BlobExpiration {
			delete(e.keyStarts, commit)
			delete(e.store, commit)

			e.log.Debug("blob pruned", "commit", commit)
		}
	}
}

// Get fetches a value from the store.
func (e *MemStore) Get(_ context.Context, commit []byte) ([]byte, error) {
	time.Sleep(e.config.GetLatency)
	e.reads++
	e.RLock()
	defer e.RUnlock()

	var cert verify.Certificate
	err := rlp.DecodeBytes(commit, &cert)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DA cert to RLP format: %w", err)
	}

	var encodedBlob []byte
	var exists bool
	if encodedBlob, exists = e.store[string(cert.BlobVerificationProof.InclusionProof)]; !exists {
		return nil, fmt.Errorf("commitment key not found")
	}

	// Don't need to do this really since it's a mock store
	err = e.verifier.VerifyCommitment(cert.BlobHeader.Commitment, encodedBlob)
	if err != nil {
		return nil, err
	}

	return e.codec.DecodeBlob(encodedBlob)
}

// Put inserts a value into the store.
func (e *MemStore) Put(_ context.Context, value []byte) ([]byte, error) {
	time.Sleep(e.config.PutLatency)
	encodedVal, err := e.codec.EncodeBlob(value)
	if err != nil {
		return nil, err
	}

	if uint64(len(encodedVal)) > e.config.MaxBlobSizeBytes {
		return nil, fmt.Errorf("%w: blob length %d, max blob size %d", common.ErrProxyOversizedBlob, len(value), e.config.MaxBlobSizeBytes)
	}

	e.Lock()
	defer e.Unlock()

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
	// construct key
	bytesKeys := cert.BlobVerificationProof.InclusionProof

	certStr := string(bytesKeys)

	if _, exists := e.store[certStr]; exists {
		return nil, fmt.Errorf("commitment key already exists")
	}

	e.store[certStr] = encodedVal
	// add expiration
	e.keyStarts[certStr] = time.Now()

	return certBytes, nil
}

func (e *MemStore) Verify(_ context.Context, _, _ []byte) error {
	return nil
}

func (e *MemStore) BackendType() common.BackendType {
	return common.MemoryBackendType
}
