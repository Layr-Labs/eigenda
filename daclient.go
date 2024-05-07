package plasma

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrNotFound is returned when the server could not find the input.
var ErrNotFound = errors.New("not found")

// ErrInvalidInput is returned when the input is not valid for posting to the DA storage.
var ErrInvalidInput = errors.New("invalid input")

// DAClient ...
type DAClient struct {
	url    string
	verify bool
}

func NewDAClient(url string, verify bool) *DAClient {
	return &DAClient{url, verify}
}

// GetInput returns the input data for the given encoded commitment bytes.
func (c *DAClient) GetInput(ctx context.Context, comm EigenDACommitment) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/get/0x%x", c.url, comm.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	defer resp.Body.Close()
	input, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// TODO: Implement verification
	if c.verify {
		if err := comm.Verify(input); err != nil {
			return nil, err
		}

	}
	return input, nil
}

// SetInput sets the input data and returns the EigenDA blob commitment as a byte array.
func (c *DAClient) SetInput(ctx context.Context, img []byte) (EigenDACommitment, error) {
	if len(img) == 0 {
		return nil, ErrInvalidInput
	}

	body := bytes.NewReader(img)
	url := fmt.Sprintf("%s/put/", c.url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to store data: %v", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	comm, err := DecodeEigenDACommitment(b)
	if err != nil {
		return nil, err
	}

	return comm, nil
}

func (c *DAClient) Health() bool {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/health", c.url), nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}
