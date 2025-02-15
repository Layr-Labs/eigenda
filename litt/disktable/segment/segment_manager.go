package segment

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"math"
	"os"
	"path"
	"sync"
	"time"
)

// SegmentManager manages a table's Segments.
type SegmentManager struct {
	// The logger for the segment manager.
	logger logging.Logger

	// The root directory for the segment manager.
	root string

	// The index of the lowest numbered segment. After initial creation, only the garbage collection
	// thread is permitted to read/write this value  for the sake of thread safety.
	lowestSegmentIndex uint32

	// The index of the highest numbered segment. All writes are applied to this segment.
	highestSegmentIndex uint32

	// All segments currently in use.
	segments map[uint32]*Segment

	// The target size for value files.
	targetFileSize uint32

	// segmentLock protects access to the segments map and highestSegmentIndex.
	// Does not protect the segments themselves.
	segmentLock sync.RWMutex
}

// NewSegmentManager creates a new SegmentManager.
func NewSegmentManager(
	logger logging.Logger,
	root string,
	targetFileSize uint32) (*SegmentManager, error) {

	manager := &SegmentManager{
		logger:         logger,
		root:           root,
		targetFileSize: targetFileSize,
	}

	err := manager.gatherSegmentFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to gather segment files: %v", err)
	}

	return manager, nil
}

// getFileIndex returns the index of the segment file. Segment files are named as <index>.<extension>.
func (s *SegmentManager) getFileIndex(fileName string) (uint32, error) {
	indexString := path.Base(fileName)
	index, err := fmt.Sscanf(indexString, "%d")
	if err != nil {
		return 0, fmt.Errorf("failed to parse index from file name %s: %v", fileName, err)
	}

	return uint32(index), nil
}

// TODO: break up this method, possibly into another file

