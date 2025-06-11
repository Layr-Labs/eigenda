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

const OnDemandQuorumID = core.QuorumID(0)

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
func (m *Meterer) Start(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(m.Config.UpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := m.ChainPaymentState.RefreshOnchainPaymentState(ctx)
				if err != nil {
					m.logger.Error("failed to refresh onchain state", "error", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// MeterRequest validates a blob header and adds it to the meterer's state
// TODO: return error if there's a rejection (with reasoning) or internal error (should be very rare)
func (m *Meterer) MeterRequest(ctx context.Context, header core.PaymentMetadata, numSymbols uint64, quorumNumbers []uint8, receivedAt time.Time) (uint64, error) {
	m.logger.Info("Validating incoming request's payment metadata", "paymentMetadata", header, "numSymbols", numSymbols, "quorumNumbers", quorumNumbers)
	symbolsCharged := SymbolsCharged(numSymbols, uint64(m.ChainPaymentState.GetMinNumSymbols(OnDemandQuorumID)))
	// Validate against the payment method
	if header.CumulativePayment.Sign() == 0 {
		reservations, err := m.ChainPaymentState.GetReservedPaymentByAccountAndQuorums(ctx, header.AccountID, quorumNumbers)
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
func (m *Meterer) ServeReservationRequest(ctx context.Context, header core.PaymentMetadata, reservations map[core.QuorumID]*core.ReservedPayment, symbolsCharged uint64, quorumNumbers []uint8, receivedAt time.Time) error {
	m.logger.Info("Recording and validating reservation usage", "header", header, "reservation", reservations)
	// Take all the quorumIDs from the reservations
	quorumIDs := make([]core.QuorumID, 0, len(reservations))
	reservationWindows := make(map[core.QuorumID]uint64, len(reservations))
	requestReservationPeriods := make(map[core.QuorumID]uint64, len(reservations))
	// Gather quorums the user had an reservations on and relevant quorum configurations
	for quorumID := range reservations {
		quorumIDs = append(quorumIDs, quorumID)
		reservationWindows[quorumID] = m.ChainPaymentState.GetReservationWindow(quorumID)
		requestReservationPeriods[quorumID] = GetReservationPeriodByNanosecond(header.Timestamp, m.ChainPaymentState.GetReservationWindow(quorumID))
	}
	// Validate quorumIDs is a subset of the quorumNumbers in the dispersal request; this should be guaranteed by GetReservedPaymentByAccountAndQuorums
	if err := ValidateQuorum(quorumNumbers, quorumIDs); err != nil {
		return fmt.Errorf("invalid quorum for reservation: %w", err)
	}

	// Validate the used reservations are active and is of valid periods
	for quorumID, reservation := range reservations {
		if !IsWithinTimeRange(reservation.StartTimestamp, reservation.EndTimestamp, header.Timestamp) {
			return fmt.Errorf("reservation not active")
		}
		if !ValidateReservationPeriod(reservation, requestReservationPeriods[quorumID], reservationWindows[quorumID], receivedAt) {
			return fmt.Errorf("invalid reservation period for reservation on quorum %d", quorumID)
		}
	}

	// Make atomic batched updates over all reservations identified by the same account and quorum
	if err := m.IncrementBinUsage(ctx, header, reservations, symbolsCharged, reservationWindows, requestReservationPeriods); err != nil {
		return fmt.Errorf("failed to increment bin usages: %w", err)
	}
	return nil
}

// IncrementBinUsage increments the bin usage atomically and checks for overflow
func (m *Meterer) IncrementBinUsage(ctx context.Context, header core.PaymentMetadata, reservations map[core.QuorumID]*core.ReservedPayment, symbolsCharged uint64, reservationWindows map[core.QuorumID]uint64, requestReservationPeriods map[core.QuorumID]uint64) error {
	charges := make(map[core.QuorumID]uint64)
	for quorumID := range reservations {
		charges[quorumID] = symbolsCharged
	}
	quorumNumbers := make([]core.QuorumID, 0, len(reservations))
	for quorumID := range reservations {
		quorumNumbers = append(quorumNumbers, quorumID)
	}

	// Batch increment all quorums for the current quorums' reservation period
	// For each quorum, increment by its specific symbolsCharged value
	updatedUsages := make(map[core.QuorumID]uint64)
	usage, err := m.MeteringStore.IncrementBinUsages(ctx, header.AccountID, quorumNumbers, requestReservationPeriods, charges)
	if err != nil {
		return err
	}
	for _, quorumID := range quorumNumbers {
		updatedUsages[quorumID] = usage[quorumID]
	}

	overflowCandidates := make(map[core.QuorumID]struct{})
	overflowAmounts := make(map[core.QuorumID]uint64)

	for quorumID, reservation := range reservations {
		reservationWindow := reservationWindows[quorumID]
		usageLimit := GetReservationBinLimit(reservation, reservationWindow)
		newUsage, ok := updatedUsages[quorumID]
		if !ok {
			return fmt.Errorf("failed to get updated usage for quorum %d", quorumID)
		}
		prevUsage := newUsage - charges[quorumID]
		if newUsage <= usageLimit {
			continue
		} else if prevUsage >= usageLimit {
			// Bin was already filled before this increment
			return fmt.Errorf("bin has already been filled for quorum %d", quorumID)
		} else if charges[quorumID] <= usageLimit && GetOverflowPeriod(requestReservationPeriods[quorumID], reservationWindow) <= GetReservationPeriod(int64(reservation.EndTimestamp), reservationWindow) {
			// Needs to go to overflow bin
			overflowCandidates[quorumID] = struct{}{}
			overflowAmounts[quorumID] = newUsage - usageLimit
		} else {
			return fmt.Errorf("overflow usage exceeds bin limit for quorum %d", quorumID)
		}
	}

	// Batch increment overflow bins for overflown reservation candidates
	for quorumID := range overflowCandidates {
		_, err := m.MeteringStore.IncrementBinUsages(ctx, header.AccountID, []core.QuorumID{quorumID}, map[core.QuorumID]uint64{quorumID: requestReservationPeriods[quorumID] + 2}, map[core.QuorumID]uint64{quorumID: overflowAmounts[quorumID]})
		if err != nil {
			// Rollback the increments for the current periods
			rollbackErr := m.MeteringStore.DecrementBinUsages(ctx, header.AccountID, quorumNumbers, requestReservationPeriods, charges)
			if rollbackErr != nil {
				return fmt.Errorf("failed to increment overflow bin for quorum %d: %w; rollback also failed: %v", quorumID, err, rollbackErr)
			}
			return fmt.Errorf("failed to increment overflow bin for quorum %d: %w; successfully rolled back increments", quorumID, err)
		}
	}

	return nil
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

	if err := ValidateQuorum(headerQuorums, quorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for On-Demand Request: %w", err)
	}

	// Verify that the claimed cumulative payment doesn't exceed the on-chain deposit
	if header.CumulativePayment.Cmp(onDemandPayment.CumulativePayment) > 0 {
		return fmt.Errorf("request claims a cumulative payment greater than the on-chain deposit")
	}

	paymentCharged := PaymentCharged(symbolsCharged, m.ChainPaymentState.GetPricePerSymbol(OnDemandQuorumID))
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

// IncrementGlobalBinUsage increments the bin usage atomically and checks for overflow
func (m *Meterer) IncrementGlobalBinUsage(ctx context.Context, symbolsCharged uint64, receivedAt time.Time) error {
	globalPeriod := GetReservationPeriod(receivedAt.Unix(), m.ChainPaymentState.GetOnDemandGlobalRatePeriodInterval(OnDemandQuorumID))

	newUsage, err := m.MeteringStore.UpdateGlobalBin(ctx, globalPeriod, symbolsCharged)
	if err != nil {
		return fmt.Errorf("failed to increment global bin usage: %w", err)
	}
	if newUsage > m.ChainPaymentState.GetOnDemandGlobalSymbolsPerSecond(OnDemandQuorumID)*uint64(m.ChainPaymentState.GetOnDemandGlobalRatePeriodInterval(OnDemandQuorumID)) {
		return fmt.Errorf("global bin usage overflows: %d > %d", newUsage, m.ChainPaymentState.GetOnDemandGlobalSymbolsPerSecond(OnDemandQuorumID)*uint64(m.ChainPaymentState.GetOnDemandGlobalRatePeriodInterval(OnDemandQuorumID)))
	}
	return nil
}

// GetReservationBinLimit returns the bin limit for a given reservation
// Note: This is called per-quorum since reservation is for a single quorum.
func GetReservationBinLimit(reservation *core.ReservedPayment, reservationWindow uint64) uint64 {
	return reservation.SymbolsPerSecond * reservationWindow
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

// GetOverflowPeriod returns the overflow period by adding the overflow offset to the current reservation period
// the offset is 2*reservationWindow, skipping the immediate next period for the period that will be used for overflow from the current period
func GetOverflowPeriod(reservationPeriod uint64, reservationWindow uint64) uint64 {
	return reservationPeriod + reservationWindow*2
}

// PaymentCharged returns the chargeable price for a given number of symbols
func PaymentCharged(numSymbols, pricePerSymbol uint64) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(numSymbols)), big.NewInt(int64(pricePerSymbol)))
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func SymbolsCharged(numSymbols uint64, minSymbols uint64) uint64 {
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

// ValidateQuorum ensures that the quorums listed in the blobHeader are present within allowedQuorums
// Note: A reservation that does not utilize all of the allowed quorums will be accepted. However, it
// will still charge against all of the allowed quorums. A on-demand requests require and only allow
// the ETH and EIGEN quorums.
func ValidateQuorum(headerQuorums []uint8, allowedQuorums []uint8) error {
	if len(headerQuorums) == 0 {
		return fmt.Errorf("no quorum numbers provided in the request")
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
// Note: This is called per-quorum since reservation is for a single quorum.
func ValidateReservationPeriod(reservation *core.ReservedPayment, requestReservationPeriod uint64, reservationWindow uint64, receivedAt time.Time) bool {
	currentReservationPeriod := GetReservationPeriod(receivedAt.Unix(), reservationWindow)
	// Valid reservation periods are either the current bin or the previous bin
	isCurrentOrPreviousPeriod := requestReservationPeriod == currentReservationPeriod || requestReservationPeriod == (currentReservationPeriod-reservationWindow)
	startPeriod := GetReservationPeriod(int64(reservation.StartTimestamp), reservationWindow)
	endPeriod := GetReservationPeriod(int64(reservation.EndTimestamp), reservationWindow)
	isWithinReservationWindow := startPeriod <= requestReservationPeriod && requestReservationPeriod < endPeriod
	if !isCurrentOrPreviousPeriod || !isWithinReservationWindow {
		return false
	}
	return true
}
