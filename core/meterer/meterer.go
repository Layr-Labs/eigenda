package meterer

import (
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Config contains network parameters that should be published on-chain. We currently configure these params through disperser env vars.
type Config struct {
	// 	GlobalBytesPerSecond is the rate limit in bytes per second for on-demand payments
	GlobalBytesPerSecond uint64
	// MinChargeableSize is the minimum size of a chargeable unit in bytes, used as a floor for on-demand payments
	MinChargeableSize uint32
	// PricePerChargeable is the price per chargeable unit in gwei, used for on-demand payments
	PricePerChargeable uint32
	// ReservationWindow is the duration of all reservations in seconds, used to calculate bin indices
	ReservationWindow uint32

	// ChainReadTimeout is the timeout for reading payment state from chain
	ChainReadTimeout time.Duration
}

// Meterer handles payment accounting across different accounts. Disperser API server receives requests from clients and each request contains a blob header
// with payments information (CumulativePayments, BinIndex, and Signature). Disperser will pass the blob header to the meterer, which will check if the
// payments information is valid.
type Meterer struct {
	Config

	// ChainState reads on-chain payment state periodically and cache it in memory
	ChainState OnchainPayment
	// OffchainStore uses DynamoDB to track metering and used to validate requests
	OffchainStore OffchainStore

	logger logging.Logger
}

func NewMeterer(
	config Config,
	paymentChainState OnchainPayment,
	offchainStore OffchainStore,
	logger logging.Logger,
) (*Meterer, error) {
	// TODO: create a separate thread to pull from the chain and update chain state
	return &Meterer{
		Config: config,

		ChainState:    paymentChainState,
		OffchainStore: offchainStore,

		logger: logger.With("component", "Meterer"),
	}, nil
}
