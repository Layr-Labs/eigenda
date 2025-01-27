package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/geth"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// PayloadDisperser provides the ability to disperse payloads to EigenDA via a Disperser grpc service.
//
// This struct is goroutine safe.
type PayloadDisperser struct {
	logger          logging.Logger
	config          PayloadDisperserConfig
	codec           codecs.BlobCodec
	disperserClient DisperserClient
	certVerifier    verification.ICertVerifier
}

// BuildPayloadDisperser builds a PayloadDisperser from config structs
func BuildPayloadDisperser(
	logger logging.Logger,
	payloadDisperserConfig PayloadDisperserConfig,
	disperserClientConfig DisperserClientConfig,
	// signer to sign blob dispersal requests
	signer core.BlobRequestSigner,
	// prover is used to compute commitments to a new blob during the dispersal process
	//
	// IMPORTANT: it is permissible for the prover parameter to be nil, but operating with this configuration
	// puts a trust assumption on the disperser. With a nil prover, the disperser is responsible for computing
	// the commitments to a blob, and the PayloadDisperser doesn't have a mechanism to verify these commitments.
	//
	// TODO: In the future, an optimized method of commitment verification using fiat shamir transformation will
	//  be implemented. This feature will allow a PayloadDisperser to offload commitment generation onto the
	//  disperser, but the disperser's commitments will be verifiable without needing a full-fledged prover
	prover encoding.Prover,
	accountant *Accountant,
	ethConfig geth.EthClientConfig,
) (*PayloadDisperser, error) {

	codec, err := codecs.CreateCodec(
		payloadDisperserConfig.PayloadPolynomialForm,
		payloadDisperserConfig.BlobEncodingVersion)
	if err != nil {
		return nil, fmt.Errorf("create codec: %w", err)
	}

	disperserClient, err := NewDisperserClient(&disperserClientConfig, signer, prover, accountant)
	if err != nil {
		return nil, fmt.Errorf("new disperser client: %s", err)
	}

	ethClient, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, logger)
	if err != nil {
		return nil, fmt.Errorf("new eth client: %w", err)
	}

	certVerifier, err := verification.NewCertVerifier(ethClient, payloadDisperserConfig.EigenDACertVerifierAddr)
	if err != nil {
		return nil, fmt.Errorf("new cert verifier: %w", err)
	}

	return &PayloadDisperser{
		logger:          logger,
		config:          payloadDisperserConfig,
		codec:           codec,
		disperserClient: disperserClient,
		certVerifier:    certVerifier,
	}, nil
}

// SendPayload executes the dispersal of a payload, with these steps:
//
//  1. Encode payload into a blob
//  2. Disperse the blob
//  3. Poll the disperser with GetBlobStatus until a terminal status is reached, or until the polling timeout is reached
//  4. Construct an EigenDACert if dispersal is successful
//  5. Verify the constructed cert with an eth_call to the EigenDACertVerifier contract
//  6. Return the valid cert
func (pd *PayloadDisperser) SendPayload(
	ctx context.Context,
	// payload is the raw data to be stored on eigenDA
	payload []byte,
	// salt is added while constructing the blob header
	// This salt should be utilized if a blob dispersal fails, in order to retry dispersing the same payload under a
	// different blob key, when using reserved bandwidth payments.
	salt uint32,
) (*verification.EigenDACert, error) {

	blobBytes, err := pd.codec.EncodeBlob(payload)
	if err != nil {
		return nil, fmt.Errorf("encode payload to blob: %w", err)
	}
	pd.logger.Debug("Payload encoded to blob")

	timeoutCtx, cancel := context.WithTimeout(ctx, pd.config.DisperseBlobTimeout)
	defer cancel()
	blobStatus, blobKey, err := pd.disperserClient.DisperseBlob(
		timeoutCtx,
		blobBytes,
		pd.config.BlobVersion,
		pd.config.Quorums,
		salt)
	if err != nil {
		return nil, fmt.Errorf("disperse blob: %w", err)
	}
	pd.logger.Debug("Successful DisperseBlob", "blobStatus", blobStatus.String(), "blobKey", blobKey)

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.BlobCertifiedTimeout)
	defer cancel()
	blobStatusReply, err := pd.pollBlobStatusUntilCertified(timeoutCtx, blobKey, blobStatus.ToProfobuf())
	if err != nil {
		return nil, fmt.Errorf("poll blob status until certified: %w", err)
	}
	pd.logger.Debug("Blob status CERTIFIED", "blobKey", blobKey)

	eigenDACert, err := pd.buildEigenDACert(ctx, blobKey, blobStatusReply)
	if err != nil {
		// error returned from method is sufficiently descriptive
		return nil, err
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	err = pd.certVerifier.VerifyCertV2(timeoutCtx, eigenDACert)
	if err != nil {
		return nil, fmt.Errorf("verify cert for blobKey %v: %w", blobKey, err)
	}
	pd.logger.Debug("EigenDACert verified", "blobKey", blobKey)

	return eigenDACert, nil
}

