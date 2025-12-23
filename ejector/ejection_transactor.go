package ejector

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/common/aws"
	contractEigenDAEjectionManager "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDAEjectionManager"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/config"
	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
//
// The ejection transactor is thread safe, although parallel calls may result in duplicate work (i.e. two calls might
// end up putting the same data in a cache).
type ejectionTransactor struct {
	logger logging.Logger

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

	// A utility for looking up which quorums a given validator is a member of at a specific reference block number.
	validatorQuorumLookup eth.ValidatorQuorumLookup

	// A utility for converting between validator IDs and addresses.
	validatorIDToAddressConverter eth.ValidatorIDToAddressConverter

	// The maximum gas limit to use for transactions.
	maxGasOverride uint64
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
	cfg *EjectorConfig,
) (EjectionTransactor, error) {

	var zeroAddress gethcommon.Address
	if selfAddress == zeroAddress {
		return nil, fmt.Errorf("selfAddress must be non-zero")
	}
	if privateKey == nil && (cfg.KmsKeyId == "" || cfg.KmsRegion == "") {
		return nil, fmt.Errorf("either privateKey or KMS configuration (keyId and region) must be provided")
	}
	if ejectionContractAddress == zeroAddress {
		return nil, fmt.Errorf("ejectionContractAddress must be non-zero")
	}
	if registryCoordinatorAddress == zeroAddress {
		return nil, fmt.Errorf("registryCoordinatorAddress must be non-zero")
	}
	if chainID.Sign() <= 0 {
		return nil, fmt.Errorf("invalid chainID: %s", chainID.String())
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

	referenceBlockProvider := eth.NewReferenceBlockProvider(logger, client, cfg.ReferenceBlockNumberOffset)
	referenceBlockProvider, err = eth.NewPeriodicReferenceBlockProvider(
		referenceBlockProvider,
		cfg.ReferenceBlockNumberPollInterval)
	if err != nil {
		return nil, fmt.Errorf("failed to create periodic reference block provider: %w", err)
	}

	validatorQuorumLookup, err := eth.NewValidatorQuorumLookup(client, registryCoordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator quorum lookup: %w", err)
	}
	validatorQuorumLookup, err = eth.NewCachedValidatorQuorumLookup(validatorQuorumLookup, int(cfg.ChainDataCacheSize))
	if err != nil {
		return nil, fmt.Errorf("failed to create cached validator quorum lookup: %w", err)
	}

	validatorIDToAddressConverter, err := eth.NewValidatorIDToAddressConverter(client, registryCoordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator ID to address converter: %w", err)
	}
	validatorIDToAddressConverter, err = eth.NewCachedValidatorIDToAddressConverter(
		validatorIDToAddressConverter,
		int(cfg.ChainDataCacheSize))
	if err != nil {
		return nil, fmt.Errorf("failed to create cached validator ID to address converter: %w", err)
	}

	signer, err := buildSigner(ctx, privateKey, chainID, cfg.KmsKeyId, cfg.KmsRegion, cfg.KmsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to build signer: %w", err)
	}

	return &ejectionTransactor{
		logger:                        logger,
		selfAddress:                   selfAddress,
		caller:                        caller,
		transactor:                    transactor,
		signer:                        signer,
		referenceBlockProvider:        referenceBlockProvider,
		validatorQuorumLookup:         validatorQuorumLookup,
		validatorIDToAddressConverter: validatorIDToAddressConverter,
		maxGasOverride:                cfg.MaxGasOverride,
	}, nil
}

// Build a signer function.
func buildSigner(
	ctx context.Context,
	privateKey *ecdsa.PrivateKey,
	chainID *big.Int,
	kmsKeyId string,
	kmsRegion string,
	kmsEndpoint string,
) (bind.SignerFn, error) {

	if privateKey != nil {
		signer, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			return nil, fmt.Errorf("failed to create transact opts: %w", err)
		}
		return signer.Signer, nil
	}

	if kmsKeyId == "" || kmsRegion == "" {
		return nil, fmt.Errorf("KMS key ID and region must be provided for KMS signing if private key is not provided")
	}

	// Load AWS config for the specified region
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(kmsRegion))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create KMS client with optional custom endpoint
	kmsClient := awskms.NewFromConfig(awsCfg, func(o *awskms.Options) {
		if kmsEndpoint != "" {
			o.BaseEndpoint = &kmsEndpoint
		}
	})

	// Load the public key from KMS
	publicKey, err := aws.LoadPublicKeyKMS(ctx, kmsClient, kmsKeyId)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key from KMS: %w", err)
	}

	// Create the EIP-155 signer for the chain
	signer := types.LatestSignerForChainID(chainID)

	// Create and return a SignerFn that uses KMS
	var signerFn bind.SignerFn = func(address gethcommon.Address, tx *types.Transaction) (*types.Transaction, error) {
		// Hash the transaction using the signer
		txHash := signer.Hash(tx)

		// Sign the hash using KMS
		signature, err := aws.SignKMS(context.Background(), kmsClient, kmsKeyId, publicKey, txHash.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to sign transaction with KMS: %w", err)
		}

		// Apply the signature to the transaction
		signedTx, err := tx.WithSignature(signer, signature)
		if err != nil {
			return nil, fmt.Errorf("failed to apply signature to transaction: %w", err)
		}

		return signedTx, nil
	}

	return signerFn, nil
}

