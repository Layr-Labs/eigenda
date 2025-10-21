package payloaddispersal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal")

// withClientSpan wraps an external call with a client span for better observability.
// It automatically handles span creation, error recording, and status setting.
func withClientSpan[T any](
	ctx context.Context,
	name string,
	fn func(context.Context) (T, error),
	attrs ...attribute.KeyValue,
) (T, error) {
	ctx, span := tracer.Start(ctx, name,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...))
	defer span.End()

	result, err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "ok")
	}
	return result, err
}

// PayloadDisperser provides the ability to disperse payloads to EigenDA via a Disperser grpc service.
//
// This struct is goroutine safe.
type PayloadDisperser struct {
	logger          logging.Logger
	config          PayloadDisperserConfig
	disperserClient *clients.DisperserClient
	blockMonitor    *verification.BlockNumberMonitor
	certBuilder     *clients.CertBuilder
	certVerifier    *verification.CertVerifier
	stageTimer      *common.StageTimer
	clientLedger    *clientledger.ClientLedger
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
	disperserClient *clients.DisperserClient,
	blockMonitor *verification.BlockNumberMonitor,
	certBuilder *clients.CertBuilder,
	certVerifier *verification.CertVerifier,
	// Manages payment state for the client. May be nil for legacy payment mode.
	clientLedger *clientledger.ClientLedger,
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
		blockMonitor:    blockMonitor,
		certBuilder:     certBuilder,
		certVerifier:    certVerifier,
		stageTimer:      stageTimer,
		clientLedger:    clientLedger,
	}, nil
}

