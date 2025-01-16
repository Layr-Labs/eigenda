package meterer

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Config contains network parameters that should be published on-chain. We currently configure these params through disperser env vars.
type Config struct {

	// ChainReadTimeout is the timeout for reading payment state from chain
	ChainReadTimeout time.Duration

	// UpdateInterval is the interval for refreshing the on-chain state
	UpdateInterval time.Duration
}

// Meterer handles payment accounting across different accounts. Disperser API server receives requests from clients and each request contains a blob header
// with payments information (CumulativePayments, ReservationPeriod, Salt, and Signature). Disperser will pass the blob header to the meterer, which will check if the
// payments information is valid.
type Meterer struct {
	Config
	// ChainPaymentState reads on-chain payment state periodically and cache it in memory
	ChainPaymentState OnchainPayment
	// OffchainStore uses DynamoDB to track metering and used to validate requests
	OffchainStore OffchainStore

	logger logging.Logger
}

func NewMeterer(
	config Config,
	paymentChainState OnchainPayment,
	offchainStore OffchainStore,
	logger logging.Logger,
) *Meterer {
	return &Meterer{
		Config: config,

		ChainPaymentState: paymentChainState,
		OffchainStore:     offchainStore,

		logger: logger.With("component", "Meterer"),
	}
}

// Start starts to periodically refreshing the on-chain state
func (m *Meterer) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(m.Config.UpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := m.ChainPaymentState.RefreshOnchainPaymentState(ctx); err != nil {
					m.logger.Error("Failed to refresh on-chain state", "error", err)
				}
				m.logger.Debug("Refreshed on-chain state")
			case <-ctx.Done():
				return
			}
		}
	}()
}

// MeterRequest validates a blob header and adds it to the meterer's state
// TODO: return error if there's a rejection (with reasoning) or internal error (should be very rare)
func (m *Meterer) MeterRequest(ctx context.Context, header core.PaymentMetadata, numSymbols uint, quorumNumbers []uint8) error {
	accountID := gethcommon.HexToAddress(header.AccountID)
	m.logger.Info("Validating incoming request's payment metadata", "paymentMetadata", header, "numSymbols", numSymbols, "quorumNumbers", quorumNumbers)
	// Validate against the payment method
	if header.CumulativePayment.Sign() == 0 {
		reservation, err := m.ChainPaymentState.GetReservedPaymentByAccount(ctx, accountID)
		if err != nil {
			return fmt.Errorf("failed to get active reservation by account: %w", err)
		}
		if err := m.ServeReservationRequest(ctx, header, reservation, numSymbols, quorumNumbers); err != nil {
			return fmt.Errorf("invalid reservation: %w", err)
		}
	} else {
		onDemandPayment, err := m.ChainPaymentState.GetOnDemandPaymentByAccount(ctx, accountID)
		if err != nil {
			return fmt.Errorf("failed to get on-demand payment by account: %w", err)
		}
		if err := m.ServeOnDemandRequest(ctx, header, onDemandPayment, numSymbols, quorumNumbers); err != nil {
			return fmt.Errorf("invalid on-demand request: %w", err)
		}
	}

	return nil
}

// ServeReservationRequest handles the rate limiting logic for incoming requests
func (m *Meterer) ServeReservationRequest(ctx context.Context, header core.PaymentMetadata, reservation *core.ReservedPayment, numSymbols uint, quorumNumbers []uint8) error {
	m.logger.Info("Recording and validating reservation usage", "header", header, "reservation", reservation)
	if !reservation.IsActive(uint64(time.Now().Unix())) {
		return fmt.Errorf("reservation not active")
	}
	if err := m.ValidateQuorum(quorumNumbers, reservation.QuorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for reservation: %w", err)
	}
	if !m.ValidateReservationPeriod(header, reservation) {
		return fmt.Errorf("invalid reservation period for reservation")
	}

	// Update bin usage atomically and check against reservation's data rate as the bin limit
	if err := m.IncrementBinUsage(ctx, header, reservation, numSymbols); err != nil {
		return fmt.Errorf("bin overflows: %w", err)
	}

	return nil
}

