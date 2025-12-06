package encoder

import (
	"fmt"
	"runtime"

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
	Version EncoderVersion

	// Port at which encoder listens for gRPC calls (default: 34000)
	GrpcPort string

	Aws        aws.ClientConfig
	BlobStore  blobstore.Config
	ChunkStore chunkstore.Config
	Kzg        kzg.KzgConfig
	Server     ServerConfig

	// MetricsPort is the port that the encoder metrics server listens on.
	MetricsPort string
	// EnableMetrics enables the metrics HTTP server for prometheus metrics collection
	EnableMetrics bool

	// LogFormat is the format of the logs: json or text
	LogFormat common.LogFormat
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
		Version:  V2,
		GrpcPort: "34000",
		Aws:      *aws.DefaultClientConfig(),
		BlobStore: blobstore.Config{
			Backend: blobstore.S3Backend,
		},
		ChunkStore: chunkstore.Config{
			Backend: string(blobstore.S3Backend),
		},
		Kzg: kzg.KzgConfig{
			SRSOrder:        268435456,
			SRSNumberToLoad: 2097152,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)),
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
		MetricsPort:   "9094",
		EnableMetrics: true,
		LogFormat:     common.JSONLogFormat,
		LogColor:      false,
		LogLevel:      "debug",
	}
}

func (c *EncoderConfig) Verify() error {
	if c.Version != V1 && c.Version != V2 {
		return fmt.Errorf("invalid encoder version: %d (must be 1 or 2)", c.Version)
	}

	if c.GrpcPort == "" {
		return fmt.Errorf("invalid gRPC port: %s", c.GrpcPort)
	}

	// For V2, bucket name is required
	if c.Version == V2 {
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

	if c.LogFormat != common.JSONLogFormat && c.LogFormat != common.TextLogFormat {
		return fmt.Errorf("invalid log format: %s (must be json or text)", c.LogFormat)
	}

	return nil
}
