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
	GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error)
	GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error)
	GetOnDemandGlobalSymbolsPerSecond(quorumID core.QuorumID) uint64
	GetOnDemandGlobalRatePeriodInterval(quorumID core.QuorumID) uint64
	GetMinNumSymbols(quorumID core.QuorumID) uint64
	GetPricePerSymbol(quorumID core.QuorumID) uint64
	GetReservationWindow(quorumID core.QuorumID) uint64
	// Return all the quorum numbers tracked by QuorumPaymentConfigs
	GetQuorumNumbers(ctx context.Context) ([]uint8, error)
}

var _ OnchainPayment = (*OnchainPaymentState)(nil)

type OnchainPaymentState struct {
	tx     *eth.Reader
	logger logging.Logger

	ReservedPayments map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment
	OnDemandPayments map[gethcommon.Address]*core.OnDemandPayment

	ReservationsLock sync.RWMutex
	OnDemandLocks    sync.RWMutex

	PaymentVaultParams atomic.Pointer[PaymentVaultParams]
}

type PaymentVaultParams struct {
	QuorumPaymentConfigs  map[core.QuorumID]*core.PaymentQuorumConfig
	QuorumProtocolConfigs map[core.QuorumID]*core.PaymentQuorumProtocolConfig
	OnDemandQuorumNumbers []uint8
}

