package meterer

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"
)

type TimeoutConfig struct {
	ChainReadTimeout    time.Duration
	ChainWriteTimeout   time.Duration
	ChainStateTimeout   time.Duration
	TxnBroadcastTimeout time.Duration
}

type Config struct {
	// PullInterval             time.Duration

	// network parameters (this should be published on-chain and read through contracts)
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

	// Metrics *Metrics
	logger logging.Logger
}

func NewMeterer(
	config Config,
	timeoutConfig TimeoutConfig,
	paymentChainState *OnchainPaymentState,
	offchainStore *OffchainStore,
	logger logging.Logger,
	// metrics *Metrics,
) (*Meterer, error) {
	// TODO: create a separate thread to pull from the chain and update chain state
	return &Meterer{
		Config:        config,
		TimeoutConfig: timeoutConfig,

		ChainState:    paymentChainState,
		OffchainStore: offchainStore,
		// Metrics:       metrics,

		logger: logger.With("component", "Meterer"),
	}, nil
}

// MeterRequest validates a blob header and adds it to the meterer's state
// TODO: return error if there's a rejection (with reasoning) or internal error (should be very rare)
func (m *Meterer) MeterRequest(ctx context.Context, header BlobHeader) error {
	if err := m.ValidateSignature(ctx, header); err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}
	//TODO: everything on chain is heavily mocked, no block number enforced but we should enforce for determinism
	blockNumber := uint(0)
	// blockNumber, err := m.ChainState.GetCurrentBlockNumber()
	// if err != nil {
	// 	return fmt.Errorf("failed to get current block number: %w", err)
	// }

	// Validate against the payment method
	if header.CumulativePayment == 0 {
		reservation, err := m.ChainState.GetActiveReservationByAccount(ctx, blockNumber, header.AccountID)
		if err != nil {
			return fmt.Errorf("failed to get active reservation by account: %w", err)
		}
		if err := m.ServeReservationRequest(ctx, header, reservation); err != nil {
			return fmt.Errorf("invalid reservation: %w", err)
		}
	} else {
		onDemandPayment, err := m.ChainState.GetOnDemandPaymentByAccount(ctx, blockNumber, header.AccountID)
		if err != nil {
			return fmt.Errorf("failed to get on-demand payment by account: %w", err)
		}
		if err := m.ServeOnDemandRequest(ctx, header, onDemandPayment); err != nil {
			return fmt.Errorf("invalid on-demand request: %w", err)
		}
	}

	// // If all checks pass, add the blob header to the meterer's state
	// if err := m.addToState(ctx, header); err != nil {
	// 	return fmt.Errorf("failed to add to meterer state: %w", err)
	// }

	return nil
}

// TODO: mocked EIP712 domain, change to the real thing when available
// ValidateSignature checks if the signature is valid against all other fields in the header
// Assuming the signature is an eip712 signature
func (m *Meterer) ValidateSignature(ctx context.Context, header BlobHeader) error {
	// Create the EIP712Signer
	//TODO: update the chainID and verifyingContract
	signer := NewEIP712Signer(big.NewInt(17000), common.HexToAddress("0x1234000000000000000000000000000000000000"))

	recoveredAddress, err := signer.RecoverSender(&header)
	if err != nil {
		return fmt.Errorf("failed to recover sender: %w", err)
	}

	accountAddress := common.HexToAddress(header.AccountID)

	if recoveredAddress != accountAddress {
		return fmt.Errorf("invalid signature: recovered address %s does not match account ID %s", recoveredAddress.Hex(), accountAddress.Hex())
	}

	return nil
}

