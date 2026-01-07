package eigenda

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/dispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/consts"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/utils"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/avast/retry-go/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Store does storage interactions and verifications for blobs with the EigenDA V2 protocol.
type Store struct {
	log logging.Logger

	// Dispersal related fields. disperser is optional, and PUT routes will return 500s if not set.
	disperser *dispersal.PayloadDisperser
	// Number of times to try blob dispersals:
	// - If > 0: Try N times total
	// - If < 0: Retry indefinitely until success
	// - If = 0: Not permitted
	putTries int

	// Verification related fields.
	certVerifier *verification.CertVerifier

	// Retrieval related fields.
	retrievers []clients.PayloadRetriever

	// Timeout used for contract calls
	contractCallTimeout time.Duration

	// offchainDerivationMap maps offchain derivation versions to their parameters.
	// offchain derivation version was introduced with EigenDA V4 certs, and is not used for earlier cert versions.
	offchainDerivationMap certs.OffchainDerivationMap
}

var _ common.EigenDAV2Store = (*Store)(nil)

func NewStore(
	log logging.Logger,
	disperser *dispersal.PayloadDisperser,
	putTries int,
	certVerifier *verification.CertVerifier,
	retrievers []clients.PayloadRetriever,
	contractCallTimeout time.Duration,
) (*Store, error) {
	if putTries == 0 {
		return nil, fmt.Errorf(
			"putTries==0 is not permitted. >0 means 'try N times', <0 means 'retry indefinitely'")
	}

	offchainDerivationMap := make(certs.OffchainDerivationMap)
	// Currently only offchain derivation version 0 exists.
	offchainDerivationMap[0] = certs.OffchainDerivationParameters{
		RBNRecencyWindowSize: consts.RBNRecencyWindowSizeV0,
	}

	return &Store{
		log:                   log,
		putTries:              putTries,
		disperser:             disperser,
		retrievers:            retrievers,
		certVerifier:          certVerifier,
		contractCallTimeout:   contractCallTimeout,
		offchainDerivationMap: offchainDerivationMap,
	}, nil
}

// Get fetches a blob from DA using certificate fields and verifies blob
// against commitment to ensure data is valid and non-tampered.
// If returnEncodedPayload is true, it returns the encoded payload without decoding.
//
// This function is bug-prone as is because it returns []byte which can either be a raw payload or an encoded payload.
// TODO: Refactor to use [coretypes.EncodedPayload] and [coretypes.Payload] instead of []byte.
func (e Store) Get(
	ctx context.Context,
	versionedCert *certs.VersionedCert,
	serializationType coretypes.CertSerializationType,
	returnEncodedPayload bool,
) ([]byte, error) {
	certTypeVersion, err := versionedCert.Version.IntoCertVersion()
	if err != nil {
		return nil, coretypes.NewCertParsingFailedError(
			hex.EncodeToString(versionedCert.SerializedCert),
			fmt.Sprintf("casting to cert type version: %v", err),
		)
	}

	cert, err := coretypes.DeserializeEigenDACert(
		versionedCert.SerializedCert,
		certTypeVersion,
		serializationType,
	)
	if err != nil {
		return nil, coretypes.NewCertParsingFailedError(
			hex.EncodeToString(versionedCert.SerializedCert),
			fmt.Sprintf("deserialize cert: %v", err),
		)
	}

	// Try each retriever in sequence until one succeeds
	var errs []error
	for _, retriever := range e.retrievers {
		if returnEncodedPayload {
			// Get encoded payload if requested
			encodedPayload, err := retriever.GetEncodedPayload(ctx, cert)
			if err == nil {
				return encodedPayload.Serialize(), nil
			}
			e.log.Debugf("Encoded payload retriever failed: %v", err)
			errs = append(errs, err)
		} else {
			// Get decoded payload (default behavior)
			payload, err := retriever.GetPayload(ctx, cert)
			if err == nil {
				return payload, nil
			}
			e.log.Debugf("Payload retriever failed: %v", err)
			errs = append(errs, err)
		}
	}

	return nil, fmt.Errorf("all retrievers failed: %w", errors.Join(errs...))
}