// gatherSegmentFiles reads the segment files on disk and populates the segments map.
func (s *SegmentManager) gatherSegmentFiles() error {

	// metadata files we've found on disk, key is the file's segment index, value is the file's path
	metadataFiles := make(map[uint32]string)
	// key files we've found on disk, key is the file's segment index, value is the file's path
	keyFiles := make(map[uint32]string)
	// value files we've found on disk, key is the file's segment index, value is the file's path
	valueFiles := make(map[uint32]string)

	s.highestSegmentIndex = uint32(0)
	s.lowestSegmentIndex = uint32(math.MaxUint32)

	files, err := os.ReadDir(s.root)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", s.root, err)
	}

	// Catalogue the segment files. While we are at it, delete rogue swap files.
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		extension := path.Ext(fileName)
		filePath := path.Join(s.root, fileName)

		switch extension {
		case MetadataFileExtension:
			index, err := s.getFileIndex(fileName)
			if err != nil {
				return fmt.Errorf("failed to get index from metadata file: %v", err)
			}
			if index > s.highestSegmentIndex {
				s.highestSegmentIndex = index
			}
			if index < s.lowestSegmentIndex {
				s.lowestSegmentIndex = index
			}
			metadataFiles[index] = filePath
		case MetadataSwapExtension:
			s.logger.Warnf("Removing rogue swap file %s", filePath)
			err = os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("failed to remove swap file %s: %v", filePath, err)
			}
		case KeysFileExtension:
			index, err := s.getFileIndex(fileName)
			if err != nil {
				return fmt.Errorf("failed to get index from keys file: %v", err)
			}
			if index > s.highestSegmentIndex {
				s.highestSegmentIndex = index
			}
			if index < s.lowestSegmentIndex {
				s.lowestSegmentIndex = index
			}
			keyFiles[index] = filePath
		case ValuesFileExtension:
			index, err := s.getFileIndex(fileName)
			if err != nil {
				return fmt.Errorf("failed to get index from values file: %v", err)
			}
			if index > s.highestSegmentIndex {
				s.highestSegmentIndex = index
			}
			if index < s.lowestSegmentIndex {
				s.lowestSegmentIndex = index
			}
			valueFiles[index] = filePath
		default:
			s.logger.Debugf("Ignoring unknown file %s", filePath)
		}
	}

	// For each segment, ensure that we have all the necessary files.
	orphanSet := make(map[uint32]struct{})
	lastSegmentOrphaned := false
	firstSegmentOrphaned := false
	for i := s.lowestSegmentIndex; i <= s.highestSegmentIndex; i++ {
		if _, ok := metadataFiles[i]; !ok {
			// We are missing a metadata file.

			if i == s.highestSegmentIndex {
				// This can happen if we crash while creating a new segment. Recoverable.
				s.logger.Warnf("Missing metadata file for last segment %d", i)
				orphanSet[i] = struct{}{}
				lastSegmentOrphaned = true
			} else if i == s.lowestSegmentIndex {
				// This can happen when deleting the oldest segment. Recoverable.
				s.logger.Warnf("Missing metadata file for first segment %d", i)
				orphanSet[i] = struct{}{}
				firstSegmentOrphaned = true
			} else {
				// Database is missing internal files. Catastrophic failure.
				return fmt.Errorf("missing metadata file for segment %d", i)
			}
		}

		if _, ok := keyFiles[i]; !ok {
			// We are missing a key file.

			if i == s.highestSegmentIndex {
				// This can happen if we crash while creating a new segment. Recoverable.
				s.logger.Warnf("Missing key file for last segment %d", i)
				orphanSet[i] = struct{}{}
				lastSegmentOrphaned = true
			} else if i == s.lowestSegmentIndex {
				// This can happen when deleting the oldest segment. Recoverable.
				s.logger.Warnf("Missing key file for first segment %d", i)
				orphanSet[i] = struct{}{}
				firstSegmentOrphaned = true
			} else {
				// Database is missing internal files. Catastrophic failure.
				return fmt.Errorf("missing key file for segment %d", i)
			}
		}

		if _, ok := valueFiles[i]; !ok {
			// We are missing a value file.

			if i == s.highestSegmentIndex {
				// This can happen if we crash while creating a new segment. Recoverable.
				s.logger.Warnf("Missing value file for last segment %d", i)
				orphanSet[i] = struct{}{}
				lastSegmentOrphaned = true
			} else if i == s.lowestSegmentIndex {
				// This can happen when deleting the oldest segment. Recoverable.
				s.logger.Warnf("Missing value file for first segment %d", i)
				orphanSet[i] = struct{}{}
				firstSegmentOrphaned = true
			} else {
				// Database is missing internal files. Catastrophic failure.
				return fmt.Errorf("missing value file for segment %d", i)
			}
		}
	}

	// Clean up any orphaned segment files.
	for orphanIndex := range orphanSet {
		metadataPath, ok := metadataFiles[orphanIndex]
		if ok {
			s.logger.Infof("Removing orphaned metadata file %s", metadataPath)
			err = os.Remove(metadataPath)
			if err != nil {
				return fmt.Errorf("failed to remove orphaned metadata file %s: %v", metadataPath, err)
			}
		}

		keyPath, ok := keyFiles[orphanIndex]
		if ok {
			s.logger.Infof("Removing orphaned key file %s", keyPath)
			err = os.Remove(keyPath)
			if err != nil {
				return fmt.Errorf("failed to remove orphaned key file %s: %v", keyPath, err)
			}
		}

		valuePath, ok := valueFiles[orphanIndex]
		if ok {
			s.logger.Infof("Removing orphaned value file %s", valuePath)
			err = os.Remove(valuePath)
			if err != nil {
				return fmt.Errorf("failed to remove orphaned value file %s: %v", valuePath, err)
			}
		}
	}

	if lastSegmentOrphaned {
		s.highestSegmentIndex--
	}
	if firstSegmentOrphaned {
		s.lowestSegmentIndex++
	}

	// Finally, load all healthy segments.
	for i := s.lowestSegmentIndex; i <= s.highestSegmentIndex; i++ {
		segment, err := NewSegment(s.logger, i, s.root, s.targetFileSize)
		if err != nil {
			return fmt.Errorf("failed to create segment %d: %v", i, err)
		}
		s.segments[i] = segment
	}

	return nil
}

// getSegment returns the segment with the given index. Segment is reserved, and it is the caller's responsibility to
// release the reservation when done.
func (s *SegmentManager) getReservedSegment(index uint32) (*Segment, error) {
	s.segmentLock.RLock()
	defer s.segmentLock.RUnlock()

	segment, ok := s.segments[index]
	if !ok {
		return nil, fmt.Errorf("segment %d does not exist", index)
	}

	ok = segment.Reserve()
	if !ok {
		// segmented was deleted out from under us
		return nil, fmt.Errorf("segment %d was deleted", index)
	}

	return segment, nil
}

