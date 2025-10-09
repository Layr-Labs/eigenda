package ejector

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	contractEigenDAEjectionManager "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDAEjectionManager"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// EjectionTransactor executes transactions related to ejections. This layer of abstraction allows for easier
// unit testing of the ejector logic.
type EjectionTransactor interface {

	// Begin ejection proceedings against the operator with the given address.
	StartEjection(ctx context.Context, addressToEject gethcommon.Address) error

	// Checks to see if an ejection is currently in progress for the operator with the given address.
	IsEjectionInProgress(ctx context.Context, addressToCheck gethcommon.Address) (bool, error)

	// Checks to see if the validator with the given address is present in any quorum.
	IsValidatorPresentInAnyQuorum(ctx context.Context, addressToCheck gethcommon.Address) (bool, error)

	// Complete the ejection proceedings against the operator with the given address.
	CompleteEjection(ctx context.Context, addressToEject gethcommon.Address) error
}

var _ EjectionTransactor = &ejectionTransactor{}

// ejectionTransactor is the production implementation of the EjectionTransactor interface.
type ejectionTransactor struct {
	// Used to execute eth reads
	caller *contractEigenDAEjectionManager.ContractIEigenDAEjectionManagerCaller
	// Used to execute eth writes
	transactor *contractEigenDAEjectionManager.ContractIEigenDAEjectionManagerTransactor
	// A function that can sign transactions from selfAddress.
	signer bind.SignerFn
}

// Create a new EjectionTransactor.
func NewEjectionTransactor(
	ctx context.Context,
	client bind.ContractBackend,
	ejectionContractAddress gethcommon.Address,
	selfAddress gethcommon.Address,
	privateKey *ecdsa.PrivateKey,
	chainID *big.Int,
) (EjectionTransactor, error) {

	var zeroAddress gethcommon.Address
	if selfAddress == zeroAddress {
		return nil, fmt.Errorf("selfAddress must be non-zero")
	}
	if privateKey == nil {
		return nil, fmt.Errorf("privateKey must be non-nil")
	}

	caller, err := contractEigenDAEjectionManager.NewContractIEigenDAEjectionManagerCaller(
		ejectionContractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create ejection manager caller: %w", err)
	}

	transactor, err := contractEigenDAEjectionManager.NewContractIEigenDAEjectionManagerTransactor(
		ejectionContractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create ejection manager transactor: %w", err)
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transact opts: %w", err)
	}

	signer := transactOpts.Signer

	return &ejectionTransactor{
		caller:     caller,
		transactor: transactor,
		signer:     signer,
	}, nil
}

// CompleteEjection implements EjectionTransactor.
func (e *ejectionTransactor) CompleteEjection(
	ctx context.Context,
	addressToEject gethcommon.Address,
) error {
	panic("unimplemented")
}

// IsEjectionInProgress implements EjectionTransactor.
func (e *ejectionTransactor) IsEjectionInProgress(
	ctx context.Context,
	addressToCheck gethcommon.Address,
) (bool, error) {
	panic("unimplemented")
}

// IsValidatorPresentInAnyQuorum implements EjectionTransactor.
func (e *ejectionTransactor) IsValidatorPresentInAnyQuorum(
	ctx context.Context,
	addressToCheck gethcommon.Address,
) (bool, error) {
	panic("unimplemented")
}

// StartEjection implements EjectionTransactor.
func (e *ejectionTransactor) StartEjection(
	ctx context.Context,
	addressToEject gethcommon.Address) error {
	panic("unimplemented")
}
