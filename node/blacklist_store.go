package node

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/leveldb"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// BlacklistStore encapsulates the database for storing blacklist entries for dispersers.
type BlacklistStore interface {
	// HasDisperserID checks if a disperser ID exists in the blacklist store
	HasDisperserID(ctx context.Context, disperserId uint32) bool

	// GetByDisperserID retrieves a blacklist by disperser ID
	GetByDisperserID(ctx context.Context, disperserId uint32) (*Blacklist, error)

	// DeleteByDisperserID deletes a blacklist by disperser ID
	DeleteByDisperserID(ctx context.Context, disperserId uint32) error

	// Get retrieves a blacklist by key
	Get(ctx context.Context, key []byte) (*Blacklist, error)

	// Put stores raw blacklist data
	Put(ctx context.Context, key []byte, value []byte) error

	// AddEntry adds or updates a blacklist entry for a disperser
	AddEntry(ctx context.Context, disperserId uint32, contextId, reason string) error

	// IsBlacklisted checks if a disperser is blacklisted
	IsBlacklisted(ctx context.Context, disperserId uint32) bool

	// blacklistDisperserFromBlobCert blacklists a disperser by retrieving the disperser's ID from the request and storing it in the blacklist store
	BlacklistDisperserFromBlobCert(request *pb.StoreChunksRequest, blobCert *corev2.BlobCertificate) error
}

type blacklistStore struct {
	mu     sync.RWMutex
	db     kvstore.Store[[]byte]
	logger logging.Logger
	time   Time
}

var _ BlacklistStore = &blacklistStore{}

// NewLevelDBBlacklistStore creates a new Store object with a levelDB at the provided path and the given logger specifically for the blacklist.
func NewLevelDBBlacklistStore(
	path string,
	logger logging.Logger,
	disableSeeksCompaction bool,
	syncWrites bool,
	time Time) (BlacklistStore, error) {

	// Create the DB at the path.
	db, err := leveldb.NewStore(logger, path, disableSeeksCompaction, syncWrites, nil)
	if err != nil {
		logger.Error("Could not create leveldb database", "err", err)
		return nil, fmt.Errorf("failed to create leveldb database: %w", err)
	}

	return &blacklistStore{
		db:     db,
		logger: logger,
		time:   time,
	}, nil
}

func (s *blacklistStore) BlacklistDisperserFromBlobCert(request *pb.StoreChunksRequest, blobCert *corev2.BlobCertificate) error {
	ctx := context.Background()
	s.logger.Info("blacklisting disperser from storeChunks request due to blobCert validation failure", "disperserID", request.DisperserID)

	// Get blob key for context
	blobKey, err := blobCert.BlobHeader.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to get blob key: %w", err)
	}

	err = s.AddEntry(ctx, request.DisperserID, fmt.Sprintf("blobKey: %x", blobKey), "blobCert validation failed")
	if err != nil {
		return fmt.Errorf("failed to add entry to blacklist: %w", err)
	}
	return nil
}

// HasDisperserID checks if a disperser ID exists in the store
func (s *blacklistStore) HasDisperserID(ctx context.Context, disperserId uint32) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// hash the disperserId and look up
	disperserIdHash := sha256.Sum256([]byte(fmt.Sprintf("%d", disperserId)))
	return s.hasKey(disperserIdHash[:])
}

// HasKey checks if a key exists in the store
func (s *blacklistStore) hasKey(key []byte) bool {
	_, err := s.db.Get(key)
	return err == nil
}

// Get retrieves a blacklist by key
func (s *blacklistStore) GetByDisperserID(ctx context.Context, disperserId uint32) (*Blacklist, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	disperserIdHash := sha256.Sum256(fmt.Appendf(nil, "%d", disperserId))
	return s.get(disperserIdHash[:])
}

// Get retrieves a blacklist by key
func (s *blacklistStore) DeleteByDisperserID(ctx context.Context, disperserId uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	disperserIdHash := sha256.Sum256(fmt.Appendf(nil, "%d", disperserId))
	return s.db.Delete(disperserIdHash[:])
}