// ValidateQuorums ensures that the quorums listed in the blobHeader are present within allowedQuorums
// Note: A reservation that does not utilize all of the allowed quorums will be accepted. However, it
// will still charge against all of the allowed quorums. A on-demand requrests require and only allow
// the ETH and EIGEN quorums.
func (m *Meterer) ValidateQuorum(headerQuorums []uint8, allowedQuorums []uint8) error {
	if len(headerQuorums) == 0 {
		return fmt.Errorf("no quorum params in blob header")
	}

	// check that all the quorum ids are in ReservedPayment's
	for _, q := range headerQuorums {
		if !slices.Contains(allowedQuorums, q) {
			// fail the entire request if there's a quorum number mismatch
			return fmt.Errorf("quorum number mismatch: %d", q)
		}
	}
	return nil
}

// ValidateReservationPeriod checks if the provided reservation period is valid
func (m *Meterer) ValidateReservationPeriod(header core.PaymentMetadata, reservation *core.ReservedPayment) bool {
	now := uint64(time.Now().Unix())
	reservationWindow := m.ChainPaymentState.GetReservationWindow()
	currentReservationPeriod := GetReservationPeriod(now, reservationWindow)
	// Valid reservation periodes are either the current bin or the previous bin
	if (header.ReservationPeriod != currentReservationPeriod && header.ReservationPeriod != (currentReservationPeriod-1)) || (GetReservationPeriod(reservation.StartTimestamp, reservationWindow) > header.ReservationPeriod || header.ReservationPeriod > GetReservationPeriod(reservation.EndTimestamp, reservationWindow)) {
		return false
	}
	return true
}

// IncrementBinUsage increments the bin usage atomically and checks for overflow
func (m *Meterer) IncrementBinUsage(ctx context.Context, header core.PaymentMetadata, reservation *core.ReservedPayment, numSymbols uint) error {
	symbolsCharged := m.SymbolsCharged(numSymbols)
	newUsage, err := m.OffchainStore.UpdateReservationBin(ctx, header.AccountID, uint64(header.ReservationPeriod), uint64(symbolsCharged))
	if err != nil {
		return fmt.Errorf("failed to increment bin usage: %w", err)
	}

	// metered usage stays within the bin limit
	usageLimit := m.GetReservationBinLimit(reservation)
	if newUsage <= usageLimit {
		return nil
	} else if newUsage-uint64(symbolsCharged) >= usageLimit {
		// metered usage before updating the size already exceeded the limit
		return fmt.Errorf("bin has already been filled")
	}
	if newUsage <= 2*usageLimit && header.ReservationPeriod+2 <= GetReservationPeriod(reservation.EndTimestamp, m.ChainPaymentState.GetReservationWindow()) {
		_, err := m.OffchainStore.UpdateReservationBin(ctx, header.AccountID, uint64(header.ReservationPeriod+2), newUsage-usageLimit)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("overflow usage exceeds bin limit")
}

// GetReservationPeriod returns the current reservation period by chunking time by the bin interval;
// bin interval used by the disperser should be public information
func GetReservationPeriod(timestamp uint64, binInterval uint32) uint32 {
	if binInterval == 0 {
		return 0
	}
	return uint32(timestamp) / binInterval
}

// ServeOnDemandRequest handles the rate limiting logic for incoming requests
// On-demand requests doesn't have additional quorum settings and should only be
// allowed by ETH and EIGEN quorums
func (m *Meterer) ServeOnDemandRequest(ctx context.Context, header core.PaymentMetadata, onDemandPayment *core.OnDemandPayment, numSymbols uint, headerQuorums []uint8) error {
	m.logger.Info("Recording and validating on-demand usage", "header", header, "onDemandPayment", onDemandPayment)
	quorumNumbers, err := m.ChainPaymentState.GetOnDemandQuorumNumbers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get on-demand quorum numbers: %w", err)
	}

	if err := m.ValidateQuorum(headerQuorums, quorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for On-Demand Request: %w", err)
	}
	// update blob header to use the miniumum chargeable size
	symbolsCharged := m.SymbolsCharged(numSymbols)
	err = m.OffchainStore.AddOnDemandPayment(ctx, header, symbolsCharged)
	if err != nil {
		return fmt.Errorf("failed to update cumulative payment: %w", err)
	}
	// Validate payments attached
	err = m.ValidatePayment(ctx, header, onDemandPayment, numSymbols)
	if err != nil {
		// No tolerance for incorrect payment amounts; no rollbacks
		return fmt.Errorf("invalid on-demand payment: %w", err)
	}

	// Update bin usage atomically and check against bin capacity
	if err := m.IncrementGlobalBinUsage(ctx, uint64(symbolsCharged)); err != nil {
		//TODO: conditionally remove the payment based on the error type (maybe if the error is store-op related)
		dbErr := m.OffchainStore.RemoveOnDemandPayment(ctx, header.AccountID, header.CumulativePayment)
		if dbErr != nil {
			return dbErr
		}
		return fmt.Errorf("failed global rate limiting: %w", err)
	}

	return nil
}

