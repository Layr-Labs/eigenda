package meterer

import (
	"context"
	"errors"
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
	m.logger.Debug("Validating incoming request's payment metadata", "paymentMetadata", header, "numSymbols", numSymbols, "quorumNumbers", quorumNumbers)
	params, err := m.ChainPaymentState.GetPaymentGlobalParams()
	if err != nil {
		return 0, fmt.Errorf("failed to get payment global params: %w", err)
	}
	// Validate against the payment method
	if !IsOnDemandPayment(&header) {
		reservations, err := m.ChainPaymentState.GetReservedPaymentByAccountAndQuorums(ctx, header.AccountID, quorumNumbers)
		if err != nil {
			return 0, fmt.Errorf("failed to get active reservation by account: %w", err)
		}
		if err := m.serveReservationRequest(ctx, params, header, reservations, numSymbols, quorumNumbers, receivedAt); err != nil {
			return 0, fmt.Errorf("invalid reservation request: %w", err)
		}
	} else {
		onDemandPayment, err := m.ChainPaymentState.GetOnDemandPaymentByAccount(ctx, header.AccountID)
		if err != nil {
			return 0, fmt.Errorf("failed to get on-demand payment by account: %w", err)
		}
		if err := m.serveOnDemandRequest(ctx, params, header, onDemandPayment, numSymbols, quorumNumbers, receivedAt); err != nil {
			return 0, fmt.Errorf("invalid on-demand request: %w", err)
		}
	}

	// TODO(hopeyen): each quorum can have different min num symbols; the returned symbolsCharged is only for used for metrics.
	// for now we simply return the charge for quorum 0, as quorums are likely to share the same min num symbols
	// we can make this more granular by adding metrics to the meterer later on
	_, protocolConfig, err := params.GetQuorumConfigs(OnDemandQuorumID)
	if err != nil {
		return 0, fmt.Errorf("failed to get on-demand quorum config: %w", err)
	}
	symbolsCharged := SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
	return symbolsCharged, nil
}

// serveReservationRequest handles the rate limiting logic for incoming requests
func (m *Meterer) serveReservationRequest(
	ctx context.Context,
	globalParams *PaymentVaultParams,
	header core.PaymentMetadata,
	reservations map[core.QuorumID]*core.ReservedPayment,
	numSymbols uint64,
	quorumNumbers []uint8,
	receivedAt time.Time,
) error {
	m.logger.Debug("Recording and validating reservation usage", "header", header, "reservation", reservations)
	if err := ValidateReservations(reservations, globalParams.QuorumProtocolConfigs, quorumNumbers, header.Timestamp, receivedAt); err != nil {
		return fmt.Errorf("invalid reservation: %w", err)
	}

	// Make atomic batched updates over all reservations identified by the same account and quorum
	if err := m.incrementBinUsage(ctx, header, reservations, globalParams, numSymbols); err != nil {
		return fmt.Errorf("failed to increment bin usages: %w", err)
	}
	return nil
}

