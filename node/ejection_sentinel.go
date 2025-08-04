package node

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	contractEigenDAEjectionManager "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAEjectionManager"
	"github.com/Layr-Labs/eigenda/contracts/directory"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// The EjectionSentinel watches for when ejection is initiated against this validator. If that happens, this utility
// may also perform an "ejection defense" in order to prevent the validator from being ejected.
type EjectionSentinel struct {
	ctx    context.Context
	logger logging.Logger

	// the time in between checks for ejection events
	period time.Duration

	// used execute read eth contract calls
	caller *contractEigenDAEjectionManager.ContractEigenDAEjectionManagerCaller

	// used to execute write eth contract calls
	transactor *contractEigenDAEjectionManager.ContractEigenDAEjectionManagerTransactor

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
}

// NewEjectionSentinel creates a new EjectionSentinel instance.
func NewEjectionSentinel(
	ctx context.Context,
	logger logging.Logger,
	contractDirectory *directory.ContractDirectory,
	ethClient common.EthClient,
	selfAddress gethcommon.Address,
	period time.Duration,
	ejectionDefenseEnabled bool,
	ignoreVersion bool,
) (*EjectionSentinel, error) {

	if period <= 0 {
		return nil, fmt.Errorf("period must be greater than 0, got %v", period)
	}

	ejectionContractAddress, err := contractDirectory.GetContractAddress(ctx, directory.EigenDAEjectionManager)
	if err != nil {
		return nil, fmt.Errorf("failed to get ejection contract address: %w", err)
	}

	caller, err := contractEigenDAEjectionManager.NewContractEigenDAEjectionManagerCaller(
		ejectionContractAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create ejection manager caller: %w", err)
	}

	var transactor *contractEigenDAEjectionManager.ContractEigenDAEjectionManagerTransactor
	if ejectionDefenseEnabled {
		logger.Info("ejection defense enabled")
		transactor, err = contractEigenDAEjectionManager.NewContractEigenDAEjectionManagerTransactor(
			ejectionContractAddress, ethClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create ejection manager transactor: %w", err)
		}
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
	}
	go sentinel.run()

	return sentinel, nil
}

// The EjectionSentinel's goroutine that watches for ejection events and performs necessary actions.
func (s *EjectionSentinel) run() {
	ticker := time.NewTicker(s.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.checkEjectionStatus()
			if err != nil {
				s.logger.Errorf("Error checking ejection status: %v", err)
			}
		case <-s.ctx.Done():
			s.logger.Info("EjectionSentinel stopped")
		}
	}
}

// checkEjectionStatus checks if the validator is being ejected and performs necessary actions based on the result.
func (s *EjectionSentinel) checkEjectionStatus() error {
	ejectionInProgress, err := s.caller.EjectionInitiated(&bind.CallOpts{Context: s.ctx}, s.selfAddress)
	if err != nil {
		return fmt.Errorf("failed to check ejection status: %w", err)
	}

	if !ejectionInProgress {
		s.logger.Debug("This validator is not currently being ejected.")
		return nil
	}

	s.logger.Warn("This validator is currently being ejected.")

	if s.transactor == nil {
		s.logger.Errorf("This validator is not configured to contest ejection. " +
			"Unless there is manual intervention, this validator may be ejected in the near future.")
		return nil
	}

	// TODO(cody-littley) check if we are running a software version that permits ejection defense

	s.logger.Info("Submitting ejection cancellation transaction.")
	txn, err := s.transactor.CancelEjection(&bind.TransactOpts{Context: s.ctx})
	if err != nil {
		return fmt.Errorf("failed to submit ejection cancellation transaction: %w", err)
	}

	s.logger.Infof("Ejection cancellation transaction submitted: %s", txn.Hash().Hex())

	return nil
}