func (e *ejectionTransactor) CompleteEjection(
	ctx context.Context,
	addressToEject gethcommon.Address,
) error {

	rbn, err := e.referenceBlockProvider.GetReferenceBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get reference block number: %w", err)
	}

	idToEject, err := e.validatorIDToAddressConverter.ValidatorAddressToID(ctx, addressToEject)
	if err != nil {
		return fmt.Errorf("failed to get validator ID from address: %w", err)
	}

	quorums, err := e.validatorQuorumLookup.GetQuorumsForValidator(ctx, idToEject, rbn)
	if err != nil {
		return fmt.Errorf("failed to get quorums for validator: %w", err)
	}

	if len(quorums) == 0 {
		return fmt.Errorf("cannot complete ejection: validator %s is not present in any quorum at RBN %d",
			idToEject.Hex(), rbn)
	}

	quorumBytes := eth.QuorumListToBytes(quorums)

	opts := &bind.TransactOpts{
		Context: ctx,
		From:    e.selfAddress,
		Signer:  e.signer,
	}

	if e.maxGasOverride != 0 {
		opts.GasLimit = e.maxGasOverride
	}

	txn, err := e.transactor.CompleteEjection(opts, addressToEject, quorumBytes)
	if err != nil {
		return fmt.Errorf("failed to complete ejection: %w", err)
	}

	e.logger.Debug("submitted CompleteEjection transaction",
		"transaction hash", txn.Hash().Hex(),
		"validator ID", idToEject.Hex(),
		"validator address", addressToEject.Hex(),
		"quorums", quorums,
		"reference block number", rbn,
	)

	return nil
}

func (e *ejectionTransactor) IsEjectionInProgress(
	ctx context.Context,
	addressToCheck gethcommon.Address,
) (bool, error) {

	opts := &bind.CallOpts{
		Context: ctx,
		From:    e.selfAddress,
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

func (e *ejectionTransactor) StartEjection(
	ctx context.Context,
	addressToEject gethcommon.Address,
) error {

	rbn, err := e.referenceBlockProvider.GetReferenceBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get reference block number: %w", err)
	}

	idToEject, err := e.validatorIDToAddressConverter.ValidatorAddressToID(ctx, addressToEject)
	if err != nil {
		return fmt.Errorf("failed to get validator ID from address: %w", err)
	}

	quorums, err := e.validatorQuorumLookup.GetQuorumsForValidator(ctx, idToEject, rbn)
	if err != nil {
		return fmt.Errorf("failed to get quorums for validator: %w", err)
	}
	quorumBytes := eth.QuorumListToBytes(quorums)

	opts := &bind.TransactOpts{
		Context: ctx,
		From:    e.selfAddress,
		Signer:  e.signer,
	}

	if e.maxGasOverride != 0 {
		opts.GasLimit = e.maxGasOverride
	}

	txn, err := e.transactor.StartEjection(opts, addressToEject, quorumBytes)
	if err != nil {
		return fmt.Errorf("failed to start ejection: %w", err)
	}

	e.logger.Debug("submitted StartEjection transaction",
		"transaction hash", txn.Hash().Hex(),
		"validator ID", idToEject.Hex(),
		"validator address", addressToEject.Hex(),
		"quorums", quorums,
		"reference block number", rbn,
	)

	return nil
}
