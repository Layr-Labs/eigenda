package encoder

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
)

type EncoderVersion uint

const (
	V1 EncoderVersion = 1
	V2 EncoderVersion = 2
)

var _ config.DocumentedConfig = (*EncoderConfig)(nil)

var _ config.VerifiableConfig = (*EncoderConfig)(nil)

// Configuration for the encoder.
type EncoderConfig struct {
	// Encoder version (1 or 2)
	EncoderVersion EncoderVersion

	// Port at which encoder listens for gRPC calls (default: 34000)
	GrpcPort string

	Aws        aws.ClientConfig
	BlobStore  blobstore.Config
	ChunkStore chunkstore.Config
	Kzg        kzg.KzgConfig
	Server     ServerConfig
	Metrics    MetricsConfig

	// LogFormat is the format of the logs: json or text
	LogFormat string
	// LogColor is a flag to enable color in the logs
	LogColor bool
	// LogLevel is the level of the logs: debug, info, warn, error
	LogLevel string
}

func (e *EncoderConfig) GetEnvVarPrefix() string {
	return "ENCODER"
}

func (e *EncoderConfig) GetName() string {
	return "Encoder"
}

func (e *EncoderConfig) GetPackagePaths() []string {
	return []string{
		"github.com/Layr-Labs/eigenda/disperser/encoder",
		"github.com/Layr-Labs/eigenda/disperser/common/blobstore",
		"github.com/Layr-Labs/eigenda/relay/chunkstore",
		"github.com/Layr-Labs/eigenda/encoding/v1/kzg",
		"github.com/Layr-Labs/eigenda/common/aws",
	}
}

// DefaultEncoderConfig returns a default configuration for the encoder.
func DefaultEncoderConfig() *EncoderConfig {
	return &EncoderConfig{
		EncoderVersion: 1,
		GrpcPort:       "34000",
		Aws: aws.ClientConfig{
			Region: "us-east-1",
		},
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
		Server: ServerConfig{
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
		Metrics: MetricsConfig{
			HTTPPort: "9100",
		},
		LogFormat: string(common.JSONLogFormat),
		LogColor:  false,
		LogLevel:  "info",
	}
}

func (c *EncoderConfig) Verify() error {
	if c.EncoderVersion != V1 && c.EncoderVersion != V2 {
		return fmt.Errorf("invalid encoder version: %d (must be 1 or 2)", c.EncoderVersion)
	}

	if c.GrpcPort == "" {
		return fmt.Errorf("invalid gRPC port: %s", c.GrpcPort)
	}

	// For V2, bucket name is required
	if c.EncoderVersion == V2 {
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
