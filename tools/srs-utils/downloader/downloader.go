package downloader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Layr-Labs/eigenda/tools/srs-utils/internal/download"
)

const (
	g1FileName         = "g1.point"
	g2FileName         = "g2.point"
	g2TrailingFileName = "g2.trailing.point"
	g2PowerOf2FileName = "g2.point.powerOf2"
)

// DownloadSRSFiles implements the CLI command for downloading SRS files and generating hash file
func DownloadSRSFiles(config DownloaderConfig) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	fmt.Println("Checking server availability and file sizes...")

	g1URL, err := download.ConstructURLPath(config.baseURL, g1FileName)
	if err != nil {
		return fmt.Errorf("construct g1.point URL: %w", err)
	}
	g1TotalSize, err := download.GetRemoteFileSize(g1URL)
	if err != nil {
		return fmt.Errorf("get remote file size: %w", err)
	}
	fmt.Printf("Total remote g1.point size: %d bytes\n", g1TotalSize)

	g2URL, err := download.ConstructURLPath(config.baseURL, g2FileName)
	if err != nil {
		return fmt.Errorf("construct g2.point URL: %w", err)
	}
	g2TotalSize, err := download.GetRemoteFileSize(g2URL)
	if err != nil {
		return fmt.Errorf("get remote file size: %w", err)
	}
	fmt.Printf("Total remote g2.point size: %d bytes\n", g2TotalSize)

	// we need to read the same number of g1 bytes as the size of the blob
	g1BytesToRead := config.blobSizeBytes
	// we need the same number of g2 points, but g2 points are twice the size of g1 points
	g2BytesToRead := config.blobSizeBytes * 2

	// Validate that our request sizes are reasonable
	if g1BytesToRead > g1TotalSize {
		return fmt.Errorf("requested blob size (%d bytes) is larger than the source g1.point file (%d bytes)",
			g1BytesToRead, g1TotalSize)
	}

	if g2BytesToRead > g2TotalSize {
		return fmt.Errorf("requested blob size *2 (%d bytes) is larger than the source g2.point file (%d bytes)",
			g2BytesToRead, g2TotalSize)
	}

	fmt.Printf("Downloading g1.point (%d bytes)...\n", g1BytesToRead)
	g1FilePath := filepath.Join(config.outputDir, g1FileName)
	if err := download.DownloadFile(
		g1URL,
		g1FilePath,
		0,
		g1BytesToRead-1,
	); err != nil {
		return err
	}

	fmt.Printf("Downloading g2.point (%d bytes)...\n", g2BytesToRead)
	if err := download.DownloadFile(
		g2URL,
		filepath.Join(config.outputDir, g2FileName),
		0,
		g2BytesToRead-1,
	); err != nil {
		return err
	}

	fmt.Printf("Downloading g2.trailing.point (%d bytes from the end of g2.point)...\n", g2BytesToRead)
	if err := download.DownloadFile(
		g2URL,
		filepath.Join(config.outputDir, g2TrailingFileName),
		g2TotalSize-g2BytesToRead,
		g2TotalSize-1,
	); err != nil {
		return err
	}

	// Download g2.point.powerOf2 if requested
	if config.includePowerOf2 {
		g2PowerOf2URL, err := download.ConstructURLPath(config.baseURL, g2PowerOf2FileName)
		if err != nil {
			return fmt.Errorf("construct g2.point.powerOf2 URL: %w", err)
		}

		g2PowerOf2TotalSize, err := download.GetRemoteFileSize(g2PowerOf2URL)
		if err != nil {
			return fmt.Errorf("get remote file size for g2.point.powerOf2: %w", err)
		}
		fmt.Printf("Total remote g2.point.powerOf2 size: %d bytes\n", g2PowerOf2TotalSize)

		fmt.Printf("Downloading g2.point.powerOf2 (full file: %d bytes)...\n", g2PowerOf2TotalSize)
		if err := download.DownloadFile(
			g2PowerOf2URL,
			filepath.Join(config.outputDir, g2PowerOf2FileName),
			0,
			g2PowerOf2TotalSize-1,
		); err != nil {
			return err
		}
	}

	fmt.Println("Calculating hashes for downloaded files...")

	srsHashFile, err := newSrsHashFile(config.blobSizeBytes, config.outputDir, config.includePowerOf2)
	if err != nil {
		return fmt.Errorf("new SRS hash file: %w", err)
	}

	err = srsHashFile.save()
	if err != nil {
		return fmt.Errorf("save hash file: %w", err)
	}

	return nil
}
