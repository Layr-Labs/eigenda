package verification

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ICertVerifier is the interface representing a CertVerifier
//
// This interface exists in order to allow verification mocking in unit tests.
type ICertVerifier interface {
	// VerifyCertV2 calls the VerifyCertV2 view function on the EigenDACertVerifier contract.
	//
	// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
	VerifyCertV2(
		ctx context.Context,
		certVerifierAddress string,
		eigenDACert *EigenDACert,
	) error

	// GetNonSignerStakesAndSignature calls the getNonSignerStakesAndSignature view function on the EigenDACertVerifier
	// contract, and returns the resulting NonSignerStakesAndSignature object.
	GetNonSignerStakesAndSignature(
		ctx context.Context,
		certVerifierAddress string,
		signedBatch *disperser.SignedBatch,
	) (*verifierBindings.NonSignerStakesAndSignature, error)

	// GetQuorumNumbersRequired queries the cert verifier contract for the configured set of quorum numbers that must
	// be set in the BlobHeader, and verified in VerifyDACertV2 and verifyDACertV2FromSignedBatch
	GetQuorumNumbersRequired(
		ctx context.Context,
		certVerifierAddress string,
	) ([]uint8, error)
}

// CertVerifier is responsible for making eth calls against the CertVerifier contract to ensure cryptographic and
// structural integrity of V2 certificates
//
// The cert verifier contract is located at https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDACertVerifier.sol
type CertVerifier struct {
	logger       logging.Logger
	ethClient    common.EthClient
	pollInterval time.Duration
	// storage shared between goroutines, containing the most recent block number observed by calling ethClient.BlockNumber()
	latestBlockNumber atomic.Uint64
	// atomic bool, so that only a single goroutine is polling the internal client with BlockNumber() calls at any given time
	pollingActive atomic.Bool
}

var _ ICertVerifier = &CertVerifier{}

// NewCertVerifier constructs a CertVerifier
func NewCertVerifier(
	logger logging.Logger,
	// the eth client, which should already be set up
	ethClient common.EthClient,
	// pollInterval is how frequently to check latest block number when waiting for the internal eth client to advance
	// to a certain block. This is needed because the RBN in a cert might be further in the future than the internal
	// eth client. In such a case, we must wait for the internal client to catch up to the block number
	// contained in the cert: otherwise, calls will fail.
	//
	// If the configured pollInterval duration is <= 0, then the block number check will be skipped, and calls that
	// rely on the client having reached a certain block number will fail if the internal client is behind.
	pollInterval time.Duration,
) (*CertVerifier, error) {
	if pollInterval <= time.Duration(0) {
		logger.Warn(
			`CertVerifier poll interval is <= 0. Therefore, any method calls made with this object that 
					rely on the internal client having reached a certain block number will fail if
					the internal client is too far behind.`,
			"pollInterval", pollInterval)
	}

	return &CertVerifier{
		logger:       logger,
		ethClient:    ethClient,
		pollInterval: pollInterval,
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
	// the hex address of the EigenDACertVerifier contract
	certVerifierAddress string,
	// The signed batch that contains the blob whose cert is being verified. This is obtained from the disperser, and
	// is used to verify that the described blob actually exists in a valid batch.
	signedBatch *disperser.SignedBatch,
	// Contains all necessary information about the blob, so that the cert can be verified.
	blobInclusionInfo *disperser.BlobInclusionInfo,
) error {
	convertedSignedBatch, err := SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return fmt.Errorf("convert signed batch: %w", err)
	}

	convertedBlobInclusionInfo, err := InclusionInfoProtoToBinding(blobInclusionInfo)
	if err != nil {
		return fmt.Errorf("convert blob inclusion info: %w", err)
	}

	err = cv.MaybeWaitForBlockNumber(ctx, signedBatch.GetHeader().GetReferenceBlockNumber())
	if err != nil {
		return fmt.Errorf("wait for block number: %w", err)
	}

	certVerifierCaller, err := verifierBindings.NewContractEigenDACertVerifierCaller(
		gethcommon.HexToAddress(certVerifierAddress),
		cv.ethClient)
	if err != nil {
		return fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
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
func (cv *CertVerifier) VerifyCertV2(
	ctx context.Context,
	// the hex address of the EigenDACertVerifier contract
	certVerifierAddress string,
	eigenDACert *EigenDACert,
) error {
	err := cv.MaybeWaitForBlockNumber(ctx, uint64(eigenDACert.BatchHeader.ReferenceBlockNumber))
	if err != nil {
		return fmt.Errorf("wait for block number: %w", err)
	}

	// don't try to bind to the address until AFTER waiting for the block number. if you try to bind too early, the
	// contract might not exist yet
	certVerifierCaller, err := verifierBindings.NewContractEigenDACertVerifierCaller(
		gethcommon.HexToAddress(certVerifierAddress),
		cv.ethClient)
	if err != nil {
		return fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
	}

	err = certVerifierCaller.VerifyDACertV2(
		&bind.CallOpts{Context: ctx},
		eigenDACert.BatchHeader,
		eigenDACert.BlobInclusionInfo,
		eigenDACert.NonSignerStakesAndSignature)

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
	// the hex address of the EigenDACertVerifier contract
	certVerifierAddress string,
	signedBatch *disperser.SignedBatch,
) (*verifierBindings.NonSignerStakesAndSignature, error) {
	signedBatchBinding, err := SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return nil, fmt.Errorf("convert signed batch: %w", err)
	}

	err = cv.MaybeWaitForBlockNumber(ctx, signedBatch.GetHeader().GetReferenceBlockNumber())
	if err != nil {
		return nil, fmt.Errorf("wait for block number: %w", err)
	}

	certVerifierCaller, err := verifierBindings.NewContractEigenDACertVerifierCaller(
		gethcommon.HexToAddress(certVerifierAddress),
		cv.ethClient)
	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
	}

	nonSignerStakesAndSignature, err := certVerifierCaller.GetNonSignerStakesAndSignature(
		&bind.CallOpts{Context: ctx},
		*signedBatchBinding)

	if err != nil {
		return nil, fmt.Errorf("get non signer stakes and signature: %w", err)
	}

	return &nonSignerStakesAndSignature, nil
}