// SendPayload executes the dispersal of a payload, with these steps:
//
//  1. Encode payload into a blob
//  2. Disperse the blob
//  3. Poll the disperser with GetBlobStatus until a terminal status is reached, or until the polling timeout is reached
//  4. Construct an EigenDACert if dispersal is successful
//  5. Verify the constructed cert via an eth_call to the EigenDACertVerifier contract
//  6. Return the valid cert
func (pd *PayloadDisperser) SendPayload(
	ctx context.Context,
	// payload is the raw data to be stored on eigenDA
	payload coretypes.Payload,
) (coretypes.EigenDACert, error) {
	ctx, span := tracer.Start(ctx, "PayloadDisperser.SendPayload",
		trace.WithAttributes(
			attribute.Int("payload_size_bytes", len(payload)),
			attribute.Int("blob_version", int(pd.config.BlobVersion)),
			attribute.Int("payload_form", int(pd.config.PayloadPolynomialForm)),
			attribute.Bool("payment.new_ledger", pd.clientLedger != nil),
		))
	defer span.End()

	probe := pd.stageTimer.NewSequence()
	defer probe.End()
	probe.SetStage("convert_to_blob")
	span.AddEvent("stage", trace.WithAttributes(attribute.String("status", "convert_to_blob")))

	// convert the payload into an EigenDA blob by interpreting the payload in polynomial form,
	// which means the encoded payload will need to be IFFT'd since EigenDA blobs are in coefficient form.
	blob, err := payload.ToBlob(pd.config.PayloadPolynomialForm)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert payload to blob")
		return nil, fmt.Errorf("failed to convert payload to blob: %w", err)
	}

	probe.SetStage("get_quorums")
	span.AddEvent("stage", trace.WithAttributes(attribute.String("status", "get_quorums")))

	timeoutCtx, cancel := context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()

	// NOTE: there is a synchronization edge case where the disperser accredits an RBN that correlates
	//       to a newly added immutable CertVerifier under the Router contract design. Resulting in
	//       potentially a few failed dispersals until the RBN advances; guaranteeing eventual consistency.
	//       This is a known issue and will be addressed with future enhancements.
	requiredQuorums, err := withClientSpan(timeoutCtx, "CertVerifier.GetQuorumNumbersRequired",
		func(c context.Context) ([]core.QuorumID, error) {
			return pd.certVerifier.GetQuorumNumbersRequired(c)
		},
		semconv.RPCSystemKey.String("jsonrpc"),
		attribute.String("component", "certVerifier"),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get quorum numbers required failed")
		return nil, fmt.Errorf("get quorum numbers required: %w", err)
	}

	span.SetAttributes(attribute.IntSlice("quorums", intSliceFromQuorums(requiredQuorums)))

	symbolCount := blob.LenSymbols()
	span.SetAttributes(attribute.Int("symbol_count", int(symbolCount)))

	var paymentMetadata *core.PaymentMetadata
	if pd.clientLedger != nil {
		// we are using the new payment system if clientLedger is non nil
		probe.SetStage("debit")
		span.AddEvent("stage", trace.WithAttributes(attribute.String("status", "debit")))

		paymentMetadata, err = withClientSpan(ctx, "ClientLedger.Debit",
			func(c context.Context) (*core.PaymentMetadata, error) {
				return pd.clientLedger.Debit(c, symbolCount, requiredQuorums)
			},
			attribute.String("component", "clientLedger"),
			attribute.Int("symbol_count", int(symbolCount)),
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "debit failed")
			return nil, fmt.Errorf("debit: %w", err)
		}

		if paymentMetadata != nil {
			span.SetAttributes(
				attribute.String("payment.account_id", paymentMetadata.AccountID.Hex()),
				attribute.Int64("payment.timestamp", paymentMetadata.Timestamp),
			)
		}
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.DisperseBlobTimeout)
	defer cancel()

	// TODO (litt3): DisperseBlob should accept an actual blob object, instead of the
	//  serialized bytes. The operations taking place in DisperseBlob require the bytes to be converted into field
	//  elements anyway, so serializing the blob here is unnecessary work. This will be a larger change that affects
	//  many areas of code, though.
	blobHeader, reply, err := pd.disperserClient.DisperseBlob(
		timeoutCtx,
		blob.Serialize(),
		pd.config.BlobVersion,
		requiredQuorums,
		probe,
		paymentMetadata)
	if err != nil {
		if pd.clientLedger != nil {
			_, revertErr := withClientSpan(ctx, "ClientLedger.RevertDebit",
				func(c context.Context) (struct{}, error) {
					return struct{}{}, pd.clientLedger.RevertDebit(c, paymentMetadata, symbolCount)
				},
				attribute.String("component", "clientLedger"),
			)
			if revertErr != nil {
				span.RecordError(errors.Join(err, revertErr))
				span.SetStatus(codes.Error, "disperse blob and revert debit failed")
				return nil, fmt.Errorf("disperse blob and revert debit: %w", errors.Join(err, revertErr))
			}
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, "disperse blob failed")
		return nil, fmt.Errorf("disperse blob: %w", err)
	}

	probe.SetStage("verify_blob_key")
	span.AddEvent("stage", trace.WithAttributes(attribute.String("status", "verify_blob_key")))

	blobKey, err := verifyReceivedBlobKey(blobHeader, reply)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "verify received blob key failed")
		return nil, fmt.Errorf("verify received blob key: %w", err)
	}

	span.SetAttributes(attribute.String("blob_key", blobKey.Hex()))

	cert, err := pd.buildEigenDACert(ctx, reply.GetResult(), blobKey, probe)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "build EigenDA cert failed")
		return nil, err
	}

	span.SetStatus(codes.Ok, "payload sent successfully")
	return cert, nil
}