func (s *SegmentManager) getMutableSegment() (*Segment, error) {
	s.segmentLock.RLock()
	defer s.segmentLock.RUnlock()

	segment := s.segments[s.highestSegmentIndex]

	ok := segment.Reserve()
	if !ok {
		// segmented was deleted out from under us. This should never happen for the mutable segment.
		return nil, fmt.Errorf("mutable segment %d was deleted", s.highestSegmentIndex)
	}

	return segment, nil
}

// createNewSegment attempts to create a new mutable segment. If multiple goroutines call this method at the same time,
// only one will succeed in creating the new segment. This method should only be called if the last segment is full.
func (s *SegmentManager) attemptSegmentCreation(previousHighestSegmentIndex uint32) error {
	s.segmentLock.Lock()
	defer s.segmentLock.Unlock()

	if s.highestSegmentIndex != previousHighestSegmentIndex {
		// another goroutine beat us to it
		s.segmentLock.Unlock()
		return nil
	}

	// Seal the previous segment.
	// TODO can we do this without holding the lock?
	now := time.Now() // TODO use time source
	err := s.segments[s.highestSegmentIndex].Seal(now)
	if err != nil {
		return fmt.Errorf("failed to seal segment %d: %v", s.highestSegmentIndex, err)
	}

	// Create a new segment.
	newSegment, err := NewSegment(s.logger, s.highestSegmentIndex+1, s.root, s.targetFileSize)
	if err != nil {
		s.segmentLock.Unlock()
		return fmt.Errorf("failed to create new segment: %v", err)
	}
	s.highestSegmentIndex++
	s.segments[s.highestSegmentIndex] = newSegment

	return nil
}

// Write records a key-value pair in the data segment, returning the resulting address of the data.
// This method does not ensure that the key-value pair is actually written to disk, only that it is recorded
// in the data segment. Flush must be called to ensure that all data previously passed to Put is written to disk.
func (s *SegmentManager) Write(key []byte, value []byte) (Address, error) {
	for {
		segment, err := s.getMutableSegment()
		if err != nil {
			return 0, fmt.Errorf("failed to get segment: %v", err)
		}

		address, ok, err := segment.Write(key, value)
		segment.Release()
		if err != nil {
			return 0, fmt.Errorf("failed to write key-value pair: %v", err)
		}

		if ok {
			// We've successfully written the key-value pair.
			return address, nil
		}

		// The segment filled up, write did not happen. Create a new segment and try again.
		err = s.attemptSegmentCreation(segment.index)
		if err != nil {
			return 0, fmt.Errorf("failed to create new segment: %v", err)
		}
	}
}

// Read fetches the data for a key from the data segment.
func (s *SegmentManager) Read(dataAddress Address) (data []byte, err error) {
	segment, err := s.getReservedSegment(dataAddress.Index())
	if err != nil {
		return nil, fmt.Errorf("failed to get segment: %v", err)
	}
	defer segment.Release()

	data, err = segment.Read(dataAddress)

	if err != nil {
		return nil, fmt.Errorf("failed to read data: %v", err)
	}

	return data, nil
}

// Flush flushes all data to disk.
func (s *SegmentManager) Flush() error {
	s.segmentLock.RLock()
	defer s.segmentLock.RUnlock()

	err := s.segments[s.highestSegmentIndex].Flush()
	if err != nil {
		return fmt.Errorf("failed to flush mutable segment: %v", err)
	}

	return nil
}

// TODO create background thread that calls this method

// DoGarbageCollection performs garbage collection on all segments, deleting old ones as necessary.
//
// Although this method is thread safe with respect to other methods in this class, it should not
// be called concurrently with itself.
func (s *SegmentManager) DoGarbageCollection(now time.Time, ttl time.Duration) {
	s.segmentLock.RLock()
	defer s.segmentLock.RUnlock()

	for index := s.lowestSegmentIndex; index <= s.highestSegmentIndex; index++ {
		segment := s.segments[index]
		if !segment.IsSealed() {
			// We can't delete an unsealed segment.
			return
		}

		segmentAge := now.Sub(segment.GetSealTime())
		if segmentAge < ttl {
			// Segment is not old enough to be deleted.
			return
		}

		// Segment is old enough to be deleted.
		// Actual deletion will happen when the segment is released by all reservation holders.
		segment.Release()
		delete(s.segments, index)

		s.lowestSegmentIndex++
	}
}
