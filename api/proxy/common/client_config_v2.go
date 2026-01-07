package common

import (
	"fmt"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/dispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
)

// ClientConfigV2 contains all non-sensitive configuration to construct V2 clients
type ClientConfigV2 struct {
	DisperserClientCfg           dispersal.DisperserClientConfig
	PayloadDisperserCfg          dispersal.PayloadDisperserConfig
	RelayPayloadRetrieverCfg     payloadretrieval.RelayPayloadRetrieverConfig
	ValidatorPayloadRetrieverCfg payloadretrieval.ValidatorPayloadRetrieverConfig

	// The following fields are not needed directly by any underlying components. Rather, these are configuration
	// values required by the proxy itself.

	// Number of times to try blob dispersals:
	// - If > 0: Try N times total
	// - If < 0: Retry indefinitely until success
	// - If = 0: Not permitted
	PutTries                           int
	MaxBlobSizeBytes                   uint64
	EigenDACertVerifierOrRouterAddress string // >= V3 cert

	// Number of GRPC connections to make to each relay
	RelayConnectionPoolSize uint

	// TODO: we should create an upstream VerifyingPayloadRetrievalClient upstream
	// that would take all of the below configs, and would verify certs before retrieving,
	// and then proceed to retrieve from its list of retrievers enabled.

	// RetrieversToEnable specifies which retrievers should be enabled
	RetrieversToEnable []RetrieverType

	// EigenDADirectory address is used to get addresses for all EigenDA contracts needed.
	EigenDADirectory string

	// The EigenDA network that is being used.
	// It is optional, and when set will be used for validating that the eth-rpc chain ID matches the network.
	EigenDANetwork EigenDANetwork

	// Determines which payment mechanism to use
	ClientLedgerMode clientledger.ClientLedgerMode

	// VaultMonitorInterval is how often to check for payment vault updates
	VaultMonitorInterval time.Duration
}

// Check checks config invariants, and returns an error if there is a problem with the config struct
func (cfg *ClientConfigV2) Check() error {
	if cfg.DisperserClientCfg.GrpcUri == "" {
		return fmt.Errorf("EigenDA disperser gRPC URI is required for using EigenDA V2 backend")
	}

	if cfg.EigenDACertVerifierOrRouterAddress == "" {
		return fmt.Errorf(`immutable v3 cert verifier address or dynamic router 
		address is required for using EigenDA V2 backend`)
	}

	if cfg.MaxBlobSizeBytes == 0 {
		return fmt.Errorf("max blob size is required for using EigenDA V2 backend")
	}

	// Check if at least one retriever is enabled
	if len(cfg.RetrieversToEnable) == 0 {
		return fmt.Errorf("at least one retriever type must be enabled for using EigenDA V2 backend")
	}

	// Check that relay retriever is not the only retriever enabled
	if slices.Contains(cfg.RetrieversToEnable, RelayRetrieverType) {
		if !slices.Contains(cfg.RetrieversToEnable, ValidatorRetrieverType) {
			return fmt.Errorf("relay retriever cannot be the only retriever enabled in EigenDA V2 backend")
		}
	}

	if slices.Contains(cfg.RetrieversToEnable, ValidatorRetrieverType) {
		if cfg.EigenDADirectory == "" {
			return fmt.Errorf("EigenDA directory is required for validator retrieval in EigenDA V2 backend")
		}
	}

	if cfg.PutTries == 0 {
		return fmt.Errorf("PutTries==0 is not permitted. >0 means 'try N times', <0 means 'retry indefinitely'")
	}

	if cfg.ClientLedgerMode == "" {
		return fmt.Errorf("client ledger mode must be specified")
	}

	if cfg.VaultMonitorInterval < 0 {
		return fmt.Errorf("vault monitor interval cannot be negative")
	}

	return nil
}

// RetrieverType defines the type of payload retriever
type RetrieverType string

const (
	RelayRetrieverType     RetrieverType = "relayRetriever"
	ValidatorRetrieverType RetrieverType = "validatorRetriever"
)
