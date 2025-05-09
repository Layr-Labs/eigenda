package verification

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	generic_verifier_binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV3"
	cert_type_binding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	opsr_binding "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenericCertVerifier is responsible for making eth calls against version agnostic CertVerifier contracts to ensure
// cryptographic and structural integrity of EigenDA certificate types.
// The V3 cert verifier contract is located at:
// https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/cert/v3/EigenDACertVerifierV3.sol
type GenericCertVerifier struct {
	logger                  logging.Logger
	ethClient               common.EthClient
	opsrCaller              *opsr_binding.ContractOperatorStateRetrieverCaller
	registryCoordinatorAddr gethcommon.Address
	addressProvider   clients.CertVerifierAddressProvider
	blockNumberMonitor *BlockNumberMonitor


	// Cache maps
	verifierCallers        sync.Map
	requiredQuorums        sync.Map
	confirmationThresholds sync.Map
	versions               sync.Map
}

// Ensure GenericCertVerifier implements the ICertVerifier interface
var _ clients.IGenericCertVerifier = &GenericCertVerifier{}

// NewGenericCertVerifier constructs a new GenericCertVerifier instance
func NewGenericCertVerifier(
	logger logging.Logger,
	ethClient common.EthClient,
	certVerifierAddressProvider clients.CertVerifierAddressProvider,
	registryCoordinatorAddr gethcommon.Address,
	opsrAddr gethcommon.Address,
) (*GenericCertVerifier, error) {
	// Create the Operator State Retriever caller
	opsrCaller, err := opsr_binding.NewContractOperatorStateRetrieverCaller(opsrAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("create operator state retriever caller: %w", err)
	}

	// Create the BlockNumberMonitor
	blockNumberMonitor, err := NewBlockNumberMonitor(logger, ethClient, time.Second * 3)
	if err != nil {
		return nil, fmt.Errorf("create block number monitor: %w", err)
	}


	return &GenericCertVerifier{
		logger:                  logger,
		ethClient:               ethClient,
		blockNumberMonitor: blockNumberMonitor,
		addressProvider:   certVerifierAddressProvider,
		registryCoordinatorAddr: registryCoordinatorAddr,
		opsrCaller:              opsrCaller,
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

// GetNonSignerStakesAndSignature constructs a NonSignerStakesAndSignature object by calling an
// onchain OperatorStateRetriever retriever to fetch necessary non-signer metadata
func (cv *GenericCertVerifier) GetNonSignerStakesAndSignature(
	ctx context.Context,
	signedBatch *disperser.SignedBatch,
) (*cert_type_binding.EigenDATypesV1NonSignerStakesAndSignature, error) {
	// 1 - Ensure that RPC node being used is synced
	//     NOTE: This check is not guaranteed when communicating with node clusters where each node can have
	//           an alternative view of the chain. Adding a retry to the operator state retriever call 
	//           can help mitigate partial async failures.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 6*time.Duration(time.Second * 5))
	defer cancel()
	err := cv.blockNumberMonitor.WaitForBlockNumber(timeoutCtx, signedBatch.GetHeader().GetReferenceBlockNumber())
	if err != nil {
		return nil, fmt.Errorf("wait for block number: %w", err)
	}

	// 2 - Pre-process inputs for operator state retriever call
	signedBatchBinding, err := coretypes.SignedBatchProtoToV2CertBinding(signedBatch)
	if err != nil {
		return nil, fmt.Errorf("convert signed batch: %w", err)
	}

	nonSignerPubKeys := signedBatch.GetAttestation().GetNonSignerPubkeys()


	// 2a - create operator IDs by hashing non-signer public keys
	nonSignerOperatorIDs := make([][32]byte, len(nonSignerPubKeys))
	for i, pubKeySet := range nonSignerPubKeys {
		nonSignerOperatorIDs[i] = crypto.Keccak256Hash(pubKeySet)
	}

	// 2b - cast []uint32 to []byte for quorum numbers
	quorumNumbers := make([]byte, len(signedBatch.GetAttestation().GetQuorumNumbers()))
	for i, qn := range signedBatch.GetAttestation().GetQuorumNumbers() {
		quorumNumbers[i] = byte(qn)
	}

	// use the reference block # from the disperser generated signed batch header
	// for referencing operator states at a specific block checkpoint
	rbn := signedBatch.GetHeader().GetReferenceBlockNumber()

	// 3 - call operator state retriever to fetch signature indices
	checkSigIndices, err := cv.opsrCaller.GetCheckSignaturesIndices(&bind.CallOpts{Context: ctx, BlockNumber: big.NewInt(int64(rbn))},
		cv.registryCoordinatorAddr, uint32(rbn), quorumNumbers, nonSignerOperatorIDs)

	if err != nil {
		return nil, fmt.Errorf("check sig indices call: %w", err)
	}

	// 4 - translate from CertVerifier binding types to cert type
	// TODO: Should probably put SignedBatch into the types directly to avoid this downstream conversion
	nonSignerPubKeysBN254 := make([]cert_type_binding.BN254G1Point, len(signedBatchBinding.Attestation.NonSignerPubkeys))
	for i, pubKeySet := range signedBatchBinding.Attestation.NonSignerPubkeys {
		nonSignerPubKeysBN254[i] = cert_type_binding.BN254G1Point{
			X: pubKeySet.X,
			Y: pubKeySet.Y,
		}
	}

	quorumApksBN254 := make([]cert_type_binding.BN254G1Point, len(signedBatchBinding.Attestation.QuorumApks))
	for i, apkSet := range signedBatchBinding.Attestation.QuorumApks {
		quorumApksBN254[i] = cert_type_binding.BN254G1Point{
			X: apkSet.X,
			Y: apkSet.Y,
		}
	}

	apkG2BN254 := cert_type_binding.BN254G2Point{
		X: signedBatchBinding.Attestation.ApkG2.X,
		Y: signedBatchBinding.Attestation.ApkG2.Y,
	}

	sigmaBN254 := cert_type_binding.BN254G1Point{
		X: signedBatchBinding.Attestation.Sigma.X,
		Y: signedBatchBinding.Attestation.Sigma.Y,
	}


	// 5 - construct non signer stakes and signature
	return &cert_type_binding.EigenDATypesV1NonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: checkSigIndices.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:             nonSignerPubKeysBN254,
		QuorumApks:                   quorumApksBN254,
		ApkG2:                        apkG2BN254,
		Sigma:                        sigmaBN254,
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
) (*generic_verifier_binding.ContractEigenDACertVerifierV3, error) {
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
) (*generic_verifier_binding.ContractEigenDACertVerifierV3, error) {
	existingCallerAny, valueExists := cv.verifierCallers.Load(certVerifierAddress)
	if valueExists {
		existingCaller, ok := existingCallerAny.(*generic_verifier_binding.ContractEigenDACertVerifierV3)
		if !ok {
			return nil, fmt.Errorf(
				"value in verifierCallers wasn't of type ContractEigenDACertVerifierV3. this should be impossible")
		}
		return existingCaller, nil
	}

	certVerifierCaller, err := generic_verifier_binding.NewContractEigenDACertVerifierV3(certVerifierAddress, cv.ethClient)
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
	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, referenceBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("get verifier caller from block number: %w", err)
	}

	version, err := certVerifierCaller.CertVersion(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get version via contract call: %w", err)
	}

	cv.versions.Store(certVerifierAddress, version)

	return version, nil
}