// Get retrieves a blacklist by key
func (s *blacklistStore) Get(ctx context.Context, key []byte) (*Blacklist, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.get(key)
}

// get is the internal unlocked version of Get
func (s *blacklistStore) get(key []byte) (*Blacklist, error) {
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
func (s *blacklistStore) Put(ctx context.Context, key []byte, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Put(key, value)
}

// AddEntry adds or updates a blacklist entry for a disperser
func (s *blacklistStore) AddEntry(ctx context.Context, disperserId uint32, contextId, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var blacklist *Blacklist
	var err error

	disperserIdHash := sha256.Sum256(fmt.Appendf(nil, "%d", disperserId))

	// Check if disperser exists (using internal unlocked hasKey)
	if s.hasKey(disperserIdHash[:]) {
		blacklist, err = s.get(disperserIdHash[:])
		if err != nil {
			return fmt.Errorf("failed to get existing blacklist: %w", err)
		}
	} else {
		blacklist = new(Blacklist)
	}

	blacklist.AddEntry(disperserId, contextId, reason)

	data, err := blacklist.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to serialize blacklist: %w", err)
	}

	s.logger.Info("Adding entry to blacklist", "disperserId", disperserId, "contextId", contextId, "reason", reason)
	return s.db.Put(disperserIdHash[:], data)
}

// IsBlacklisted checks if a disperser is blacklisted
func (s *blacklistStore) IsBlacklisted(ctx context.Context, disperserId uint32) bool {
	// Start with read lock for the common case (no deletion needed)
	s.mu.RLock()

	disperserIdHash := sha256.Sum256(fmt.Appendf(nil, "%d", disperserId))
	blacklist, err := s.get(disperserIdHash[:])
	if err != nil {
		s.mu.RUnlock()
		return false
	}

	// The number of entries determines the number of times a disperser has offended
	// We exponentially increase the amount of time the disperser is blacklisted
	// So it is 1 hour, 1 day and 1 week for the 3 offences and is based on the LastUpdated timestamp
	// Uses the mockable time interface
	lastUpdated := blacklist.LastUpdated
	timeSinceLastUpdate := s.time.Since(s.time.Unix(int64(lastUpdated), 0))

	// Check if currently blacklisted based on entry count and time
	isCurrentlyBlacklisted := false
	if len(blacklist.Entries) == 1 {
		isCurrentlyBlacklisted = timeSinceLastUpdate < time.Hour
	} else if len(blacklist.Entries) == 2 {
		isCurrentlyBlacklisted = timeSinceLastUpdate < time.Hour*24
	} else if len(blacklist.Entries) >= 3 {
		isCurrentlyBlacklisted = timeSinceLastUpdate < time.Hour*24*7
	}

	// If currently blacklisted, return early with read unlock
	if isCurrentlyBlacklisted {
		s.mu.RUnlock()
		return true
	}

	// Check if we need to delete expired entry (14+ days old)
	needsCleanup := timeSinceLastUpdate >= time.Hour*24*14

	if needsCleanup {
		// Upgrade to write lock for deletion
		s.mu.RUnlock()
		s.mu.Lock()
		defer s.mu.Unlock()

		// Double-check the entry still exists after lock upgrade (avoid race condition)
		blacklist, err = s.get(disperserIdHash[:])
		if err != nil {
			// Entry was already deleted by another goroutine
			return false
		}

		// Verify it still needs cleanup (time could have changed or entry could have been updated)
		if s.time.Since(s.time.Unix(int64(blacklist.LastUpdated), 0)) >= time.Hour*24*14 {
			err = s.db.Delete(disperserIdHash[:])
			if err != nil {
				s.logger.Error("failed to delete disperser from blacklist", "disperserId", disperserId, "err", err)
			}
		}
	} else {
		s.mu.RUnlock()
	}

	// if the disperser is not blacklisted, return false
	return false
}
