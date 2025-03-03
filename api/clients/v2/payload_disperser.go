package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// PayloadDisperser provides the ability to disperse payloads to EigenDA via a Disperser grpc service.
//
// This struct is goroutine safe.
type PayloadDisperser struct {
	logger               logging.Logger
	config               PayloadDisperserConfig
	disperserClient      DisperserClient
	certVerifier         verification.ICertVerifier
	requiredQuorumsStore *RequiredQuorumsStore
}

// NewPayloadDisperser creates a PayloadDisperser from subcomponents that have already been constructed and initialized.
func NewPayloadDisperser(
	logger logging.Logger,
	payloadDisperserConfig PayloadDisperserConfig,
	// IMPORTANT: it is permissible for the disperserClient to be configured without a prover, but operating with this
	// configuration puts a trust assumption on the disperser. With a nil prover, the disperser is responsible for computing
	// the commitments to a blob, and the PayloadDisperser doesn't have a mechanism to verify these commitments.
	//
	// TODO: In the future, an optimized method of commitment verification using fiat shamir transformation will
	//  be implemented. This feature will allow a PayloadDisperser to offload commitment generation onto the
	//  disperser, but the disperser's commitments will be verifiable without needing a full-fledged prover
	disperserClient DisperserClient,
	certVerifier verification.ICertVerifier,
) (*PayloadDisperser, error) {

	err := payloadDisperserConfig.checkAndSetDefaults()
	if err != nil {
		return nil, fmt.Errorf("check and set PayloadDisperserConfig defaults: %w", err)
	}

	requiredQuorumsStore, err := NewRequiredQuorumsStore(certVerifier)
	if err != nil {
		return nil, fmt.Errorf("new required quorums store: %w", err)
	}

	return &PayloadDisperser{
		logger:               logger,
		config:               payloadDisperserConfig,
		disperserClient:      disperserClient,
		certVerifier:         certVerifier,
		requiredQuorumsStore: requiredQuorumsStore,
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
	certVerifierAddress string,
	// payload is the raw data to be stored on eigenDA
	payload *coretypes.Payload,
) (*verification.EigenDACert, error) {
	blob, err := payload.ToBlob(pd.config.PayloadPolynomialForm)
	if err != nil {
		return nil, fmt.Errorf("convert payload to blob: %w", err)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	requiredQuorums, err := pd.requiredQuorumsStore.GetQuorumNumbersRequired(timeoutCtx, certVerifierAddress)
	if err != nil {
		return nil, fmt.Errorf("get quorum numbers required: %w", err)
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.DisperseBlobTimeout)
	defer cancel()

	// TODO (litt3): eventually, we should consider making DisperseBlob accept an actual blob object, instead of the
	//  serialized bytes. The operations taking place in DisperseBlob require the bytes to be converted into field
	//  elements anyway, so serializing the blob here is unnecessary work. This will be a larger change that affects
	//  many areas of code, though.
	blobStatus, blobKey, err := pd.disperserClient.DisperseBlob(
		timeoutCtx,
		blob.Serialize(),
		pd.config.BlobVersion,
		requiredQuorums,
	)
	if err != nil {
		return nil, fmt.Errorf("disperse blob: %w", err)
	}
	pd.logger.Debug("Successful DisperseBlob", "blobStatus", blobStatus.String(), "blobKey", blobKey.Hex())

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.BlobCertifiedTimeout)
	defer cancel()
	blobStatusReply, err := pd.pollBlobStatusUntilCertified(timeoutCtx, blobKey, blobStatus.ToProfobuf())
	if err != nil {
		return nil, fmt.Errorf("poll blob status until certified: %w", err)
	}
	pd.logger.Debug("Blob status CERTIFIED", "blobKey", blobKey.Hex())

	eigenDACert, err := pd.buildEigenDACert(ctx, certVerifierAddress, blobKey, blobStatusReply)
	if err != nil {
		// error returned from method is sufficiently descriptive
		return nil, err
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	err = pd.certVerifier.VerifyCertV2(timeoutCtx, certVerifierAddress, eigenDACert)
	if err != nil {
		return nil, fmt.Errorf("verify cert for blobKey %v: %w", blobKey.Hex(), err)
	}
	pd.logger.Debug("EigenDACert verified", "blobKey", blobKey.Hex())

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
				dispgrpc.BlobStatus_COMPLETE.String(),
				previousStatus.String(),
				ctx.Err())
		case <-ticker.C:
			// This call to the disperser doesn't have a dedicated timeout configured.
			// If this call fails to return in a timely fashion, the timeout configured for the poll loop will trigger
			blobStatusReply, err := pd.disperserClient.GetBlobStatus(ctx, blobKey)
			if err != nil {
				pd.logger.Warn("get blob status", "err", err, "blobKey", blobKey.Hex())
				continue
			}

			newStatus := blobStatusReply.Status
			if newStatus != previousStatus {
				pd.logger.Debug(
					"Blob status changed",
					"blob key", blobKey.Hex(),
					"previous status", previousStatus.String(),
					"new status", newStatus.String())
				previousStatus = newStatus
			}

			// TODO: we'll need to add more in-depth response status processing to derive failover errors
			switch newStatus {
			case dispgrpc.BlobStatus_COMPLETE:
				return blobStatusReply, nil
			case dispgrpc.BlobStatus_QUEUED, dispgrpc.BlobStatus_ENCODED, dispgrpc.BlobStatus_GATHERING_SIGNATURES:
				// TODO (litt): check signing percentage when we are gathering signatures, potentially break
				//  out of this loop early if we have enough signatures
				continue
			default:
				return nil, fmt.Errorf(
					"terminal dispersal failure for blobKey %v. blob status: %v",
					blobKey.Hex(),
					newStatus.String())
			}
		}
	}
}

// buildEigenDACert makes a call to the getNonSignerStakesAndSignature view function on the EigenDACertVerifier
// contract, and then assembles an EigenDACert
func (pd *PayloadDisperser) buildEigenDACert(
	ctx context.Context,
	certVerifierAddress string,
	blobKey core.BlobKey,
	blobStatusReply *dispgrpc.BlobStatusReply,
) (*verification.EigenDACert, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	nonSignerStakesAndSignature, err := pd.certVerifier.GetNonSignerStakesAndSignature(
		timeoutCtx, certVerifierAddress, blobStatusReply.GetSignedBatch())
	if err != nil {
		return nil, fmt.Errorf("get non signer stake and signature: %w", err)
	}
	pd.logger.Debug("Retrieved NonSignerStakesAndSignature", "blobKey", blobKey.Hex())

	eigenDACert, err := verification.BuildEigenDACert(blobStatusReply, nonSignerStakesAndSignature)
	if err != nil {
		return nil, fmt.Errorf("build eigen da cert: %w", err)
	}
	pd.logger.Debug("Constructed EigenDACert", "blobKey", blobKey.Hex())

	return eigenDACert, nil
}
