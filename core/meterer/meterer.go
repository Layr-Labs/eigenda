package meterer

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer/paymentlogic"
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
	m.logger.Info("Validating incoming request's payment metadata", "paymentMetadata", header, "numSymbols", numSymbols, "quorumNumbers", quorumNumbers)

	params, err := m.ChainPaymentState.GetPaymentGlobalParams()
	if err != nil {
		return 0, fmt.Errorf("failed to get payment global params: %w", err)
	}
	// Validate against the payment method
	if !paymentlogic.IsOnDemandPayment(&header) {
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
	symbolsCharged := paymentlogic.SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
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
	if err := paymentlogic.ValidateReservations(reservations, globalParams.QuorumProtocolConfigs, quorumNumbers, header.Timestamp, receivedAt.UnixNano()); err != nil {
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
		charges[quorumID] = paymentlogic.SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
		quorumNumbers = append(quorumNumbers, quorumID)
		reservationWindows[quorumID] = protocolConfig.ReservationRateLimitWindow
		requestReservationPeriods[quorumID] = paymentlogic.GetReservationPeriodByNanosecond(header.Timestamp, protocolConfig.ReservationRateLimitWindow)
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
		usageLimit := paymentlogic.GetBinLimit(reservation.SymbolsPerSecond, reservationWindow)
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
		overflowPeriod := paymentlogic.GetOverflowPeriod(requestReservationPeriod, reservationWindow)
		if charges[quorumID] <= usageLimit && overflowPeriod <= paymentlogic.GetReservationPeriod(int64(reservation.EndTimestamp), reservationWindow) {
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

	if err := paymentlogic.ValidateQuorum(headerQuorums, globalParams.OnDemandQuorumNumbers); err != nil {
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

	symbolsCharged = paymentlogic.SymbolsCharged(symbolsCharged, protocolConfig.MinNumSymbols)
	paymentCharged := paymentlogic.PaymentCharged(symbolsCharged, paymentConfig.OnDemandPricePerSymbol)
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

	globalPeriod := paymentlogic.GetReservationPeriod(receivedAt.Unix(), protocolConfig.OnDemandRateLimitWindow)

	newUsage, err := m.MeteringStore.UpdateGlobalBin(ctx, globalPeriod, symbolsCharged)
	if err != nil {
		return fmt.Errorf("failed to increment global bin usage: %w", err)
	}
	if newUsage > paymentlogic.GetBinLimit(paymentConfig.OnDemandSymbolsPerSecond, protocolConfig.OnDemandRateLimitWindow) {
		return fmt.Errorf("global bin usage overflows")
	}
	return nil
}
