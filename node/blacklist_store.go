package node

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/leveldb"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type BlacklistStore struct {
	db     kvstore.Store[[]byte]
	logger logging.Logger
}

// NewLevelDBBlacklistStore creates a new Store object with a levelDB at the provided path and the given logger specifically for the blacklist.
func NewLevelDBBlacklistStore(
	path string,
	logger logging.Logger,
	disableSeeksCompaction bool,
	syncWrites bool) (*BlacklistStore, error) {

	// Create the DB at the path.
	db, err := leveldb.NewStore(logger, path, disableSeeksCompaction, syncWrites, nil)
	if err != nil {
		logger.Error("Could not create leveldb database", "err", err)
		return nil, fmt.Errorf("failed to create leveldb database: %w", err)
	}

	return &BlacklistStore{
		db:     db,
		logger: logger,
	}, nil
}

// HasKey checks if a key exists in the store
func (s *BlacklistStore) HasDisperserID(ctx context.Context, disperserId uint32) bool {
	// hash the disperserId and look up
	disperserIdHash := sha256.Sum256([]byte(fmt.Sprintf("%d", disperserId)))
	return s.HasKey(ctx, disperserIdHash[:])
}

// HasKey checks if a key exists in the store
func (s *BlacklistStore) HasKey(ctx context.Context, key []byte) bool {
	_, err := s.db.Get(key)
	return err == nil
}

// Get retrieves a blacklist by key
func (s *BlacklistStore) GetByDisperserID(ctx context.Context, disperserId uint32) (*Blacklist, error) {
	disperserIdHash := sha256.Sum256([]byte(fmt.Sprintf("%d", disperserId)))
	return s.Get(ctx, disperserIdHash[:])
}

// Get retrieves a blacklist by key
func (s *BlacklistStore) Get(ctx context.Context, key []byte) (*Blacklist, error) {
	rawBlackList, err := s.db.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get blacklist data: %w", err)
	}

	blacklist := new(Blacklist)
	if err = blacklist.FromBytes(rawBlackList); err != nil {
		return nil, fmt.Errorf("failed to parse blacklist data: %w", err)
	}

	return blacklist, nil
}

// Put stores raw blacklist data
func (s *BlacklistStore) Put(ctx context.Context, key []byte, value []byte) error {
	return s.db.Put(key, value)
}

// AddEntry adds or updates a blacklist entry for a disperser
func (s *BlacklistStore) AddEntry(ctx context.Context, disperserId uint32, contextReason, reason string) error {
	var blacklist *Blacklist
	var err error

	if s.HasDisperserID(ctx, disperserId) {
		blacklist, err = s.GetByDisperserID(ctx, disperserId)
		if err != nil {
			return fmt.Errorf("failed to get existing blacklist: %w", err)
		}
	} else {
		blacklist = new(Blacklist)
	}

	blacklist.AddEntry(disperserId, contextReason, reason)

	data, err := blacklist.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to serialize blacklist: %w", err)
	}

	s.logger.Info("Adding entry to blacklist", "disperserId", disperserId, "contextReason", contextReason, "reason", reason)
	disperserIdHash := sha256.Sum256([]byte(fmt.Sprintf("%d", disperserId)))
	return s.db.Put(disperserIdHash[:], data)
}

// IsBlacklisted checks if a disperser is blacklisted
func (s *BlacklistStore) IsBlacklisted(ctx context.Context, disperserId uint32) bool {
	blacklist, err := s.GetByDisperserID(ctx, disperserId)
	if err != nil {
		return false
	}

	// The number of entries determines the number of times a disperser has offended
	// We exponentially increase the amount of time the disperser is blacklisted
	// So it is 1 hour, 1 day and 1 week for the 3 offences and is based on the LastUpdated timestamp

	lastUpdated := blacklist.LastUpdated
	if len(blacklist.Entries) == 1 {
		if time.Since(time.Unix(int64(lastUpdated), 0)) < time.Hour {
			return true
		}
	} else if len(blacklist.Entries) == 2 {
		if time.Since(time.Unix(int64(lastUpdated), 0)) < time.Hour*24 {
			return true
		}
	} else if len(blacklist.Entries) >= 3 {
		if time.Since(time.Unix(int64(lastUpdated), 0)) < time.Hour*24*7 {
			return true
		}
	}

	// is blacklisted if number of entries is >= 3
	return len(blacklist.Entries) >= 3
}
