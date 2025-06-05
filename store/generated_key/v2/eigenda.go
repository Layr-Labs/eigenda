package eigenda

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/utils"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/rlp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Store does storage interactions and verifications for blobs with the EigenDA V2 protocol.
type Store struct {
	log logging.Logger

	// Number of times to try blob dispersals:
	// - If > 0: Try N times total
	// - If < 0: Retry indefinitely until success
	// - If = 0: Not permitted
	putTries int
	// Allowed distance (in L1 blocks) between the eigenDA cert's reference block number (RBN)
	// and the L1 block number at which the cert was included in the rollup's batch inbox.
	// If cert.L1InclusionBlock > batch.RBN + rbnRecencyWindowSize, an
	// [RBNRecencyCheckFailedError] is returned.
	// This check is optional and will be skipped when rbnRecencyWindowSize is set to 0.
	rbnRecencyWindowSize uint64

	disperser    *payloaddispersal.PayloadDisperser
	retrievers   []clients.PayloadRetriever
	certVerifier *verification.CertVerifier
}

var _ common.EigenDAV2Store = (*Store)(nil)

func NewStore(
	log logging.Logger,
	putTries int,
	rbnRecencyWindowSize uint64,
	disperser *payloaddispersal.PayloadDisperser,
	retrievers []clients.PayloadRetriever,
	certVerifier *verification.CertVerifier,

) (*Store, error) {
	if putTries == 0 {
		return nil, fmt.Errorf(
			"putTries==0 is not permitted. >0 means 'try N times', <0 means 'retry indefinitely'")
	}

	return &Store{
		log:                  log,
		putTries:             putTries,
		rbnRecencyWindowSize: rbnRecencyWindowSize,
		disperser:            disperser,
		retrievers:           retrievers,
		certVerifier:         certVerifier,
	}, nil
}

// Get fetches a blob from DA using certificate fields and verifies blob
// against commitment to ensure data is valid and non-tampered.
func (e Store) Get(ctx context.Context, versionedCert certs.VersionedCert) ([]byte, error) {
	var cert coretypes.RetrievableEigenDACert

	switch versionedCert.Version {
	case certs.V0VersionByte, certs.V1VersionByte:
		var v2Cert coretypes.EigenDACertV2
		err := rlp.DecodeBytes(versionedCert.SerializedCert, &v2Cert)
		if err != nil {
			return nil, fmt.Errorf("RLP decoding EigenDA v2 cert: %w", err)
		}
		cert = &v2Cert

	case certs.V2VersionByte:
		var v3Cert coretypes.EigenDACertV3
		err := rlp.DecodeBytes(versionedCert.SerializedCert, &v3Cert)
		if err != nil {
			return nil, fmt.Errorf("RLP decoding EigenDA v3 cert: %w", err)
		}
		cert = &v3Cert

	default:
		return nil, fmt.Errorf("unknown certificate version: %d", versionedCert.Version)
	}

	// Try each retriever in sequence until one succeeds
	var errs []error
	for _, retriever := range e.retrievers {
		payload, err := retriever.GetPayload(ctx, cert)
		if err == nil {
			return payload.Serialize(), nil
		}

		e.log.Debugf("Payload retriever failed: %v", err)
	}

	return nil, fmt.Errorf("all retrievers failed: %w", errors.Join(errs...))
}

