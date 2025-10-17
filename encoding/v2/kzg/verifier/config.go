package verifier

import (
	"github.com/Layr-Labs/eigenda/encoding/kzgconfig"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover"
)

// Config holds configuration for the V2 KZG verifier.
type Config struct {
	// Number of G1 points to be loaded from the G1 SRS file located at G1Path.
	// This number times 32 bytes will be loaded from G1Path.
	SRSNumberToLoad uint64

	// G1Path is the path to the G1 SRS file.
	G1Path string

	// NumWorker is the number of goroutines used to read and parse the G1 SRS file.
	NumWorker uint64
}

// The v2 verifier's KzgConfig is a strict subset of the prover's config,
// since it doesn't need the SRSTable information which is only used for proving.
func ConfigFromProverV2Config(v2Prover *prover.KzgConfig) *Config {
	return &Config{
		SRSNumberToLoad: v2Prover.SRSNumberToLoad,
		G1Path:          v2Prover.G1Path,
		NumWorker:       v2Prover.NumWorker,
	}
}

// ConfigFromV1KzgConfig converts a v1 KzgConfig to a v2 verifier KzgConfig.
// The V1 KzgConfig is used all over the place in multiple different structs,
// making it very hard to update, optimize, change, or remove unused fields.
// The V2 verifier has its own KzgConfig, which is a very small subset of the V1 KzgConfig.
func ConfigFromV1KzgConfig(v1 *kzgconfig.Config) *Config {
	return &Config{
		SRSNumberToLoad: v1.SRSNumberToLoad,
		G1Path:          v1.G1Path,
		NumWorker:       v1.NumWorker,
	}
}
