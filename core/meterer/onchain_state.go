package meterer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// PaymentAccounts (For reservations and on-demand payments)

// OnchainPaymentState is an interface for getting information about the current chain state for payments.
type OnchainPayment interface {
	RefreshOnchainPaymentState(ctx context.Context) error

	// Reservation payment methods
	GetReservedPaymentByAccountAndQuorum(ctx context.Context, accountID gethcommon.Address, quorumId uint8) (*core.ReservedPayment, error)
	GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumIds []uint8) (map[uint8]*core.ReservedPayment, error)

	// On-demand payment methods
	GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error)
	GetOnDemandPaymentByAccountAndQuorum(ctx context.Context, accountID gethcommon.Address, quorumId uint64) (*core.OnDemandPayment, error)

	// Configuration methods
	GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error)
	GetQuorumPaymentConfig(ctx context.Context, quorumId uint64) (*core.QuorumConfig, error)
	GetQuorumProtocolConfig(ctx context.Context, quorumId uint64) (*core.QuorumProtocolConfig, error)

	// Legacy global parameter methods
	GetOnDemandSymbolsPerSecond() uint64
	GetOnDemandRatePeriodInterval() uint64
	GetMinNumSymbols() uint64
	GetPricePerSymbol() uint64
	GetReservationWindow() uint64
}

var _ OnchainPayment = (*OnchainPaymentState)(nil)

type OnchainPaymentState struct {
	tx     *eth.Reader
	logger logging.Logger

	ReservedPayments map[gethcommon.Address]map[uint8]*core.ReservedPayment
	OnDemandPayments map[gethcommon.Address]*core.OnDemandPayment

	ReservationsLock sync.RWMutex
	OnDemandLocks    sync.RWMutex

	PaymentVaultParams atomic.Pointer[PaymentVaultParams]
}

type PaymentVaultParams struct {
	GlobalSymbolsPerSecond   uint64
	GlobalRatePeriodInterval uint64
	MinNumSymbols            uint64
	PricePerSymbol           uint64
	ReservationWindow        uint64
	OnDemandQuorumNumbers    []uint8
}

func NewOnchainPaymentState(ctx context.Context, tx *eth.Reader, logger logging.Logger) (*OnchainPaymentState, error) {
	state := OnchainPaymentState{
		tx:                 tx,
		logger:             logger.With("component", "OnchainPaymentState"),
		ReservedPayments:   make(map[gethcommon.Address]map[uint8]*core.ReservedPayment),
		OnDemandPayments:   make(map[gethcommon.Address]*core.OnDemandPayment),
		PaymentVaultParams: atomic.Pointer[PaymentVaultParams]{},
	}

	paymentVaultParams, err := state.GetPaymentVaultParams(ctx)
	if err != nil {
		return nil, err
	}

	state.PaymentVaultParams.Store(paymentVaultParams)

	return &state, nil
}

