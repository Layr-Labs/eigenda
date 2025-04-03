package meterer

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/leveldb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

const MinNumBins int32 = 3

// Key prefixes for different tables
const (
	reservationPrefix = "reservation:"
	onDemandPrefix    = "ondemand:"
	globalBinPrefix   = "globalbin:"
)

type OffchainStore struct {
	db     kvstore.Store[[]byte]
	logger logging.Logger
}

func NewOffchainStore(
	path string,
	logger logging.Logger,
) (OffchainStore, error) {
	db, err := leveldb.NewStore(logger, path, false, true, nil)
	if err != nil {
		return OffchainStore{}, fmt.Errorf("failed to create leveldb store: %w", err)
	}

	return OffchainStore{
		db:     db,
		logger: logger,
	}, nil
}

// buildReservationKey builds a key for the reservation table
func buildReservationKey(accountID gethcommon.Address, reservationPeriod uint64) []byte {
	key := make([]byte, 0, len(reservationPrefix)+20+8)
	key = append(key, []byte(reservationPrefix)...)
	key = append(key, accountID.Bytes()...)
	periodBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(periodBytes, reservationPeriod)
	key = append(key, periodBytes...)
	return key
}

// buildOnDemandKey builds a key for the on-demand table
func buildOnDemandKey(accountID gethcommon.Address) []byte {
	key := make([]byte, 0, len(onDemandPrefix)+20)
	key = append(key, []byte(onDemandPrefix)...)
	key = append(key, accountID.Bytes()...)
	return key
}

// buildGlobalBinKey builds a key for the global bin table
func buildGlobalBinKey(reservationPeriod uint64) []byte {
	key := make([]byte, 0, len(globalBinPrefix)+8)
	key = append(key, []byte(globalBinPrefix)...)
	periodBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(periodBytes, reservationPeriod)
	key = append(key, periodBytes...)
	return key
}

func (s *OffchainStore) UpdateReservationBin(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64, size uint64) (uint64, error) {
	key := buildReservationKey(accountID, reservationPeriod)

	// Get current value
	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			// If not found, create new entry
			value = make([]byte, 8)
			binary.BigEndian.PutUint64(value, size)
			if err := s.db.Put(key, value); err != nil {
				return 0, fmt.Errorf("failed to create new reservation bin: %w", err)
			}
			return size, nil
		}
		return 0, fmt.Errorf("failed to get reservation bin: %w", err)
	}

	// Update existing value
	currentSize := binary.BigEndian.Uint64(value)
	newSize := currentSize + size
	newValue := make([]byte, 8)
	binary.BigEndian.PutUint64(newValue, newSize)

	if err := s.db.Put(key, newValue); err != nil {
		return 0, fmt.Errorf("failed to update reservation bin: %w", err)
	}

	return newSize, nil
}

func (s *OffchainStore) UpdateGlobalBin(ctx context.Context, reservationPeriod uint64, size uint64) (uint64, error) {
	key := buildGlobalBinKey(reservationPeriod)

	// Get current value
	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			// If not found, create new entry
			value = make([]byte, 8)
			binary.BigEndian.PutUint64(value, size)
			if err := s.db.Put(key, value); err != nil {
				return 0, fmt.Errorf("failed to create new global bin: %w", err)
			}
			return size, nil
		}
		return 0, fmt.Errorf("failed to get global bin: %w", err)
	}

	// Update existing value
	currentSize := binary.BigEndian.Uint64(value)
	newSize := currentSize + size
	if newSize < currentSize {
		return 0, fmt.Errorf("global bin usage overflows")
	}
	newValue := make([]byte, 8)
	binary.BigEndian.PutUint64(newValue, newSize)

	if err := s.db.Put(key, newValue); err != nil {
		return 0, fmt.Errorf("failed to update global bin: %w", err)
	}

	return newSize, nil
}