// incrementBinUsage increments the bin usage atomically and checks for overflow
func (m *Meterer) incrementBinUsage(
	ctx context.Context, header core.PaymentMetadata,
	reservations map[core.QuorumID]*core.ReservedPayment,
	globalParams *PaymentVaultParams,
	numSymbols uint64,
) error {
	charges := make(map[core.QuorumID]uint64)
	quorumNumbers := make([]core.QuorumID, 0, len(reservations))
	reservationWindows := make(map[core.QuorumID]uint64, len(reservations))
	requestReservationPeriods := make(map[core.QuorumID]uint64, len(reservations))
	for quorumID := range reservations {
		_, protocolConfig, err := globalParams.GetQuorumConfigs(quorumID)
		if err != nil {
			return fmt.Errorf("failed to get quorum config for quorum %d: %w", quorumID, err)
		}
		charges[quorumID] = SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
		quorumNumbers = append(quorumNumbers, quorumID)
		reservationWindows[quorumID] = protocolConfig.ReservationRateLimitWindow
		requestReservationPeriods[quorumID] = GetReservationPeriodByNanosecond(header.Timestamp, protocolConfig.ReservationRateLimitWindow)
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
	overflowAmounts := make(map[core.QuorumID]uint64)
	overflowPeriods := make(map[core.QuorumID]uint64)

	for quorumID, reservation := range reservations {
		reservationWindow := reservationWindows[quorumID]
		requestReservationPeriod := requestReservationPeriods[quorumID]
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
		}
		overflowPeriod := GetOverflowPeriod(requestReservationPeriod, reservationWindow)
		if charges[quorumID] <= usageLimit && overflowPeriod <= GetReservationPeriod(int64(reservation.EndTimestamp), reservationWindow) {
			// Needs to go to overflow bin
			overflowAmounts[quorumID] = newUsage - usageLimit
			overflowPeriods[quorumID] = overflowPeriod
		} else {
			return fmt.Errorf("overflow usage exceeds bin limit for quorum %d", quorumID)
		}
	}
	if len(overflowAmounts) != len(overflowPeriods) {
		return fmt.Errorf("overflow amount and period mismatch")
	}
	// Batch increment overflow bins for all overflown reservation candidates
	if len(overflowAmounts) > 0 {
		m.logger.Debug("Utilizing reservation overflow period", "overflowAmounts", overflowAmounts, "overflowPeriods", overflowPeriods)
		overflowQuorums := make([]core.QuorumID, 0, len(overflowAmounts))
		for quorumID := range overflowAmounts {
			overflowQuorums = append(overflowQuorums, quorumID)
		}
		_, err := m.MeteringStore.IncrementBinUsages(ctx, header.AccountID, overflowQuorums, overflowPeriods, overflowAmounts)
		if err != nil {
			// Rollback the increments for the current periods
			rollbackErr := m.MeteringStore.DecrementBinUsages(ctx, header.AccountID, quorumNumbers, requestReservationPeriods, charges)
			if rollbackErr != nil {
				return fmt.Errorf("failed to increment overflow bins: %w; rollback also failed: %v", err, rollbackErr)
			}
			return fmt.Errorf("failed to increment overflow bins: %w; successfully rolled back increments", err)
		}
	}

	return nil
}

