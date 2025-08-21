package metadata

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// A utility for providing a reference block number (RBN) for the creation of new batches. Ensures that the reference
// block number never goes backwards, regardless of whatever the chain is doing. (Note that this invariant is not
// guaranteed after the software is restarted.)
//
// This utility is not thread safe.
type ReferenceBlockProvider interface {
	// GetReferenceBlockNumber returns a reference block number, based on the current chain height and the
	// configured offset. Value returned will only go forwards, never backwards.
	GetReferenceBlockNumber(ctx context.Context) (uint64, error)
}

var _ ReferenceBlockProvider = (*referenceBlockProvider)(nil)

// A standard implementation of the ReferenceBlockProvider interface.
type referenceBlockProvider struct {
	logger logging.Logger

	// The handle for interacting with the blockchain.
	contractBackend bind.ContractBackend

	// The offset to use when calculating the reference block number. This is the number of blocks in the past
	// that we want to use as the reference block number. This is a hedge against forking.
	offset uint64

	// Used to prevent the reference block number from going backwards.
	previousReferenceBlockNumber uint64
}

// NewReferenceBlockProvider creates a new ReferenceBlockProvider instance.
func NewReferenceBlockProvider(
	logger logging.Logger,
	contractBackend bind.ContractBackend,
	offset uint64,
) ReferenceBlockProvider {

	return &referenceBlockProvider{
		logger:          logger,
		contractBackend: contractBackend,
		offset:          offset,
	}
}

func (r *referenceBlockProvider) GetReferenceBlockNumber(ctx context.Context) (uint64, error) {
	latestHeader, err := r.contractBackend.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block header: %w", err)
	}
	latestBlockNumber := latestHeader.Number.Uint64()

	if latestBlockNumber < r.offset {
		return 0, fmt.Errorf("latest block number is less than RBN offset: %d < %d",
			latestBlockNumber, r.offset)
	}

	newReferenceBlockNumber := latestBlockNumber - r.offset

	if newReferenceBlockNumber < r.previousReferenceBlockNumber {
		r.logger.Warnf("Reference block number is going backwards: %d < %d... was there a fork? "+
			"Using previous value %d instead.",
			newReferenceBlockNumber,
			r.previousReferenceBlockNumber,
			r.previousReferenceBlockNumber)

		return r.previousReferenceBlockNumber, nil
	}

	r.previousReferenceBlockNumber = newReferenceBlockNumber
	return newReferenceBlockNumber, nil
}
