package main

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
)

var _ config.DocumentedConfig = (*RootEncoderConfig)(nil)

// The root configuration for the encoder service. This config should be discarded after parsing
// and only the sub-configs should be used. This is a safety mechanism to make it harder to
// accidentally print/log the secret config.
type RootEncoderConfig struct {
	Config *EncoderConfig
	Secret *EncoderSecretConfig
}

var _ config.VerifiableConfig = (*EncoderConfig)(nil)

// Configuration for the encoder.
type EncoderConfig struct {
	// Encoder version (1 or 2)
	EncoderVersion uint `docs:"required"`

	// Port at which encoder listens for gRPC calls
	GrpcPort string `docs:"required"`

	// Object storage configuration
	BlobStore  blobstore.Config
	ChunkStore chunkstore.Config

	// KZG configuration
	Kzg kzg.KzgConfig

	// Server configuration
	Server encoder.ServerConfig

	// Metrics configuration
	Metrics encoder.MetricsConfig

	// Logger configuration
	LogFormat string
	LogColor  bool
	LogLevel  string

	// AWS client configuration
	Aws aws.ClientConfig
}

// Create a new root encoder config with default values.
func DefaultRootEncoderConfig() *RootEncoderConfig {
	return &RootEncoderConfig{
		Config: DefaultEncoderConfig(),
		Secret: &EncoderSecretConfig{},
	}
}

func (e *RootEncoderConfig) GetEnvVarPrefix() string {
	return "ENCODER"
}

func (e *RootEncoderConfig) GetName() string {
	return "Encoder"
}

func (e *RootEncoderConfig) GetPackagePaths() []string {
	return []string{
		"github.com/Layr-Labs/eigenda/disperser/cmd/encoder",
	}
}

func (e *RootEncoderConfig) Verify() error {
	err := e.Config.Verify()
	if err != nil {
		return fmt.Errorf("invalid encoder config: %w", err)
	}
	err = e.Secret.Verify()
	if err != nil {
		return fmt.Errorf("invalid encoder secret config: %w", err)
	}
	return nil
}

var _ config.VerifiableConfig = (*EncoderSecretConfig)(nil)

// Configuration for secrets used by the encoder.
// Currently empty as AWS credentials are handled through aws.ClientConfig,
// but this structure is kept for consistency with the ejector pattern
// and potential future secret fields.
type EncoderSecretConfig struct {
}

func (c *EncoderSecretConfig) Verify() error {
	// No secrets to verify currently
	return nil
}

// DefaultEncoderConfig returns a default configuration for the encoder.
func DefaultEncoderConfig() *EncoderConfig {
	return &EncoderConfig{
		EncoderVersion: 1,
		GrpcPort:       "34000",
		BlobStore: blobstore.Config{
			Backend: blobstore.S3Backend,
		},
		ChunkStore: chunkstore.Config{
			Backend: string(blobstore.S3Backend),
		},
		Kzg: kzg.KzgConfig{
			SRSOrder:        10000,
			SRSNumberToLoad: 10000,
			NumWorker:       12,
			PreloadEncoder:  false,
			Verbose:         false,
		},
		Server: encoder.ServerConfig{
			MaxConcurrentRequestsDangerous: 16,
			RequestPoolSize:                32,
			RequestQueueSize:               32,
			EnableGnarkChunkEncoding:       false,
			PreventReencoding:              true,
			Backend:                        string(encoding.GnarkBackend),
			GPUEnable:                      false,
			PprofHttpPort:                  "6060",
			EnablePprof:                    false,
		},
		Metrics: encoder.MetricsConfig{
			HTTPPort:      "9100",
			EnableMetrics: false,
		},
		LogFormat: string(common.JSONLogFormat),
		LogColor:  false,
		Aws: aws.ClientConfig{
			Region: "us-east-1",
		},
	}
}

func (c *EncoderConfig) Verify() error {
	if c.EncoderVersion != 1 && c.EncoderVersion != 2 {
		return fmt.Errorf("invalid encoder version: %d (must be 1 or 2)", c.EncoderVersion)
	}

	if c.GrpcPort == "" {
		return fmt.Errorf("invalid gRPC port: %s", c.GrpcPort)
	}

	// For V2, bucket name is required
	if c.EncoderVersion == 2 {
		if c.BlobStore.BucketName == "" {
			return fmt.Errorf("blob store bucket name is required for encoder v2")
		}
		if c.ChunkStore.BucketName == "" {
			return fmt.Errorf("chunk store bucket name is required for encoder v2")
		}
	}

	// Verify KZG config
	if c.Kzg.G1Path == "" {
		return fmt.Errorf("G1 path is required")
	}

	if c.Kzg.SRSNumberToLoad == 0 {
		return fmt.Errorf("SRS number to load must be greater than 0")
	}

	if c.Kzg.NumWorker == 0 {
		return fmt.Errorf("number of workers must be greater than 0")
	}

	if c.Server.MaxConcurrentRequestsDangerous <= 0 {
		return fmt.Errorf("max concurrent requests must be greater than 0")
	}

	if c.Server.RequestPoolSize <= 0 {
		return fmt.Errorf("request pool size must be greater than 0")
	}

	if c.Server.RequestQueueSize <= 0 {
		return fmt.Errorf("request queue size must be greater than 0")
	}

	return nil
}
