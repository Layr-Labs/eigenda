package meterer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// OnchainPaymentState is an interface for getting information about the current chain state for payments.
type OnchainPayment interface {
	// State management
	RefreshOnchainPaymentState(ctx context.Context) error

	// Account queries
	GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error)

	// Config access
	GetPaymentGlobalParams() (*PaymentVaultParams, error)
}

var _ OnchainPayment = (*OnchainPaymentState)(nil)

// OnchainPaymentState manages the state of on-chain payments including reservations and on-demand payments
type OnchainPaymentState struct {
	tx     *eth.Reader
	logger logging.Logger

	ReservedPayments map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment
	OnDemandPayments map[gethcommon.Address]*core.OnDemandPayment

	ReservationsLock sync.RWMutex
	OnDemandLocks    sync.RWMutex

	PaymentVaultParams atomic.Pointer[PaymentVaultParams]
}

// PaymentVaultParams contains all configuration parameters for the payment vault
type PaymentVaultParams struct {
	QuorumPaymentConfigs  map[core.QuorumID]*core.PaymentQuorumConfig
	QuorumProtocolConfigs map[core.QuorumID]*core.PaymentQuorumProtocolConfig
	OnDemandQuorumNumbers []core.QuorumID
}

// NewOnchainPaymentState creates a new OnchainPaymentState instance and initializes it with current chain state
func NewOnchainPaymentState(ctx context.Context, tx *eth.Reader, logger logging.Logger) (*OnchainPaymentState, error) {
	state := OnchainPaymentState{
		tx:                 tx,
		logger:             logger.With("component", "OnchainPaymentState"),
		ReservedPayments:   make(map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment),
		OnDemandPayments:   make(map[gethcommon.Address]*core.OnDemandPayment),
		PaymentVaultParams: atomic.Pointer[PaymentVaultParams]{},
	}

	err := state.RefreshOnchainPaymentState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize payment state: %w", err)
	}

	return &state, nil
}

// NewOnchainPaymentStateEmpty creates a new OnchainPaymentState instance without initializing chain state
func NewOnchainPaymentStateEmpty(ctx context.Context, tx *eth.Reader, logger logging.Logger) (*OnchainPaymentState, error) {
	state := OnchainPaymentState{
		tx:                 tx,
		logger:             logger.With("component", "OnchainPaymentState"),
		ReservedPayments:   make(map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment),
		OnDemandPayments:   make(map[gethcommon.Address]*core.OnDemandPayment),
		PaymentVaultParams: atomic.Pointer[PaymentVaultParams]{},
	}

	return &state, nil
}

// GetPaymentVaultParams retrieves the current payment vault parameters from the chain
// TODO(hopeyen): this function will be updated with the new UsageAuthorizationRegistry interface updates
func (pcs *OnchainPaymentState) GetPaymentVaultParams(ctx context.Context) (*PaymentVaultParams, error) {
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}

	requiredQuorumNumbers, err := pcs.tx.GetRequiredQuorumNumbers(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get required quorum numbers: %w", err)
	}

	quorumCount, err := pcs.tx.GetQuorumCount(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get quorum count: %w", err)
	}

	quorumNumbers := make([]uint8, quorumCount)
	for i := range quorumNumbers {
		quorumNumbers[i] = uint8(i)
	}

	// TODO(hopeyen): the construction of quorum configs will be updated with payment vault interface updates
	globalSymbolsPerSecond, err := pcs.tx.GetOnDemandGlobalSymbolsPerSecond(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get global symbols per second: %w", err)
	}
	globalRatePeriodInterval, err := pcs.tx.GetOnDemandGlobalRatePeriodInterval(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get global rate period interval: %w", err)
	}
	minNumSymbols, err := pcs.tx.GetMinNumSymbols(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get min num symbols: %w", err)
	}
	pricePerSymbol, err := pcs.tx.GetPricePerSymbol(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get price per symbol: %w", err)
	}
	reservationWindow, err := pcs.tx.GetReservationWindow(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation window: %w", err)
	}

	// Initialize config maps
	quorumPaymentConfigs := make(map[core.QuorumID]*core.PaymentQuorumConfig)
	quorumProtocolConfigs := make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig)

	// Populate configs for each quorum
	for _, quorumNumber := range quorumNumbers {
		quorumPaymentConfigs[quorumNumber] = &core.PaymentQuorumConfig{
			OnDemandSymbolsPerSecond:    globalSymbolsPerSecond,
			OnDemandPricePerSymbol:      pricePerSymbol,
			ReservationSymbolsPerSecond: 0, // placeholder
		}

		quorumProtocolConfigs[quorumNumber] = &core.PaymentQuorumProtocolConfig{
			ReservationRateLimitWindow: reservationWindow,
			OnDemandRateLimitWindow:    globalRatePeriodInterval,
			MinNumSymbols:              minNumSymbols,
			OnDemandEnabled:            false, // placeholder
			ReservationAdvanceWindow:   0,     // placeholder
		}
	}

	// Enable on-demand for Quorum 0
	quorumProtocolConfigs[OnDemandQuorumID].OnDemandEnabled = true

	return &PaymentVaultParams{
		OnDemandQuorumNumbers: requiredQuorumNumbers,
		QuorumPaymentConfigs:  quorumPaymentConfigs,
		QuorumProtocolConfigs: quorumProtocolConfigs,
	}, nil
}

