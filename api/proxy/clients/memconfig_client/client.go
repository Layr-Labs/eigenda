// Package memconfig_client provides a client for interacting with the eigenda-proxy's memstore configuration API.
// It is used in tests to drive memstore behavior such as causing proxy to start returning 503 failover errors.
package memconfig_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	memConfigEndpoint = "/memstore/config"
)

type Config struct {
	URL string // EigenDA proxy REST API URL
}

// copy derivation error to avoid cyclic deps
// see full implementation at
// https://github.com/Layr-Labs/eigenda/blob/e5f489aae34a1f68eb750e0da7ded52c200d7c36/api/clients/v2/coretypes/derivation_errors.go#L20
type DerivationError struct {
	StatusCode uint8
	Msg        string
}

// MemConfig ... contains properties that are used to configure the MemStore's behavior.
// this is copied directly from /store/generated_key/memstore/memconfig.
// importing the struct isn't possible since it'd create cyclic dependency loop
// with core proxy's go.mod
type MemConfig struct {
	MaxBlobSizeBytes                 uint64
	BlobExpiration                   time.Duration
	PutLatency                       time.Duration
	GetLatency                       time.Duration
	PutReturnsFailoverError          bool
	PutWithGetReturnsDerivationError *DerivationError
}

// MarshalJSON implements custom JSON marshaling for Config.
// This is needed because time.Duration is serialized to nanoseconds,
// which is hard to read.
func (c MemConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(intermediaryCfg{
		MaxBlobSizeBytes:                 c.MaxBlobSizeBytes,
		BlobExpiration:                   c.BlobExpiration.String(),
		PutLatency:                       c.PutLatency.String(),
		GetLatency:                       c.GetLatency.String(),
		PutReturnsFailoverError:          c.PutReturnsFailoverError,
		PutWithGetReturnsDerivationError: c.PutWithGetReturnsDerivationError,
	})
}

// intermediaryCfg ... used for decoding into a less rich type before
// translating to a structured MemConfig
type intermediaryCfg struct {
	MaxBlobSizeBytes                 uint64
	BlobExpiration                   string
	PutLatency                       string
	GetLatency                       string
	PutReturnsFailoverError          bool
	PutWithGetReturnsDerivationError *DerivationError
}

// IntoMemConfig ... converts an intermediary config into a memconfig
// with structured type definitions
func (cfg *intermediaryCfg) IntoMemConfig() (*MemConfig, error) {
	getLatency, err := time.ParseDuration(cfg.GetLatency)
	if err != nil {
		return nil, fmt.Errorf("failed to parse getLatency: %w", err)
	}

	putLatency, err := time.ParseDuration(cfg.PutLatency)
	if err != nil {
		return nil, fmt.Errorf("failed to parse putLatency: %w", err)
	}

	blobExpiration, err := time.ParseDuration(cfg.BlobExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to parse blobExpiration: %w", err)
	}

	return &MemConfig{
		MaxBlobSizeBytes:                 cfg.MaxBlobSizeBytes,
		BlobExpiration:                   blobExpiration,
		PutLatency:                       putLatency,
		GetLatency:                       getLatency,
		PutReturnsFailoverError:          cfg.PutReturnsFailoverError,
		PutWithGetReturnsDerivationError: cfg.PutWithGetReturnsDerivationError,
	}, nil
}

// Client implements a standard client for the eigenda-proxy
// that can be used for updating a memstore configuration in real-time
// this is useful for API driven tests in protocol forks that leverage
// custom integrations with EigenDA
type Client struct {
	cfg        *Config
	httpClient *http.Client
}

// New ... memconfig client constructor
func New(cfg *Config) *Client {
	cfg.URL = cfg.URL + memConfigEndpoint // initialize once to avoid unnecessary ops when patch/get

	scc := &Client{
		cfg:        cfg,
		httpClient: http.DefaultClient,
	}

	return scc
}

// decodeResponseToMemCfg ... converts http response to structured MemConfig
func decodeResponseToMemCfg(resp *http.Response) (*MemConfig, error) {
	var cfg intermediaryCfg
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("could not decode response body to intermediary cfg: %w", err)
	}
	return cfg.IntoMemConfig()
}

// GetConfig retrieves the current configuration.
func (c *Client) GetConfig(ctx context.Context) (*MemConfig, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.cfg.URL, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to read config. expected status code 200, got: %d", resp.StatusCode)
	}

	return decodeResponseToMemCfg(resp)
}

// UpdateConfig updates the configuration using the new MemConfig instance
// Despite the API using a PATH method, this function treats the "update" config
// as a POST and modifies every associated field. This could present issues if
// misused in a testing framework which imports it.
func (c *Client) UpdateConfig(ctx context.Context, update *MemConfig) (*MemConfig, error) {
	fmt.Printf("update %v\n", update)
	body, err := update.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config update to json bytes: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.cfg.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("resp.Status %v\n", resp.Header)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update config, status code: %d", resp.StatusCode)
	}

	return decodeResponseToMemCfg(resp)
}
