package verification

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// CertVerifier is responsible for making eth calls against the CertVerifier contract to ensure cryptographic and
// structural integrity of V2 certificates
//
// The cert verifier contract is located at https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDACertVerifier.sol
type CertVerifier struct {
	logger                      logging.Logger
	ethClient                   common.EthClient
	certVerifierAddressProvider clients.CertVerifierAddressProvider
	// maps contract address to a ContractEigenDACertVerifierCaller object
	verifierCallers sync.Map
	// maps contract address to set of required quorums specified in the contract at that address
	requiredQuorums sync.Map
	// maps contract address to the confirmation threshold required by that address
	confirmationThresholds sync.Map
}

var _ clients.ICertVerifier = &CertVerifier{}

// NewCertVerifier constructs a CertVerifier
func NewCertVerifier(
	logger logging.Logger,
	// the eth client, which should already be set up
	ethClient common.EthClient,
	certVerifierAddressProvider clients.CertVerifierAddressProvider,
) (*CertVerifier, error) {
	return &CertVerifier{
		logger:                      logger,
		ethClient:                   ethClient,
		certVerifierAddressProvider: certVerifierAddressProvider,
	}, nil
}

// VerifyCertV2FromSignedBatch calls the verifyDACertV2FromSignedBatch view function on the EigenDACertVerifier contract.
//
// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
func (cv *CertVerifier) VerifyCertV2FromSignedBatch(
	ctx context.Context,
	// The signed batch that contains the blob whose cert is being verified. This is obtained from the disperser, and
	// is used to verify that the described blob actually exists in a valid batch.
	signedBatch *disperser.SignedBatch,
	// Contains all necessary information about the blob, so that the cert can be verified.
	blobInclusionInfo *disperser.BlobInclusionInfo,
) error {
	convertedSignedBatch, err := coretypes.SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return fmt.Errorf("convert signed batch: %w", err)
	}

	convertedBlobInclusionInfo, err := coretypes.InclusionInfoProtoToBinding(blobInclusionInfo)
	if err != nil {
		return fmt.Errorf("convert blob inclusion info: %w", err)
	}

	referenceBlockNumber := signedBatch.GetHeader().GetReferenceBlockNumber()

	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, referenceBlockNumber)
	if err != nil {
		return fmt.Errorf("get verifier caller: %w", err)
	}

	err = certVerifierCaller.VerifyDACertV2FromSignedBatch(
		&bind.CallOpts{Context: ctx},
		*convertedSignedBatch,
		*convertedBlobInclusionInfo)

	if err != nil {
		return fmt.Errorf("verify cert v2 from signed batch: %w", err)
	}

	return nil
}

// VerifyCertV2 calls the VerifyCertV2 view function on the EigenDACertVerifier contract.
//
// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
func (cv *CertVerifier) VerifyCertV2(ctx context.Context, eigenDACert *coretypes.EigenDACert) error {
	referenceBlockNumber := uint64(eigenDACert.BatchHeader.ReferenceBlockNumber)

	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, referenceBlockNumber)
	if err != nil {
		return fmt.Errorf("get verifier caller: %w", err)
	}

	err = certVerifierCaller.VerifyDACertV2(
		&bind.CallOpts{Context: ctx},
		eigenDACert.BatchHeader,
		eigenDACert.BlobInclusionInfo,
		eigenDACert.NonSignerStakesAndSignature,
		eigenDACert.SignedQuorumNumbers)

	if err != nil {
		return fmt.Errorf("verify cert v2: %w", err)
	}

	return nil
}

// GetNonSignerStakesAndSignature calls the getNonSignerStakesAndSignature view function on the EigenDACertVerifier
// contract, and returns the resulting NonSignerStakesAndSignature object.
func (cv *CertVerifier) GetNonSignerStakesAndSignature(
	ctx context.Context,
	signedBatch *disperser.SignedBatch,
) (*verifierBindings.NonSignerStakesAndSignature, error) {
	signedBatchBinding, err := coretypes.SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return nil, fmt.Errorf("convert signed batch: %w", err)
	}

	referenceBlockNumber := signedBatch.GetHeader().GetReferenceBlockNumber()

	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get verifier caller: %w", err)
	}

	nonSignerStakesAndSignature, err := certVerifierCaller.GetNonSignerStakesAndSignature(
		&bind.CallOpts{Context: ctx},
		*signedBatchBinding)

	if err != nil {
		return nil, fmt.Errorf("get non signer stakes and signature: %w", err)
	}

	return &nonSignerStakesAndSignature, nil
}

// GetQuorumNumbersRequired returns the set of quorum numbers that must be set in the BlobHeader, and verified in
// VerifyDACertV2 and verifyDACertV2FromSignedBatch.
//
// This method will return required quorum numbers from an internal cache if they are already known for the currently
// active cert verifier. Otherwise, this method will query the required quorum numbers from the currently active
// cert verifier, and cache the result for future use.
func (cv *CertVerifier) GetQuorumNumbersRequired(ctx context.Context) ([]uint8, error) {
	blockNumber, err := cv.ethClient.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch block number from eth client: %w", err)
	}

	certVerifierAddress, err := cv.certVerifierAddressProvider.GetCertVerifierAddress(ctx, blockNumber)
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

	quorumNumbersRequired, err := certVerifierCaller.QuorumNumbersRequiredV2(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("get quorum numbers required: %w", err)
	}

	cv.requiredQuorums.Store(certVerifierAddress, quorumNumbersRequired)

	return quorumNumbersRequired, nil
}

// getVerifierCallerFromBlockNumber returns a ContractEigenDACertVerifierCaller that corresponds to the input reference
// block number.
//
// This method caches ContractEigenDACertVerifierCaller instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *CertVerifier) getVerifierCallerFromBlockNumber(
	ctx context.Context,
	referenceBlockNumber uint64,
) (*verifierBindings.ContractEigenDACertVerifierV2Caller, error) {
	certVerifierAddress, err := cv.certVerifierAddressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get cert verifier address: %w", err)
	}

	return cv.getVerifierCallerFromAddress(certVerifierAddress)
}

// getVerifierCallerFromAddress returns a ContractEigenDACertVerifierCaller that corresponds to the input contract
// address
//
// This method caches ContractEigenDACertVerifierCaller instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *CertVerifier) getVerifierCallerFromAddress(
	certVerifierAddress gethcommon.Address,
) (*verifierBindings.ContractEigenDACertVerifierV2Caller, error) {
	existingCallerAny, valueExists := cv.verifierCallers.Load(certVerifierAddress)
	if valueExists {
		existingCaller, ok := existingCallerAny.(*verifierBindings.ContractEigenDACertVerifierV2Caller)
		if !ok {
			return nil, fmt.Errorf(
				"value in verifierCallers wasn't of type ContractEigenDACertVerifierCaller. this should be impossible")
		}
		return existingCaller, nil
	}

	certVerifierCaller, err := verifierBindings.NewContractEigenDACertVerifierV2Caller(certVerifierAddress, cv.ethClient)
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
	certVerifierAddress, err := cv.certVerifierAddressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
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
	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, referenceBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("get verifier caller from block number: %w", err)
	}

	securityThresholds, err := certVerifierCaller.SecurityThresholdsV2(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get security thresholds via contract call: %w", err)
	}

	cv.confirmationThresholds.Store(certVerifierAddress, securityThresholds.ConfirmationThreshold)

	return securityThresholds.ConfirmationThreshold, nil
}
