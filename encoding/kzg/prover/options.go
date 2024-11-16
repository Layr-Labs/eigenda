package prover

import (
	"errors"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

// ProverOption defines a function that configures a Prover
type ProverOption func(*Prover) error

// WithBackend sets the backend type for the prover
func WithBackend(backend encoding.BackendType) ProverOption {
	return func(p *Prover) error {
		p.Config.BackendType = backend
		return nil
	}
}

// WithGPU enables or disables GPU usage
func WithGPU(enable bool) ProverOption {
	return func(e *Prover) error {
		e.Config.GPUEnable = enable
		return nil
	}
}

// WithKZGConfig sets the KZG configuration
func WithKZGConfig(config *kzg.KzgConfig) ProverOption {
	return func(p *Prover) error {
		if config.SRSNumberToLoad > config.SRSOrder {
			return errors.New("SRSOrder is less than srsNumberToLoad")
		}
		p.KzgConfig = config
		return nil
	}
}

// WithRSEncoder sets a custom RS encoder
func WithRSEncoder(encoder *rs.Encoder) ProverOption {
	return func(p *Prover) error {
		p.Encoder = encoder
		return nil
	}
}

// WithLoadG2Points enables or disables G2 points loading
func WithLoadG2Points(load bool) ProverOption {
	return func(p *Prover) error {
		p.LoadG2Points = load
		return nil
	}
}

// WithPreloadEncoder enables or disables encoder preloading
func WithPreloadEncoder(preload bool) ProverOption {
	return func(p *Prover) error {
		if !preload {
			return nil
		}

		if p.KzgConfig == nil {
			return errors.New("KZG config must be set before enabling preload encoder")
		}

		// Create table dir if not exist
		err := os.MkdirAll(p.KzgConfig.CacheDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot make CacheDir: %w", err)
		}

		return p.PreloadAllEncoders()
	}
}

// WithVerbose enables or disables verbose logging
func WithVerbose(verbose bool) ProverOption {
	return func(p *Prover) error {
		p.Config.Verbose = verbose
		return nil
	}
}
