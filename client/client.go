package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

var (
	// 503 error type informing rollup to failover to other DA location
	ErrServiceUnavailable = fmt.Errorf("eigenda service is temporarily unavailable")
)

type Config struct {
	URL string // EigenDA proxy REST API URL
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientOption func(c *Client)

// WithHTTPClient ... Embeds custom http client type
func WithHTTPClient(client HTTPClient) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// Client implements a standard client for the eigenda-proxy
// that can put/get standard commitment data and query the health endpoint.
// Currently it is meant to be used by Arbitrum nitro integrations but can be extended to others in the future.
// Optimism has its own client: https://github.com/ethereum-optimism/optimism/blob/develop/op-alt-da/daclient.go
// so clients wanting to send op commitment mode data should use that client.
type Client struct {
	cfg        *Config
	httpClient HTTPClient
}

// New ... constructor
func New(cfg *Config, opts ...ClientOption) *Client {
	scc := &Client{
		cfg,
		http.DefaultClient,
	}

	for _, opt := range opts {
		opt(scc)
	}

	return scc
}

// Health indicates if the server is operational; useful for event based awaits
// when integration testing
func (c *Client) Health() error {
	url := c.cfg.URL + "/health"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received bad status code: %d", resp.StatusCode)
	}

	return nil
}

// GetData fetches blob data associated with a DA certificate
func (c *Client) GetData(ctx context.Context, comm []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/get/0x%x?commitment_mode=standard", c.cfg.URL, comm)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received error response when reading from eigenda-proxy, code=%d, msg = %s", resp.StatusCode, string(b))
	}

	return b, nil
}

// SetData writes raw byte data to DA and returns the associated certificate
// which should be verified within the proxy
func (c *Client) SetData(ctx context.Context, b []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/put?commitment_mode=standard", c.cfg.URL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// failover signal
	if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, ErrServiceUnavailable
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received error response when dispersing to eigenda-proxy, code=%d, err = %s", resp.StatusCode, string(b))
	}

	if len(b) == 0 {
		return nil, fmt.Errorf("received an empty certificate")
	}

	return b, err
}
