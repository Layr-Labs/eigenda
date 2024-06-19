package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/eigenda"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	// NOTE: this will need to be updated as plasma's implementation changes
	decodingOffset = 3
)

// TODO: Add support for custom http client option
type Config struct {
	URL string
}

// ProxyClient is an interface for communicating with the EigenDA proxy server
type ProxyClient interface {
	Health() error
	GetData(ctx context.Context, cert *common.Certificate, domain common.DomainType) ([]byte, error)
	SetData(ctx context.Context, b []byte) (*common.Certificate, error)
}

// client is the implementation of ProxyClient
type client struct {
	cfg        *Config
	httpClient *http.Client
}

func New(cfg *Config) ProxyClient {
	return &client{
		cfg,
		http.DefaultClient,
	}
}

// Health indicates if server is operational; useful for event based awaits
// when integration testing
func (c *client) Health() error {
	url := c.cfg.URL + "/health"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received bad status code: %d", resp.StatusCode)
	}

	return nil
}

// GetData fetches blob data associated with a DA certificate
func (c *client) GetData(ctx context.Context, cert *common.Certificate, domain common.DomainType) ([]byte, error) {
	b, err := rlp.EncodeToBytes(cert)
	if err != nil {
		return nil, err
	}

	// encode prefix bytes
	b = eigenda.Commitment(b).Encode()

	url := fmt.Sprintf("%s/get/0x%x?domain=%s", c.cfg.URL, b, domain.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct http request: %e", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received unexpected response code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// SetData writes raw byte data to DA and returns the respective certificate
func (c *client) SetData(ctx context.Context, b []byte) (*common.Certificate, error) {
	url := fmt.Sprintf("%s/put/", c.cfg.URL)
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to store data: %v", resp.StatusCode)
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(b) < decodingOffset {
		return nil, fmt.Errorf("read certificate is of invalid length: %d", len(b))
	}

	var cert *common.Certificate
	if err = rlp.DecodeBytes(b[decodingOffset:], &cert); err != nil {
		return nil, err
	}

	return cert, err
}
