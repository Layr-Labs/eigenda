package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

const (
	g1FileName         = "g1.point"
	g2FileName         = "g2.point"
	g2TrailingFileName = "g2.trailing.point"
)

// DownloadSRSFiles implements the CLI command for downloading SRS files and generating hash file
func DownloadSRSFiles(config DownloaderConfig) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	fmt.Println("Checking server availability and file sizes...")

	g1URL, err := constructURLPath(config.baseURL, g1FileName)
	if err != nil {
		return fmt.Errorf("construct g1.point URL: %w", err)
	}
	g1TotalSize, err := getRemoteFileSize(g1URL)
	if err != nil {
		return fmt.Errorf("get remote file size: %w", err)
	}
	fmt.Printf("Total remote g1.point size: %d bytes\n", g1TotalSize)

	g2URL, err := constructURLPath(config.baseURL, g2FileName)
	if err != nil {
		return fmt.Errorf("construct g2.point URL: %w", err)
	}
	g2TotalSize, err := getRemoteFileSize(g2URL)
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
	if err := downloadFile(
		g1URL,
		g1FilePath,
		0,
		g1BytesToRead-1,
	); err != nil {
		return err
	}

	fmt.Printf("Downloading g2.point (%d bytes)...\n", g2BytesToRead)
	if err := downloadFile(
		g2URL,
		filepath.Join(config.outputDir, g2FileName),
		0,
		g2BytesToRead-1,
	); err != nil {
		return err
	}

	fmt.Printf("Downloading g2.trailing.point (%d bytes from the end of g2.point)...\n", g2BytesToRead)
	if err := downloadFile(
		g2URL,
		filepath.Join(config.outputDir, g2TrailingFileName),
		g2TotalSize-g2BytesToRead,
		g2TotalSize-1,
	); err != nil {
		return err
	}

	fmt.Println("Calculating hashes for downloaded files...")

	srsHashFile, err := newSrsHashFile(config.blobSizeBytes, config.outputDir)
	if err != nil {
		return fmt.Errorf("new SRS hash file: %w", err)
	}

	err = srsHashFile.save()
	if err != nil {
		return fmt.Errorf("save hash file: %w", err)
	}

	return nil
}

// constructURLPath creates a proper URL for SRS file downloading
func constructURLPath(baseURL string, filename string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path.Join(u.Path, filename)
	return u.String(), nil
}

// getRemoteFileSize retrieves the size of a file from the server via a HEAD request
func getRemoteFileSize(url string) (uint64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("failed to access %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	if resp.ContentLength < 0 {
		return 0, fmt.Errorf("could not determine file size for %s", url)
	}

	return uint64(resp.ContentLength), nil
}

// downloadFile downloads a file from the given URL
func downloadFile(url string, outputPath string, rangeStart uint64, rangeEnd uint64) error {
	// Create parent directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", outputPath, err)
	}
	defer file.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}

	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("save downloaded data: %w", err)
	}

	return nil
}
