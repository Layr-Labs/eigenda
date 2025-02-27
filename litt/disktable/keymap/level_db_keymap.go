package keymap

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/syndtr/goleveldb/leveldb"
)

var _ KeyMap = &LevelDBKeyMap{}

// LevelDBKeyMap is a key map that uses LevelDB as the underlying storage.
type LevelDBKeyMap struct {
	logger logging.Logger
	db     *leveldb.DB
	path   string
	alive  atomic.Bool
}

// NewLevelDBKeyMap creates a new LevelDBKeyMap instance.
func NewLevelDBKeyMap(logger logging.Logger, path string) (*LevelDBKeyMap, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB: %w", err)
	}

	kmap := &LevelDBKeyMap{
		logger: logger,
		db:     db,
		path:   path,
	}
	kmap.alive.Store(true)

	return kmap, nil
}

func (l *LevelDBKeyMap) Put(pairs []*types.KAPair) error {
	batch := new(leveldb.Batch)
	for _, pair := range pairs {
		batch.Put(pair.Key, pair.Address.Serialize())
	}

	err := l.db.Write(batch, nil)
	if err != nil {
		return fmt.Errorf("failed to put batch to LevelDB: %w", err)
	}
	return nil
}

func (l *LevelDBKeyMap) Get(key []byte) (types.Address, bool, error) {
	addressBytes, err := l.db.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get key from LevelDB: %w", err)
	}

	address, err := types.DeserializeAddress(addressBytes)
	if err != nil {
		return 0, false, fmt.Errorf("failed to deserialize address: %w", err)
	}

	return address, true, nil
}

func (l *LevelDBKeyMap) Delete(keys []*types.KAPair) error {
	batch := new(leveldb.Batch)
	for _, key := range keys {
		batch.Delete(key.Key)
	}

	err := l.db.Write(batch, nil)
	if err != nil {
		return fmt.Errorf("failed to delete keys from LevelDB: %w", err)
	}

	return nil
}

func (l *LevelDBKeyMap) Stop() error {
	alive := l.alive.Swap(false)
	if !alive {
		return nil
	}

	err := l.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close LevelDB: %w", err)
	}
	return nil
}

func (l *LevelDBKeyMap) Destroy() error {
	err := l.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop LevelDB: %w", err)
	}

	l.logger.Info(fmt.Sprintf("deleting LevelDB key map at path: %s", l.path))
	err = os.RemoveAll(l.path)
	if err != nil {
		return err
	}
	return nil
}

func (l *LevelDBKeyMap) LoadFromSegments(
	segments map[uint32]*segment.Segment,
	lowestSegmentIndex uint32,
	highestSegmentIndex uint32) error {

	// All data is persisted on disk via levelDB, so no need to do anything here.
	return nil
}
