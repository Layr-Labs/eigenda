package payloaddispersal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
)

// PayloadDisperser provides the ability to disperse payloads to EigenDA via a Disperser grpc service.
//
// This struct is goroutine safe.
type PayloadDisperser struct {
	logger          logging.Logger
	config          PayloadDisperserConfig
	disperserClient clients.DisperserClient
	certVerifier    clients.ICertVerifier
	stageTimer      *common.StageTimer
}

// NewPayloadDisperser creates a PayloadDisperser from subcomponents that have already been constructed and initialized.
// If the registry is nil then no metrics will be collected.
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
	disperserClient clients.DisperserClient,
	certVerifier clients.ICertVerifier,
	// if nil, then no metrics will be collected
	registry *prometheus.Registry,
) (*PayloadDisperser, error) {

	err := payloadDisperserConfig.checkAndSetDefaults()
	if err != nil {
		return nil, fmt.Errorf("check and set PayloadDisperserConfig defaults: %w", err)
	}

	stageTimer := common.NewStageTimer(registry, "PayloadDisperser", "SendPayload", false)
	return &PayloadDisperser{
		logger:          logger,
		config:          payloadDisperserConfig,
		disperserClient: disperserClient,
		certVerifier:    certVerifier,
		stageTimer:      stageTimer,
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
	payload *coretypes.Payload,
) (*coretypes.EigenDACert, error) {

	probe := pd.stageTimer.NewSequence()
	defer probe.End()
	probe.SetStage("convert_to_blob")

	blob, err := payload.ToBlob(pd.config.PayloadPolynomialForm)
	if err != nil {
		return nil, fmt.Errorf("convert payload to blob: %w", err)
	}

	probe.SetStage("get_quorums")

	timeoutCtx, cancel := context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	requiredQuorums, err := pd.certVerifier.GetQuorumNumbersRequired(timeoutCtx)
	if err != nil {
		return nil, fmt.Errorf("get quorum numbers required: %w", err)
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.DisperseBlobTimeout)
	defer cancel()

	// TODO (litt3): eventually, we should consider making DisperseBlob accept an actual blob object, instead of the
	//  serialized bytes. The operations taking place in DisperseBlob require the bytes to be converted into field
	//  elements anyway, so serializing the blob here is unnecessary work. This will be a larger change that affects
	//  many areas of code, though.
	blobStatus, blobKey, err := pd.disperserClient.DisperseBlobWithProbe(
		timeoutCtx,
		blob.Serialize(),
		pd.config.BlobVersion,
		requiredQuorums,
		probe)
	if err != nil {
		return nil, fmt.Errorf("disperse blob: %w", err)
	}
	pd.logger.Debug("Successful DisperseBlob", "blobStatus", blobStatus.String(), "blobKey", blobKey.Hex())

	probe.SetStage("QUEUED")

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.BlobCompleteTimeout)
	defer cancel()
	blobStatusReply, err := pd.pollBlobStatusUntilSigned(timeoutCtx, blobKey, blobStatus.ToProfobuf(), probe)
	if err != nil {
		return nil, fmt.Errorf("poll blob status until signed: %w", err)
	}
	pd.logger.Debug("Blob status COMPLETE", "blobKey", blobKey.Hex())

	probe.SetStage("build_cert")

	eigenDACert, err := pd.buildEigenDACert(ctx, blobKey, blobStatusReply)
	if err != nil {
		// error returned from method is sufficiently descriptive
		return nil, err
	}

	probe.SetStage("verify_cert")

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	err = pd.certVerifier.VerifyCertV2(timeoutCtx, eigenDACert)
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

// pollBlobStatusUntilSigned polls the disperser for the status of a blob that has been dispersed
//
// This method will only return a non-nil BlobStatusReply if all quorums meet the required confirmation threshold prior
// to timeout. In all other cases, this method will return a nil BlobStatusReply, along with an error describing the
// failure.
func (pd *PayloadDisperser) pollBlobStatusUntilSigned(
	ctx context.Context,
	blobKey core.BlobKey,
	initialStatus dispgrpc.BlobStatus,
	probe *common.SequenceProbe,
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
				err := checkThresholds(ctx, pd.certVerifier, blobStatusReply, blobKey.Hex())
				if err != nil {
					// returned error is verbose enough, no need to wrap it with additional context
					return nil, err
				}

				return blobStatusReply, nil
			case dispgrpc.BlobStatus_QUEUED, dispgrpc.BlobStatus_ENCODED:
				// Report all non-terminal statuses to the probe. Repeat reports are no-ops.
				probe.SetStage(newStatus.String())
				continue
			case dispgrpc.BlobStatus_GATHERING_SIGNATURES:
				// Report all non-terminal statuses to the probe. Repeat reports are no-ops.
				probe.SetStage(newStatus.String())

				err := checkThresholds(ctx, pd.certVerifier, blobStatusReply, blobKey.Hex())
				if err == nil {
					// If there's no error, then all thresholds are met, so we can stop polling
					return blobStatusReply, nil
				}

				var thresholdNotMetErr *thresholdNotMetError
				if !errors.As(err, &thresholdNotMetErr) {
					// an error occurred which was unrelated to an unmet threshold: something went wrong while checking!
					pd.logger.Warnf("error checking thresholds: %v", err)
				}

				// thresholds weren't met yet. that's ok, since signature gathering is still in progress
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
	blobKey core.BlobKey,
	blobStatusReply *dispgrpc.BlobStatusReply,
) (*coretypes.EigenDACert, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	nonSignerStakesAndSignature, err := pd.certVerifier.GetNonSignerStakesAndSignature(
		timeoutCtx, blobStatusReply.GetSignedBatch())
	if err != nil {
		return nil, fmt.Errorf("get non signer stake and signature: %w", err)
	}
	pd.logger.Debug("Retrieved NonSignerStakesAndSignature", "blobKey", blobKey.Hex())

	eigenDACert, err := coretypes.BuildEigenDACert(blobStatusReply, nonSignerStakesAndSignature)
	if err != nil {
		return nil, fmt.Errorf("build eigen da cert: %w", err)
	}
	pd.logger.Debug("Constructed EigenDACert", "blobKey", blobKey.Hex())

	return eigenDACert, nil
}