// Waits for a blob to be signed, and builds the EigenDA cert with the operator signatures
//
// If the blob does not become fully signed before the BlobCompleteTimeout timeout elapses, returns an error
func (pd *PayloadDisperser) buildEigenDACert(
	ctx context.Context,
	initialBlobStatus dispgrpc.BlobStatus,
	blobKey corev2.BlobKey,
	probe *common.SequenceProbe,
) (coretypes.EigenDACert, error) {
	ctx, span := tracer.Start(ctx, "PayloadDisperser.buildEigenDACert",
		trace.WithAttributes(
			attribute.String("blob_key", blobKey.Hex()),
			attribute.String("initial_status", initialBlobStatus.String()),
		))
	defer span.End()

	probe.SetStage("QUEUED")
	span.AddEvent("stage", trace.WithAttributes(attribute.String("status", "QUEUED")))

	// poll the disperser for the status of the blob until it's received adequate signatures in regards to
	// confirmation thresholds, a terminal error, or a timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, pd.config.BlobCompleteTimeout)
	defer cancel()
	blobStatusReply, err := pd.pollBlobStatusUntilSigned(timeoutCtx, blobKey, initialBlobStatus, probe)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "poll blob status until signed failed")
		return nil, fmt.Errorf("poll blob status until signed: %w", err)
	}

	pd.logSigningPercentages(blobKey, blobStatusReply)
	addSigningPercentagesToSpan(span, blobStatusReply)

	rbn := blobStatusReply.GetSignedBatch().GetHeader().GetReferenceBlockNumber()
	span.SetAttributes(attribute.Int64("reference_block_number", int64(rbn)))

	probe.SetStage("wait_for_block_number")
	span.AddEvent("stage", trace.WithAttributes(attribute.String("status", "wait_for_block_number")))
	// TODO: given the repeated context timeout declaration in this method we should consider creating some
	// generic function or helper to enhance DRY
	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	_, waitErr := withClientSpan(timeoutCtx, "BlockMonitor.WaitForBlockNumber",
		func(c context.Context) (struct{}, error) {
			return struct{}{}, pd.blockMonitor.WaitForBlockNumber(c, rbn)
		},
		attribute.String("component", "blockMonitor"),
		attribute.Int64("block_number", int64(rbn)),
	)
	if waitErr != nil {
		span.RecordError(waitErr)
		span.SetStatus(codes.Error, "wait for block number failed")
		return nil, fmt.Errorf("wait for block number: %w", waitErr)
	}

	certVersion, err := withClientSpan(ctx, "CertVerifier.GetCertVersion",
		func(c context.Context) (coretypes.CertificateVersion, error) {
			return pd.certVerifier.GetCertVersion(c, rbn)
		},
		semconv.RPCSystemKey.String("jsonrpc"),
		attribute.String("component", "certVerifier"),
		attribute.Int64("block_number", int64(rbn)),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get certificate version failed")
		return nil, fmt.Errorf("get certificate version: %w", err)
	}

	span.SetAttributes(attribute.Int("cert_version", int(certVersion)))

	probe.SetStage("build_cert")
	span.AddEvent("stage", trace.WithAttributes(attribute.String("status", "build_cert")))
	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()
	eigenDACert, err := withClientSpan(timeoutCtx, "CertBuilder.BuildCert",
		func(c context.Context) (coretypes.EigenDACert, error) {
			return pd.certBuilder.BuildCert(c, certVersion, blobStatusReply)
		},
		attribute.String("component", "certBuilder"),
		attribute.Int("cert_version", int(certVersion)),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "build cert failed")
		return nil, fmt.Errorf("build cert: %w", err)
	}
	pd.logger.Debug("EigenDACert built", "blobKey", blobKey.Hex(), "certVersion", certVersion)

	probe.SetStage("verify_cert")
	span.AddEvent("stage", trace.WithAttributes(attribute.String("status", "verify_cert")))

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.ContractCallTimeout)
	defer cancel()

	_, checkErr := withClientSpan(timeoutCtx, "CertVerifier.CheckDACert",
		func(c context.Context) (struct{}, error) {
			return struct{}{}, pd.certVerifier.CheckDACert(c, eigenDACert)
		},
		semconv.RPCSystemKey.String("jsonrpc"),
		attribute.String("component", "certVerifier"),
	)
	if checkErr != nil {
		var errInvalidCert *verification.CertVerifierInvalidCertError
		if errors.As(checkErr, &errInvalidCert) {
			// Regardless of whether the cert is invalid (400) or certVerifier contract has a bug (500),
			// we send a failover signal. If we can't construct a valid cert after retrying a few times (proxy retry
			// policy), then its safer for the rollup to failover to another DA layer.
			span.RecordError(checkErr)
			span.SetStatus(codes.Error, "check DA cert invalid")
			return nil, api.NewErrorFailover(fmt.Errorf("checkDACert failed with blobKey %v: %w", blobKey.Hex(), checkErr))
		}
		span.RecordError(checkErr)
		span.SetStatus(codes.Error, "verify cert failed")
		return nil, fmt.Errorf("verify cert for blobKey %v: %w", blobKey.Hex(), checkErr)
	}

	pd.logger.Debug("EigenDACert verified", "blobKey", blobKey.Hex())

	span.SetStatus(codes.Ok, "cert built and verified successfully")
	return eigenDACert, nil
}

