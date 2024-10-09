package meterer

import (
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Config contains network parameters that should be published on-chain. We currently configure these params through disperser env vars.
type Config struct {
	// for rate limiting 2^64 ~= 18 exabytes per second; 2^32 ~= 4GB/s
	// for payments      2^64 ~= 18M Eth;                2^32 ~= 4ETH
	GlobalBytesPerSecond uint64 // Global rate limit in bytes per second for on-demand payments
	MinChargeableSize    uint32 // Minimum size of a chargeable unit in bytes, used as a floor for on-demand payments
	PricePerChargeable   uint32 // Price per chargeable unit in gwei, used for on-demand payments
	ReservationWindow    uint32 // Duration of all reservations in seconds, used to calculate bin indices
}

// disperser API server will receive requests from clients. these requests will be with a blobHeader with payments information (CumulativePayments, BinIndex, and Signature)
// Disperser will pass the blob header to the meterer, which will check if the payments information is valid. if it is, it will be added to the meterer's state.
// To check if the payment is valid, the meterer will:
//  1. check if the signature is valid
//     (against the CumulativePayments and BinIndex fields ;
//     maybe need something else to secure against using this appraoch for reservations when rev request comes in same bin interval; say that nonce is signed over as well)
//  2. For reservations, check offchain bin state as demonstrated in pseudocode, also check onchain state before rejecting (since onchain data is pulled)
//  3. For on-demand, check against payments and the global rates, similar to the reservation case
//
// If the payment is valid, the meterer will add the blob header to its state and return a success response to the disperser API server.
// if any of the checks fail, the meterer will return a failure response to the disperser API server.
var OnDemandQuorumNumbers = []uint8{0, 1}

type Meterer struct {
	Config

	ChainState    *OnchainPaymentState
	OffchainStore *OffchainStore

	logger logging.Logger
}

func NewMeterer(
	config Config,
	paymentChainState *OnchainPaymentState,
	offchainStore *OffchainStore,
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