// GetQuorumNumbersRequired queries the cert verifier contract for the configured set of quorum numbers that must
// be set in the BlobHeader, and verified in VerifyDACertV2 and verifyDACertV2FromSignedBatch
func (cv *CertVerifier) GetQuorumNumbersRequired(ctx context.Context, certVerifierAddress string) ([]uint8, error) {
	certVerifierCaller, err := verifierBindings.NewContractEigenDACertVerifierCaller(
		gethcommon.HexToAddress(certVerifierAddress),
		cv.ethClient)
	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
	}

	quorumNumbersRequired, err := certVerifierCaller.QuorumNumbersRequiredV2(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("get quorum numbers required: %w", err)
	}

	return quorumNumbersRequired, nil
}

// MaybeWaitForBlockNumber waits until the internal eth client has advanced to a certain targetBlockNumber, unless
// configured pollInterval is <= 0, in which case this method will NOT wait for the internal client to advance.
//
// This method will check the current block number of the internal client every CertVerifier.pollInterval duration.
// It will return nil if the internal client advances to (or past) the targetBlockNumber. It will return an error
// if the input context times out, or if any error occurs when checking the block number of the internal client.
//
// This method is synchronized in a way that, if called by multiple goroutines, only a single goroutine will actually
// poll the internal eth client for most recent block number. The goroutine responsible for polling at a given time
// updates an atomic integer, so that all goroutines may check the most recent block without duplicating work.
func (cv *CertVerifier) MaybeWaitForBlockNumber(ctx context.Context, targetBlockNumber uint64) error {
	if cv.pollInterval <= 0 {
		// don't wait for the internal client to advance
		return nil
	}

	if cv.latestBlockNumber.Load() >= targetBlockNumber {
		// immediately return if the local client isn't behind the target block number
		return nil
	}

	ticker := time.NewTicker(cv.pollInterval)
	defer ticker.Stop()

	polling := false
	if cv.pollingActive.CompareAndSwap(false, true) {
		// no other goroutine is currently polling, so assume responsibility
		polling = true
		defer cv.pollingActive.Store(false)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf(
				"timed out waiting for block number %d (latest block number observed was %d): %w",
				targetBlockNumber, cv.latestBlockNumber.Load(), ctx.Err())
		case <-ticker.C:
			if cv.latestBlockNumber.Load() >= targetBlockNumber {
				return nil
			}

			if cv.pollingActive.CompareAndSwap(false, true) {
				// no other goroutine is currently polling, so assume responsibility
				polling = true
				defer cv.pollingActive.Store(false)
			}

			if polling {
				actualBlockNumber, err := cv.ethClient.BlockNumber(ctx)
				if err != nil {
					cv.logger.Debug(
						"ethClient.BlockNumber returned an error",
						"targetBlockNumber", targetBlockNumber,
						"latestBlockNumber", cv.latestBlockNumber.Load(),
						"error", err)

					// tolerate some failures here. if failure continues for too long, it will be caught by the timeout
					continue
				}

				cv.latestBlockNumber.Store(actualBlockNumber)
				if actualBlockNumber >= targetBlockNumber {
					return nil
				}
			}

			cv.logger.Debug(
				"local client is behind the reference block number",
				"targetBlockNumber", targetBlockNumber,
				"actualBlockNumber", cv.latestBlockNumber.Load())
		}
	}
}
