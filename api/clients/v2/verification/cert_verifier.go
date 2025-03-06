package verification

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
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
	blockNumberProvider         *BlockNumberProvider
	certVerifierAddressProvider clients.CertVerifierAddressProvider
	// maps contract address to a ContractEigenDACertVerifierCaller object
	verifierCallers sync.Map
	// maps cert verifier address string to set of required quorums specified in the contract at that address
	requiredQuorums sync.Map
}

var _ clients.ICertVerifier = &CertVerifier{}

// NewCertVerifier constructs a CertVerifier
func NewCertVerifier(
	logger logging.Logger,
	// the eth client, which should already be set up
	ethClient common.EthClient,
	certVerifierAddressProvider clients.CertVerifierAddressProvider,
	// pollInterval is how frequently to check latest block number when waiting for the internal eth client to advance
	// to a certain block. This is needed because the RBN in a cert might be further in the future than the internal
	// eth client. In such a case, we must wait for the internal client to catch up to the block number
	// contained in the cert: otherwise, calls will fail.
	//
	// If the configured pollInterval duration is <= 0, then the block number check will be skipped, and calls that
	// rely on the client having reached a certain block number will fail if the internal client is behind.
	pollInterval time.Duration,
) (*CertVerifier, error) {
	blockNumberProvider := NewBlockNumberProvider(logger, ethClient, pollInterval)

	return &CertVerifier{
		logger:                      logger,
		ethClient:                   ethClient,
		certVerifierAddressProvider: certVerifierAddressProvider,
		blockNumberProvider:         blockNumberProvider,
	}, nil
}

// VerifyCertV2FromSignedBatch calls the verifyDACertV2FromSignedBatch view function on the EigenDACertVerifier contract.
//
// Before verifying the cert, this method will wait for the internal client to advance to a sufficient block height.
// This wait will time out if the duration exceeds the timeout configured for the input ctx parameter. If
// CertVerifier.pollInterval is configured to be <= 0, then this method will *not* wait for the internal client to
// advance, and will instead simply fail verification if the internal client is behind.
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

	blockNumber := signedBatch.GetHeader().GetReferenceBlockNumber()

	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, blockNumber)
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
// Before verifying the cert, this method will wait for the internal client to advance to a sufficient block height.
// This wait will time out if the duration exceeds the timeout configured for the input ctx parameter. If
// CertVerifier.pollInterval is configured to be <= 0, then this method will *not* wait for the internal client to
// advance, and will instead simply fail verification if the internal client is behind.
//
// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
func (cv *CertVerifier) VerifyCertV2(ctx context.Context, eigenDACert *coretypes.EigenDACert) error {
	blockNumber := uint64(eigenDACert.BatchHeader.ReferenceBlockNumber)

	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, blockNumber)
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
//
// Before getting the NonSignerStakesAndSignature, this method will wait for the internal client to advance to a
// sufficient block height. This wait will time out if the duration exceeds the timeout configured for the input ctx
// parameter. If CertVerifier.pollInterval is configured to be <= 0, then this method will *not* wait for the internal
// client to advance, and will instead simply fail to get the NonSignerStakesAndSignature if the internal client is
// behind.
func (cv *CertVerifier) GetNonSignerStakesAndSignature(
	ctx context.Context,
	signedBatch *disperser.SignedBatch,
) (*verifierBindings.NonSignerStakesAndSignature, error) {
	signedBatchBinding, err := coretypes.SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return nil, fmt.Errorf("convert signed batch: %w", err)
	}

	blockNumber := signedBatch.GetHeader().GetReferenceBlockNumber()

	certVerifierCaller, err := cv.getVerifierCallerFromBlockNumber(ctx, blockNumber)
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
	blockNumber, err := cv.blockNumberProvider.FetchLatestBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch latest block number: %w", err)
	}

	certVerifierAddress, err := cv.certVerifierAddressProvider.GetCertVerifierAddress(blockNumber)
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

// getVerifierCallerFromBlockNumber returns a ContractEigenDACertVerifierCaller that corresponds to the input block
// number
//
// If the eth node hasn't yet advanced to the input block number, this method will wait until that block is reached
// before attempting to get the verifier caller.
//
// This method caches ContractEigenDACertVerifierCaller instances, since their construction requires acquiring a lock
// and parsing json, and is therefore not trivially inexpensive.
func (cv *CertVerifier) getVerifierCallerFromBlockNumber(
	ctx context.Context,
	blockNumber uint64,
) (*verifierBindings.ContractEigenDACertVerifierCaller, error) {
	err := cv.blockNumberProvider.MaybeWaitForBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("wait for block number: %w", err)
	}

	certVerifierAddress, err := cv.certVerifierAddressProvider.GetCertVerifierAddress(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("get cert verifier address: %w", err)
	}

	return cv.getVerifierCallerFromAddress(certVerifierAddress)
}

// getVerifierCallerFromAddress returns a ContractEigenDACertVerifierCaller that corresponds to the input contract
// address
//
// This method caches ContractEigenDACertVerifierCaller instances, since their construction requires acquiring a lock
// and parsing json, and is therefore not trivially inexpensive.
func (cv *CertVerifier) getVerifierCallerFromAddress(
	certVerifierAddress string,
) (*verifierBindings.ContractEigenDACertVerifierCaller, error) {
	existingCallerAny, valueExists := cv.verifierCallers.Load(certVerifierAddress)
	if valueExists {
		existingCaller, ok := existingCallerAny.(*verifierBindings.ContractEigenDACertVerifierCaller)
		if !ok {
			return nil, fmt.Errorf(
				"value in verifierCallers wasn't of type ContractEigenDACertVerifierCaller. this should be impossible")
		}
		return existingCaller, nil
	}

	certVerifierCaller, err := verifierBindings.NewContractEigenDACertVerifierCaller(
		gethcommon.HexToAddress(certVerifierAddress), cv.ethClient)
	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
	}

	cv.verifierCallers.Store(certVerifierAddress, certVerifierCaller)
	return certVerifierCaller, nil
}