// Put disperses a blob for some pre-image and returns the associated certificate commit.
func (e Store) Put(
	ctx context.Context, value []byte, serializationType coretypes.CertSerializationType,
) (*certs.VersionedCert, error) {
	if e.disperser == nil {
		return nil, fmt.Errorf("PUT routes are disabled, did you provide a signer private key?")
	}
	e.log.Debug("Dispersing payload to EigenDA V2 network")

	// TODO: https://github.com/Layr-Labs/eigenda/issues/1271

	// We attempt to disperse the blob to EigenDA up to PutRetries times, unless we get a 400 error on any attempt.

	payload := coretypes.Payload(value)

	cert, err := retry.DoWithData(
		func() (coretypes.EigenDACert, error) {
			return e.disperser.SendPayload(ctx, payload)
		},
		retry.RetryIf(
			func(err error) bool {
				if err == nil {
					// This should never happen since RetryIf function should only be called err != nil.
					// But returning false since if no error happened... then don't need to retry,
					// unless there's a bug in the RetryIf library...
					return false
				}
				if errors.Is(err, &api.ErrorFailover{}) {
					// Failover errors should be retried before failing over.
					return true
				}
				grpcStatus, isGRPCError := status.FromError(err)
				if !isGRPCError {
					e.log.Warn("Received non-grpc error, retrying", "err", err)
					return true
				}
				//nolint:exhaustive // we only care about a few grpc error codes
				switch grpcStatus.Code() {
				case codes.InvalidArgument:
					// we don't retry 400 errors because there is no point, we are passing invalid data
					e.log.Warn("Received InvalidArgument status code, not retrying", "err", err)
					return false
				case codes.ResourceExhausted:
					// we retry on 429s because it *can* mean we are being rate limited
					// we sleep 1 second... very arbitrarily, because we don't have more info.
					// grpc error itself should return a backoff time,
					// see https://github.com/Layr-Labs/eigenda/issues/845 for more details
					e.log.Warn("Received ResourceExhausted status code, retrying after 1 second", "err", err)
					time.Sleep(1 * time.Second)
					return true
				default:
					e.log.Warn("Received gRPC error, retrying", "err", err, "code", grpcStatus.Code())
					return true
				}
			}),
		// only return the last error. If it is an api.ErrorFailover, then the handler will convert
		// it to an http 503 to signify to the client (batcher) to failover to ethda b/c eigenda is temporarily down.
		retry.LastErrorOnly(true),
		// retry.Attempts uses different semantics than our config field. ConvertToRetryGoAttempts converts between
		// these two semantics.
		retry.Attempts(utils.ConvertToRetryGoAttempts(e.putTries)),
	)
	if err != nil {
		// TODO: we will want to filter for errors here and return a 503 when needed, i.e. when dispersal itself failed,
		//  or that we timed out waiting for batch to land onchain
		return nil, err
	}

	switch cert := cert.(type) {
	case *coretypes.EigenDACertV2:
		return nil, fmt.Errorf("EigenDA V2 certs are not supported anymore, use V3 instead")
	case *coretypes.EigenDACertV3:
		serializedCert, err := cert.Serialize(serializationType)
		if err != nil {
			return nil, fmt.Errorf("serialize cert: %w", err)
		}
		return certs.NewVersionedCert(serializedCert, certs.V2VersionByte), nil
	case *coretypes.EigenDACertV4:
		serializedCert, err := cert.Serialize(serializationType)
		if err != nil {
			return nil, fmt.Errorf("serialize cert: %w", err)
		}
		return certs.NewVersionedCert(serializedCert, certs.V3VersionByte), nil
	default:
		return nil, fmt.Errorf("unsupported cert version: %T", cert)
	}
}

// BackendType returns the backend type for EigenDA Store
func (e Store) BackendType() common.BackendType {
	return common.EigenDAV2BackendType
}