// RefreshOnchainPaymentState updates the payment state with current chain data
func (pcs *OnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context) error {
	paymentVaultParams, err := pcs.GetPaymentVaultParams(ctx)
	if err != nil {
		return fmt.Errorf("failed to get payment vault params: %w", err)
	}
	pcs.PaymentVaultParams.Store(paymentVaultParams)

	// Refresh reserved and on-demand payments
	var refreshErr error
	if reservedPaymentsErr := pcs.refreshReservedPayments(ctx); reservedPaymentsErr != nil {
		pcs.logger.Error("failed to refresh reserved payments", "error", reservedPaymentsErr)
		refreshErr = errors.Join(refreshErr, reservedPaymentsErr)
	}

	if ondemandPaymentsErr := pcs.refreshOnDemandPayments(ctx); ondemandPaymentsErr != nil {
		pcs.logger.Error("failed to refresh on-demand payments", "error", ondemandPaymentsErr)
		refreshErr = errors.Join(refreshErr, ondemandPaymentsErr)
	}

	if refreshErr != nil {
		return fmt.Errorf("failed to refresh payment state: %w", refreshErr)
	}

	return nil
}

func (pcs *OnchainPaymentState) refreshReservedPayments(ctx context.Context) error {
	pcs.ReservationsLock.Lock()
	defer pcs.ReservationsLock.Unlock()

	if len(pcs.ReservedPayments) == 0 {
		pcs.logger.Info("No reserved payments to refresh")
		return nil
	}

	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return err
	}

	quorumCount, err := pcs.tx.GetQuorumCount(ctx, blockNumber)
	if err != nil {
		return err
	}
	quorumNumbers := make([]uint8, quorumCount)
	for i := range quorumNumbers {
		quorumNumbers[i] = uint8(i)
	}

	accountIDs := make([]gethcommon.Address, 0, len(pcs.ReservedPayments))
	for accountID := range pcs.ReservedPayments {
		accountIDs = append(accountIDs, accountID)
	}

	// Updated to use Usage Authorization Registry - read reservations for all quorums
	reservedPaymentsByQuorum := make(map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment)
	for _, accountID := range accountIDs {
		newRes, err := pcs.getReservedPayments(ctx, accountID, quorumNumbers)
		if err != nil {
			return fmt.Errorf("failed to get reserved payments: %w", err)
		}
		reservedPaymentsByQuorum[accountID] = newRes
	}
	pcs.ReservedPayments = reservedPaymentsByQuorum
	return nil
}

func (pcs *OnchainPaymentState) refreshOnDemandPayments(ctx context.Context) error {
	pcs.OnDemandLocks.Lock()
	defer pcs.OnDemandLocks.Unlock()

	if len(pcs.OnDemandPayments) == 0 {
		pcs.logger.Info("No on-demand payments to refresh")
		return nil
	}

	accountIDs := make([]gethcommon.Address, 0, len(pcs.OnDemandPayments))
	for accountID := range pcs.OnDemandPayments {
		accountIDs = append(accountIDs, accountID)
	}

	// Updated to use Usage Authorization Registry - read on-demand deposits for all quorums
	onDemandPayments := make(map[gethcommon.Address]*core.OnDemandPayment)
	for _, accountID := range accountIDs {
		deposit, err := pcs.tx.GetUsageAuthOnDemandDeposit(ctx, core.QuorumID(OnDemandQuorumID), accountID)
		if err != nil {
			// Log but continue for other accounts
			pcs.logger.Debug("Failed to get usage auth on-demand deposit", "account", accountID, "err", err)
			continue
		}

		onDemandPayments[accountID] = &core.OnDemandPayment{
			CumulativePayment: deposit,
		}
	}
	pcs.OnDemandPayments = onDemandPayments
	return nil
}

