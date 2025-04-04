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

	// Create a batch for atomic updates
	batch := s.db.NewBatch()

	// Get current value
	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			// If not found, create new entry
			value = make([]byte, 8)
			binary.BigEndian.PutUint64(value, size)
			batch.Put(key, value)
			if err := batch.Apply(); err != nil {
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

	batch.Put(key, newValue)
	if err := batch.Apply(); err != nil {
		return 0, fmt.Errorf("failed to update reservation bin: %w", err)
	}

	return newSize, nil
}

func (s *OffchainStore) UpdateGlobalBin(ctx context.Context, reservationPeriod uint64, size uint64) (uint64, error) {
	key := buildGlobalBinKey(reservationPeriod)

	// Create a batch for atomic updates
	batch := s.db.NewBatch()

	// Get current value
	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			// If not found, create new entry
			value = make([]byte, 8)
			binary.BigEndian.PutUint64(value, size)
			batch.Put(key, value)
			if err := batch.Apply(); err != nil {
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

	batch.Put(key, newValue)
	if err := batch.Apply(); err != nil {
		return 0, fmt.Errorf("failed to update global bin: %w", err)
	}

	return newSize, nil
}

func (s *OffchainStore) AddOnDemandPayment(ctx context.Context, paymentMetadata core.PaymentMetadata, paymentCharged *big.Int) (*big.Int, error) {
	key := buildOnDemandKey(paymentMetadata.AccountID)

	// Create a batch for atomic updates
	batch := s.db.NewBatch()

	// Get current value
	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			// If not found, create new entry
			batch.Put(key, paymentMetadata.CumulativePayment.Bytes())
			if err := batch.Apply(); err != nil {
				return nil, fmt.Errorf("failed to apply batch: %w", err)
			}
			value = make([]byte, 32)
			copy(value, big.NewInt(0).Bytes())
			return big.NewInt(0), nil
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Validate payment increment first
	oldPayment := new(big.Int).SetBytes(value)
	paymentCheckpoint := new(big.Int).Sub(paymentMetadata.CumulativePayment, paymentCharged)

	if paymentCheckpoint.Sign() < 0 {
		if oldPayment.Cmp(big.NewInt(0)) == 0 {
			batch.Delete(key)
			if err := batch.Apply(); err != nil {
				return nil, fmt.Errorf("failed to apply batch: %w", err)
			}
		}
		return nil, fmt.Errorf("payment validation failed: payment charged is greater than cumulative payment")
	}

	if oldPayment.Cmp(paymentCheckpoint) > 0 {
		if oldPayment.Cmp(big.NewInt(0)) == 0 {
			batch.Delete(key)
			if err := batch.Apply(); err != nil {
				return nil, fmt.Errorf("failed to apply batch: %w", err)
			}
		}
		return nil, fmt.Errorf("insufficient cumulative payment increment")
	}

	// Only store after validation
	batch.Put(key, paymentMetadata.CumulativePayment.Bytes())

	// Apply the batch atomically
	if err := batch.Apply(); err != nil {
		return nil, fmt.Errorf("failed to apply batch: %w", err)
	}

	return oldPayment, nil
}

func (s *OffchainStore) RollbackOnDemandPayment(ctx context.Context, accountID gethcommon.Address, newPayment, oldPayment *big.Int) error {
	key := buildOnDemandKey(accountID)

	// Create a batch for atomic updates
	batch := s.db.NewBatch()

	// Get current value
	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			// If not found, create new entry with oldPayment (which might be nil)
			if oldPayment == nil {
				oldPayment = big.NewInt(0)
			}
			batch.Put(key, oldPayment.Bytes())
			if err := batch.Apply(); err != nil {
				return fmt.Errorf("failed to create new payment: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get payment: %w", err)
	}

	currentPayment := new(big.Int).SetBytes(value)

	if currentPayment.Cmp(newPayment) != 0 {
		return nil
	}

	// Update payment
	if oldPayment == nil {
		oldPayment = big.NewInt(0)
	}
	batch.Put(key, oldPayment.Bytes())
	if err := batch.Apply(); err != nil {
		return fmt.Errorf("failed to rollback payment: %w", err)
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

func (s *OffchainStore) GetReservationBin(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64) (uint64, error) {
	key := buildReservationKey(accountID, reservationPeriod)

	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			return 0, nil // Return 0 for non-existent keys
		}
		return 0, fmt.Errorf("failed to get reservation bin: %w", err)
	}

	return binary.BigEndian.Uint64(value), nil
}

func (s *OffchainStore) GetOnDemandPayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	key := buildOnDemandKey(accountID)

	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			return big.NewInt(0), nil // Return 0 for non-existent keys
		}
		return nil, fmt.Errorf("failed to get on-demand payment: %w", err)
	}

	return new(big.Int).SetBytes(value), nil
}

func (s *OffchainStore) GetGlobalBin(ctx context.Context, reservationPeriod uint64) (uint64, error) {
	key := buildGlobalBinKey(reservationPeriod)

	value, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			return 0, nil // Return 0 for non-existent keys
		}
		return 0, fmt.Errorf("failed to get global bin: %w", err)
	}

	return binary.BigEndian.Uint64(value), nil
}
