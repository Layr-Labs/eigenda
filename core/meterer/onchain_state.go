package meterer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// PaymentAccounts (For reservations and on-demand payments)

// OnchainPaymentState is an interface for getting information about the current chain state for payments.
type OnchainPayment interface {
	// Core state refresh and access methods
	RefreshOnchainPaymentState(ctx context.Context) error
	GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error)
	GetQuorumNumbers(ctx context.Context) ([]core.QuorumID, error)

	// Quorum config getters
	GetQuorumPaymentConfig(quorumID core.QuorumID) (*core.PaymentQuorumConfig, error)
	GetQuorumProtocolConfig(quorumID core.QuorumID) (*core.PaymentQuorumProtocolConfig, error)
	GetQuorumPaymentConfigs() map[core.QuorumID]*core.PaymentQuorumConfig
	GetQuorumProtocolConfigs() map[core.QuorumID]*core.PaymentQuorumProtocolConfig

	// On-demand specific getters
	GetOnDemandQuorumNumbers(ctx context.Context) ([]core.QuorumID, error)
	GetOnDemandGlobalSymbolsPerSecond(quorumID core.QuorumID) (uint64, error)
	GetOnDemandGlobalRatePeriodInterval(quorumID core.QuorumID) (uint64, error)

	// Payment parameter getters
	GetMinNumSymbols(quorumID core.QuorumID) (uint64, error)
	GetPricePerSymbol(quorumID core.QuorumID) (uint64, error)
	GetReservationWindow(quorumID core.QuorumID) (uint64, error)
}

var _ OnchainPayment = (*OnchainPaymentState)(nil)

// OnchainPaymentState manages the state of on-chain payments including reservations and on-demand payments
type OnchainPaymentState struct {
	tx     *eth.Reader
	logger logging.Logger

	// Maps to store payment state
	ReservedPayments map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment
	OnDemandPayments map[gethcommon.Address]*core.OnDemandPayment

	// Locks for thread-safe access
	ReservationsLock sync.RWMutex
	OnDemandLocks    sync.RWMutex

	// Atomic pointer to store payment vault parameters
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

// GetPaymentVaultParams retrieves the current payment vault parameters from the chain
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

	// Get global parameters
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

// GetReservedPaymentByAccountAndQuorums retrieves reserved payments for an account across specified quorums
func (pcs *OnchainPaymentState) GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error) {
	pcs.ReservationsLock.RLock()
	if quorumReservations, ok := pcs.ReservedPayments[accountID]; ok {
		// Check if all quorums are present
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

	// Fetch from chain if not found in cache
	res, err := pcs.tx.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reserved payment: %w", err)
	}

	pcs.ReservationsLock.Lock()
	defer pcs.ReservationsLock.Unlock()

	// Initialize map if needed
	if _, ok := pcs.ReservedPayments[accountID]; !ok {
		pcs.ReservedPayments[accountID] = make(map[core.QuorumID]*core.ReservedPayment)
	}

	// Update cache with new data
	for _, quorumNumber := range quorumNumbers {
		if reservation, ok := res[quorumNumber]; ok {
			pcs.ReservedPayments[accountID][quorumNumber] = reservation
		}
	}

	return res, nil
}

// GetOnDemandPaymentByAccount retrieves on-demand payment information for an account
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	pcs.OnDemandLocks.RLock()
	if payment, ok := pcs.OnDemandPayments[accountID]; ok {
		pcs.OnDemandLocks.RUnlock()
		return payment, nil
	}
	pcs.OnDemandLocks.RUnlock()

	// Fetch from chain if not found in cache
	res, err := pcs.tx.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get on-demand payment: %w", err)
	}

	pcs.OnDemandLocks.Lock()
	pcs.OnDemandPayments[accountID] = res
	pcs.OnDemandLocks.Unlock()

	return res, nil
}

// GetQuorumPaymentConfig retrieves payment configuration for a specific quorum
func (pcs *OnchainPaymentState) GetQuorumPaymentConfig(quorumID core.QuorumID) (*core.PaymentQuorumConfig, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		return nil, fmt.Errorf("payment vault params not initialized")
	}
	return params.GetQuorumPaymentConfig(quorumID)
}