// GetReservedPaymentByAccountAndQuorums retrieves reserved payments for an account across specified quorums
func (pcs *OnchainPaymentState) GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error) {
	pcs.ReservationsLock.RLock()
	// Initialize map if needed
	if _, ok := pcs.ReservedPayments[accountID]; !ok {
		pcs.ReservedPayments[accountID] = make(map[core.QuorumID]*core.ReservedPayment)
	}

	// first read from the periodically refreshed cache
	notCachedQuorums := make([]core.QuorumID, 0)
	accountReservations := make(map[core.QuorumID]*core.ReservedPayment)
	if quorumReservations, ok := pcs.ReservedPayments[accountID]; ok {
		for _, quorumNumber := range quorumNumbers {
			if reservation, cached := quorumReservations[quorumNumber]; cached {
				accountReservations[quorumNumber] = reservation
			} else {
				notCachedQuorums = append(notCachedQuorums, quorumNumber)
			}
		}

		// If notCachedQuorums are empty, return the existing cache
		if len(notCachedQuorums) == 0 {
			pcs.ReservationsLock.RUnlock()
			return accountReservations, nil
		}
	} else {
		// No cached data for this account, need to fetch all quorums
		notCachedQuorums = quorumNumbers
	}
	pcs.ReservationsLock.RUnlock()

	// pulls the chain state using Usage Authorization Registry - only for notCachedQuorums
	newRes, err := pcs.getReservedPayments(ctx, accountID, notCachedQuorums)
	if err != nil {
		return nil, fmt.Errorf("failed to get reserved payments: %w", err)
	}

	pcs.ReservationsLock.Lock()
	defer pcs.ReservationsLock.Unlock()

	// Combine cached and newly fetched results; update cache as well
	for quorumNumber, reservation := range newRes {
		pcs.ReservedPayments[accountID][quorumNumber] = reservation
		accountReservations[quorumNumber] = reservation
	}

	return accountReservations, nil
}

// getReservedPayments retrieves reserved payments for an account across specified quorums directly from the chain
// Zero-valued reservations are not included in the returned map.
func (pcs *OnchainPaymentState) getReservedPayments(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error) {
	allRes := make(map[core.QuorumID]*core.ReservedPayment)
	for _, quorumNumber := range quorumNumbers {
		reservation, err := pcs.tx.GetUsageAuthReservation(ctx, quorumNumber, accountID)
		if err != nil {
			// Log but continue for other quorums
			pcs.logger.Debug("Failed to get usage auth reservation", "account", accountID, "quorum", quorumNumber, "err", err)
			continue
		}

		// Convert to ReservedPayment and check if zero-valued
		reservedPayment := &core.ReservedPayment{
			SymbolsPerSecond: reservation.SymbolsPerSecond,
			StartTimestamp:   reservation.StartTimestamp,
			EndTimestamp:     reservation.EndTimestamp,
		}

		if !eth.IsZeroValuedReservation(reservedPayment) {
			allRes[quorumNumber] = reservedPayment
		}
	}

	return allRes, nil
}

// GetOnDemandPaymentByAccount retrieves on-demand payment information for an account
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	pcs.OnDemandLocks.RLock()
	if payment, ok := pcs.OnDemandPayments[accountID]; ok {
		pcs.OnDemandLocks.RUnlock()
		return payment, nil
	}
	pcs.OnDemandLocks.RUnlock()
	deposit, err := pcs.tx.GetUsageAuthOnDemandDeposit(ctx, OnDemandQuorumID, accountID)
	if err != nil {
		return nil, err
	}

	res := &core.OnDemandPayment{
		CumulativePayment: deposit,
	}

	pcs.OnDemandLocks.Lock()
	pcs.OnDemandPayments[accountID] = res
	pcs.OnDemandLocks.Unlock()

	return res, nil
}

