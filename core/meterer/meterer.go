package meterer

import (
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

type TimeoutConfig struct {
	ChainReadTimeout    time.Duration
	ChainWriteTimeout   time.Duration
	ChainStateTimeout   time.Duration
	TxnBroadcastTimeout time.Duration
}

// network parameters (this should be published on-chain and read through contracts)
type Config struct {
	GlobalBytesPerSecond uint64 // 2^64 bytes ~= 18 exabytes per second; if we use uint32, that's ~4GB/s
	PricePerChargeable   uint32 // 2^64 gwei ~= 18M Eth; uint32 => ~4ETH
	MinChargeableSize    uint32
	ReservationWindow    uint32
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
	TimeoutConfig

	ChainState    *OnchainPaymentState
	OffchainStore *OffchainStore

	logger logging.Logger
}

func NewMeterer(
	config Config,
	timeoutConfig TimeoutConfig,
	paymentChainState *OnchainPaymentState,
	offchainStore *OffchainStore,
	logger logging.Logger,
) (*Meterer, error) {
	// TODO: create a separate thread to pull from the chain and update chain state
	return &Meterer{
		Config:        config,
		TimeoutConfig: timeoutConfig,

		ChainState:    paymentChainState,
		OffchainStore: offchainStore,

		logger: logger.With("component", "Meterer"),
	}, nil
}
