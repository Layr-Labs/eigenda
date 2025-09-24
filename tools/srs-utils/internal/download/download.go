package download

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/Layr-Labs/eigenda/core"
)

// ConstructURLPath creates a proper URL for SRS file downloading
func ConstructURLPath(baseURL string, filename string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path.Join(u.Path, filename)
	return u.String(), nil
}

// GetRemoteFileSize retrieves the size of a file from the server via a HEAD request
func GetRemoteFileSize(url string) (uint64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("failed to access %s: %w", url, err)
	}
	defer core.CloseLogOnError(resp.Body, "downloader: close response body", nil)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	if resp.ContentLength < 0 {
		return 0, fmt.Errorf("could not determine file size for %s", url)
	}

	return uint64(resp.ContentLength), nil
}

// DownloadFile downloads a file from the given URL
func DownloadFile(url string, outputPath string, rangeStart uint64, rangeEnd uint64) error {
	// Create parent directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", outputPath, err)
	}
	defer core.CloseLogOnError(file, file.Name(), nil)

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
	defer core.CloseLogOnError(resp.Body, "downloader: close response body", nil)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("save downloaded data: %w", err)
	}

	return nil
}

// FormatBytes converts bytes to a human-readable string
func FormatBytes(bytes uint64) string {
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
