package verification

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	cv_v3_binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV3"
	cert_type_binding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	opsr_binding "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// CertVerifierV3 is responsible for making eth calls against the CertVerifier contract to ensure
// cryptographic and structural integrity of EigenDA V3 certificate types.
// The V3 cert verifier contract is located at:
// https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/cert/v3/EigenDACertVerifierV3.sol
type CertVerifierV3 struct {
	logger                  logging.Logger
	ethClient               common.EthClient
	opsrCaller              opsr_binding.ContractOperatorStateRetrieverCaller
	registryCoordinatorAddr gethcommon.Address
	routerAddressProvider   clients.CertVerifierAddressProvider

	// Cache maps
	verifierCallers        sync.Map // maps address -> ContractEigenDACertVerifierV3Caller
	requiredQuorums        sync.Map // maps address -> []uint8 (quorum numbers)
	confirmationThresholds sync.Map // maps address -> uint8 (threshold)
	versions               sync.Map // maps address -> uint8 (version)
}

// Ensure CertVerifierV3 implements the ICertVerifier interface
var _ clients.IV3CertVerifier = &CertVerifierV3{}

// NewV3CertVerifier constructs a new CertVerifierV3 instance
func NewV3CertVerifier(
	logger logging.Logger,
	ethClient common.EthClient,
	certVerifierAddressProvider clients.CertVerifierAddressProvider,
	registryCoordinatorAddr gethcommon.Address,
	opsrAddr gethcommon.Address,
) (*CertVerifierV3, error) {
	// Create the Operator State Retriever caller
	opsrCaller, err := opsr_binding.NewContractOperatorStateRetrieverCaller(opsrAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create OPSR caller: %w", err)
	}

	return &CertVerifierV3{
		logger:                  logger,
		ethClient:               ethClient,
		routerAddressProvider:   certVerifierAddressProvider,
		registryCoordinatorAddr: registryCoordinatorAddr,
		opsrCaller:              opsrCaller,
	}, nil
}

// CheckDACert calls the CheckDACert view function on the EigenDACertVerifier contract.
// This method returns nil if the certificate is successfully verified; otherwise, it returns an error.
func (cv *CertVerifierV3) CheckDACert(
	ctx context.Context,
	referenceBlockNumber uint64,
	certBytes []byte,
) error {
	// Get verifier caller for the block number
	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, referenceBlockNumber)
	if err != nil {
		return fmt.Errorf("get verifier caller: %w", err)
	}

	// Call the contract method
	result, err := certVerifierCaller.CheckDACert(
		&bind.CallOpts{Context: ctx},
		certBytes,
	)
	if err != nil {
		return fmt.Errorf("verify cert: %w", err)
	}

	// Check the result - 1 means success
	if result != 1 {
		return fmt.Errorf("cert verification failed with status code: %d", result)
	}

	return nil
}

// GetNonSignerStakesAndSignature constructs a NonSignerStakesAndSignature object by calling an
// onchain OperatorStateRetriever retriever to fetch necessary nosigner metadata
func (cv *CertVerifierV3) GetNonSignerStakesAndSignature(
	ctx context.Context,
	signedBatch *disperser.SignedBatch,
) (*cert_type_binding.NonSignerStakesAndSignature, error) {
	signedBatchBinding, err := coretypes.SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return nil, fmt.Errorf("convert signed batch: %w", err)
	}

	nonSignerPubKeys := signedBatch.GetAttestation().GetNonSignerPubkeys()

	// 1 - Pre-process inputs for operator state retriever call
	nonSignerOperatorIDs := make([][32]byte, len(nonSignerPubKeys))
	for i, pubKeySet := range nonSignerPubKeys {
		nonSignerOperatorIDs[i] = crypto.Keccak256Hash(pubKeySet)
	}

	quorumNumbers := make([]byte, len(signedBatch.GetAttestation().GetQuorumNumbers()))
	for i, qn := range signedBatch.GetAttestation().GetQuorumNumbers() {
		quorumNumbers[i] = byte(qn)
	}

	referenceBlockNumber := signedBatch.GetHeader().GetReferenceBlockNumber()

	// 2 - call operator state retriever to fetch signature indices
	checkSigIndices, err := cv.opsrCaller.GetCheckSignaturesIndices(&bind.CallOpts{Context: ctx},
		cv.registryCoordinatorAddr, uint32(referenceBlockNumber), quorumNumbers, nonSignerOperatorIDs)

	if err != nil {
		return nil, fmt.Errorf("check sig indices call: %w", err)
	}

	// 3 - construct non signer stakes and signature
	return &cert_type_binding.EigenDATypesV1NonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: checkSigIndices.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:             signedBatchBinding.Attestation.NonSignerPubkeys,
		QuorumApks:                   signedBatchBinding.Attestation.QuorumApks,
		ApkG2:                        signedBatchBinding.Attestation.ApkG2,
		Sigma:                        signedBatchBinding.Attestation.Sigma,
		QuorumApkIndices:             checkSigIndices.QuorumApkIndices,
		TotalStakeIndices:            checkSigIndices.TotalStakeIndices,
		NonSignerStakeIndices:        checkSigIndices.NonSignerStakeIndices,
	}, nil
}