// ServeReservationRequest handles the rate limiting logic for incoming requests
func (m *Meterer) ServeReservationRequest(ctx context.Context, blobHeader BlobHeader, reservation *ActiveReservation) error {
	if err := m.ValidateQuorum(blobHeader, reservation.QuorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for reservation: %w", err)
	}
	if !m.ValidateBinIndex(blobHeader, reservation) {
		return fmt.Errorf("invalid bin index for reservation")
	}

	// Update bin usage atomically and check against reservation's data rate as the bin limit
	if err := m.IncrementBinUsage(ctx, blobHeader, reservation); err != nil {
		return fmt.Errorf("bin overflows: %w", err)
	}

	return nil
}

func (m *Meterer) ValidateQuorum(blobHeader BlobHeader, allowedQuoroms []uint8) error {
	if len(blobHeader.QuorumNumbers) == 0 {
		return fmt.Errorf("no quorum params in blob header")
	}

	// check that all the quorum ids are in ActiveReservation's
	for _, q := range blobHeader.QuorumNumbers {
		if !slices.Contains(allowedQuoroms, q) {
			// fail the entire request if there's a quorum number mismatch
			return fmt.Errorf("quorum number mismatch: %d", q)
		}
	}
	return nil
}

// ValidateBinIndex checks if the provided bin index is valid
func (m *Meterer) ValidateBinIndex(blobHeader BlobHeader, reservation *ActiveReservation) bool {
	currentBinIndex := GetCurrentBinIndex(m.ReservationWindow)
	// Valid bin indexes are either the current bin or the previous bin
	if (blobHeader.BinIndex != currentBinIndex && blobHeader.BinIndex != (currentBinIndex-1)) || (reservation.StartEpoch > blobHeader.BinIndex || blobHeader.BinIndex > reservation.EndEpoch) {
		return false
	}
	return true
}

// IncrementBinUsage increments the bin usage atomically and checks for overflow
// TODO: Bin limit should be direct write to the Store
func (m *Meterer) IncrementBinUsage(ctx context.Context, blobHeader BlobHeader, reservation *ActiveReservation) error {
	//todo: sizes use uint64?
	recordedSize := max(blobHeader.BlobSize, uint32(m.MinChargeableSize))
	newUsage, err := m.OffchainStore.UpdateReservationBin(ctx, blobHeader.AccountID, uint64(blobHeader.BinIndex), recordedSize)
	if err != nil {
		return fmt.Errorf("failed to increment bin usage: %w", err)
	}

	// metered usage stays within the bin limit
	if newUsage <= reservation.DataRate {
		return nil
	} else if newUsage-recordedSize >= reservation.DataRate {
		// metered usage before updating the size already exceeded the limit
		return fmt.Errorf("Bin has already been filled")
	}
	if newUsage <= 2*reservation.DataRate && blobHeader.BinIndex+2 <= reservation.EndEpoch {
		m.OffchainStore.UpdateReservationBin(ctx, blobHeader.AccountID, uint64(blobHeader.BinIndex+2), newUsage-reservation.DataRate)
		return nil
	}
	return fmt.Errorf("Overflow usage exceeds bin limit")
}

// GetCurrentBinIndex returns the current bin index based on time
func GetCurrentBinIndex(binInterval uint32) uint32 {
	currentTime := time.Now().Unix()
	return uint32(currentTime) / binInterval
}

//TODO: should we track some number of blobHeaders in the meterer state and expose an API? is it stored somewhere else?
// func (m *Meterer) addToState(header BlobHeader) error {
// 	return nil
// }

