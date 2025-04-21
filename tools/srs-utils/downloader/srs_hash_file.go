package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// srsHashFile represents a file containing SRS file hashes
type srsHashFile struct {
	blobSizeBytes uint64
	generatedAt   time.Time
	srsFileInfo   []*fileHashInfo
	filePath      string
}

// fileHashInfo holds information about a file and its hash
type fileHashInfo struct {
	filename string
	hash     string
}

// newSrsHashFile creates a new srsHashFile
func newSrsHashFile(blobSizeBytes uint64, outputDir string) (*srsHashFile, error) {
	var srsFileInfo []*fileHashInfo

	fileNames := []string{g1FileName, g2FileName, g2TrailingFileName}
	for _, fileName := range fileNames {
		hashInfo, err := getFileHashInfo(outputDir, fileName)
		if err != nil {
			return nil, fmt.Errorf("get file hash info for %s: %w", fileName, err)
		}

		srsFileInfo = append(srsFileInfo, hashInfo)
	}

	return &srsHashFile{
		blobSizeBytes: blobSizeBytes,
		generatedAt:   time.Now().UTC(),
		srsFileInfo:   srsFileInfo,
		filePath:      filepath.Join(outputDir, fmt.Sprintf("srs-files-%d.sha256", blobSizeBytes)),
	}, nil
}

// save writes the srsHashFile to the specified path
func (sf *srsHashFile) save() error {
	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(sf.filePath), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	// Create the hash file
	file, err := os.Create(sf.filePath)
	if err != nil {
		return fmt.Errorf("creating hash file: %w", err)
	}
	defer file.Close()

	// Write header
	timeStr := sf.generatedAt.Format("2006-01-02 15:04:05 UTC")
	header := fmt.Sprintf(
		"# SRS files hashes for blob size %d bytes\n"+
			"# Generated on %s\n"+
			"# Format: SHA256 (filename)\n\n",
		sf.blobSizeBytes, timeStr)

	_, err = file.WriteString(header)
	if err != nil {
		return fmt.Errorf("writing header to hash file: %w", err)
	}

	// Write file hashes
	for _, fileInfo := range sf.srsFileInfo {
		_, err = fmt.Fprintf(file, "%s  %s\n", fileInfo.hash, fileInfo.filename)
		if err != nil {
			return fmt.Errorf("writing hash to file: %w", err)
		}
	}

	return nil
}

// getFileHashInfo computes SHA-256 hash of a file
func getFileHashInfo(outputDir string, fileName string) (*fileHashInfo, error) {
	filePath := filepath.Join(outputDir, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s not found", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file for hashing: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, fmt.Errorf("calculating hash: %w", err)
	}

	return &fileHashInfo{
		fileName,
		hex.EncodeToString(hasher.Sum(nil)),
	}, nil
}
