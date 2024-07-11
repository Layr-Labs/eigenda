package server

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	DefaultPruneInterval = 500 * time.Millisecond
)

// MemStore is a simple in-memory store for blobs which uses an expiration
// time to evict blobs to best emulate the ephemeral nature of blobs dispersed to
// EigenDA operators.
type MemStore struct {
	sync.RWMutex

	l         log.Logger
	keyStarts map[string]time.Time
	store     map[string][]byte
	verifier  *verify.Verifier
	codec     codecs.BlobCodec

	maxBlobSizeBytes uint64
	blobExpiration   time.Duration
	reads            int
}

var _ Store = (*MemStore)(nil)

// NewMemStore ... constructor
func NewMemStore(ctx context.Context, verifier *verify.Verifier, l log.Logger, maxBlobSizeBytes uint64, blobExpiration time.Duration) (*MemStore, error) {
	store := &MemStore{
		l:                l,
		keyStarts:        make(map[string]time.Time),
		store:            make(map[string][]byte),
		verifier:         verifier,
		codec:            codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec()),
		maxBlobSizeBytes: maxBlobSizeBytes,
		blobExpiration:   blobExpiration,
	}

	if store.blobExpiration != 0 {
		l.Info("memstore expiration enabled", "time", store.blobExpiration)
		go store.EventLoop(ctx)
	}

	return store, nil
}

func (e *MemStore) EventLoop(ctx context.Context) {
	timer := time.NewTicker(DefaultPruneInterval)

	for {
		select {
		case <-ctx.Done():
			return

		case <-timer.C:
			e.l.Debug("pruning expired blobs")
			e.pruneExpired()
		}
	}
}

func (e *MemStore) pruneExpired() {
	e.Lock()
	defer e.Unlock()

	for commit, dur := range e.keyStarts {
		if time.Since(dur) >= e.blobExpiration {
			delete(e.keyStarts, commit)
			delete(e.store, commit)

			e.l.Info("blob pruned", "commit", commit)
		}
	}

}

// Get fetches a value from the store.
func (e *MemStore) Get(ctx context.Context, commit []byte, domain DomainType) ([]byte, error) {
	e.reads += 1
	e.RLock()
	defer e.RUnlock()

	var cert verify.Certificate
	err := rlp.DecodeBytes(commit, &cert)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DA cert to RLP format: %w", err)
	}

	var encodedBlob []byte
	var exists bool
	if encodedBlob, exists = e.store[string(commit)]; !exists {
		return nil, fmt.Errorf("commitment key not found")
	}

	// Don't need to do this really since it's a mock store
	err = e.verifier.VerifyCommitment(cert.BlobHeader.Commitment, encodedBlob)
	if err != nil {
		return nil, err
	}

	switch domain {
	case BinaryDomain:
		return e.codec.DecodeBlob(encodedBlob)
	case PolyDomain:
		return encodedBlob, nil
	default:
		return nil, fmt.Errorf("unexpected domain type: %d", domain)
	}
}

// Put inserts a value into the store.
func (e *MemStore) Put(ctx context.Context, value []byte) ([]byte, error) {
	if uint64(len(value)) > e.maxBlobSizeBytes {
		return nil, fmt.Errorf("blob is larger than max blob size: blob length %d, max blob size %d", len(value), e.maxBlobSizeBytes)
	}

	e.Lock()
	defer e.Unlock()

	encodedVal, err := e.codec.EncodeBlob(value)
	if err != nil {
		return nil, err
	}

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

	num := uint32(blockNum.Uint64())

	cert := &verify.Certificate{
		BlobHeader: &disperser.BlobHeader{
			Commitment: &common.G1Commitment{
				X: commitment.X.Marshal(),
				Y: commitment.Y.Marshal(),
			},
			DataLength: uint32(len(encodedVal)),
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

func (e *MemStore) Stats() *Stats {
	e.RLock()
	defer e.RUnlock()
	return &Stats{
		Entries: len(e.store),
		Reads:   e.reads,
	}
}