// ValidatePayment checks if the provided payment header is valid against the local accounting
// prevPmt is the largest  cumulative payment strictly less    than PaymentMetadata.cumulativePayment if exists
// nextPmt is the smallest cumulative payment strictly greater than PaymentMetadata.cumulativePayment if exists
// nextPmtnumSymbols is the numSymbols of corresponding to nextPmt if exists
// prevPmt + PaymentMetadata.numSymbols * m.FixedFeePerByte
// <= PaymentMetadata.CumulativePayment
// <= nextPmt - nextPmtnumSymbols * m.FixedFeePerByte > nextPmt
func (m *Meterer) ValidatePayment(ctx context.Context, header core.PaymentMetadata, onDemandPayment *core.OnDemandPayment, numSymbols uint) error {
	if header.CumulativePayment.Cmp(onDemandPayment.CumulativePayment) > 0 {
		return fmt.Errorf("request claims a cumulative payment greater than the on-chain deposit")
	}

	prevPmt, nextPmt, nextPmtnumSymbols, err := m.OffchainStore.GetRelevantOnDemandRecords(ctx, header.AccountID, header.CumulativePayment) // zero if DNE
	if err != nil {
		return fmt.Errorf("failed to get relevant on-demand records: %w", err)
	}
	// the current request must increment cumulative payment by a magnitude sufficient to cover the blob size
	if prevPmt.Add(prevPmt, m.PaymentCharged(numSymbols)).Cmp(header.CumulativePayment) > 0 {
		return fmt.Errorf("insufficient cumulative payment increment")
	}
	// the current request must not break the payment magnitude for the next payment if the two requests were delivered out-of-order
	if nextPmt.Cmp(big.NewInt(0)) != 0 && header.CumulativePayment.Add(header.CumulativePayment, m.PaymentCharged(uint(nextPmtnumSymbols))).Cmp(nextPmt) > 0 {
		return fmt.Errorf("breaking cumulative payment invariants")
	}
	// check passed: blob can be safely inserted into the set of payments
	return nil
}

// PaymentCharged returns the chargeable price for a given data length
func (m *Meterer) PaymentCharged(numSymbols uint) *big.Int {
	symbolsCharged := big.NewInt(int64(m.SymbolsCharged(numSymbols)))
	pricePerSymbol := big.NewInt(int64(m.ChainPaymentState.GetPricePerSymbol()))
	return symbolsCharged.Mul(symbolsCharged, pricePerSymbol)
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func (m *Meterer) SymbolsCharged(numSymbols uint) uint32 {
	if numSymbols <= uint(m.ChainPaymentState.GetMinNumSymbols()) {
		return m.ChainPaymentState.GetMinNumSymbols()
	}
	// Round up to the nearest multiple of MinNumSymbols
	return uint32(core.RoundUpDivide(uint(numSymbols), uint(m.ChainPaymentState.GetMinNumSymbols()))) * m.ChainPaymentState.GetMinNumSymbols()
}

// IncrementBinUsage increments the bin usage atomically and checks for overflow
func (m *Meterer) IncrementGlobalBinUsage(ctx context.Context, symbolsCharged uint64) error {
	globalPeriod := GetReservationPeriod(uint64(time.Now().Unix()), m.ChainPaymentState.GetGlobalRatePeriodInterval())

	newUsage, err := m.OffchainStore.UpdateGlobalBin(ctx, globalPeriod, symbolsCharged)
	if err != nil {
		return fmt.Errorf("failed to increment global bin usage: %w", err)
	}
	if newUsage > m.ChainPaymentState.GetGlobalSymbolsPerSecond()*uint64(m.ChainPaymentState.GetGlobalRatePeriodInterval()) {
		return fmt.Errorf("global bin usage overflows")
	}
	return nil
}

// GetReservationBinLimit returns the bin limit for a given reservation
func (m *Meterer) GetReservationBinLimit(reservation *core.ReservedPayment) uint64 {
	return reservation.SymbolsPerSecond * uint64(m.ChainPaymentState.GetReservationWindow())
}