// Put disperses a blob for some pre-image and returns the associated RLP encoded certificate commit.
// TODO: Client polling for different status codes, Mapping status codes to 503 failover
func (e Store) Put(ctx context.Context, value []byte) ([]byte, error) {
	e.log.Debug("Dispersing payload to EigenDA V2 network")

	// TODO: https://github.com/Layr-Labs/eigenda/issues/1271

	// We attempt to disperse the blob to EigenDA up to PutRetries times, unless we get a 400 error on any attempt.

	payload := coretypes.NewPayload(value)

	cert, err := retry.DoWithData(
		func() (coretypes.EigenDACert, error) {
			return e.disperser.SendPayload(ctx, payload)
		},
		retry.RetryIf(
			func(err error) bool {
				grpcStatus, isGRPCError := status.FromError(err)
				if !isGRPCError {
					// api.ErrorFailover is returned, so we should retry
					return true
				}
				//nolint:exhaustive // we only care about a few grpc error codes
				switch grpcStatus.Code() {
				case codes.InvalidArgument:
					// we don't retry 400 errors because there is no point, we are passing invalid data
					return false
				case codes.ResourceExhausted:
					// we retry on 429s because *can* mean we are being rate limited
					// we sleep 1 second... very arbitrarily, because we don't have more info.
					// grpc error itself should return a backoff time,
					// see https://github.com/Layr-Labs/eigenda/issues/845 for more details
					time.Sleep(1 * time.Second)
					return true
				default:
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

	switch cert.Version() {
	case coretypes.VersionTwoCert:
		return nil, fmt.Errorf("EigenDA V2 certs are not supported anymore, use V3 instead")

	case coretypes.VersionThreeCert:
		eigenDACertV3, ok := cert.(*coretypes.EigenDACertV3)
		if !ok {
			return nil, fmt.Errorf("expected EigenDACertV3, got %T", cert)
		}

		return eigenDACertV3.Serialize(coretypes.CertSerializationRLP)

	default:
		return nil, fmt.Errorf("unsupported EigenDA cert version: %d", cert.Version())
	}
}

// BackendType returns the backend type for EigenDA Store
func (e Store) BackendType() common.BackendType {
	return common.EigenDAV2BackendType
}

// TODO: this whole function should be upstreamed to a new eigenda VerifyingPayloadRetrieval client
// that would verify certs, and then retrieve the payloads (from relay with fallback to eigenda validators if needed).
// Then proxy could remain a very thing server wrapper around eigenda clients.
// Verify verifies an EigenDACert by calling the verifyEigenDACertV2 view function
//
// Since v2 methods for fetching a payload are responsible for verifying the received bytes against the certificate,
// this Verify method only needs to check the cert on chain. That is why the third parameter is ignored.
func (e Store) Verify(ctx context.Context, versionedCert certs.VersionedCert, opts common.CertVerificationOpts) error {
	switch versionedCert.Version {
	case certs.V0VersionByte, certs.V1VersionByte:
		e.log.Warn("EigenDA V2 certs are not supported anymore more for verification. Defaulting to successful verification.")
		return nil

	case certs.V2VersionByte:
		var eigenDACert coretypes.EigenDACertV3
		err := rlp.DecodeBytes(versionedCert.SerializedCert, &eigenDACert)
		if err != nil {
			return fmt.Errorf("RLP decoding EigenDA v3 cert: %w", err)
		}

		err = verifyCertRBNRecencyCheck(eigenDACert.ReferenceBlockNumber(), opts.L1InclusionBlockNum, e.rbnRecencyWindowSize)
		if err != nil {
			return fmt.Errorf("rbn recency check failed: %w", err)
		}

		err = e.certVerifier.CheckDACert(ctx, &eigenDACert)
		if err != nil {
			return fmt.Errorf("verify v3 cert: %w", err)
		}

	default:
		return fmt.Errorf("unsupported EigenDA cert version: %d", versionedCert.Version)
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
	if certRBN == 0 {
		return fmt.Errorf("certRBN should never be 0, this is likely a bug")
	}
	if certL1IBN <= certRBN {
		return fmt.Errorf(
			"cert's l1 inclusion block number (%d) <= cert reference block number (%d), but this is physically impossible "+
				"since the cert has to be signed by all eigenda validators before being submitted to the batcher inbox, "+
				"and validators will only sign a batchRoot (contained in certs) if the RBN is in the past. "+
				"This is a serious bug, please report it",
			certRBN,
			certL1IBN,
		)
	}

	// Actual Recency Check
	if !(certL1IBN <= certRBN+rbnRecencyWindowSize) {
		//nolint:gosec // disable G115 // We checked the length of thresholds above
		return RBNRecencyCheckFailedError{
			certRBN:              uint32(certRBN),
			certL1IBN:            certL1IBN,
			rbnRecencyWindowSize: rbnRecencyWindowSize,
		}
	}
	return nil
}