// logSigningPercentages logs the signing percentage of each quorum for a blob that has been dispersed and satisfied
// required signing thresholds
func (pd *PayloadDisperser) logSigningPercentages(blobKey corev2.BlobKey, blobStatusReply *dispgrpc.BlobStatusReply) {
	attestation := blobStatusReply.GetSignedBatch().GetAttestation()
	if len(attestation.GetQuorumNumbers()) != len(attestation.GetQuorumSignedPercentages()) {
		pd.logger.Error("quorum number count and signed percentage count don't match. This should never happen",
			"blobKey", blobKey.Hex(),
			"quorumNumberCount", len(attestation.GetQuorumNumbers()),
			"signedPercentageCount", len(attestation.GetQuorumSignedPercentages()))
	}

	quorumPercentagesBuilder := strings.Builder{}
	quorumPercentagesBuilder.WriteString("(")

	for index, quorumNumber := range attestation.GetQuorumNumbers() {
		quorumPercentagesBuilder.WriteString(
			fmt.Sprintf("quorum_%d: %d%%, ", quorumNumber, attestation.GetQuorumSignedPercentages()[index]))
	}
	quorumPercentagesBuilder.WriteString(")")

	pd.logger.Debug("Blob signed",
		"blobKey", blobKey.Hex(), "quorumPercentages", quorumPercentagesBuilder.String())
}

