package segment

import (
	"fmt"
	"math"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// getFileIndex returns the index of the segment file. Segment files are named as <index>.<extension>.
func getFileIndex(fileName string) (uint32, error) {
	extension := path.Ext(fileName)
	indexString := path.Base(fileName)[:len(fileName)-len(extension)]
	index, err := strconv.Atoi(indexString)
	if err != nil {
		return 0, fmt.Errorf("failed to parse index from file name %s: %v", fileName, err)
	}

	return uint32(index), nil
}

// scanDirectory scans a directory for segment files and returns a map of metadata, key, and value files.
// Also returns a list of garbage files that should be deleted. Does not do anything to files with unrecognized
// extensions.
func scanDirectory(logger logging.Logger, rootDirectory string) (
	metadataFiles map[uint32]string,
	keyFiles map[uint32]string,
	valueFiles map[uint32]string,
	garbageFiles []string,
	highestSegmentIndex uint32,
	lowestSegmentIndex uint32,
	isEmpty bool,
	err error) {

	highestSegmentIndex = uint32(0)
	lowestSegmentIndex = uint32(math.MaxUint32)

	// key is the file's segment index, value is the file's path
	metadataFiles = make(map[uint32]string)
	keyFiles = make(map[uint32]string)
	valueFiles = make(map[uint32]string)

	garbageFiles = make([]string, 0)

	files, err := os.ReadDir(rootDirectory)
	if err != nil {
		return nil, nil, nil, nil,
			0, 0, false,
			fmt.Errorf("failed to read directory %s: %v", rootDirectory, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		extension := path.Ext(fileName)
		filePath := path.Join(rootDirectory, fileName)
		var index uint32
		var targetMap map[uint32]string

		switch extension {
		case MetadataSwapExtension:
			garbageFiles = append(garbageFiles, filePath)
			continue
		case MetadataFileExtension:
			targetMap = metadataFiles
		case KeysFileExtension:
			targetMap = keyFiles
		case ValuesFileExtension:
			targetMap = valueFiles
		default:
			logger.Debugf("Ignoring unknown file %s", filePath)
			continue
		}

		index, err = getFileIndex(fileName)
		if err != nil {
			return nil, nil, nil, nil,
				0, 0, false,
				fmt.Errorf("failed to get file index: %v", err)
		}
		targetMap[index] = filePath

		if index > highestSegmentIndex {
			highestSegmentIndex = index
		}
		if index < lowestSegmentIndex {
			lowestSegmentIndex = index
		}
	}

	if lowestSegmentIndex == math.MaxUint32 {
		// No segments found, fix the index.
		lowestSegmentIndex = 0
	}

	isEmpty = len(metadataFiles) == 0 && len(keyFiles) == 0 && len(valueFiles) == 0

	return metadataFiles, keyFiles, valueFiles, garbageFiles, highestSegmentIndex, lowestSegmentIndex, isEmpty, nil
}

func checkForMissingFile(
	logger logging.Logger,
	index uint32,
	lowestFileIndex uint32,
	highestFileIndex uint32,
	files map[uint32]string,
	fileType string,
	orphanSet map[uint32]struct{}) error {

	if _, ok := files[index]; !ok {
		if index == highestFileIndex {
			// This can happen if we crash while creating a new segment. Recoverable.
			logger.Warnf("Missing %s file for last segment %d", fileType, index)
			orphanSet[index] = struct{}{}
		} else if index == lowestFileIndex {
			// This can happen when deleting the oldest segment. Recoverable.
			logger.Warnf("Missing %s file for first segment %d", fileType, index)
			orphanSet[index] = struct{}{}
		} else {
			// Database is missing internal files. Catastrophic failure.
			return fmt.Errorf("missing %s file for segment %d", fileType, index)
		}
	}

	return nil
}

// lookForMissingFiles ensures that all files that should be present are actually present. Returns an error
// if files are missing in a way that cannot be recovered. If recoverable, returns a set of segments that
// have orphaned files.
func lookForMissingFiles(
	logger logging.Logger,
	lowestSegmentIndex uint32,
	highestSegmentIndex uint32,
	metadataFiles map[uint32]string,
	keyFiles map[uint32]string,
	valueFiles map[uint32]string) (map[uint32]struct{}, error) {

	orphanSet := make(map[uint32]struct{})

	for i := lowestSegmentIndex; i <= highestSegmentIndex; i++ {
		err := checkForMissingFile(
			logger,
			i,
			lowestSegmentIndex,
			highestSegmentIndex,
			metadataFiles,
			"metadata",
			orphanSet)
		if err != nil {
			return nil, err
		}

		err = checkForMissingFile(
			logger,
			i,
			lowestSegmentIndex,
			highestSegmentIndex,
			keyFiles,
			"key",
			orphanSet)
		if err != nil {
			return nil, err
		}

		err = checkForMissingFile(
			logger,
			i,
			lowestSegmentIndex,
			highestSegmentIndex,
			valueFiles,
			"value",
			orphanSet)
		if err != nil {
			return nil, err
		}
	}

	return orphanSet, nil
}

// deleteOrphanedFiles deletes any files that are in the orphan set.
func deleteOrphanedFiles(
	logger logging.Logger,
	orphanSet map[uint32]struct{},
	metadataFiles map[uint32]string,
	keyFiles map[uint32]string,
	valueFiles map[uint32]string) error {

	for orphanIndex := range orphanSet {
		metadataPath, ok := metadataFiles[orphanIndex]
		if ok {
			logger.Infof("Removing orphaned metadata file %s", metadataPath)
			err := os.Remove(metadataPath)
			if err != nil {
				return fmt.Errorf("failed to remove orphaned metadata file %s: %v", metadataPath, err)
			}
		}

		keyPath, ok := keyFiles[orphanIndex]
		if ok {
			logger.Infof("Removing orphaned key file %s", keyPath)
			err := os.Remove(keyPath)
			if err != nil {
				return fmt.Errorf("failed to remove orphaned key file %s: %v", keyPath, err)
			}
		}

		valuePath, ok := valueFiles[orphanIndex]
		if ok {
			logger.Infof("Removing orphaned value file %s", valuePath)
			err := os.Remove(valuePath)
			if err != nil {
				return fmt.Errorf("failed to remove orphaned value file %s: %v", valuePath, err)
			}
		}
	}

	return nil
}

// linkSegments links together adjacent segments via SetNextSegment().
func linkSegments(lowestSegmentIndex uint32, highestSegmentIndex uint32, segments map[uint32]*Segment) error {
	if lowestSegmentIndex == highestSegmentIndex {
		// Only one segment, nothing to link. This is checked explicitly to avoid 0-1 underflow.
		return nil
	}

	for i := lowestSegmentIndex; i <= highestSegmentIndex-1; i++ {
		first, ok := segments[i]
		if !ok {
			return fmt.Errorf("missing segment %d", i)
		}
		second, ok := segments[i+1]
		if !ok {
			return fmt.Errorf("missing segment %d", i+1)
		}
		first.SetNextSegment(second)
	}
	return nil
}

// GatherSegmentFiles scans a directory for segment files and loads them into memory. It also deletes
// orphaned files and checks for corrupted files. It creates a new mutable segment at the end.
func GatherSegmentFiles(
	logger logging.Logger,
	rootDirectory string,
	now time.Time,
	createMutableSegment bool,
) (lowestSegmentIndex uint32, highestSegmentIndex uint32, segments map[uint32]*Segment, err error) {

	// Scan the root directory for segment files.
	metadataFiles, keyFiles, valueFiles, garbageFiles, highestSegmentIndex, lowestSegmentIndex, isEmpty, err :=
		scanDirectory(logger, rootDirectory)
	if err != nil {
		return 0, 0, nil,
			fmt.Errorf("failed to scan directory: %v", err)
	}

	segments = make(map[uint32]*Segment)

	// Delete any garbage files. Ignore files with unrecognized extensions.
	if !isEmpty {
		for _, garbageFile := range garbageFiles {
			logger.Infof("deleting file %s", garbageFile)
			err = os.Remove(garbageFile)
			if err != nil {
				return 0, 0, nil,
					fmt.Errorf("failed to remove garbage file %s: %v", garbageFile, err)
			}
		}

		// Check for missing files.
		orphanSet, err := lookForMissingFiles(
			logger,
			lowestSegmentIndex,
			highestSegmentIndex,
			metadataFiles,
			keyFiles,
			valueFiles)
		if err != nil {
			return 0, 0, nil,
				fmt.Errorf("there are one or more missing files: %v", err)
		}

		// Clean up any orphaned segment files.
		err = deleteOrphanedFiles(
			logger,
			orphanSet,
			metadataFiles,
			keyFiles,
			valueFiles)
		if err != nil {
			return 0, 0, nil,
				fmt.Errorf("failed to delete orphaned files: %v", err)
		}

		// Adjust the segment range to exclude orphaned segments.
		if _, ok := orphanSet[highestSegmentIndex]; ok {
			highestSegmentIndex--
		}
		if _, ok := orphanSet[lowestSegmentIndex]; ok {
			lowestSegmentIndex++
		}

		// Load all healthy segments.
		for i := lowestSegmentIndex; i <= highestSegmentIndex; i++ {
			segment, err := NewSegment(logger, i, rootDirectory, now, true)
			if err != nil {
				return 0, 0, nil,
					fmt.Errorf("failed to create segment %d: %v", i, err)
			}
			segments[i] = segment
		}
	}

	if createMutableSegment {
		// Create a new mutable segment at the end.
		if isEmpty {
			segment, err := NewSegment(logger, lowestSegmentIndex, rootDirectory, now, false)
			if err != nil {
				return 0, 0, nil,
					fmt.Errorf("failed to create segment %d: %v", lowestSegmentIndex, err)
			}

			segments[0] = segment
		} else {
			segment, err := NewSegment(logger, highestSegmentIndex+1, rootDirectory, now, false)
			if err != nil {
				return 0, 0, nil,
					fmt.Errorf("failed to create segment %d: %v", highestSegmentIndex+1, err)
			}

			segments[highestSegmentIndex+1] = segment
			highestSegmentIndex++
		}
	}

	// Stitch together the segments.
	err = linkSegments(lowestSegmentIndex, highestSegmentIndex, segments)
	if err != nil {
		return 0, 0, nil,
			fmt.Errorf("failed to link segments: %v", err)
	}

	return lowestSegmentIndex, highestSegmentIndex, segments, nil
}
