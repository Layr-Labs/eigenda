package verification

import (
	"context"
	"fmt"
	"sync"

	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	genericVerifierBinding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV3"
	opsrbinding "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// GenericCertVerifier is responsible for making eth calls against version agnostic CertVerifier contracts to ensure
// cryptographic and structural integrity of EigenDA certificate types.
// The V3 cert verifier contract is located at:
// https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/cert/v3/EigenDACertVerifierV3.sol
type GenericCertVerifier struct {
	logger            logging.Logger
	ethClient         common.EthClient
	opsrCaller        *opsrbinding.ContractOperatorStateRetrieverCaller
	addressProvider   clients.CertVerifierAddressProvider

	// Cache maps
	verifierCallers        sync.Map
	requiredQuorums        sync.Map
	confirmationThresholds sync.Map
	versions               sync.Map
}

// NewGenericCertVerifier constructs a new GenericCertVerifier instance
func NewGenericCertVerifier(
	logger logging.Logger,
	ethClient common.EthClient,
	certVerifierAddressProvider clients.CertVerifierAddressProvider,
) (*GenericCertVerifier, error) {
	return &GenericCertVerifier{
		logger:          logger,
		ethClient:       ethClient,
		addressProvider: certVerifierAddressProvider,
	}, nil
}

// CheckDACert calls the CheckDACert view function on the EigenDACertVerifier contract.
// This method returns nil if the certificate is successfully verified; otherwise, it returns an error.
func (cv *GenericCertVerifier) CheckDACert(
	ctx context.Context,
	referenceBlockNumber uint64,
	certBytes []byte,
) error {
	// 1 - Get verifier caller for the block number
	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, referenceBlockNumber)
	if err != nil {
		return fmt.Errorf("get verifier caller: %w", err)
	}

	// 2 - Call the contract method
	result, err := certVerifierCaller.CheckDACert(
		&bind.CallOpts{Context: ctx},
		certBytes,
	)
	if err != nil {
		return fmt.Errorf("verify cert: %w", err)
	}

	// 3 - Check the result, 1 means success while anything else indicates failure
	// TODO: Structured error responses by translating response codes
	if result != 1 {
		return fmt.Errorf("cert verification failed with status code: %d", result)
	}

	return nil
}



// GetQuorumNumbersRequired returns the set of quorum numbers that must be set in the BlobHeader, and verified in
// VerifyCert and CheckDACert.
//
// This method will return required quorum numbers from an internal cache if they are already known for the currently
// active cert verifier. Otherwise, this method will query the required quorum numbers from the currently active
// cert verifier, and cache the result for future use.
func (cv *GenericCertVerifier) GetQuorumNumbersRequired(ctx context.Context, referenceBlockNumber uint64) ([]uint8, error) {
	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get cert verifier address: %w", err)
	}

	// if the quorum numbers for the active cert verifier address have already been cached, return them immediately
	cachedQuorumNumbers, ok := cv.requiredQuorums.Load(certVerifierAddress)
	if ok {
		castQuorums, ok := cachedQuorumNumbers.([]uint8)
		if !ok {
			return nil, fmt.Errorf("expected quorum numbers to be []uint8")
		}
		return castQuorums, nil
	}

	// quorum numbers weren't cached, so proceed to fetch them
	certVerifierCaller, err := cv.getVerifierCallerFromAddress(certVerifierAddress)
	if err != nil {
		return nil, fmt.Errorf("get verifier caller from address: %w", err)
	}

	quorumNumbersRequired, err := certVerifierCaller.QuorumNumbersRequired(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("get quorum numbers required: %w", err)
	}

	cv.requiredQuorums.Store(certVerifierAddress, quorumNumbersRequired)

	return quorumNumbersRequired, nil
}