// GetQuorumNumbersRequired returns the set of quorum numbers that must be set in the BlobHeader, and verified in
// VerifyCert and CheckDACert.
//
// This method will return required quorum numbers from an internal cache if they are already known for the currently
// active cert verifier. Otherwise, this method will query the required quorum numbers from the currently active
// cert verifier, and cache the result for future use.
func (cv *CertVerifierV3) GetQuorumNumbersRequired(ctx context.Context) ([]uint8, error) {
	blockNumber, err := cv.ethClient.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch block number from eth client: %w", err)
	}

	certVerifierAddress, err := cv.routerAddressProvider.GetCertVerifierAddress(ctx, blockNumber)
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

// getVerifierCallerFromBlockNumber returns a ContractEigenDACertVerifierV3Caller that corresponds to the input reference
// block number.
//
// This method caches ContractEigenDACertVerifierV3Caller instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *CertVerifierV3) getVerifierCallerFromBlockNumber(
	ctx context.Context,
	referenceBlockNumber uint64,
) (*cv_v3_binding.ContractEigenDACertVerifierV3Caller, error) {
	certVerifierAddress, err := cv.routerAddressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get cert verifier address: %w", err)
	}

	return cv.getVerifierCallerFromAddress(certVerifierAddress)
}

// getVerifierCallerFromAddress returns a ContractEigenDACertVerifierV3Caller that corresponds to the input contract
// address
//
// This method caches ContractEigenDACertVerifierV3Caller instances, since their construction requires acquiring a lock
// and parsing json, and is therefore non-trivially expensive.
func (cv *CertVerifierV3) getVerifierCallerFromAddress(
	certVerifierAddress gethcommon.Address,
) (*cv_v3_binding.ContractEigenDACertVerifierV3Caller, error) {
	existingCallerAny, valueExists := cv.verifierCallers.Load(certVerifierAddress)
	if valueExists {
		existingCaller, ok := existingCallerAny.(*cv_v3_binding.ContractEigenDACertVerifierV3Caller)
		if !ok {
			return nil, fmt.Errorf(
				"value in verifierCallers wasn't of type ContractEigenDACertVerifierV3Caller. this should be impossible")
		}
		return existingCaller, nil
	}

	certVerifierCaller, err := cv_v3_binding.NewContractEigenDACertVerifierV3Caller(certVerifierAddress, cv.ethClient)
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
func (cv *CertVerifierV3) GetConfirmationThreshold(ctx context.Context, referenceBlockNumber uint64) (uint8, error) {
	certVerifierAddress, err := cv.routerAddressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
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

	securityThresholds, err := certVerifierCaller.SecurityThresholds(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get security thresholds via contract call: %w", err)
	}

	cv.confirmationThresholds.Store(certVerifierAddress, securityThresholds.ConfirmationThreshold)

	return securityThresholds.ConfirmationThreshold, nil
}

// GetVersion returns the Version that corresponds to the input reference block number.
//
// This method will return the version from an internal cache if it is already known for the cert
// verifier which corresponds to the input reference block number. Otherwise, this method will query the version
// and cache the result for future use.
func (cv *CertVerifierV3) GetVersion(ctx context.Context, referenceBlockNumber uint64) (uint8, error) {
	certVerifierAddress, err := cv.routerAddressProvider.GetCertVerifierAddress(ctx, referenceBlockNumber)
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
	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, referenceBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("get verifier caller from block number: %w", err)
	}

	version, err := certVerifierCaller.Version(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get version via contract call: %w", err)
	}

	cv.versions.Store(certVerifierAddress, version)

	return version, nil
}
