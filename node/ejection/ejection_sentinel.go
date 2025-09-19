package ejection

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	contractEigenDAEjectionManager "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDAEjectionManager"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// The EjectionSentinel watches for when ejection is initiated against this validator. If that happens, this utility
// may also perform an "ejection defense" in order to prevent the validator from being ejected.
type EjectionSentinel struct {
	ctx    context.Context
	logger logging.Logger

	// the time in between checks for ejection events
	period time.Duration

	// used to execute eth reads
	caller *contractEigenDAEjectionManager.ContractIEigenDAEjectionManagerCaller

	// used to execute eth writes
	transactor *contractEigenDAEjectionManager.ContractIEigenDAEjectionManagerTransactor

	// the address of this validator
	selfAddress gethcommon.Address

	// If true, the sentinel will attempt to contest ejection by sending a transaction to cancel the ejection.
	ejectionDefenseEnabled bool

	// Normally, the sentinel will check the software version of the validator before deciding whether to contest
	// ejection. Under normal circumstances, an honest validator should not contest ejection if it is running software
	// that does not meet the minimum version number. However, if the governing body in control of setting the minimum
	// version number goes rogue, honest validators may want to contest ejection regardless of the claimed minimum
	// version number.
	ignoreVersion bool

	// A function that can sign transactions from selfAddress. nil if ejectionDefenseEnabled is false.
	signer func(address gethcommon.Address, tx *types.Transaction) (*types.Transaction, error)
}

// NewEjectionSentinel creates a new EjectionSentinel instance.
func NewEjectionSentinel(
	ctx context.Context,
	logger logging.Logger,
	ejectionContractAddress gethcommon.Address,
	ethClient common.EthClient,
	privateKey *ecdsa.PrivateKey,
	selfAddress gethcommon.Address,
	period time.Duration,
	ejectionDefenseEnabled bool,
	ignoreVersion bool,
) (*EjectionSentinel, error) {

	if period <= 0 {
		return nil, fmt.Errorf("period must be greater than 0, got %v", period)
	}

	var zeroAddress gethcommon.Address
	if selfAddress == zeroAddress {
		return nil, fmt.Errorf("selfAddress must be non-zero")
	}

	caller, err := contractEigenDAEjectionManager.NewContractIEigenDAEjectionManagerCaller(
		ejectionContractAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create ejection manager caller: %w", err)
	}

	var transactor *contractEigenDAEjectionManager.ContractIEigenDAEjectionManagerTransactor
	var signer func(address gethcommon.Address, tx *types.Transaction) (*types.Transaction, error)
	if ejectionDefenseEnabled {
		if privateKey == nil {
			return nil, fmt.Errorf("privateKey must be provided if ejection defense is enabled")
		}

		logger.Info("ejection defense enabled")

		transactor, err = contractEigenDAEjectionManager.NewContractIEigenDAEjectionManagerTransactor(
			ejectionContractAddress, ethClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create ejection manager transactor: %w", err)
		}

		chainID, err := ethClient.ChainID(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get chain ID: %w", err)
		}

		transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			return nil, fmt.Errorf("failed to create transact opts: %w", err)
		}

		signer = transactOpts.Signer
	} else {
		logger.Info("ejection defense not enabled")
	}

	sentinel := &EjectionSentinel{
		ctx:                    ctx,
		logger:                 logger,
		period:                 period,
		selfAddress:            selfAddress,
		caller:                 caller,
		transactor:             transactor,
		ejectionDefenseEnabled: ejectionDefenseEnabled,
		ignoreVersion:          ignoreVersion,
		signer:                 signer,
	}
	go sentinel.run()

	return sentinel, nil
}

// The EjectionSentinel's goroutine that watches for ejection events and performs necessary actions.
func (s *EjectionSentinel) run() {
	ticker := time.NewTicker(s.period)
	defer ticker.Stop()

	s.logger.Debugf("Ejection Sentinel is running with a period of %s", s.period)

	for {
		select {
		case <-ticker.C:
			err := s.checkEjectionStatus()
			if err != nil {
				s.logger.Errorf("Error checking ejection status: %v", err)
			}
		case <-s.ctx.Done():
			s.logger.Info("EjectionSentinel stopped")
			return
		}
	}
}

// checkEjectionStatus checks if the validator is being ejected and performs necessary actions based on the result.
func (s *EjectionSentinel) checkEjectionStatus() error {

	// This method will return the ID of the entity attempting an ejection if an ejection is in progress,
	// or the zero address if no ejection is in progress.
	ejector, err := s.caller.GetEjector(&bind.CallOpts{Context: s.ctx}, s.selfAddress)
	if err != nil {
		return fmt.Errorf("failed to check ejection status: %w", err)
	}

	var zeroAddress gethcommon.Address
	ejectionInProgress := ejector != zeroAddress
	if !ejectionInProgress {
		s.logger.Debug("This validator is not currently being ejected.")
		return nil
	}

	s.logger.Warnf("This validator is currently being ejected by %s", ejector.Hex())

	if s.transactor == nil {
		// TODO(cody.littley) Talk to Lulu about the "special log" we need to do to support validators
		//  who want to sign cancellation with external key management systems. That log should happen here.

		s.logger.Errorf("This validator is not configured to contest ejection. " +
			"Unless there is manual intervention, this validator may be ejected in the near future.")
		return nil
	}

	// TODO(cody.littley) check if we are running a software version that permits ejection defense
	//  Minimum software version is not currently written onchain so we can't write the offchain logic yet.

	s.logger.Info("Submitting ejection cancellation transaction.")

	txn, err := s.transactor.CancelEjection(&bind.TransactOpts{
		From:    s.selfAddress,
		Context: s.ctx,
		Signer:  s.signer,
	})
	if err != nil {
		return fmt.Errorf("failed to submit ejection cancellation transaction: %w", err)
	}

	s.logger.Infof("Ejection cancellation transaction submitted: %s", txn.Hash().Hex())

	return nil
}
