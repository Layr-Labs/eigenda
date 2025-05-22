package meterer

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Config contains network parameters that should be published on-chain. We currently configure these params through disperser env vars.
type Config struct {

	// ChainReadTimeout is the timeout for reading payment state from chain
	ChainReadTimeout time.Duration

	// UpdateInterval is the interval for refreshing the on-chain state
	UpdateInterval time.Duration
}

// Meterer handles payment accounting across different accounts. Disperser API server receives requests from clients and each request contains a blob header
// with payments information (CumulativePayments, Timestamp, and Signature). Disperser will pass the blob header to the meterer, which will check if the
// payments information is valid.
type Meterer struct {
	Config

	// ChainPaymentState reads on-chain payment state periodically and caches it in memory
	ChainPaymentState OnchainPayment

	// MeteringStore tracks usage and payments in a storage backend
	MeteringStore MeteringStore

	logger logging.Logger
}

func NewMeterer(
	config Config,
	paymentChainState OnchainPayment,
	meteringStore MeteringStore,
	logger logging.Logger,
) *Meterer {
	return &Meterer{
		Config:            config,
		ChainPaymentState: paymentChainState,
		MeteringStore:     meteringStore,
		logger:            logger.With("component", "Meterer"),
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
func (m *Meterer) MeterRequest(ctx context.Context, header core.PaymentMetadata, numSymbols uint64, quorumNumbers []uint8, receivedAt time.Time) (uint64, error) {
	symbolsCharged := m.SymbolsCharged(numSymbols)
	m.logger.Info("Validating incoming request's payment metadata", "paymentMetadata", header, "numSymbols", numSymbols, "quorumNumbers", quorumNumbers)
	// Validate against the payment method
	if header.CumulativePayment.Sign() == 0 {
		reservations, err := m.ChainPaymentState.GetReservedPaymentByAccount(ctx, header.AccountID)
		if err != nil {
			return 0, fmt.Errorf("failed to get active reservation by account: %w", err)
		}
		if err := m.ServeReservationRequest(ctx, header, reservations, symbolsCharged, quorumNumbers, receivedAt); err != nil {
			return 0, fmt.Errorf("invalid reservation: %w", err)
		}
	} else {
		onDemandPayment, err := m.ChainPaymentState.GetOnDemandPaymentByAccount(ctx, header.AccountID)
		if err != nil {
			return 0, fmt.Errorf("failed to get on-demand payment by account: %w", err)
		}
		if err := m.ServeOnDemandRequest(ctx, header, onDemandPayment, symbolsCharged, quorumNumbers, receivedAt); err != nil {
			return 0, fmt.Errorf("invalid on-demand request: %w", err)
		}
	}

	return symbolsCharged, nil
}

// ServeReservationRequest handles the rate limiting logic for incoming requests
func (m *Meterer) ServeReservationRequest(ctx context.Context, header core.PaymentMetadata, reservations map[uint8]*core.ReservedPayment, symbolsCharged uint64, quorumNumbers []uint8, receivedAt time.Time) error {
	m.logger.Info("Recording and validating reservation usage", "header", header, "reservations", reservations)
	reservedQuorumNumbers := make([]uint8, 0, len(reservations))
	for quorum := range reservations {
		reservedQuorumNumbers = append(reservedQuorumNumbers, quorum)
	}
	fmt.Println("reservedQuorumNumbers", reservedQuorumNumbers, "header quorumNumbers", quorumNumbers)
	if err := m.ValidateQuorum(quorumNumbers, reservedQuorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for reservation: %w", err)
	}
	reservationWindow := m.ChainPaymentState.GetReservationWindow()
	requestReservationPeriod := GetReservationPeriodByNanosecond(header.Timestamp, reservationWindow)
	for _, quorum := range quorumNumbers {
		reservation, ok := reservations[quorum]
		if !ok {
			return fmt.Errorf("no reservation found for quorum %d", quorum)
		}
		if !reservation.IsActiveByNanosecond(header.Timestamp) {
			return fmt.Errorf("reservation not active for quorum %d", quorum)
		}
		if !m.ValidateReservationPeriod(reservation, requestReservationPeriod, reservationWindow, receivedAt) {
			return fmt.Errorf("invalid reservation period for reservation quorum %d", quorum)
		}
	}
	// TODO(hopeyen): uncomment the loop once the offchain store PR is merged
	// Update bin usage atomically for each quorum and check against reservation's data rate as the bin limit
	// for _, quorum := range quorumNumbers {
	// 	if err := m.IncrementBinUsage(ctx, header, reservations[quorum], quorum, symbolsCharged, reservationWindow, requestReservationPeriod); err != nil {
	// 		return fmt.Errorf("bin overflows for quorum %d reservation period %d: %w", quorum, requestReservationPeriod, err)
	// 	}
	// }
	if err := m.IncrementBinUsage(ctx, header, reservations[0], 0, symbolsCharged, reservationWindow, requestReservationPeriod); err != nil {
		return fmt.Errorf("bin overflows for quorum %d reservation period %d (quorum is fixed to the first one as testing placeholder): %w", 0, requestReservationPeriod, err)
	}
	return nil
}

// ValidateQuorums ensures that the quorums listed in the blobHeader are present within allowedQuorums
// Note: A reservation that does not utilize all of the allowed quorums will be accepted. However, it
// will still charge against all of the allowed quorums. A on-demand requests require and only allow
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
// Note: This is called per-quorum, so reservation is for a single quorum.
func (m *Meterer) ValidateReservationPeriod(reservation *core.ReservedPayment, requestReservationPeriod uint64, reservationWindow uint64, receivedAt time.Time) bool {
	currentReservationPeriod := GetReservationPeriod(receivedAt.Unix(), reservationWindow)
	// Valid reservation periods are either the current bin or the previous bin
	isCurrentOrPreviousPeriod := requestReservationPeriod == currentReservationPeriod || requestReservationPeriod == (currentReservationPeriod-reservationWindow)
	startPeriod := GetReservationPeriod(int64(reservation.StartTimestamp), reservationWindow)
	endPeriod := GetReservationPeriod(int64(reservation.EndTimestamp), reservationWindow)
	fmt.Println("startPeriod", startPeriod, "endPeriod", endPeriod, "requestReservationPeriod", requestReservationPeriod, "currentReservationPeriod", currentReservationPeriod, "isCurrentOrPreviousPeriod", isCurrentOrPreviousPeriod)
	isWithinReservationWindow := startPeriod <= requestReservationPeriod && requestReservationPeriod < endPeriod
	if !isCurrentOrPreviousPeriod || !isWithinReservationWindow {
		return false
	}
	return true
}

// IncrementBinUsage increments the bin usage atomically and checks for overflow
// Note: This is called per-quorum, so reservation is for a single quorum.
func (m *Meterer) IncrementBinUsage(ctx context.Context, header core.PaymentMetadata, reservation *core.ReservedPayment, quorumNumber uint8, symbolsCharged uint64, reservationWindow uint64, requestReservationPeriod uint64) error {
	newUsage, err := m.MeteringStore.UpdateReservationBin(ctx, header.AccountID, requestReservationPeriod, quorumNumber, symbolsCharged)
	if err != nil {
		return fmt.Errorf("failed to increment bin usage: %w", err)
	}

	// metered usage stays within the bin limit
	usageLimit := m.GetReservationBinLimit(reservation, reservationWindow)
	if newUsage <= usageLimit {
		return nil
	} else if newUsage-symbolsCharged >= usageLimit {
		// metered usage before updating the size already exceeded the limit
		return fmt.Errorf("bin has already been filled")
	}
	if newUsage <= 2*usageLimit && requestReservationPeriod+2 <= GetReservationPeriod(int64(reservation.EndTimestamp), reservationWindow) {
		_, err := m.MeteringStore.UpdateReservationBin(ctx, header.AccountID, uint64(requestReservationPeriod+2), quorumNumber, newUsage-usageLimit)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("overflow usage exceeds bin limit")
}

// GetReservationPeriodByNanosecondTimestamp returns the current reservation period by finding the nearest lower multiple of the bin interval;
// bin interval used by the disperser is publicly recorded on-chain at the payment vault contract
func GetReservationPeriodByNanosecond(nanosecondTimestamp int64, binInterval uint64) uint64 {
	if nanosecondTimestamp < 0 {
		return 0
	}
	return GetReservationPeriod(int64((time.Duration(nanosecondTimestamp) * time.Nanosecond).Seconds()), binInterval)
}

// GetReservationPeriod returns the current reservation period by finding the nearest lower multiple of the bin interval;
// bin interval used by the disperser is publicly recorded on-chain at the payment vault contract
func GetReservationPeriod(timestamp int64, binInterval uint64) uint64 {
	if binInterval == 0 {
		return 0
	}
	return uint64(timestamp) / binInterval * binInterval
}

// ServeOnDemandRequest handles the rate limiting logic for incoming requests
// On-demand requests doesn't have additional quorum settings and should only be
// allowed by ETH and EIGEN quorums
func (m *Meterer) ServeOnDemandRequest(ctx context.Context, header core.PaymentMetadata, onDemandPayment *core.OnDemandPayment, symbolsCharged uint64, headerQuorums []uint8, receivedAt time.Time) error {
	m.logger.Debug("Recording and validating on-demand usage", "header", header, "onDemandPayment", onDemandPayment)
	quorumNumbers, err := m.ChainPaymentState.GetOnDemandQuorumNumbers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get on-demand quorum numbers: %w", err)
	}

	if err := m.ValidateQuorum(headerQuorums, quorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for On-Demand Request: %w", err)
	}

	// Verify that the claimed cumulative payment doesn't exceed the on-chain deposit
	if header.CumulativePayment.Cmp(onDemandPayment.CumulativePayment) > 0 {
		return fmt.Errorf("request claims a cumulative payment greater than the on-chain deposit")
	}

	paymentCharged := PaymentCharged(symbolsCharged, m.ChainPaymentState.GetPricePerSymbol())
	oldPayment, err := m.MeteringStore.AddOnDemandPayment(ctx, header, paymentCharged)
	if err != nil {
		return fmt.Errorf("failed to update cumulative payment: %w", err)
	}

	// Update bin usage atomically and check against bin capacity
	if err := m.IncrementGlobalBinUsage(ctx, uint64(symbolsCharged), receivedAt); err != nil {
		// If global bin usage update fails, roll back the payment to its previous value
		// The rollback will only happen if the current payment value still matches what we just wrote
		// This ensures we don't accidentally roll back a newer payment that might have been processed
		dbErr := m.MeteringStore.RollbackOnDemandPayment(ctx, header.AccountID, header.CumulativePayment, oldPayment)
		if dbErr != nil {
			return dbErr
		}
		return fmt.Errorf("failed global rate limiting: %w", err)
	}

	return nil
}

// PaymentCharged returns the chargeable price for a given number of symbols
func PaymentCharged(numSymbols, pricePerSymbol uint64) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(numSymbols)), big.NewInt(int64(pricePerSymbol)))
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func (m *Meterer) SymbolsCharged(numSymbols uint64) uint64 {
	minSymbols := uint64(m.ChainPaymentState.GetMinNumSymbols())
	if numSymbols <= minSymbols {
		return minSymbols
	}
	// Round up to the nearest multiple of MinNumSymbols
	roundedUp := core.RoundUpDivide(numSymbols, minSymbols) * minSymbols
	// Check for overflow; this case should never happen
	if roundedUp < numSymbols {
		return math.MaxUint64
	}
	return roundedUp
}

// IncrementGlobalBinUsage increments the bin usage atomically and checks for overflow
func (m *Meterer) IncrementGlobalBinUsage(ctx context.Context, symbolsCharged uint64, receivedAt time.Time) error {
	globalPeriod := GetReservationPeriod(receivedAt.Unix(), m.ChainPaymentState.GetGlobalRatePeriodInterval())

	newUsage, err := m.MeteringStore.UpdateGlobalBin(ctx, globalPeriod, symbolsCharged)
	if err != nil {
		return fmt.Errorf("failed to increment global bin usage: %w", err)
	}
	if newUsage > m.ChainPaymentState.GetGlobalSymbolsPerSecond()*uint64(m.ChainPaymentState.GetGlobalRatePeriodInterval()) {
		return fmt.Errorf("global bin usage overflows")
	}
	return nil
}

// GetReservationBinLimit returns the bin limit for a given reservation
// Note: This is called per-quorum, so reservation is for a single quorum.
func (m *Meterer) GetReservationBinLimit(reservation *core.ReservedPayment, reservationWindow uint64) uint64 {
	return reservation.SymbolsPerSecond * reservationWindow
}
