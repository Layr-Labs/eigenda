package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	dispgrpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/geth"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
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

// BuildPayloadDisperser builds a PayloadDisperser from config structs.
func BuildPayloadDisperser(log logging.Logger, payloadDispCfg PayloadDisperserConfig,
	dispClientCfg *DisperserClientConfig,
	ethCfg *geth.EthClientConfig,
	kzgConfig *kzg.KzgConfig, encoderCfg *encoding.Config) (*PayloadDisperser, error) {

	// 1 - verify key semantics and create signer
	signer, err := auth.NewLocalBlobRequestSigner(payloadDispCfg.SignerPaymentKey)
	if err != nil {
		return nil, fmt.Errorf("new local blob request signer: %w", err)
	}

	// 2 - create prover (if applicable)

	var kzgProver encoding.Prover
	if kzgConfig != nil {
		if encoderCfg == nil {
			encoderCfg = encoding.DefaultConfig()
		}

		kzgProver, err = prover.NewProver(kzgConfig, encoderCfg)
		if err != nil {
			return nil, fmt.Errorf("new kzg prover: %w", err)
		}
	} else {
		log.Warn("No prover provided, using disperser for blob commitment generation")
	}

	// 3 - create disperser client & set accountant to nil
	// to then populate using signer field via in-place method
	// which queries disperser directly for payment states
	disperserClient, err := NewDisperserClient(dispClientCfg, signer, kzgProver, nil)
	if err != nil {
		return nil, fmt.Errorf("new disperser client: %w", err)
	}

	err = disperserClient.PopulateAccountant(context.Background())
	if err != nil {
		return nil, fmt.Errorf("populating accountant in disperser client: %w", err)
	}

	// 4 - construct eth client to wire up cert verifier
	ethClient, err := geth.NewClient(*ethCfg, gethcommon.Address{}, 0, log)
	if err != nil {
		return nil, fmt.Errorf("new eth client: %w", err)
	}

	certVerifier, err := verification.NewCertVerifier(
		log,
		ethClient,
		payloadDispCfg.EigenDACertVerifierAddr,
		payloadDispCfg.BlockNumberPollInterval,
	)

	if err != nil {
		return nil, fmt.Errorf("new cert verifier: %w", err)
	}

	// 5 - create codec
	codec, err := codecs.CreateCodec(payloadDispCfg.PayloadPolynomialForm, payloadDispCfg.BlobEncodingVersion)
	if err != nil {
		return nil, err
	}

	return NewPayloadDisperser(log, payloadDispCfg, codec, disperserClient, certVerifier)
}

// NewPayloadDisperser creates a PayloadDisperser from subcomponents that have already been constructed and initialized.
func NewPayloadDisperser(
	logger logging.Logger,
	payloadDisperserConfig PayloadDisperserConfig,
	codec codecs.BlobCodec,
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
	pd.logger.Debug("Successful DisperseBlob", "blobStatus", blobStatus.String(), "blobKey", blobKey.Hex())

	timeoutCtx, cancel = context.WithTimeout(ctx, pd.config.BlobCertifiedTimeout)
	defer cancel()
	blobStatusReply, err := pd.pollBlobStatusUntilCertified(timeoutCtx, blobKey, blobStatus.ToProfobuf())
	if err != nil {
		return nil, fmt.Errorf("poll blob status until certified: %w", err)
	}
	pd.logger.Debug("Blob status CERTIFIED", "blobKey", blobKey.Hex())

	eigenDACert, err := pd.buildEigenDACert(ctx, blobKey, blobStatusReply)
	if err != nil {
		// error returned from method is sufficiently descriptive
		return nil, err
	}

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
	pd.logger.Debug("Retrieved NonSignerStakesAndSignature", "blobKey", blobKey.Hex())

	eigenDACert, err := verification.BuildEigenDACert(blobStatusReply, nonSignerStakesAndSignature)
	if err != nil {
		return nil, fmt.Errorf("build eigen da cert: %w", err)
	}
	pd.logger.Debug("Constructed EigenDACert", "blobKey", blobKey.Hex())

	return eigenDACert, nil
}
