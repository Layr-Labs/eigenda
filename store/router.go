package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type Router struct {
	log log.Logger
	eigenda *EigenDAStore
	mem *MemStore
	s3 *S3Store
}

func NewRouter(eigenda *EigenDAStore, mem *MemStore, s3 *S3Store, l log.Logger) (*Router, error) {
	return &Router{
		log: l,
		eigenda: eigenda,
		mem: mem,
		s3: s3,
	}, nil
}

// The general pseudo-opinionated way of processing commitments is as follows:
// 1. generic --> EigenDA
// 2. Optimism --> EigenDA || S3  (if not EigenDA) 
func (r *Router) Get(ctx context.Context, key []byte, cm commitments.CommitmentMode) ([]byte, error) {

	switch cm {
	case commitments.OptimismGeneric:

		if r.s3 == nil {
			return nil, errors.New("expected S3 backend for OP keccak256 commitment type, but none configured")
		}

		r.log.Debug("Fetching data from S3 backend")
		return r.s3.Get(ctx, key)

	case commitments.SimpleCommitmentMode, commitments.OptimismAltDA:
		if r.mem != nil {
			return r.mem.Get(ctx, key)
		}
		
		return r.eigenda.Get(ctx, key)

	default:
		return nil, errors.New("could not determine which storage backend to route to based on unknown commitment mode")

	}

}

func (r *Router) Put(ctx context.Context, cm commitments.CommitmentMode, key, value []byte) ([]byte, error) {
	switch cm {
	case commitments.OptimismGeneric:
		return r.PutWithKey(ctx, key, value)

	
	case commitments.OptimismAltDA, commitments.SimpleCommitmentMode:
		return r.PutWithoutKey(ctx, value)

	default:
		return nil, fmt.Errorf("unknown commitment mode")
	}

}

// PutWithoutKey ... 
func (r *Router) PutWithoutKey(ctx context.Context, value []byte) (key []byte, err error) {
	if r.mem != nil {
		r.log.Debug("Storing data to memstore")
		return r.mem.Put(ctx, value)
	}

	if r.eigenda != nil {
		r.log.Debug("Storing data to eigenda backend")
		return r.eigenda.Put(ctx, value)
	}

	if r.s3 != nil {
		r.log.Debug("Storing data to S3 backend")
		commitment := crypto.Keccak256(value)

		err = r.s3.Put(ctx, commitment, value)
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("no DA storage backend found")

}


// PutWithKey is only supported for S3 storage backends using OP's alt-da keccak256 commitment type 
func (r *Router) PutWithKey(ctx context.Context, key []byte, value []byte) ([]byte, error) {
	println("Storing to s3 with key")
	if r.s3 == nil {
		return nil, errors.New("S3 is disabled but is only supported for posting known commitment keys")
	}
	// key should be a hash of the preimage value
	if actualHash := crypto.Keccak256(value); !utils.EqualSlices(actualHash, key) {
		return nil, fmt.Errorf("provided key isn't the result of Keccak256(preimage); expected: %s, actual: %s", hexutil.Encode(key), crypto.Keccak256(value))
	}

	return key, r.s3.Put(ctx, key, value)
}

func (r *Router) GetMemStore() *MemStore {
	return r.mem
}

func (r *Router) GetS3Store() *S3Store {
	return r.s3
}