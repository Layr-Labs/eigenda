package meterer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
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
	GetQuorumPaymentConfig(quorumID core.QuorumID) (*core.PaymentQuorumConfig, error)
	GetQuorumProtocolConfig(quorumID core.QuorumID) (*core.PaymentQuorumProtocolConfig, error)
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

	// TODO(hopeyen): the construction of quorum configs will be updated with payment vault interface updates
	globalSymbolsPerSecond, err := pcs.tx.GetOnDemandGlobalSymbolsPerSecond(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	globalRatePeriodInterval, err := pcs.tx.GetOnDemandGlobalRatePeriodInterval(ctx, blockNumber)
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
	quorumPaymentConfig := &core.PaymentQuorumConfig{
		OnDemandSymbolsPerSecond: globalSymbolsPerSecond,
		OnDemandPricePerSymbol:   pricePerSymbol,
		// These two fields are not used in the offchain state
		ReservationSymbolsPerSecond: uint64(0),
	}
	quorumProtocolConfig := &core.PaymentQuorumProtocolConfig{
		MinNumSymbols:              minNumSymbols,
		ReservationRateLimitWindow: reservationWindow,
		OnDemandRateLimitWindow:    globalRatePeriodInterval,
		// These two fields are not used in the offchain state
		ReservationAdvanceWindow: uint64(0),
		OnDemandEnabled:          false,
	}
	for _, quorumNumber := range quorumNumbers {
		quorumPaymentConfigs[core.QuorumID(quorumNumber)] = quorumPaymentConfig
		quorumProtocolConfigs[core.QuorumID(quorumNumber)] = quorumProtocolConfig
	}
	// OnDemand is initially only enabled on Quorum 0
	quorumProtocolConfigs[OnDemandQuorumID].OnDemandEnabled = true

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

	// TODO(hopeyen): with payment vault update, this function will take quorum numbers;
	// Currently we just build the same reservation for each quorum
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
	// TODO(hopeyen): use specific quorum IDs from the chain when payment vault is updated
	res, err := pcs.tx.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	pcs.ReservationsLock.Lock()
	if (pcs.ReservedPayments)[accountID] == nil {
		(pcs.ReservedPayments)[accountID] = make(map[core.QuorumID]*core.ReservedPayment)
	}
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

// GetQuorumPaymentConfig safely retrieves a quorum payment config
func (pcs *OnchainPaymentState) GetQuorumPaymentConfig(quorumID core.QuorumID) (*core.PaymentQuorumConfig, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		pcs.logger.Error("PaymentVaultParams is nil")
		return nil, fmt.Errorf("payment vault params is nil")
	}
	config, ok := params.QuorumPaymentConfigs[quorumID]
	if !ok {
		pcs.logger.Error("Quorum payment config not found", "quorumID", quorumID)
		return nil, fmt.Errorf("quorum payment config not found for quorum %d", quorumID)
	}
	return config, nil
}

// GetQuorumProtocolConfig safely retrieves a quorum protocol config
func (pcs *OnchainPaymentState) GetQuorumProtocolConfig(quorumID core.QuorumID) (*core.PaymentQuorumProtocolConfig, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		pcs.logger.Error("PaymentVaultParams is nil")
		return nil, fmt.Errorf("payment vault params is nil")
	}
	config, ok := params.QuorumProtocolConfigs[quorumID]
	if !ok {
		pcs.logger.Error("Quorum protocol config not found", "quorumID", quorumID)
		return nil, fmt.Errorf("quorum protocol config not found for quorum %d", quorumID)
	}
	return config, nil
}

func (pcs *OnchainPaymentState) GetOnDemandGlobalSymbolsPerSecond(quorumID core.QuorumID) uint64 {
	config, err := pcs.GetQuorumPaymentConfig(quorumID)
	if err != nil {
		return 0
	}
	return config.OnDemandSymbolsPerSecond
}

func (pcs *OnchainPaymentState) GetOnDemandGlobalRatePeriodInterval(quorumID core.QuorumID) uint64 {
	config, err := pcs.GetQuorumProtocolConfig(quorumID)
	if err != nil {
		return 0
	}
	return config.OnDemandRateLimitWindow
}

func (pcs *OnchainPaymentState) GetMinNumSymbols(quorumID core.QuorumID) uint64 {
	config, err := pcs.GetQuorumProtocolConfig(quorumID)
	if err != nil {
		return math.MaxUint64
	}
	return config.MinNumSymbols
}

func (pcs *OnchainPaymentState) GetPricePerSymbol(quorumID core.QuorumID) uint64 {
	config, err := pcs.GetQuorumPaymentConfig(quorumID)
	if err != nil {
		return math.MaxUint64
	}
	return config.OnDemandPricePerSymbol
}

func (pcs *OnchainPaymentState) GetReservationWindow(quorumID core.QuorumID) uint64 {
	config, err := pcs.GetQuorumProtocolConfig(quorumID)
	if err != nil {
		return 0
	}
	return config.ReservationRateLimitWindow
}

// return the key of QuorumPaymentConfigs
func (pcs *OnchainPaymentState) GetQuorumNumbers(ctx context.Context) ([]uint8, error) {
	quorumNumbers := make([]uint8, 0, len(pcs.PaymentVaultParams.Load().QuorumPaymentConfigs))
	for quorumNumber := range pcs.PaymentVaultParams.Load().QuorumPaymentConfigs {
		quorumNumbers = append(quorumNumbers, uint8(quorumNumber))
	}
	return quorumNumbers, nil
}
