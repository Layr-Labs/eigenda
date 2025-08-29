package vault

import (
	"context"
	"errors"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	paymentvault "github.com/Layr-Labs/eigenda/contracts/bindings/PaymentVault"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// PaymentVault provides access to payment vault contract operations
type PaymentVault struct {
	ethClient           common.EthClient
	logger              logging.Logger
	paymentVaultBinding *paymentvault.ContractPaymentVault
}

// NewPaymentVault creates a new PaymentVault instance
func NewPaymentVault(
	logger logging.Logger,
	ethClient common.EthClient,
	paymentVaultAddress gethcommon.Address,
) (*PaymentVault, error) {
	if ethClient == nil {
		return nil, errors.New("ethClient cannot be nil")
	}

	paymentVaultBinding, err := paymentvault.NewContractPaymentVault(paymentVaultAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &PaymentVault{
		ethClient:           ethClient,
		logger:              logger.With("component", "PaymentVault"),
		paymentVaultBinding: paymentVaultBinding,
	}, nil
}

// GetTotalDeposits retrieves on-demand payment information for multiple accounts
func (pv *PaymentVault) GetTotalDeposits(ctx context.Context, accountIDs []gethcommon.Address) (map[gethcommon.Address]*big.Int, error) {
	if pv.paymentVaultBinding == nil {
		return nil, errors.New("payment vault not deployed")
	}
	paymentsMap := make(map[gethcommon.Address]*big.Int)
	payments, err := pv.paymentVaultBinding.GetOnDemandTotalDeposits(&bind.CallOpts{
		Context: ctx}, accountIDs)
	if err != nil {
		return nil, err
	}

	// since payments are returned in the same order as the accountIDs, we can directly map them
	for i, payment := range payments {
		if payment.Cmp(big.NewInt(0)) == 0 {
			pv.logger.Warn("failed to get on demand payment for account", "account", accountIDs[i])
			continue
		}
		paymentsMap[accountIDs[i]] = payment
	}

	return paymentsMap, nil
}

func (pv *PaymentVault) GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	if pv.paymentVaultBinding == nil {
		return nil, errors.New("payment vault not deployed")
	}
	onDemandPayment, err := pv.paymentVaultBinding.GetOnDemandTotalDeposit(&bind.CallOpts{
		Context: ctx,
	}, accountID)
	if err != nil {
		return nil, err
	}
	if onDemandPayment.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("ondemand payment does not exist for given account")
	}
	return onDemandPayment, nil
}

// GetGlobalSymbolsPerSecond retrieves the global symbols per second parameter
func (pv *PaymentVault) GetGlobalSymbolsPerSecond(ctx context.Context, blockNumber uint32) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	globalSymbolsPerSecond, err := pv.paymentVaultBinding.GlobalSymbolsPerPeriod(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return globalSymbolsPerSecond, nil
}

// GetGlobalRatePeriodInterval retrieves the global rate period interval parameter
func (pv *PaymentVault) GetGlobalRatePeriodInterval(ctx context.Context, blockNumber uint32) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	globalRateBinInterval, err := pv.paymentVaultBinding.GlobalRatePeriodInterval(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return globalRateBinInterval, nil
}

// GetMinNumSymbols retrieves the minimum number of symbols parameter
func (pv *PaymentVault) GetMinNumSymbols(ctx context.Context, blockNumber uint32) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	minNumSymbols, err := pv.paymentVaultBinding.MinNumSymbols(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return minNumSymbols, nil
}

// GetPricePerSymbol retrieves the price per symbol parameter
func (pv *PaymentVault) GetPricePerSymbol(ctx context.Context, blockNumber uint32) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	pricePerSymbol, err := pv.paymentVaultBinding.PricePerSymbol(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return pricePerSymbol, nil
}

// GetReservationWindow retrieves the reservation window parameter
func (pv *PaymentVault) GetReservationWindow(ctx context.Context, blockNumber uint32) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	reservationWindow, err := pv.paymentVaultBinding.ReservationPeriodInterval(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return reservationWindow, nil
}

// GetCurrentBlockNumber retrieves the current block number from the blockchain
func (pv *PaymentVault) GetCurrentBlockNumber(ctx context.Context) (uint32, error) {
	blockNumber, err := pv.ethClient.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return uint32(blockNumber), nil
}