// GetPaymentGlobalParams retrieves all payment vault parameters
func (pcs *OnchainPaymentState) GetPaymentGlobalParams() (*PaymentVaultParams, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		return nil, fmt.Errorf("payment vault params not initialized")
	}
	return params, nil
}

// GetQuorumConfigs retrieves payment and protocol configurations for a specific quorum
func (pvp *PaymentVaultParams) GetQuorumConfigs(quorumID core.QuorumID) (*core.PaymentQuorumConfig, *core.PaymentQuorumProtocolConfig, error) {
	paymentConfig, ok := pvp.QuorumPaymentConfigs[quorumID]
	if !ok {
		return nil, nil, fmt.Errorf("payment config not found for quorum %d", quorumID)
	}
	protocolConfig, ok := pvp.QuorumProtocolConfigs[quorumID]
	if !ok {
		return nil, nil, fmt.Errorf("protocol config not found for quorum %d", quorumID)
	}
	return paymentConfig, protocolConfig, nil
}

// PaymentVaultParamsFromProtobuf converts a protobuf payment vault params to a core payment vault params
func PaymentVaultParamsFromProtobuf(vaultParams *disperser_rpc.PaymentVaultParams) (*PaymentVaultParams, error) {
	if vaultParams == nil {
		return nil, fmt.Errorf("payment vault params cannot be nil")
	}

	if vaultParams.GetQuorumPaymentConfigs() == nil {
		return nil, fmt.Errorf("payment quorum configs cannot be nil")
	}

	if vaultParams.GetQuorumProtocolConfigs() == nil {
		return nil, fmt.Errorf("payment quorum protocol configs cannot be nil")
	}

	// Convert protobuf configs to core types
	quorumPaymentConfigs := make(map[core.QuorumID]*core.PaymentQuorumConfig)
	quorumProtocolConfigs := make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig)

	for quorumID, pbPaymentConfig := range vaultParams.GetQuorumPaymentConfigs() {
		quorumPaymentConfigs[core.QuorumID(quorumID)] = &core.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: pbPaymentConfig.GetReservationSymbolsPerSecond(),
			OnDemandSymbolsPerSecond:    pbPaymentConfig.GetOnDemandSymbolsPerSecond(),
			OnDemandPricePerSymbol:      pbPaymentConfig.GetOnDemandPricePerSymbol(),
		}
	}

	for quorumID, pbProtocolConfig := range vaultParams.GetQuorumProtocolConfigs() {
		quorumProtocolConfigs[core.QuorumID(quorumID)] = &core.PaymentQuorumProtocolConfig{
			MinNumSymbols:              pbProtocolConfig.GetMinNumSymbols(),
			ReservationAdvanceWindow:   pbProtocolConfig.GetReservationAdvanceWindow(),
			ReservationRateLimitWindow: pbProtocolConfig.GetReservationRateLimitWindow(),
			OnDemandRateLimitWindow:    pbProtocolConfig.GetOnDemandRateLimitWindow(),
			OnDemandEnabled:            pbProtocolConfig.GetOnDemandEnabled(),
		}
	}
	// Convert uint32 slice to core.QuorumID slice
	onDemandQuorumNumbers := make([]core.QuorumID, len(vaultParams.GetOnDemandQuorumNumbers()))
	for i, num := range vaultParams.GetOnDemandQuorumNumbers() {
		onDemandQuorumNumbers[i] = core.QuorumID(num)
	}
	return &PaymentVaultParams{
		QuorumPaymentConfigs:  quorumPaymentConfigs,
		QuorumProtocolConfigs: quorumProtocolConfigs,
		OnDemandQuorumNumbers: onDemandQuorumNumbers,
	}, nil
}