func (s *OffchainStore) AddOnDemandPayment(ctx context.Context, paymentMetadata core.PaymentMetadata, paymentCharged *big.Int) (*big.Int, error) {
	key := buildOnDemandKey(paymentMetadata.AccountID)

	// Get current value
	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			// If not found, create new entry
			if err := s.db.Put(key, paymentMetadata.CumulativePayment.Bytes()); err != nil {
				return nil, fmt.Errorf("failed to create new payment: %w", err)
			}
			return big.NewInt(0), nil // Return 0 for new accounts to match original behavior
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Validate payment increment first
	oldPayment := new(big.Int).SetBytes(value)
	paymentCheckpoint := new(big.Int).Sub(paymentMetadata.CumulativePayment, paymentCharged)
	if paymentCheckpoint.Sign() < 0 {
		return nil, fmt.Errorf("payment validation failed: payment charged is greater than cumulative payment")
	}

	if oldPayment.Cmp(paymentCheckpoint) > 0 {
		return nil, fmt.Errorf("insufficient cumulative payment increment")
	}

	// Only store after validation
	if err := s.db.Put(key, paymentMetadata.CumulativePayment.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return oldPayment, nil
}

func (s *OffchainStore) RollbackOnDemandPayment(ctx context.Context, accountID gethcommon.Address, newPayment, oldPayment *big.Int) error {
	key := buildOnDemandKey(accountID)

	// Get current value
	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			// If not found, create new entry with oldPayment (which might be nil)
			if oldPayment == nil {
				oldPayment = big.NewInt(0)
			}
			if err := s.db.Put(key, oldPayment.Bytes()); err != nil {
				return fmt.Errorf("failed to create new payment: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get payment: %w", err)
	}

	currentPayment := new(big.Int).SetBytes(value)
	if currentPayment.Cmp(newPayment) != 0 {
		if s.logger != nil {
			s.logger.Debug("Skipping rollback as current payment doesn't match the expected value",
				"accountID", accountID.Hex(),
				"expectedPayment", newPayment.String())
		}
		return nil
	}

	// Update payment
	if oldPayment == nil {
		oldPayment = big.NewInt(0)
	}
	if err := s.db.Put(key, oldPayment.Bytes()); err != nil {
		return fmt.Errorf("failed to rollback payment: %w", err)
	}

	if s.logger != nil {
		s.logger.Debug("Successfully rolled back payment to previous value",
			"accountID", accountID.Hex(),
			"rolledBackFrom", newPayment.String(),
			"rolledBackTo", oldPayment.String())
	}

	return nil
}

func (s *OffchainStore) GetPeriodRecords(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64) ([MinNumBins]*pb.PeriodRecord, error) {
	records := [MinNumBins]*pb.PeriodRecord{}

	// Get records for the next MinNumBins periods
	for i := 0; i < int(MinNumBins); i++ {
		period := reservationPeriod + uint64(i)
		key := buildReservationKey(accountID, period)

		value, err := s.db.Get(key)
		if err != nil {
			if errors.Is(err, kvstore.ErrNotFound) {
				continue
			}
			return [MinNumBins]*pb.PeriodRecord{}, fmt.Errorf("failed to get period record: %w", err)
		}

		usage := binary.BigEndian.Uint64(value)
		records[i] = &pb.PeriodRecord{
			Index: uint32(period),
			Usage: usage,
		}
	}

	return records, nil
}

func (s *OffchainStore) GetLargestCumulativePayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	key := buildOnDemandKey(accountID)

	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			return big.NewInt(0), nil // Return 0 for non-existent keys to match original behavior
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return new(big.Int).SetBytes(value), nil
}

func (s *OffchainStore) GetGlobalBinUsage(ctx context.Context, reservationPeriod uint64) (uint64, error) {
	key := buildGlobalBinKey(reservationPeriod)

	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			return 0, nil // Return 0 for non-existent keys
		}
		return 0, fmt.Errorf("failed to get global bin usage: %w", err)
	}

	return binary.BigEndian.Uint64(value), nil
}