func NewOnchainPaymentState(ctx context.Context, tx *eth.Reader, logger logging.Logger) (*OnchainPaymentState, error) {
	state := OnchainPaymentState{
		tx:                 tx,
		logger:             logger.With("component", "OnchainPaymentState"),
		ReservedPayments:   make(map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment),
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
	requiredQuorumNumbers, err := pcs.tx.GetRequiredQuorumNumbers(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	quorumCount, err := pcs.tx.GetQuorumCount(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	quorumNumbers := make([]uint8, quorumCount)
	for i := range quorumNumbers {
		quorumNumbers[i] = uint8(i)
	}

	// TODO(hopeyen): these will be replaced in a later PR with payment vault interface updates
	globalSymbolsPerSecond, err := pcs.tx.GetOnDemandGlobalSymbolsPerSecond(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	globalRatePeriodInterval, err := pcs.tx.GetGlobalRatePeriodInterval(ctx, blockNumber)
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

	quorumPaymentConfigs := make(map[core.QuorumID]*core.PaymentQuorumConfig)
	quorumProtocolConfigs := make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig)
	// TODO(hopeyen): update to quorum specific values when payment vault is updated
	quorumPaymentConfig := &core.PaymentQuorumConfig{
		OnDemandSymbolsPerSecond:    globalSymbolsPerSecond,
		OnDemandPricePerSymbol:      pricePerSymbol,
		ReservationSymbolsPerSecond: uint64(0),
	}
	quorumProtocolConfig := &core.PaymentQuorumProtocolConfig{
		MinNumSymbols:              minNumSymbols,
		ReservationAdvanceWindow:   reservationWindow,
		ReservationRateLimitWindow: uint64(0),
		// OnDemand is initially only enabled on Quorum 0
		OnDemandRateLimitWindow: globalRatePeriodInterval,
		OnDemandEnabled:         false,
	}
	for _, quorumNumber := range quorumNumbers {
		quorumPaymentConfigs[core.QuorumID(quorumNumber)] = quorumPaymentConfig
		quorumProtocolConfigs[core.QuorumID(quorumNumber)] = quorumProtocolConfig
	}

	return &PaymentVaultParams{
		OnDemandQuorumNumbers: requiredQuorumNumbers,
		QuorumPaymentConfigs:  quorumPaymentConfigs,
		QuorumProtocolConfigs: quorumProtocolConfigs,
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

	// TODO(hopeyen): get all quorum numbers and use quorum specific calls when there's updated payment vault interface
	reservedPayments, err := pcs.tx.GetReservedPayments(ctx, accountIDs)
	if err != nil {
		return err
	}

	reservedPaymentsByQuorum := make(map[gethcommon.Address]map[uint8]*core.ReservedPayment)
	for accountID, payments := range reservedPayments {
		reservedPaymentsByQuorum[accountID] = make(map[uint8]*core.ReservedPayment)
		for quorumNumber, reservation := range payments {
			reservedPaymentsByQuorum[accountID][uint8(quorumNumber)] = reservation
		}
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

	onDemandPayments, err := pcs.tx.GetOnDemandPayments(ctx, accountIDs)
	if err != nil {
		return err
	}
	pcs.OnDemandPayments = onDemandPayments
	return nil
}

// GetReservedPaymentByAccountAndQuorums returns a pointer to the active reservation for the given account ID; no writes will be made to the reservation
func (pcs *OnchainPaymentState) GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error) {
	pcs.ReservationsLock.RLock()
	// defer pcs.ReservationsLock.RUnlock()
	if quorumReservations, ok := (pcs.ReservedPayments)[accountID]; ok {
		// Check if all the quorums are present; pull the chain state if not
		allFound := true
		for _, quorumNumber := range quorumNumbers {
			if _, ok := quorumReservations[quorumNumber]; !ok {
				allFound = false
				break
			}
		}
		if allFound {
			pcs.ReservationsLock.RUnlock()
			return quorumReservations, nil
		}
	}
	pcs.ReservationsLock.RUnlock()

	// pulls the chain state
	// TODO(hopeyen): update this to be pulling specific quorum IDs from the chain when payment vault is updated
	res, err := pcs.tx.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	pcs.ReservationsLock.Lock()
	// update specific quorum reservations
	for _, quorumNumber := range quorumNumbers {
		if _, ok := res[quorumNumber]; ok {
			(pcs.ReservedPayments)[accountID][quorumNumber] = res[quorumNumber]
		}
	}
	pcs.ReservationsLock.Unlock()

	return res, nil
}

// GetReservedPaymentByAccount returns a pointer to all quorums' reservation for the given account ID; no writes will be made to the reservation
func (pcs *OnchainPaymentState) GetReservedPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (map[core.QuorumID]*core.ReservedPayment, error) {
	pcs.ReservationsLock.RLock()
	if reservation, ok := (pcs.ReservedPayments)[accountID]; ok {
		pcs.ReservationsLock.RUnlock()
		return reservation, nil
	}
	pcs.ReservationsLock.RUnlock()

	// pulls the chain state
	res, err := pcs.tx.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	reservedPaymentsByQuorum := make(map[core.QuorumID]*core.ReservedPayment)
	for quorumNumber, reservation := range res {
		reservedPaymentsByQuorum[uint8(quorumNumber)] = reservation
	}
	pcs.ReservationsLock.Lock()
	(pcs.ReservedPayments)[accountID] = reservedPaymentsByQuorum
	pcs.ReservationsLock.Unlock()

	return reservedPaymentsByQuorum, nil
}

// GetOnDemandPaymentByAccount returns a pointer to the on-demand payment for the given account ID; no writes will be made to the payment
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	pcs.OnDemandLocks.RLock()
	if payment, ok := (pcs.OnDemandPayments)[accountID]; ok {
		pcs.OnDemandLocks.RUnlock()
		return payment, nil
	}
	pcs.OnDemandLocks.RUnlock()

	// pulls the chain state
	res, err := pcs.tx.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

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

func (pcs *OnchainPaymentState) GetOnDemandGlobalSymbolsPerSecond(quorumID core.QuorumID) uint64 {
	return pcs.PaymentVaultParams.Load().QuorumPaymentConfigs[quorumID].OnDemandSymbolsPerSecond
}

func (pcs *OnchainPaymentState) GetOnDemandGlobalRatePeriodInterval(quorumID core.QuorumID) uint64 {
	return pcs.PaymentVaultParams.Load().QuorumProtocolConfigs[quorumID].OnDemandRateLimitWindow
}

func (pcs *OnchainPaymentState) GetMinNumSymbols(quorumID core.QuorumID) uint64 {
	return pcs.PaymentVaultParams.Load().QuorumProtocolConfigs[quorumID].MinNumSymbols
}

func (pcs *OnchainPaymentState) GetPricePerSymbol(quorumID core.QuorumID) uint64 {
	return pcs.PaymentVaultParams.Load().QuorumPaymentConfigs[quorumID].OnDemandPricePerSymbol
}

func (pcs *OnchainPaymentState) GetReservationWindow(quorumID core.QuorumID) uint64 {
	return pcs.PaymentVaultParams.Load().QuorumProtocolConfigs[quorumID].ReservationRateLimitWindow
}

// return the key of QuorumPaymentConfigs
func (pcs *OnchainPaymentState) GetQuorumNumbers(ctx context.Context) ([]uint8, error) {
	quorumNumbers := make([]uint8, 0, len(pcs.PaymentVaultParams.Load().QuorumPaymentConfigs))
	for quorumNumber := range pcs.PaymentVaultParams.Load().QuorumPaymentConfigs {
		quorumNumbers = append(quorumNumbers, uint8(quorumNumber))
	}
	return quorumNumbers, nil
}