// PaymentVaultParamsToProtobuf converts core payment vault params to protobuf format
func (pvp *PaymentVaultParams) PaymentVaultParamsToProtobuf() (*disperser_rpc.PaymentVaultParams, error) {
	if pvp == nil {
		return nil, fmt.Errorf("payment vault params cannot be nil")
	}

	if pvp.QuorumPaymentConfigs == nil {
		return nil, fmt.Errorf("payment quorum configs cannot be nil")
	}

	if pvp.QuorumProtocolConfigs == nil {
		return nil, fmt.Errorf("payment quorum protocol configs cannot be nil")
	}

	quorumPaymentConfigs := make(map[uint32]*disperser_rpc.PaymentQuorumConfig)
	quorumProtocolConfigs := make(map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig)

	for quorumID, paymentConfig := range pvp.QuorumPaymentConfigs {
		quorumPaymentConfigs[uint32(quorumID)] = &disperser_rpc.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: paymentConfig.ReservationSymbolsPerSecond,
			OnDemandSymbolsPerSecond:    paymentConfig.OnDemandSymbolsPerSecond,
			OnDemandPricePerSymbol:      paymentConfig.OnDemandPricePerSymbol,
		}
	}

	for quorumID, protocolConfig := range pvp.QuorumProtocolConfigs {
		quorumProtocolConfigs[uint32(quorumID)] = &disperser_rpc.PaymentQuorumProtocolConfig{
			MinNumSymbols:              protocolConfig.MinNumSymbols,
			ReservationAdvanceWindow:   protocolConfig.ReservationAdvanceWindow,
			ReservationRateLimitWindow: protocolConfig.ReservationRateLimitWindow,
			OnDemandRateLimitWindow:    protocolConfig.OnDemandRateLimitWindow,
			OnDemandEnabled:            protocolConfig.OnDemandEnabled,
		}
	}

	onDemandQuorumNumbers := make([]uint32, len(pvp.OnDemandQuorumNumbers))
	for i, num := range pvp.OnDemandQuorumNumbers {
		onDemandQuorumNumbers[i] = uint32(num)
	}

	return &disperser_rpc.PaymentVaultParams{
		QuorumPaymentConfigs:  quorumPaymentConfigs,
		QuorumProtocolConfigs: quorumProtocolConfigs,
		OnDemandQuorumNumbers: onDemandQuorumNumbers,
	}, nil
}

// ReservationsFromProtobuf converts protobuf reservations to native types
func ReservationsFromProtobuf(pbReservations map[uint32]*disperser_rpc.QuorumReservation) map[core.QuorumID]*core.ReservedPayment {
	if pbReservations == nil {
		return nil
	}

	reservations := make(map[core.QuorumID]*core.ReservedPayment)
	for quorumNumber, reservation := range pbReservations {
		if reservation == nil {
			continue
		}
		quorumID := core.QuorumID(quorumNumber)
		reservations[quorumID] = &core.ReservedPayment{
			SymbolsPerSecond: reservation.GetSymbolsPerSecond(),
			StartTimestamp:   uint64(reservation.GetStartTimestamp()),
			EndTimestamp:     uint64(reservation.GetEndTimestamp()),
		}
	}
	return reservations
}

// CumulativePaymentFromProtobuf converts protobuf payment bytes to *big.Int
func CumulativePaymentFromProtobuf(paymentBytes []byte) *big.Int {
	if paymentBytes == nil {
		return nil
	}
	return new(big.Int).SetBytes(paymentBytes)
}

// ConvertPaymentStateFromProtobuf converts a protobuf GetPaymentStateForAllQuorumsReply to native types
func ConvertPaymentStateFromProtobuf(paymentStateProto *disperser_rpc.GetPaymentStateForAllQuorumsReply) (
	*PaymentVaultParams,
	map[core.QuorumID]*core.ReservedPayment,
	*big.Int,
	*big.Int,
	QuorumPeriodRecords,
	error,
) {
	if paymentStateProto == nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("payment state cannot be nil")
	}

	paymentVaultParams, err := PaymentVaultParamsFromProtobuf(paymentStateProto.GetPaymentVaultParams())
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error converting payment vault params: %w", err)
	}

	reservations := ReservationsFromProtobuf(paymentStateProto.GetReservations())

	cumulativePayment := CumulativePaymentFromProtobuf(paymentStateProto.GetCumulativePayment())
	onchainCumulativePayment := CumulativePaymentFromProtobuf(paymentStateProto.GetOnchainCumulativePayment())

	var periodRecords QuorumPeriodRecords
	if paymentStateProto.GetPeriodRecords() != nil {
		periodRecords = FromProtoRecords(paymentStateProto.GetPeriodRecords())
	}

	return paymentVaultParams, reservations, cumulativePayment, onchainCumulativePayment, periodRecords, nil
}