// addSigningPercentagesToSpan adds the signing percentages of each quorum to the span attributes
func addSigningPercentagesToSpan(span trace.Span, blobStatusReply *dispgrpc.BlobStatusReply) {
	attestation := blobStatusReply.GetSignedBatch().GetAttestation()
	if len(attestation.GetQuorumNumbers()) != len(attestation.GetQuorumSignedPercentages()) {
		return
	}

	for index, quorumNumber := range attestation.GetQuorumNumbers() {
		span.SetAttributes(attribute.Int(
			fmt.Sprintf("signing.q%d_pct", quorumNumber),
			int(attestation.GetQuorumSignedPercentages()[index]),
		))
	}
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
	blobKey corev2.BlobKey,
	initialStatus dispgrpc.BlobStatus,
	probe *common.SequenceProbe,
) (*dispgrpc.BlobStatusReply, error) {
	ctx, span := tracer.Start(ctx, "PayloadDisperser.pollBlobStatusUntilSigned",
		trace.WithAttributes(
			attribute.String("blob_key", blobKey.Hex()),
			attribute.String("initial_status", initialStatus.String()),
		))
	defer span.End()

	previousStatus := initialStatus

	ticker := time.NewTicker(pd.config.BlobStatusPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Failover to another DA layer because EigenDA is not completing its signing duty in time.
			err := api.NewErrorFailover(fmt.Errorf(
				"timed out waiting for %v blob status, final status was %v: %w",
				dispgrpc.BlobStatus_COMPLETE.String(),
				previousStatus.String(),
				ctx.Err()))
			span.RecordError(err)
			span.SetStatus(codes.Error, "timeout waiting for blob status")
			return nil, err
		case <-ticker.C:
			// This call to the disperser doesn't have a dedicated timeout configured.
			// If this call fails to return in a timely fashion, the timeout configured for the poll loop will trigger
			blobStatusReply, err := pd.disperserClient.GetBlobStatus(ctx, blobKey)
			if err != nil {
				// this is expected to fail multiple times before we get a valid response, so only do a Debug log
				pd.logger.Debug("get blob status", "err", err, "blobKey", blobKey.Hex())
				continue
			}

			newStatus := blobStatusReply.GetStatus()
			if newStatus != previousStatus {
				pd.logger.Debug(
					"Blob status changed",
					"blob key", blobKey.Hex(),
					"previous status", previousStatus.String(),
					"new status", newStatus.String())
				span.AddEvent("status_change", trace.WithAttributes(
					attribute.String("previous_status", previousStatus.String()),
					attribute.String("new_status", newStatus.String()),
				))
				previousStatus = newStatus
			}

			switch newStatus {
			case dispgrpc.BlobStatus_COMPLETE:
				err := checkThresholds(ctx, pd.certVerifier, blobStatusReply, blobKey.Hex())
				if err != nil {
					// TODO(samlaf): checkThresholds should return more fine-grained errors
					// For now, we only failover if thresholds were unmet, not anything else.
					// The risk of failing over for everything is that eth-rpc calls could fail
					// for networking reasons, which we don't want to failover to eth for!
					var thresholdNotMetErr *thresholdNotMetError
					if errors.As(err, &thresholdNotMetErr) {
						span.RecordError(err)
						span.SetStatus(codes.Error, "threshold not met")
						return nil, api.NewErrorFailover(fmt.Errorf("check thresholds: %w", err))
					}
					span.RecordError(err)
					span.SetStatus(codes.Error, "check thresholds failed")
					return nil, fmt.Errorf("check thresholds: %w", err)
				}

				span.SetAttributes(attribute.String("final_status", newStatus.String()))
				span.SetStatus(codes.Ok, "blob signed successfully")
				return blobStatusReply, nil
			case dispgrpc.BlobStatus_QUEUED, dispgrpc.BlobStatus_ENCODED:
				// Report all non-terminal statuses to the probe. Repeat reports are no-ops.
				probe.SetStage(newStatus.String())
				span.AddEvent("stage", trace.WithAttributes(attribute.String("status", newStatus.String())))
				continue
			case dispgrpc.BlobStatus_GATHERING_SIGNATURES:
				// Report all non-terminal statuses to the probe. Repeat reports are no-ops.
				probe.SetStage(newStatus.String())
				span.AddEvent("stage", trace.WithAttributes(attribute.String("status", newStatus.String())))

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
				// Failover to another DA layer because something is wrong with EigenDA.
				err := api.NewErrorFailover(
					fmt.Errorf("terminal dispersal failure for blobKey %v. blob status: %v",
						blobKey.Hex(),
						newStatus.String()))
				span.RecordError(err)
				span.SetStatus(codes.Error, "terminal dispersal failure")
				return nil, err
			}
		}
	}
}

// verifyReceivedBlobKey computes the BlobKey from the BlobHeader which was sent to the disperser, and compares it with
// the BlobKey which was returned by the disperser in the DisperseBlobReply
//
// A successful verification guarantees that the disperser didn't make any modifications to the BlobHeader that it
// received from this client.
//
// This function returns the verified blob key if the verification succeeds, and otherwise returns an error describing
// the failure
func verifyReceivedBlobKey(
	// the blob header which was constructed locally and sent to the disperser
	blobHeader *corev2.BlobHeader,
	// the reply received back from the disperser
	disperserReply *dispgrpc.DisperseBlobReply,
) (corev2.BlobKey, error) {

	actualBlobKey, err := blobHeader.BlobKey()
	enforce.NilError(err, "compute blob key")

	blobKeyFromDisperser, err := corev2.BytesToBlobKey(disperserReply.GetBlobKey())
	if err != nil {
		return corev2.BlobKey{}, fmt.Errorf("converting returned bytes to blob key: %w", err)
	}

	if actualBlobKey != blobKeyFromDisperser {
		return corev2.BlobKey{}, fmt.Errorf(
			"blob key returned by disperser (%v) doesn't match blob which was dispersed (%v)",
			blobKeyFromDisperser, actualBlobKey)
	}

	return blobKeyFromDisperser, nil
}

// intSliceFromQuorums converts a slice of QuorumIDs to a slice of ints for tracing attributes
func intSliceFromQuorums(quorums []core.QuorumID) []int {
	result := make([]int, len(quorums))
	for i, q := range quorums {
		result[i] = int(q)
	}
	return result
}