func (pcs *OnchainPaymentState) GetPaymentVaultParams(ctx context.Context) (*PaymentVaultParams, error) {
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	quorumNumbers, err := pcs.tx.GetRequiredQuorumNumbers(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	globalSymbolsPerSecond, err := pcs.tx.GetOnDemandSymbolsPerSecond(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	globalRatePeriodInterval, err := pcs.tx.GetOnDemandRatePeriodInterval(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	minNumSymbols, err := pcs.tx.GetMinNumSymbols(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	pricePerSymbol, err := pcs.tx.GetPricePerSymbol(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	reservationWindow, err := pcs.tx.GetReservationWindow(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	return &PaymentVaultParams{
		OnDemandQuorumNumbers:    quorumNumbers,
		GlobalSymbolsPerSecond:   globalSymbolsPerSecond,
		GlobalRatePeriodInterval: globalRatePeriodInterval,
		MinNumSymbols:            minNumSymbols,
		PricePerSymbol:           pricePerSymbol,
		ReservationWindow:        reservationWindow,
	}, nil
}

// RefreshOnchainPaymentState returns the current onchain payment state
func (pcs *OnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context) error {
	paymentVaultParams, err := pcs.GetPaymentVaultParams(ctx)
	if err != nil {
		return err
	}
	// These parameters should be rarely updated, but we refresh them anyway
	pcs.PaymentVaultParams.Store(paymentVaultParams)

	var refreshErr error
	if reservedPaymentsErr := pcs.refreshReservedPayments(ctx); reservedPaymentsErr != nil {
		pcs.logger.Error("failed to refresh reserved payments", "error", reservedPaymentsErr)
		refreshErr = errors.Join(refreshErr, reservedPaymentsErr)
	}

	if ondemandPaymentsErr := pcs.refreshOnDemandPayments(ctx); ondemandPaymentsErr != nil {
		pcs.logger.Error("failed to refresh on-demand payments", "error", ondemandPaymentsErr)
		refreshErr = errors.Join(refreshErr, ondemandPaymentsErr)
	}

	return refreshErr
}

func (pcs *OnchainPaymentState) refreshReservedPayments(ctx context.Context) error {
	pcs.ReservationsLock.Lock()
	defer pcs.ReservationsLock.Unlock()

	if len(pcs.ReservedPayments) == 0 {
		pcs.logger.Info("No reserved payments to refresh")
		return nil
	}

	accountIDs := make([]gethcommon.Address, 0, len(pcs.ReservedPayments))
	quorumIdsMap := make(map[uint8]struct{})

	for accountID, quorums := range pcs.ReservedPayments {
		accountIDs = append(accountIDs, accountID)
		for quorumId := range quorums {
			quorumIdsMap[quorumId] = struct{}{}
		}
	}

	quorumIds := make([]uint8, 0, len(quorumIdsMap))
	for quorumId := range quorumIdsMap {
		quorumIds = append(quorumIds, quorumId)
	}

	reservedPaymentsMap, err := pcs.tx.GetReservedPayments(ctx, accountIDs, quorumIds)
	if err != nil {
		return err
	}

	// Update the cache with first reservation found for each account
	newReservedPayments := make(map[gethcommon.Address]map[uint8]*core.ReservedPayment)
	for accountID, quorumMap := range reservedPaymentsMap {
		if len(quorumMap) > 0 {
			newReservedPayments[accountID] = quorumMap
		}
	}

	pcs.ReservedPayments = newReservedPayments
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

	onDemandPayments, err := pcs.tx.GetOnDemandPayments(ctx, accountIDs)
	if err != nil {
		return err
	}
	pcs.OnDemandPayments = onDemandPayments
	return nil
}

// GetReservedPaymentByAccountAndQuorum returns a pointer to the active reservation for the given account ID and quorum
func (pcs *OnchainPaymentState) GetReservedPaymentByAccountAndQuorum(ctx context.Context, accountID gethcommon.Address, quorumId uint8) (*core.ReservedPayment, error) {
	// Check cache
	pcs.ReservationsLock.RLock()
	cachedReservations, hasCachedAccount := (pcs.ReservedPayments)[accountID]
	if hasCachedAccount && cachedReservations != nil {
		pcs.ReservationsLock.RUnlock()
		return cachedReservations[quorumId], nil
	}
	pcs.ReservationsLock.RUnlock()

	// Pull the chain state
	res, err := pcs.tx.GetReservedPaymentByAccountAndQuorum(ctx, accountID, quorumId)
	if err != nil {
		return nil, err
	}

	// Update cache
	pcs.ReservationsLock.Lock()
	(pcs.ReservedPayments)[accountID][quorumId] = res
	pcs.ReservationsLock.Unlock()

	return res, nil
}

// GetReservedPaymentByAccountAndQuorums returns a map of quorum ID to ReservedPayment for each quorum.
// All requested quorumIds must be gathered before returning, either from cache if available
// or pulled from the chain and cached.
func (pcs *OnchainPaymentState) GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumIds []uint8) (map[uint8]*core.ReservedPayment, error) {
	if len(quorumIds) == 0 {
		return make(map[uint8]*core.ReservedPayment), nil
	}

	result := make(map[uint8]*core.ReservedPayment)
	var quorumsToFetch []uint8

	// First check if we have any quorums already cached per-quorum
	pcs.ReservationsLock.RLock()
	// Check for per-quorum cached reservations
	cachedAccount, hasCachedAccount := (pcs.ReservedPayments)[accountID]
	pcs.ReservationsLock.RUnlock()

	// For each requested quorum, either use the cached value or mark for fetching
	for _, quorumId := range quorumIds {
		if hasCachedAccount && cachedAccount != nil {
			result[quorumId] = cachedAccount[quorumId]
		} else {
			quorumsToFetch = append(quorumsToFetch, quorumId)
		}
	}

	// If we already have all quorums cached, return immediately
	if len(quorumsToFetch) == 0 {
		return result, nil
	}

	// Pull the chain state for the quorums we need to fetch
	chainResult, err := pcs.tx.GetReservedPaymentsByAccountAndQuorums(ctx, accountID, quorumIds)
	if err != nil {
		return nil, err
	}

	// Update cache with fetched reservations
	pcs.ReservationsLock.Lock()
	for quorumId, reservation := range chainResult {
		if reservation != nil {
			result[quorumId] = reservation
		}
	}
	pcs.ReservationsLock.Unlock()

	// Verify that we have all requested quorums
	for _, quorumId := range quorumIds {
		if _, ok := result[quorumId]; !ok {
			return nil, fmt.Errorf("failed to get reservation for quorum %d", quorumId)
		}
	}

	return result, nil
}

// GetOnDemandPaymentByAccount returns a pointer to the on-demand payment for the given account ID; no writes will be made to the payment
// Uses quorum 0 for backwards compatibility
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	return pcs.GetOnDemandPaymentByAccountAndQuorum(ctx, accountID, 0)
}

// GetOnDemandPaymentByAccountAndQuorum returns a pointer to the on-demand payment for the given account ID and quorum; no writes will be made to the payment
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccountAndQuorum(ctx context.Context, accountID gethcommon.Address, quorumId uint64) (*core.OnDemandPayment, error) {
	// Check if we have a cached value
	pcs.OnDemandLocks.RLock()
	if payment, ok := (pcs.OnDemandPayments)[accountID]; ok {
		pcs.OnDemandLocks.RUnlock()
		return payment, nil
	}
	pcs.OnDemandLocks.RUnlock()

	// Pull the chain state
	res, err := pcs.tx.GetOnDemandPaymentByAccountAndQuorum(ctx, accountID, quorumId)
	if err != nil {
		return nil, err
	}

	// Update cache
	pcs.OnDemandLocks.Lock()
	(pcs.OnDemandPayments)[accountID] = res
	pcs.OnDemandLocks.Unlock()

	return res, nil
}

func (pcs *OnchainPaymentState) GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error) {
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	quorumNumbers, err := pcs.tx.GetRequiredQuorumNumbers(ctx, blockNumber)
	if err != nil {
		// On demand required quorum is unlikely to change, so we are comfortable using the cached value
		// in case the contract read fails
		log.Println("Failed to get required quorum numbers, read from cache", "error", err)
		params := pcs.PaymentVaultParams.Load()
		if params == nil {
			log.Println("Failed to get required quorum numbers and no cached params")
			return nil, fmt.Errorf("failed to get required quorum numbers and no cached params")
		}
		// params.OnDemandQuorumNumbers could be empty if set by the protocol
		return params.OnDemandQuorumNumbers, nil
	}
	return quorumNumbers, nil
}

// GetQuorumPaymentConfig retrieves the payment configuration for a specific quorum
func (pcs *OnchainPaymentState) GetQuorumPaymentConfig(ctx context.Context, quorumId uint64) (*core.QuorumConfig, error) {
	return pcs.tx.GetQuorumPaymentConfig(ctx, quorumId)
}

// GetQuorumProtocolConfig retrieves the protocol configuration for a specific quorum
func (pcs *OnchainPaymentState) GetQuorumProtocolConfig(ctx context.Context, quorumId uint64) (*core.QuorumProtocolConfig, error) {
	return pcs.tx.GetQuorumProtocolConfig(ctx, quorumId)
}

func (pcs *OnchainPaymentState) GetOnDemandSymbolsPerSecond() uint64 {
	return pcs.PaymentVaultParams.Load().GlobalSymbolsPerSecond
}

func (pcs *OnchainPaymentState) GetOnDemandRatePeriodInterval() uint64 {
	return pcs.PaymentVaultParams.Load().GlobalRatePeriodInterval
}

func (pcs *OnchainPaymentState) GetMinNumSymbols() uint64 {
	return pcs.PaymentVaultParams.Load().MinNumSymbols
}

func (pcs *OnchainPaymentState) GetPricePerSymbol() uint64 {
	return pcs.PaymentVaultParams.Load().PricePerSymbol
}

func (pcs *OnchainPaymentState) GetReservationWindow() uint64 {
	return pcs.PaymentVaultParams.Load().ReservationWindow
}