// getVerifierCallerFromBlockNumber returns a ContractEigenDACertVerifierV3 that corresponds to the input reference
// block number.
//
// This method caches ContractEigenDACertVerifierV3 instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *GenericCertVerifier) getVerifierCallerFromBlockNumber(
	ctx context.Context,
	referenceBlockNumber uint64,
) (*genericVerifierBinding.ContractEigenDACertVerifierV3, error) {
	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get cert verifier address: %w", err)
	}

	return cv.getVerifierCallerFromAddress(certVerifierAddress)
}

// getVerifierCallerFromAddress returns a ContractEigenDACertVerifierV3 that corresponds to the input contract
// address
//
// This method caches ContractEigenDACertVerifierV3 instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *GenericCertVerifier) getVerifierCallerFromAddress(
	certVerifierAddress gethcommon.Address,
) (*genericVerifierBinding.ContractEigenDACertVerifierV3, error) {
	existingCallerAny, valueExists := cv.verifierCallers.Load(certVerifierAddress)
	if valueExists {
		existingCaller, ok := existingCallerAny.(*genericVerifierBinding.ContractEigenDACertVerifierV3)
		if !ok {
			return nil, fmt.Errorf(
				"value in verifierCallers wasn't of type ContractEigenDACertVerifierV3. this should be impossible")
		}
		return existingCaller, nil
	}

	certVerifierCaller, err := genericVerifierBinding.NewContractEigenDACertVerifierV3(certVerifierAddress, cv.ethClient)
	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
	}

	cv.verifierCallers.Store(certVerifierAddress, certVerifierCaller)
	return certVerifierCaller, nil
}

// GetConfirmationThreshold returns the ConfirmationThreshold that corresponds to the input reference block number.
//
// This method will return the confirmation threshold from an internal cache if it is already known for the cert
// verifier which corresponds to the input reference block number. Otherwise, this method will query the confirmation
// threshold and cache the result for future use.
func (cv *GenericCertVerifier) GetConfirmationThreshold(ctx context.Context, referenceBlockNumber uint64) (uint8, error) {
	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("get cert verifier address: %w", err)
	}

	// if the confirmation threshold for the active cert verifier address has already been cached, return it immediately
	cachedThreshold, ok := cv.confirmationThresholds.Load(certVerifierAddress)
	if ok {
		castThreshold, ok := cachedThreshold.(uint8)
		if !ok {
			return 0, fmt.Errorf("expected confirmation threshold to be uint8")
		}
		return castThreshold, nil
	}

	// confirmation threshold wasn't cached, so proceed to fetch it
	certVerifierCaller, err := cv.getVerifierCallerFromAddress(certVerifierAddress)
	if err != nil {
		return 0, fmt.Errorf("get verifier caller from address: %w", err)
	}

	securityThresholds, err := certVerifierCaller.SecurityThresholds(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get security thresholds via contract call: %w", err)
	}

	cv.confirmationThresholds.Store(certVerifierAddress, securityThresholds.ConfirmationThreshold)

	return securityThresholds.ConfirmationThreshold, nil
}

// GetCertVersion returns the CertVersion that corresponds to the input reference block number.
//
// This method will return the version from an internal cache if it is already known for the cert
// verifier which corresponds to the input reference block number. Otherwise, this method will query the version
// and cache the result for future use.
func (cv *GenericCertVerifier) GetCertVersion(ctx context.Context, referenceBlockNumber uint64) (uint8, error) {
	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("get cert verifier address: %w", err)
	}

	// if the version for the active cert verifier address has already been cached, return it immediately
	cachedVersion, ok := cv.versions.Load(certVerifierAddress)
	if ok {
		castVersion, ok := cachedVersion.(uint8)
		if !ok {
			return 0, fmt.Errorf("expected version to be uint8")
		}
		return castVersion, nil
	}

	// version wasn't cached, so proceed to fetch it
	certVerifierCaller, err := cv.getVerifierCallerFromAddress(certVerifierAddress)
	if err != nil {
		return 0, fmt.Errorf("get verifier caller from address: %w", err)
	}

	version, err := certVerifierCaller.CertVersion(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get version via contract call: %w", err)
	}

	cv.versions.Store(certVerifierAddress, version)

	return version, nil
}