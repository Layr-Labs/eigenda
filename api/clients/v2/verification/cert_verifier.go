package verification

import (
	"context"
	"fmt"
	"time"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/geth"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ICertVerifier is the interface representing a CertVerifier
//
// This interface exists in order to allow verification mocking in unit tests.
type ICertVerifier interface {
	VerifyCertV2(
		ctx context.Context,
		eigenDACert *EigenDACert,
	) error

	GetNonSignerStakesAndSignature(
		ctx context.Context,
		signedBatch *disperser.SignedBatch,
	) (*verifierBindings.NonSignerStakesAndSignature, error)
}

// CertVerifier is responsible for making eth calls against the CertVerifier contract to ensure cryptographic and
// structural integrity of V2 certificates
//
// The cert verifier contract is located at https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDACertVerifier.sol
type CertVerifier struct {
	logger logging.Logger
	// go binding around the EigenDACertVerifier ethereum contract
	certVerifierCaller *verifierBindings.ContractEigenDACertVerifierCaller
	ethClient *geth.EthClient
	pollInterval       time.Duration
}

var _ ICertVerifier = &CertVerifier{}

// NewCertVerifier constructs a CertVerifier
func NewCertVerifier(
	logger logging.Logger,
	// the eth client, which should already be set up
	ethClient *geth.EthClient,
	// the hex address of the EigenDACertVerifier contract
	certVerifierAddress string,
	// pollInterval is how frequently to check latest block number when waiting for the internal eth client to advance
	// to a certain block. This is needed because the RBN in a cert might be further in the future than the internal
	// eth client. In such a case, we must wait for the internal client to catch up to the block number
	// contained in the cert: otherwise, calls will fail.
	//
	// If the configured pollInterval duration is <= 0, then the block number check will be skipped, and calls that
	// rely on the client having reached a certain block number will fail if the internal client is behind.
	pollInterval time.Duration,
) (*CertVerifier, error) {
	verifierCaller, err := verifierBindings.NewContractEigenDACertVerifierCaller(
		gethcommon.HexToAddress(certVerifierAddress),
		ethClient)

	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
	}

	if pollInterval <= time.Duration(0) {
		logger.Warn(
			`CertVerifier poll interval is <= 0. Therefore, any method calls made with this object that 
					rely on the internal client having reached a certain block number will fail if
					the internal client is too far behind.`,
			"pollInterval", pollInterval)
	}

	return &CertVerifier{
		logger:             logger,
		certVerifierCaller: verifierCaller,
		ethClient:          ethClient,
		pollInterval:       pollInterval,
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
	convertedSignedBatch, err := SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return fmt.Errorf("convert signed batch: %w", err)
	}

	convertedBlobInclusionInfo, err := InclusionInfoProtoToBinding(blobInclusionInfo)
	if err != nil {
		return fmt.Errorf("convert blob inclusion info: %w", err)
	}

	err = cv.waitForBlockNumber(ctx, signedBatch.GetHeader().GetReferenceBlockNumber())
	if err != nil {
		return fmt.Errorf("wait for block number: %w", err)
	}

	err = cv.certVerifierCaller.VerifyDACertV2FromSignedBatch(
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
	eigenDACert *EigenDACert,
) error {
	err := cv.waitForBlockNumber(ctx, uint64(eigenDACert.BatchHeader.ReferenceBlockNumber))
	if err != nil {
		return fmt.Errorf("wait for block number: %w", err)
	}

	err = cv.certVerifierCaller.VerifyDACertV2(
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
	signedBatch *disperser.SignedBatch,
) (*verifierBindings.NonSignerStakesAndSignature, error) {

	signedBatchBinding, err := SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return nil, fmt.Errorf("convert signed batch: %w", err)
	}

	err = cv.waitForBlockNumber(ctx, signedBatch.GetHeader().GetReferenceBlockNumber())
	if err != nil {
		return nil, fmt.Errorf("wait for block number: %w", err)
	}

	nonSignerStakesAndSignature, err := cv.certVerifierCaller.GetNonSignerStakesAndSignature(
		&bind.CallOpts{Context: ctx},
		*signedBatchBinding)

	if err != nil {
		return nil, fmt.Errorf("get non signer stakes and signature: %w", err)
	}

	return &nonSignerStakesAndSignature, nil
}

// waitForBlockNumber waits until the internal eth client has advanced to a certain targetBlockNumber.
//
// This method will check the current block number of the internal client every CertVerifier.pollInterval duration.
// It will return nil if the internal client advances to (or past) the targetBlockNumber. It will return an error
// if the input context times out, or if any error occurs when checking the block number of the internal client.
func (cv *CertVerifier) waitForBlockNumber(ctx context.Context, targetBlockNumber uint64) error {
	if cv.pollInterval <= 0 {
		// don't wait for the internal client to advance
		return nil
	}

	ticker := time.NewTicker(cv.pollInterval)
	defer ticker.Stop()

	var actualBlockNumber uint64
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf(
				"timed out waiting for block number %d (latest block number observed was %d): %w",
				targetBlockNumber, actualBlockNumber, ctx.Err())
		case <-ticker.C:
			actualBlockNumber, err := cv.ethClient.BlockNumber(ctx)
			if err != nil {
				return fmt.Errorf("get block number: %w", err)
			}

			if actualBlockNumber >= targetBlockNumber {
				return nil
			}

			cv.logger.Debug(
				"local client is behind the target block number",
				"targetBlockNumber", targetBlockNumber,
				"actualBlockNumber", actualBlockNumber)
		}
	}
}
