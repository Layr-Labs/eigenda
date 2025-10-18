package eth

import (
	"context"
	"fmt"
	"time"

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

var _ ReferenceBlockProvider = (*periodicReferenceBlockProvider)(nil)

// A ReferenceBlockProvider implementation that periodically updates the reference block number once in a while,
// but otherwise just returns the last value it saw.
type periodicReferenceBlockProvider struct {
	base ReferenceBlockProvider

	// The most recently fetched reference block number.
	currentReferenceBlockNumber uint64

	// The time between updates to the reference block number.
	updatePeriod time.Duration

	// The last time we updated the reference block number.
	lastUpdate time.Time
}

// NewPeriodicReferenceBlockProvider creates a new ReferenceBlockProvider that wraps the given base
// ReferenceBlockProvider. The returned implementation will only call the base provider once every updatePeriod, and
// will return the last value it saw in between updates.
func NewPeriodicReferenceBlockProvider(
	base ReferenceBlockProvider,
	updatePeriod time.Duration,
) (ReferenceBlockProvider, error) {

	if updatePeriod < 0 {
		return nil, fmt.Errorf("updatePeriod must be positive")
	}

	return &periodicReferenceBlockProvider{
		base:         base,
		updatePeriod: updatePeriod,
		lastUpdate:   time.Time{},
	}, nil
}

func (p *periodicReferenceBlockProvider) GetReferenceBlockNumber(ctx context.Context) (uint64, error) {
	if time.Since(p.lastUpdate) >= p.updatePeriod {
		rbn, err := p.base.GetReferenceBlockNumber(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to get reference block number: %w", err)
		}
		p.currentReferenceBlockNumber = rbn
		p.lastUpdate = time.Now()
	}
	return p.currentReferenceBlockNumber, nil
}