// serveOnDemandRequest handles the rate limiting logic for incoming requests
// On-demand requests doesn't have additional quorum settings and should only be
// allowed by ETH and EIGEN quorums
func (m *Meterer) serveOnDemandRequest(ctx context.Context, globalParams *PaymentVaultParams, header core.PaymentMetadata, onDemandPayment *core.OnDemandPayment, symbolsCharged uint64, headerQuorums []uint8, receivedAt time.Time) error {
	m.logger.Debug("Recording and validating on-demand usage", "header", header, "onDemandPayment", onDemandPayment)

	if err := ValidateQuorum(headerQuorums, globalParams.OnDemandQuorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for On-Demand Request: %w", err)
	}

	// Verify that the claimed cumulative payment doesn't exceed the on-chain deposit
	if header.CumulativePayment.Cmp(onDemandPayment.CumulativePayment) > 0 {
		return fmt.Errorf("request claims a cumulative payment greater than the on-chain deposit")
	}

	paymentConfig, protocolConfig, err := globalParams.GetQuorumConfigs(OnDemandQuorumID)
	if err != nil {
		return fmt.Errorf("failed to get payment config for on-demand quorum: %w", err)
	}

	symbolsCharged = SymbolsCharged(symbolsCharged, protocolConfig.MinNumSymbols)
	paymentCharged := PaymentCharged(symbolsCharged, paymentConfig.OnDemandPricePerSymbol)
	oldPayment, err := m.MeteringStore.AddOnDemandPayment(ctx, header, paymentCharged)
	if err != nil {
		return fmt.Errorf("failed to update cumulative payment: %w", err)
	}

	// Update bin usage atomically and check against bin capacity
	if err := m.incrementGlobalBinUsage(ctx, globalParams, uint64(symbolsCharged), receivedAt); err != nil {
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
func (m *Meterer) incrementGlobalBinUsage(ctx context.Context, params *PaymentVaultParams, symbolsCharged uint64, receivedAt time.Time) error {
	paymentConfig, protocolConfig, err := params.GetQuorumConfigs(OnDemandQuorumID)
	if err != nil {
		return fmt.Errorf("failed to get quorum configs for on-demand quorum: %w", err)
	}

	globalPeriod := GetReservationPeriod(receivedAt.Unix(), protocolConfig.OnDemandRateLimitWindow)

	newUsage, err := m.MeteringStore.UpdateGlobalBin(ctx, globalPeriod, symbolsCharged)
	if err != nil {
		return fmt.Errorf("failed to increment global bin usage: %w", err)
	}
	if newUsage > GetBinLimit(paymentConfig.OnDemandSymbolsPerSecond, protocolConfig.OnDemandRateLimitWindow) {
		return fmt.Errorf("global bin usage overflows")
	}
	return nil
}

// GetReservationBinLimit returns the bin limit for a given reservation
// Note: This is called per-quorum since reservation is for a single quorum.
func GetReservationBinLimit(reservation *core.ReservedPayment, reservationWindow uint64) uint64 {
	return GetBinLimit(reservation.SymbolsPerSecond, reservationWindow)
}

// GetBinLimit returns the bin limit given the bin interval and the symbols per second
func GetBinLimit(symbolsPerSecond uint64, binInterval uint64) uint64 {
	return symbolsPerSecond * binInterval
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
	// directly convert to uint64 to avoid overflow
	numSymbolsInt := new(big.Int).SetUint64(numSymbols)
	pricePerSymbolInt := new(big.Int).SetUint64(pricePerSymbol)
	return new(big.Int).Mul(numSymbolsInt, pricePerSymbolInt)
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func SymbolsCharged(numSymbols uint64, minSymbols uint64) uint64 {
	if numSymbols <= minSymbols {
		return minSymbols
	}
	if minSymbols == 0 {
		return numSymbols
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

// ValidateReservations ensures that the quorums listed in the blobHeader are present within allowedQuorums.
//
// Parameters:
//   - timestamp: time in nanoseconds
//
// Notes:
//   - Reservations that don't use all allowed quorums are still accepted
//   - Charges apply to ALL allowed quorums, even if not all are used
//   - On-demand requests have special requirements: they must use ETH and EIGEN quorums only
func ValidateReservations(
	reservations map[core.QuorumID]*core.ReservedPayment,
	quorumConfigs map[core.QuorumID]*core.PaymentQuorumProtocolConfig,
	quorumNumbers []uint8,
	timestamp int64,
	receivedAt time.Time,
) error {
	reservationQuorums := make([]uint8, 0, len(reservations))
	reservationWindows := make(map[core.QuorumID]uint64, len(reservations))
	requestReservationPeriods := make(map[core.QuorumID]uint64, len(reservations))

	// Gather quorums the user had an reservations on and relevant quorum configurations
	for quorumID := range reservations {
		reservationQuorums = append(reservationQuorums, uint8(quorumID))
		_, ok := quorumConfigs[quorumID]
		if !ok {
			return fmt.Errorf("quorum config not found for quorum %d", quorumID)
		}
		reservationWindows[quorumID] = quorumConfigs[quorumID].ReservationRateLimitWindow
		requestReservationPeriods[quorumID] = GetReservationPeriodByNanosecond(timestamp, quorumConfigs[quorumID].ReservationRateLimitWindow)
	}
	if err := ValidateQuorum(quorumNumbers, reservationQuorums); err != nil {
		return err
	}
	// Validate the used reservations are active and is of valid periods
	for _, quorumID := range quorumNumbers {
		reservation := reservations[core.QuorumID(quorumID)]
		if !reservation.IsActiveByNanosecond(timestamp) {
			return errors.New("reservation not active")
		}
		if !ValidateReservationPeriod(reservation, requestReservationPeriods[quorumID], reservationWindows[quorumID], receivedAt) {
			return fmt.Errorf("invalid reservation period for reservation on quorum %d", quorumID)
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

func IsOnDemandPayment(paymentMetadata *core.PaymentMetadata) bool {
	return paymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) > 0
}