// GetQuorumProtocolConfig retrieves protocol configuration for a specific quorum
func (pcs *OnchainPaymentState) GetQuorumProtocolConfig(quorumID core.QuorumID) (*core.PaymentQuorumProtocolConfig, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		return nil, fmt.Errorf("payment vault params not initialized")
	}
	return params.GetQuorumProtocolConfig(quorumID)
}

// GetQuorumPaymentConfigs retrieves all quorum payment configurations
func (pcs *OnchainPaymentState) GetQuorumPaymentConfigs() map[core.QuorumID]*core.PaymentQuorumConfig {
	return pcs.PaymentVaultParams.Load().QuorumPaymentConfigs
}

// GetQuorumProtocolConfigs retrieves all quorum protocol configurations
func (pcs *OnchainPaymentState) GetQuorumProtocolConfigs() map[core.QuorumID]*core.PaymentQuorumProtocolConfig {
	return pcs.PaymentVaultParams.Load().QuorumProtocolConfigs
}

// GetOnDemandQuorumNumbers retrieves the list of quorums enabled for on-demand payments
func (pcs *OnchainPaymentState) GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error) {
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}

	quorumNumbers, err := pcs.tx.GetRequiredQuorumNumbers(ctx, blockNumber)
	if err != nil {
		// Fallback to cached value if chain read fails
		pcs.logger.Warn("failed to get required quorum numbers, using cached value", "error", err)
		params := pcs.PaymentVaultParams.Load()
		if params == nil {
			return nil, fmt.Errorf("failed to get quorum numbers and no cached params available")
		}
		return params.OnDemandQuorumNumbers, nil
	}
	return quorumNumbers, nil
}

// GetOnDemandGlobalSymbolsPerSecond retrieves the global symbols per second rate for on-demand payments
func (pcs *OnchainPaymentState) GetOnDemandGlobalSymbolsPerSecond(quorumID core.QuorumID) (uint64, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		return 0, fmt.Errorf("payment vault params not initialized")
	}
	return params.GetOnDemandGlobalSymbolsPerSecond(quorumID)
}

// GetOnDemandGlobalRatePeriodInterval retrieves the rate period interval for on-demand payments
func (pcs *OnchainPaymentState) GetOnDemandGlobalRatePeriodInterval(quorumID core.QuorumID) (uint64, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		return 0, fmt.Errorf("payment vault params not initialized")
	}
	return params.GetOnDemandGlobalRatePeriodInterval(quorumID)
}

// GetMinNumSymbols retrieves the minimum number of symbols required for a quorum
func (pcs *OnchainPaymentState) GetMinNumSymbols(quorumID core.QuorumID) (uint64, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		return math.MaxUint64, fmt.Errorf("payment vault params not initialized")
	}
	return params.GetMinNumSymbols(quorumID)
}

// GetPricePerSymbol retrieves the price per symbol for a quorum
func (pcs *OnchainPaymentState) GetPricePerSymbol(quorumID core.QuorumID) (uint64, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		return math.MaxUint64, fmt.Errorf("payment vault params not initialized")
	}
	return params.GetPricePerSymbol(quorumID)
}

// GetReservationWindow retrieves the reservation window duration for a quorum
func (pcs *OnchainPaymentState) GetReservationWindow(quorumID core.QuorumID) (uint64, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		return 0, fmt.Errorf("payment vault params not initialized")
	}
	return params.GetReservationWindow(quorumID)
}

// GetQuorumNumbers retrieves all quorum numbers tracked by the payment system
func (pcs *OnchainPaymentState) GetQuorumNumbers(ctx context.Context) ([]core.QuorumID, error) {
	params := pcs.PaymentVaultParams.Load()
	if params == nil {
		return nil, fmt.Errorf("payment vault params not initialized")
	}
	return params.GetQuorumNumbers(), nil
}

