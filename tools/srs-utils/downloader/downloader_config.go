package downloader

import "fmt"

const (
	// maxSizeBytes is the maximum allowed SRS file size (16GB)
	maxSizeBytes = 16 * 1024 * 1024 * 1024
	// minSizeBytes is the minimum allowed SRS file size (32 bytes)
	minSizeBytes = 32
	// defaultBaseURL is the default URL for SRS files
	defaultBaseURL = "https://srs-mainnet.s3.amazonaws.com/kzg"
	// defaultOutputDir is the default directory for downloaded files
	defaultOutputDir = "srs-files"
)

// DownloaderConfig holds configuration for the SRS file download
type DownloaderConfig struct {
	blobSizeBytes uint64
	outputDir     string
	baseURL       string
}

// NewDownloaderConfig creates a new config with the specified parameters,
// applies defaults to empty fields, and validates the configuration
func NewDownloaderConfig(
	blobSizeBytes uint64,
	outputDir string,
	baseURL string,
) (DownloaderConfig, error) {
	// Apply defaults
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if outputDir == "" {
		outputDir = defaultOutputDir
	}

	if blobSizeBytes < minSizeBytes {
		return DownloaderConfig{}, fmt.Errorf("blob size must be at least %d bytes", minSizeBytes)
	}
	if blobSizeBytes > maxSizeBytes {
		return DownloaderConfig{}, fmt.Errorf("blob size must be less than %d bytes (16GB)", maxSizeBytes)
	}

	return DownloaderConfig{
		blobSizeBytes: blobSizeBytes,
		outputDir:     outputDir,
		baseURL:       baseURL,
	}, nil
}