// VerifyCert verifies an EigenDACert by calling the verifyEigenDACertV2 view function
//
// Since v2 methods for fetching a payload are responsible for verifying the received bytes against the certificate,
// this VerifyCert method only needs to check the cert on chain. That is why the third parameter is ignored.
//
// TODO: this whole function should be upstreamed to a new eigenda VerifyingPayloadRetrieval client
// that would verify certs, and then retrieve the payloads (from relay with fallback to eigenda validators if needed).
// Then proxy could remain a very thin server wrapper around eigenda clients.
func (e Store) VerifyCert(ctx context.Context, versionedCert *certs.VersionedCert,
	serializationType coretypes.CertSerializationType, l1InclusionBlockNum uint64) error {
	var sumDACert coretypes.EigenDACert
	var certVersion coretypes.CertificateVersion

	switch versionedCert.Version {
	case certs.V0VersionByte:
		return coretypes.NewCertParsingFailedError(
			hex.EncodeToString(versionedCert.SerializedCert),
			"version 0 byte certs should never be verified by the EigenDA V2 store",
		)

	case certs.V1VersionByte, certs.V2VersionByte, certs.V3VersionByte:
		certTypeVersion, err := versionedCert.Version.IntoCertVersion()
		if err != nil {
			return coretypes.NewCertParsingFailedError(
				hex.EncodeToString(versionedCert.SerializedCert), fmt.Sprintf("casting to cert type version: %v", err))
		}

		cert, err := coretypes.DeserializeEigenDACert(
			versionedCert.SerializedCert,
			certTypeVersion,
			serializationType,
		)
		if err != nil {
			return coretypes.NewCertParsingFailedError(
				hex.EncodeToString(versionedCert.SerializedCert),
				fmt.Sprintf("deserialize EigenDA cert: %v", err))
		}
		certVersion = certTypeVersion
		sumDACert = cert

	default:
		return coretypes.NewCertParsingFailedError(
			hex.EncodeToString(versionedCert.SerializedCert),
			fmt.Sprintf("unknown EigenDA cert version: %d", versionedCert.Version))
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, e.contractCallTimeout)
	defer cancel()

	// verify cert via simulation call to verifier contract
	err := e.certVerifier.CheckDACert(timeoutCtx, sumDACert)
	if err != nil {
		var certVerifierInvalidCertErr *verification.CertVerifierInvalidCertError
		if errors.As(err, &certVerifierInvalidCertErr) {
			// We convert the cert verifier failure error, which contains the low-level detailed status code,
			// into the higher-level CertDerivationError which will get converted to a 418 HTTP error by the error middleware.
			return coretypes.ErrInvalidCertDerivationError.WithMessage(certVerifierInvalidCertErr.Error())
		}
		// Other errors are internal proxy errors, so we just wrap them for extra context.
		// They will be converted to 500 HTTP errors by the error middleware.
		return fmt.Errorf("eth-call to CertVerifier.checkDACert: %w", err)
	}

	// For cert versions that support offchain derivation versioning (v4+),
	// we need to fetch the offchain derivation version from the contract,
	// and then use that to get the relevant offchain derivation parameters
	// (e.g. RBN recency window size) to perform additional checks.
	//
	// For cert versions that do not support offchain derivation versioning (v3 and below),
	// we skip these additional checks.
	//
	// Note: offchain derivation versioning was introduced in cert version 4.
	if certVersion >= coretypes.VersionFourCert {
		// The CheckDACert call above has already verified the cert's onchain validity,
		// including that the cert's offchain derivation version is supported onchain.
		// So we can safely cast to V4 here.
		certV4 := sumDACert.(*coretypes.EigenDACertV4)
		offchainDerivationVersion := certV4.OffchainDerivationVersion

		offchainDerivationParams, exists := e.offchainDerivationMap[offchainDerivationVersion]
		if !exists {
			// Note: If we encounter this error, we've updated the derivation version onchain and not updated the
			// hardcoded offchain map. This should never happen in practice unless there's a misconfiguration.
			return coretypes.NewCertParsingFailedError(
				hex.EncodeToString(versionedCert.SerializedCert),
				fmt.Sprintf("unsupported offchain derivation version: %d", offchainDerivationVersion),
			)
		}

		err = verifyCertRBNRecencyCheck(
			certV4.ReferenceBlockNumber(),
			l1InclusionBlockNum,
			offchainDerivationParams.RBNRecencyWindowSize,
		)
		if err != nil {
			// Already a structured error converted to a 418 HTTP error by the error middleware.
			return err
		}
	}

	return nil
}

// verifyCertRBNRecencyCheck arguments:
//   - certRBN: ReferenceBlockNumber included in the cert itself at which operator stakes are referenced
//     when verifying that a cert's signature meets the required quorum thresholds.
//   - certL1IBN: InclusionBlockNumber at which the EigenDA cert was included in the rollup batcher inbox.
//     The IBN is not part of the cert. It is received as an optional query param on GET requests.
//     0 means to skip the check (return nil).
//   - rbnRecencyWindowSize: distance allowed between the RBN and IBN. See below for more details.
//     Value should be set by proxy operator as a flag. 0 means to skip the check (return nil).
//
// Certs in the rollup batcher-inbox that do not respect the below equation are discarded.
//
//	certRBN < certL1IBN <= certRBN + RBNRecencyWindowSize
//
// This check serves 2 purposes:
//  1. liveness: prevents derivation pipeline from stalling on blobs that are no longer available on the DA layer
//  2. safety: prevents a malicious EigenDA sequencer from using a very stale RBN whose operator distribution
//     does not represent the actual stake distribution. Operators that withdrew a lot of stake would
//     not be slashable anymore, even though because of the old RBN their signature would count for a lot of stake.
//
// Note that for a secure integration, this same check needs to be verified onchain.
// There are 2 approaches to doing this:
//  1. Pessimistic approach: use a smart batcher inbox to disallow stale blobs from even being included
//     in the batcher inbox (see https://github.com/ethereum-optimism/design-docs/pull/229)
//  2. Optimistic approach: verify the check in op-program or hokulea (kona)'s derivation pipeline. See
//     https://github.com/Layr-Labs/hokulea/blob/8c4c89bc4f/crates/eigenda/src/eigenda.rs#L90
func verifyCertRBNRecencyCheck(certRBN uint64, certL1IBN uint64, rbnRecencyWindowSize uint64) error {
	// Input Validation
	if certL1IBN == 0 || rbnRecencyWindowSize == 0 {
		return nil
	}

	// Actual Recency Check
	if !(certL1IBN <= certRBN+rbnRecencyWindowSize) { //nolint:staticcheck // inequality is clearer as is
		return coretypes.NewRBNRecencyCheckFailedError(certRBN, certL1IBN, rbnRecencyWindowSize)
	}
	return nil
}