// GetQuorumPaymentConfig retrieves payment configuration for a specific quorum
func (pvp *PaymentVaultParams) GetQuorumPaymentConfig(quorumID core.QuorumID) (*core.PaymentQuorumConfig, error) {
	config, ok := pvp.QuorumPaymentConfigs[quorumID]
	if !ok {
		return nil, fmt.Errorf("payment config not found for quorum %d", quorumID)
	}
	return config, nil
}

// GetQuorumProtocolConfig retrieves protocol configuration for a specific quorum
func (pvp *PaymentVaultParams) GetQuorumProtocolConfig(quorumID core.QuorumID) (*core.PaymentQuorumProtocolConfig, error) {
	config, ok := pvp.QuorumProtocolConfigs[quorumID]
	if !ok {
		return nil, fmt.Errorf("protocol config not found for quorum %d", quorumID)
	}
	return config, nil
}

// GetOnDemandGlobalSymbolsPerSecond retrieves the global symbols per second rate for on-demand payments
func (pvp *PaymentVaultParams) GetOnDemandGlobalSymbolsPerSecond(quorumID core.QuorumID) (uint64, error) {
	config, err := pvp.GetQuorumPaymentConfig(quorumID)
	if err != nil {
		return 0, err
	}
	return config.OnDemandSymbolsPerSecond, nil
}

// GetOnDemandGlobalRatePeriodInterval retrieves the rate period interval for on-demand payments
func (pvp *PaymentVaultParams) GetOnDemandGlobalRatePeriodInterval(quorumID core.QuorumID) (uint64, error) {
	config, err := pvp.GetQuorumProtocolConfig(quorumID)
	if err != nil {
		return 0, err
	}
	return config.OnDemandRateLimitWindow, nil
}

// GetMinNumSymbols retrieves the minimum number of symbols required for a quorum
func (pvp *PaymentVaultParams) GetMinNumSymbols(quorumID core.QuorumID) (uint64, error) {
	config, err := pvp.GetQuorumProtocolConfig(quorumID)
	if err != nil {
		return 0, err
	}
	return config.MinNumSymbols, nil
}

// GetPricePerSymbol retrieves the price per symbol for a quorum
func (pvp *PaymentVaultParams) GetPricePerSymbol(quorumID core.QuorumID) (uint64, error) {
	config, err := pvp.GetQuorumPaymentConfig(quorumID)
	if err != nil {
		return 0, err
	}
	return config.OnDemandPricePerSymbol, nil
}

// GetReservationWindow retrieves the reservation window duration for a quorum
func (pvp *PaymentVaultParams) GetReservationWindow(quorumID core.QuorumID) (uint64, error) {
	config, err := pvp.GetQuorumProtocolConfig(quorumID)
	if err != nil {
		return 0, err
	}
	return config.ReservationRateLimitWindow, nil
}

// GetQuorumNumbers retrieves all quorum numbers tracked by the payment system
func (pvp *PaymentVaultParams) GetQuorumNumbers() []core.QuorumID {
	quorumNumbers := make([]core.QuorumID, 0, len(pvp.QuorumPaymentConfigs))
	for quorumNumber := range pvp.QuorumPaymentConfigs {
		quorumNumbers = append(quorumNumbers, core.QuorumID(quorumNumber))
	}
	return quorumNumbers
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

func (pvp *PaymentVaultParams) GetConfigs(quorumNumber core.QuorumID) (*core.PaymentQuorumConfig, *core.PaymentQuorumProtocolConfig, error) {
	if pvp == nil {
		return nil, nil, fmt.Errorf("payment vault params is nil")
	}
	paymentQuorumConfig, ok := pvp.QuorumPaymentConfigs[quorumNumber]
	if !ok {
		return nil, nil, fmt.Errorf("payment quorum config not found for quorum %d", quorumNumber)
	}
	protocolConfig, ok := pvp.QuorumProtocolConfigs[quorumNumber]
	if !ok {
		return nil, nil, fmt.Errorf("payment quorum protocol config not found for quorum %d", quorumNumber)
	}
	return paymentQuorumConfig, protocolConfig, nil
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