// Close is responsible for calling close on all internal clients. This method will do its best to close all internal
// clients, even if some closes fail.
//
// Any and all errors returned from closing internal clients will be joined and returned.
//
// This method should only be called once.
func (pd *PayloadDisperser) Close() error {
	err := pd.disperserClient.Close()
	if err != nil {
		return fmt.Errorf("close disperser client: %w", err)
	}

	return nil
}

// pollBlobStatusUntilCertified polls the disperser for the status of a blob that has been dispersed
//
// This method will only return a non-nil BlobStatusReply if the blob is reported to be CERTIFIED prior to the timeout.
// In all other cases, this method will return a nil BlobStatusReply, along with an error describing the failure.
func (pd *PayloadDisperser) pollBlobStatusUntilCertified(
	ctx context.Context,
	blobKey core.BlobKey,
	initialStatus dispgrpc.BlobStatus,
) (*dispgrpc.BlobStatusReply, error) {

	previousStatus := initialStatus

	ticker := time.NewTicker(pd.config.BlobStatusPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf(
				"timed out waiting for %v blob status, final status was %v: %w",
				dispgrpc.BlobStatus_CERTIFIED.Descriptor(),
				previousStatus.Descriptor(),
				ctx.Err())
		case <-ticker.C:
			// This call to the disperser doesn't have a dedicated timeout configured.
			// If this call fails to return in a timely fashion, the timeout configured for the poll loop will trigger
			blobStatusReply, err := pd.disperserClient.GetBlobStatus(ctx, blobKey)
			if err != nil {
				pd.logger.Warn("get blob status", "err", err, "blobKey", blobKey)
				continue
			}

			newStatus := blobStatusReply.Status
			if newStatus != previousStatus {
				pd.logger.Debug(
					"Blob status changed",
					"blob key", blobKey,
					"previous status", previousStatus.Descriptor(),
					"new status", newStatus.Descriptor())
				previousStatus = newStatus
			}

			switch newStatus {
			case dispgrpc.BlobStatus_CERTIFIED:
				return blobStatusReply, nil
			case dispgrpc.BlobStatus_QUEUED, dispgrpc.BlobStatus_ENCODED:
				continue
			default:
				return nil, fmt.Errorf(
					"terminal dispersal failure for blobKey %v. blob status: %v",
					blobKey,
					newStatus.Descriptor())
			}
		}
	}
}

// buildEigenDACert makes a call to the getNonSignerStakesAndSignature view function on the EigenDACertVerifier
// contract, and then assembles an EigenDACert
func (pd *PayloadDisperser) buildEigenDACert(
	ctx context.Context,
	blobKey core.BlobKey,
	blobStatusReply *dispgrpc.BlobStatusReply,
) (*verification.EigenDACert, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	nonSignerStakesAndSignature, err := pd.certVerifier.GetNonSignerStakesAndSignature(
		timeoutCtx, blobStatusReply.GetSignedBatch())
	if err != nil {
		return nil, fmt.Errorf("get non signer stake and signature: %w", err)
	}
	pd.logger.Debug("Retrieved NonSignerStakesAndSignature", "blobKey", blobKey)

	eigenDACert, err := verification.BuildEigenDACert(blobStatusReply, nonSignerStakesAndSignature)
	if err != nil {
		return nil, fmt.Errorf("build eigen da cert: %w", err)
	}
	pd.logger.Debug("Constructed EigenDACert", "blobKey", blobKey)

	return eigenDACert, nil
}
