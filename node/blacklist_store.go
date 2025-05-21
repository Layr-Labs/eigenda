package node

import (
	"context"
	"fmt"

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
func (s *BlacklistStore) HasKey(ctx context.Context, key []byte) bool {
	_, err := s.db.Get(key)
	return err == nil
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
func (s *BlacklistStore) AddEntry(ctx context.Context, disperserAddr []byte, contextID, reason string) error {
	var blacklist *Blacklist
	var err error

	if s.HasKey(ctx, disperserAddr) {
		blacklist, err = s.Get(ctx, disperserAddr)
		if err != nil {
			return fmt.Errorf("failed to get existing blacklist: %w", err)
		}
	} else {
		blacklist = new(Blacklist)
	}

	blacklist.AddEntry(disperserAddr, contextID, reason)

	data, err := blacklist.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to serialize blacklist: %w", err)
	}

	return s.db.Put(disperserAddr, data)
}
