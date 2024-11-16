package verifier

import (
	"errors"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

// VerifierOption defines a function that configures a Verifier
type VerifierOption func(*Verifier) error

// WithKZGConfig sets the KZG configuration
func WithKZGConfig(config *kzg.KzgConfig) VerifierOption {
	return func(v *Verifier) error {
		if config.SRSNumberToLoad > config.SRSOrder {
			return errors.New("SRSOrder is less than srsNumberToLoad")
		}
		v.kzgConfig = config
		return nil
	}
}

// WithRSEncoder sets a custom RS encoder
func WithRSEncoder(encoder *rs.Encoder) VerifierOption {
	return func(v *Verifier) error {
		v.Encoder = encoder
		return nil
	}
}

// WithLoadG2Points enables or disables G2 points loading
func WithLoadG2Points(load bool) VerifierOption {
	return func(v *Verifier) error {
		v.LoadG2Points = load
		return nil
	}
}

// WithVerbose enables or disables verbose logging
func WithVerbose(verbose bool) VerifierOption {
	return func(v *Verifier) error {
		v.config.Verbose = verbose
		return nil
	}
}
