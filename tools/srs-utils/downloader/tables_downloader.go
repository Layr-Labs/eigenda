package downloader

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultTablesBaseURL  = "https://srs-mainnet.s3.amazonaws.com/kzg/SRSTables"
	defaultTablesOutputDir = "resources/srs/SRSTables"
	defaultDimension      = "dimE8192"
)

// TablesDownloaderConfig holds configuration for SRS table files download
type TablesDownloaderConfig struct {
	dimension  string
	outputDir  string
	baseURL    string
	cosetSizes []int
}

// NewTablesDownloaderConfig creates a new config with the specified parameters,
// applies defaults to empty fields, and validates the configuration
func NewTablesDownloaderConfig(
	dimension string,
	outputDir string,
	baseURL string,
	cosetSizes []int,
) (TablesDownloaderConfig, error) {
	// Apply defaults
	if dimension == "" {
		dimension = defaultDimension
	}
	if outputDir == "" {
		outputDir = defaultTablesOutputDir
	}
	if baseURL == "" {
		baseURL = defaultTablesBaseURL
	}
	if len(cosetSizes) == 0 {
		cosetSizes = []int{4, 8, 16, 32, 64, 128, 256, 512, 1024}
	}

	return TablesDownloaderConfig{
		dimension:  dimension,
		outputDir:  outputDir,
		baseURL:    baseURL,
		cosetSizes: cosetSizes,
	}, nil
}

// DownloadSRSTables implements the CLI command for downloading SRS table files
func DownloadSRSTables(config TablesDownloaderConfig) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	fmt.Printf("Downloading SRS tables for %s from %s/\n", config.dimension, config.baseURL)
	fmt.Println("Checking server availability and file sizes...")

	totalBytes := uint64(0)
	downloadedFiles := 0

	for _, cosetSize := range config.cosetSizes {
		fileName := fmt.Sprintf("%s.coset%d", config.dimension, cosetSize)
		fileURL, err := constructURLPath(config.baseURL, fileName)
		if err != nil {
			return fmt.Errorf("construct URL for %s: %w", fileName, err)
		}

		// Get file size
		fileSize, err := getRemoteFileSize(fileURL)
		if err != nil {
			fmt.Printf("Warning: Could not get size for %s: %v (skipping)\n", fileName, err)
			continue
		}

		fmt.Printf("Downloading %s (%d MB)...\n", fileName, fileSize/(1024*1024))

		outputPath := filepath.Join(config.outputDir, fileName)
		if err := downloadFile(fileURL, outputPath, 0, fileSize-1); err != nil {
			return fmt.Errorf("download %s: %w", fileName, err)
		}

		totalBytes += fileSize
		downloadedFiles++
		fmt.Printf("  Downloaded %s\n", fileName)
	}

	if downloadedFiles == 0 {
		return fmt.Errorf("no files were downloaded")
	}

	fmt.Printf("\nSuccessfully downloaded %d files (%s) to %s\n",
		downloadedFiles, formatBytes(totalBytes), config.outputDir)

	return nil
}

// formatBytes converts bytes to a human-readable string
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}