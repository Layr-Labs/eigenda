package verification

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common"
	certVerifierBinding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// CertVerifier is responsible for making eth calls against version agnostic CertVerifier contracts to ensure
// cryptographic and structural integrity of EigenDA certificate types.
// The V3 cert verifier contract is located at:
// https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/cert/EigenDACertVerifier.sol
type CertVerifier struct {
	logger          logging.Logger
	ethClient       common.EthClient
	addressProvider clients.CertVerifierAddressProvider

	// maps contract address to a ContractEigenDACertVerifierCaller object
	verifierCallers sync.Map
	// maps contract address to set of required quorums specified in the contract at that address
	requiredQuorums sync.Map
	// maps contract address to the confirmation threshold required by that address
	confirmationThresholds sync.Map
	// maps contract address to the cert version specified in the contract at that address
	versions sync.Map
}

// NewCertVerifier constructs a new CertVerifier instance
func NewCertVerifier(
	logger logging.Logger,
	ethClient common.EthClient,
	certVerifierAddressProvider clients.CertVerifierAddressProvider,
) (*CertVerifier, error) {
	return &CertVerifier{
		logger:          logger,
		ethClient:       ethClient,
		addressProvider: certVerifierAddressProvider,
	}, nil
}

// CheckDACert calls the CheckDACert view function on the EigenDACertVerifier contract.
// This method returns nil if the certificate is successfully verified; otherwise, it returns one of
// [CertVerifierInputError], [CertVerifierInvalidCertError], or [CertVerifierInternalError] errors.
func (cv *CertVerifier) CheckDACert(
	ctx context.Context,
	cert coretypes.EigenDACert,
) error {
	// 1 - switch on the certificate version to determine which underlying type to decode into
	//     and which contract to call

	// EigenDACertV3 is the only version that is supported by the CheckDACert function
	var certV3 *coretypes.EigenDACertV3
	var err error
	switch c := cert.(type) {
	case *coretypes.EigenDACertV3:
		certV3 = c
	case *coretypes.EigenDACertV2:
		certV3, err = c.ToV3()
		if err != nil {
			return &CertVerifierInternalError{Msg: "convert V2 cert to V3", Err: err}
		}
	default:
		return &CertVerifierInputError{Msg: fmt.Sprintf("unsupported cert version: %T", cert)}
	}

	// 2 - Call the contract method CheckDACert to verify the certificate
	// TODO: Determine adequate future proofing strategy for EigenDACertVerifierRouter to be compliant
	//       with future reference timestamp change which deprecates the reference block number
	//       used for quorum stake check-pointing.
	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, certV3.ReferenceBlockNumber())
	if err != nil {
		return &CertVerifierInternalError{Msg: "get verifier caller", Err: err}
	}

	certBytes, err := certV3.Serialize(coretypes.CertSerializationABI)
	if err != nil {
		return &CertVerifierInternalError{Msg: "serialize cert", Err: err}
	}

	// TODO: determine if there's any merit in passing call options to impose better determinism and
	// safety on the operation
	result, err := certVerifierCaller.CheckDACert(
		&bind.CallOpts{Context: ctx},
		certBytes,
	)
	if err != nil {
		return &CertVerifierInternalError{Msg: "checkDACert eth call", Err: err}
	}

	// 3 - Cast result to structured enum type and check for success
	verifyResultCode := coretypes.VerificationStatusCode(result)
	if verifyResultCode == coretypes.StatusNullError {
		return &CertVerifierInternalError{Msg: fmt.Sprintf("checkDACert eth-call bug: %s", verifyResultCode.String())}
	} else if verifyResultCode != coretypes.StatusSuccess {
		return &CertVerifierInvalidCertError{
			StatusCode: verifyResultCode,
			Msg:        verifyResultCode.String(),
		}
	}
	return nil
}

// GetQuorumNumbersRequired returns the set of quorum numbers that must be set in the BlobHeader, and verified in
// VerifyCert and CheckDACert.
//
// This method will return required quorum numbers from an internal cache if they are already known for the currently
// active cert verifier. Otherwise, this method will query the required quorum numbers from the currently active
// cert verifier, and cache the result for future use.
func (cv *CertVerifier) GetQuorumNumbersRequired(ctx context.Context) ([]uint8, error) {
	// get the latest cert verifier address from the address provider

	blockNum, err := cv.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("get latest block number: %w", err)
	}

	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, blockNum.NumberU64())
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

// getVerifierCallerFromBlockNumber returns a ContractEigenDACertVerifier that corresponds to the input reference
// block number.
//
// This method caches ContractEigenDACertVerifier instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *CertVerifier) getVerifierCallerFromBlockNumber(
	ctx context.Context,
	referenceBlockNumber uint64,
) (*certVerifierBinding.ContractEigenDACertVerifier, error) {
	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get cert verifier address: %w", err)
	}

	return cv.getVerifierCallerFromAddress(certVerifierAddress)
}

// getVerifierCallerFromAddress returns a ContractEigenDACertVerifier that corresponds to the input contract
// address
//
// This method caches ContractEigenDACertVerifier instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *CertVerifier) getVerifierCallerFromAddress(
	certVerifierAddress gethcommon.Address,
) (*certVerifierBinding.ContractEigenDACertVerifier, error) {
	existingCallerAny, valueExists := cv.verifierCallers.Load(certVerifierAddress)
	if valueExists {
		existingCaller, ok := existingCallerAny.(*certVerifierBinding.ContractEigenDACertVerifier)
		if !ok {
			return nil, fmt.Errorf(
				"value in verifierCallers wasn't of type ContractEigenDACertVerifier. this should be impossible")
		}
		return existingCaller, nil
	}

	certVerifierCaller, err := certVerifierBinding.NewContractEigenDACertVerifier(certVerifierAddress, cv.ethClient)
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
func (cv *CertVerifier) GetConfirmationThreshold(ctx context.Context, referenceBlockNumber uint64) (uint8, error) {
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
func (cv *CertVerifier) GetCertVersion(ctx context.Context, referenceBlockNumber uint64) (uint8, error) {
	certVerifierAddress, err := cv.addressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("get cert verifier address: %w", err)
	}

	// if the version for the active cert verifier address has already been cached, return it immediately
	cachedVersion, ok := cv.versions.Load(certVerifierAddress)
	if ok {
		castVersion, ok := cachedVersion.(uint8)
		if !ok {
			return 0, fmt.Errorf("expected version to be uint64")
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