// ServeOnDemandRequest handles the rate limiting logic for incoming requests
func (m *Meterer) ServeOnDemandRequest(ctx context.Context, blobHeader BlobHeader, onDemandPayment *OnDemandPayment) error {
	if err := m.ValidateQuorum(blobHeader, OnDemandQuorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for On-Demand Request: %w", err)
	}
	// update blob header to use the miniumum chargeable size
	blobHeader.BlobSize = max(blobHeader.BlobSize, uint32(m.MinChargeableSize))
	err := m.OffchainStore.AddOnDemandPayment(ctx, blobHeader)
	if err != nil {
		return fmt.Errorf("failed to update cumulative payment: %w", err)
	}
	// Validate payments attached
	err = m.ValidatePayment(ctx, blobHeader, onDemandPayment)
	if err != nil {
		// No tolerance for incorrect payment amounts; no rollbacks
		return fmt.Errorf("invalid on-demand payment: %w", err)
	}

	// Update bin usage atomically and check against bin capacity
	if err := m.IncrementGlobalBinUsage(ctx, blobHeader); err != nil {
		//TODO: conditionally remove the payment based on the error type (maybe if the error is store-op related)
		m.OffchainStore.RemoveOnDemandPayment(ctx, blobHeader.AccountID, blobHeader.CumulativePayment)
		return fmt.Errorf("failed global rate limiting")
	}

	return nil
}

// ValidatePayment checks if the provided payment header is valid against the local accounting
// prevPmt is the largest  cumulative payment strictly less    than blobHeader.cumulativePayment if exists
// nextPmt is the smallest cumulative payment strictly greater than blobHeader.cumulativePayment if exists
// nextPmtBlobSize is the blobSize of corresponding to nextPmt if exists
func (m *Meterer) ValidatePayment(ctx context.Context, blobHeader BlobHeader, onDemandPayment *OnDemandPayment) error {
	if blobHeader.CumulativePayment > uint64(onDemandPayment.CumulativePayment) {
		return fmt.Errorf("request claims a cumulative payment greater than the on-chain deposit")
	}

	prevPmt, nextPmt, nextPmtBlobSize, err := m.OffchainStore.GetRelevantOnDemandRecords(ctx, blobHeader.AccountID, blobHeader.CumulativePayment) // zero if DNE
	if err != nil {
		return fmt.Errorf("failed to get relevant on-demand records: %w", err)
	}
	// the current request must increment cumulative payment by a magnitude sufficient to cover the blob size
	if prevPmt+uint64(blobHeader.BlobSize*m.Config.PricePerChargeable/m.Config.MinChargeableSize) > blobHeader.CumulativePayment {
		return fmt.Errorf("insufficient cumulative payment increment")
	}
	// the current request must not break the payment magnitude for the next payment if the two requests were delivered out-of-order
	if nextPmt != 0 && blobHeader.CumulativePayment+uint64(nextPmtBlobSize*m.Config.PricePerChargeable/m.Config.MinChargeableSize) > nextPmt {
		return fmt.Errorf("breaking cumulative payment invariants")
	}
	// check passed: blob can be safely inserted into the set of payments
	// prevPmt + blobHeader.BlobSize * m.FixedFeePerByte <= blobHeader.CumulativePayment
	//                                                   <= nextPmt - nextPmtBlobSize * m.FixedFeePerByte > nextPmt
	return nil
}

// ValidateBinIndex checks if the provided bin index is valid
func (m *Meterer) ValidateGlobalBinIndex(blobHeader BlobHeader) (uint32, error) {
	// Deterministic function: local clock -> index (1second intervals)
	currentBinIndex := uint32(time.Now().Unix())

	// Valid bin indexes are either the current bin or the previous bin (allow this second or prev sec)
	if blobHeader.BinIndex != currentBinIndex && blobHeader.BinIndex != (currentBinIndex-1) {
		return 0, fmt.Errorf("invalid bin index for on-demand request")
	}
	return currentBinIndex, nil
}

// IncrementBinUsage increments the bin usage atomically and checks for overflow
func (m *Meterer) IncrementGlobalBinUsage(ctx context.Context, blobHeader BlobHeader) error {
	globalIndex := uint64(time.Now().Unix())
	newUsage, err := m.OffchainStore.UpdateGlobalBin(ctx, globalIndex, blobHeader.BlobSize)
	if err != nil {
		return fmt.Errorf("failed to increment global bin usage: %w", err)
	}
	if newUsage > m.GlobalBytesPerSecond {
		return fmt.Errorf("global bin usage overflows")
	}
	return nil
}
