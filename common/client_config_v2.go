package common

import (
	"fmt"
	"slices"

	clients_v2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
)

// ClientConfigV2 contains all non-sensitive configuration to construct V2 clients
type ClientConfigV2 struct {
	DisperserClientCfg           clients_v2.DisperserClientConfig
	PayloadDisperserCfg          payloaddispersal.PayloadDisperserConfig
	RelayPayloadRetrieverCfg     payloadretrieval.RelayPayloadRetrieverConfig
	ValidatorPayloadRetrieverCfg payloadretrieval.ValidatorPayloadRetrieverConfig

	// The following fields are not needed directly by any underlying components. Rather, these are configuration
	// values required by the proxy itself.

	// Number of times to try blob dispersals:
	// - If > 0: Try N times total
	// - If < 0: Retry indefinitely until success
	// - If = 0: Not permitted
	PutTries                   int
	MaxBlobSizeBytes           uint64
	EigenDACertVerifierAddress string

	// Fields required for validator payload retrieval
	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string

	// RetrieversToEnable specifies which retrievers should be enabled
	RetrieversToEnable []RetrieverType
}

// Check checks config invariants, and returns an error if there is a problem with the config struct
func (cfg *ClientConfigV2) Check() error {
	if cfg.DisperserClientCfg.Hostname == "" {
		return fmt.Errorf("disperser hostname is required for using EigenDA V2 backend")
	}

	if cfg.DisperserClientCfg.Port == "" {
		return fmt.Errorf("disperser port is required for using EigenDA V2 backend")
	}

	if cfg.EigenDACertVerifierAddress == "" {
		return fmt.Errorf("cert verifier address is required for using EigenDA V2 backend")
	}

	if cfg.MaxBlobSizeBytes == 0 {
		return fmt.Errorf("max blob size is required for using EigenDA V2 backend")
	}

	// Check if at least one retriever is enabled
	if len(cfg.RetrieversToEnable) == 0 {
		return fmt.Errorf("at least one retriever type must be enabled for using EigenDA V2 backend")
	}

	if slices.Contains(cfg.RetrieversToEnable, ValidatorRetrieverType) {
		if cfg.BLSOperatorStateRetrieverAddr == "" {
			return fmt.Errorf(
				"BLS operator state retriever address is required for validator retrieval in EigenDA V2 backend")
		}

		if cfg.EigenDAServiceManagerAddr == "" {
			return fmt.Errorf(
				"EigenDA service manager address is required for validator retrieval in EigenDA V2 backend")
		}
	}

	if cfg.PutTries == 0 {
		return fmt.Errorf("PutTries==0 is not permitted. >0 means 'try N times', <0 means 'retry indefinitely'")
	}

	return nil
}

// RetrieverType defines the type of payload retriever
type RetrieverType string

const (
	RelayRetrieverType     RetrieverType = "relayRetriever"
	ValidatorRetrieverType RetrieverType = "validatorRetriever"
)
