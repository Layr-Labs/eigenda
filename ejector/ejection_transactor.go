package ejector

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	contractEigenDAEjectionManager "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDAEjectionManager"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
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
	// The address of this ejector instance.
	selfAddress gethcommon.Address

	// Used to execute eth reads
	caller *contractEigenDAEjectionManager.ContractIEigenDAEjectionManagerCaller

	// Used to execute eth writes
	transactor *contractEigenDAEjectionManager.ContractIEigenDAEjectionManagerTransactor

	// A function that can sign transactions from selfAddress.
	signer bind.SignerFn

	// A utility for getting the reference block number.
	referenceBlockProvider eth.ReferenceBlockProvider

	// A utility for getting a list of all quorums.
	quorumScanner eth.QuorumScanner

	// A utility for looking up which quorums a given validator is a member of at a specific reference block number.
	validatorQuorumLookup eth.ValidatorQuorumLookup

	// A utility for converting between validator IDs and addresses.
	validatorIDToAddressConverter eth.ValidatorIDToAddressConverter
}

// Create a new EjectionTransactor.
func NewEjectionTransactor(
	ctx context.Context,
	logger logging.Logger,
	client bind.ContractBackend,
	ejectionContractAddress gethcommon.Address,
	registryCoordinatorAddress gethcommon.Address,
	selfAddress gethcommon.Address,
	privateKey *ecdsa.PrivateKey,
	chainID *big.Int,
	referenceBlockNumberOffset uint64,
	referenceBlockNumberPollInterval time.Duration,
	ethCacheSize int,
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

	referenceBlockProvider := eth.NewReferenceBlockProvider(logger, client, referenceBlockNumberOffset)
	referenceBlockProvider = eth.NewPeriodicReferenceBlockProvider(
		referenceBlockProvider,
		referenceBlockNumberPollInterval)

	quorumScanner, err := eth.NewQuorumScanner(client, registryCoordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create quorum scanner: %w", err)
	}
	quorumScanner, err = eth.NewCachedQuorumScanner(quorumScanner, ethCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create cached quorum scanner: %w", err)
	}

	validatorQuorumLookup, err := eth.NewValidatorQuorumLookup(client, registryCoordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator quorum lookup: %w", err)
	}
	validatorQuorumLookup, err = eth.NewCachedValidatorQuorumLookup(validatorQuorumLookup, ethCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create cached validator quorum lookup: %w", err)
	}

	validatorIDToAddressConverter, err := eth.NewValidatorIDToAddressConverter(client, registryCoordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator ID to address converter: %w", err)
	}
	validatorIDToAddressConverter, err = eth.NewCachedValidatorIDToAddressConverter(
		validatorIDToAddressConverter,
		ethCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create cached validator ID to address converter: %w", err)
	}

	return &ejectionTransactor{
		selfAddress:                   selfAddress,
		caller:                        caller,
		transactor:                    transactor,
		signer:                        transactOpts.Signer,
		referenceBlockProvider:        referenceBlockProvider,
		quorumScanner:                 quorumScanner,
		validatorQuorumLookup:         validatorQuorumLookup,
		validatorIDToAddressConverter: validatorIDToAddressConverter,
	}, nil
}

// CompleteEjection implements EjectionTransactor.
func (e *ejectionTransactor) CompleteEjection(
	ctx context.Context,
	addressToEject gethcommon.Address,
) error {

	rbn, err := e.referenceBlockProvider.GetReferenceBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get reference block number: %w", err)
	}

	quorums, err := e.quorumScanner.GetQuorums(ctx, rbn)
	if err != nil {
		return fmt.Errorf("failed to get quorums: %w", err)
	}
	quorumBytes := eth.QuorumListToBytes(quorums)

	opts := &bind.TransactOpts{ // TODO make sure these are correct
		From:   e.selfAddress,
		Signer: e.signer,
	}

	_, err = e.transactor.CompleteEjection(opts, addressToEject, quorumBytes)
	if err != nil {
		return fmt.Errorf("failed to complete ejection: %w", err)
	}
	return nil
}

// IsEjectionInProgress implements EjectionTransactor.
func (e *ejectionTransactor) IsEjectionInProgress(
	ctx context.Context,
	addressToCheck gethcommon.Address,
) (bool, error) {

	opts := &bind.CallOpts{ // TODO make sure these are correct
		From:    e.selfAddress,
		Context: ctx,
	}

	// This method returns the zero address if no ejection is in progress.
	ejector, err := e.caller.GetEjector(opts, addressToCheck)
	if err != nil {
		return false, fmt.Errorf("failed to check if ejection is in progress: %w", err)
	}

	var zeroAddress gethcommon.Address
	if ejector != zeroAddress {
		return true, nil
	}
	return false, nil
}

// IsValidatorPresentInAnyQuorum implements EjectionTransactor.
func (e *ejectionTransactor) IsValidatorPresentInAnyQuorum(
	ctx context.Context,
	addressToCheck gethcommon.Address,
) (bool, error) {

	rbn, err := e.referenceBlockProvider.GetReferenceBlockNumber(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get reference block number: %w", err)
	}

	validatorID, err := e.validatorIDToAddressConverter.ValidatorAddressToID(ctx, addressToCheck)
	if err != nil {
		return false, fmt.Errorf("failed to get validator ID from address: %w", err)
	}

	quorums, err := e.validatorQuorumLookup.GetQuorumsForValidator(ctx, validatorID, rbn)
	if err != nil {
		return false, fmt.Errorf("failed to get quorums for validator: %w", err)
	}

	return len(quorums) > 0, nil
}

// StartEjection implements EjectionTransactor.
func (e *ejectionTransactor) StartEjection(
	ctx context.Context,
	addressToEject gethcommon.Address) error {

	rbn, err := e.referenceBlockProvider.GetReferenceBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get reference block number: %w", err)
	}

	quorums, err := e.quorumScanner.GetQuorums(ctx, rbn)
	if err != nil {
		return fmt.Errorf("failed to get quorums: %w", err)
	}
	quorumBytes := eth.QuorumListToBytes(quorums)

	opts := &bind.TransactOpts{ // TODO make sure these are correct
		From:   e.selfAddress,
		Signer: e.signer,
	}

	_, err = e.transactor.StartEjection(opts, addressToEject, quorumBytes)
	if err != nil {
		return fmt.Errorf("failed to start ejection: %w", err)
	}
	return nil
}
