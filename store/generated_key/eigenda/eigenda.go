package eigenda

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/utils"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/rlp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StoreConfig struct {
	MaxBlobSizeBytes uint64
	// the # of Ethereum blocks to wait after the EigenDA L1BlockReference # before attempting to verify
	// & accredit a blob
	EthConfirmationDepth uint64

	// total duration time that client waits for blob to confirm
	StatusQueryTimeout time.Duration

	// Number of times to try blob dispersals:
	// - If > 0: Try N times total
	// - If < 0: Retry indefinitely until success
	// - If = 0: Not permitted
	PutTries int
}

// Store does storage interactions and verifications for blobs with DA.
type Store struct {
	client   *clients.EigenDAClient
	verifier *verify.Verifier
	cfg      *StoreConfig
	log      logging.Logger
}

var _ common.EigenDAStore = (*Store)(nil)

// NewStoreConfig creates a new StoreConfig with validation for PutTries.
// PutTries==0 is not permitted, and will return an error
func NewStoreConfig(
	maxBlobSizeBytes uint64,
	ethConfirmationDepth uint64,
	statusQueryTimeout time.Duration,
	putTries int,
) (*StoreConfig, error) {
	if putTries == 0 {
		return nil, fmt.Errorf(
			"putTries==0 is not permitted. >0 means 'try N times', <0 means 'retry indefinitely'")
	}

	return &StoreConfig{
		MaxBlobSizeBytes:     maxBlobSizeBytes,
		EthConfirmationDepth: ethConfirmationDepth,
		StatusQueryTimeout:   statusQueryTimeout,
		PutTries:             putTries,
	}, nil
}

func NewStore(
	client *clients.EigenDAClient,
	v *verify.Verifier,
	log logging.Logger,
	cfg *StoreConfig,
) (*Store, error) {
	return &Store{
		client:   client,
		verifier: v,
		log:      log,
		cfg:      cfg,
	}, nil
}

// Get fetches a blob from DA using certificate fields and verifies blob
// against commitment to ensure data is valid and non-tampered.
func (e Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	var cert verify.Certificate
	err := rlp.DecodeBytes(key, &cert)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DA cert to RLP format: %w", err)
	}

	err = cert.NoNilFields()
	if err != nil {
		return nil, fmt.Errorf("failed to verify DA cert: %w", err)
	}

	decodedBlob, err := e.client.GetBlob(
		ctx,
		cert.BlobVerificationProof.GetBatchMetadata().GetBatchHeaderHash(),
		cert.BlobVerificationProof.GetBlobIndex(),
	)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to retrieve decoded blob: %w", err)
	}

	return decodedBlob, nil
}

// Put disperses a blob for some pre-image and returns the associated RLP encoded certificate commit.
func (e Store) Put(ctx context.Context, value []byte) ([]byte, error) {
	encodedBlob, err := e.client.GetCodec().EncodeBlob(value)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to re-encode blob: %w", err)
	}
	// TODO: We should move this length check inside PutBlob
	if uint64(len(encodedBlob)) > e.cfg.MaxBlobSizeBytes {
		return nil, fmt.Errorf(
			"%w: blob length %d, max blob size %d",
			common.ErrProxyOversizedBlob,
			len(value),
			e.cfg.MaxBlobSizeBytes,
		)
	}

	// We attempt to disperse the blob to EigenDA up to e.cfg.PutTries times total,
	// unless we get a 400 error which aborts retries.
	blobInfo, err := retry.DoWithData(
		func() (*disperser.BlobInfo, error) {
			return e.client.PutBlob(ctx, value)
		},
		retry.RetryIf(func(err error) bool {
			st, isGRPCError := status.FromError(err)
			if !isGRPCError {
				// api.ErrorFailover is returned, so we should retry
				return true
			}
			//nolint:exhaustive // we only care about a few grpc error codes
			switch st.Code() {
			case codes.InvalidArgument:
				// we don't retry 400 errors because there is no point,
				// we are passing invalid data
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
		// it to an http 503 to signify to the client (batcher) to failover to ethda
		// b/c eigenda is temporarily down.
		retry.LastErrorOnly(true),
		// retry.Attempts uses different semantics than our config field. ConvertToRetryGoAttempts converts between
		// these two semantics.
		retry.Attempts(utils.ConvertToRetryGoAttempts(e.cfg.PutTries)),
	)
	if err != nil {
		// TODO: we will want to filter for errors here and return a 503 when needed
		// ie when dispersal itself failed, or that we timed out waiting for batch to land onchain
		return nil, err
	}
	cert := (*verify.Certificate)(blobInfo)

	err = cert.NoNilFields()
	if err != nil {
		return nil, fmt.Errorf("failed to verify DA cert due to nil fields: %w", err)
	}

	err = cert.ValidFieldLengths()
	if err != nil {
		return nil, fmt.Errorf("failed to verify DA cert due to invalid field lengths: %w", err)
	}

	err = e.verifier.VerifyCommitment(cert.BlobHeader.GetCommitment(), encodedBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to verify commitment: %w", err)
	}

	// The cert will only be included in the batcher's inbox after
	// the proxy returns the verified cert to the batcher.
	err = e.verifier.VerifyCert(ctx, cert)
	if err != nil {
		if errors.Is(err, verify.ErrBatchMetadataHashMismatch) {
			// This error might have been caused by an L1 reorg.
			// See
			// https://github.com/Layr-Labs/eigenda-proxy/blob/main/docs/troubleshooting_v1.md#batch-hash-mismatch-error
			// We could try to update the cert's confirmation block number here,
			// but it's currently not possible because the client.PutBlob call doesn't return the
			// request_id to be able to query the GetBlobStatus endpoint...
			// V2 already solves this by using the blob_header_hash as the request_id,
			// but also having moved verification to the client.
			// V1 feels doomed..... maybe just accept this sad state of affairs
			// and try to get people to migrate to V2 asap.
			return nil, fmt.Errorf("failed to verify DA cert: %w", err)
		}
		return nil, fmt.Errorf("failed to verify DA cert: %w", err)
	}

	bytes, err := rlp.EncodeToBytes(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to encode DA cert to RLP format: %w", err)
	}

	return bytes, nil
}

// BackendType returns the backend type for EigenDA Store
func (e Store) BackendType() common.BackendType {
	return common.EigenDABackendType
}

// Key is used to recover certificate fields and that verifies blob
// against commitment to ensure data is valid and non-tampered.
func (e Store) Verify(ctx context.Context, serializedCert []byte, payload []byte, _ common.CertVerificationOpts) error {
	var cert verify.Certificate
	err := rlp.DecodeBytes(serializedCert, &cert)
	if err != nil {
		return fmt.Errorf("failed to decode DA cert to RLP format: %w", err)
	}

	// re-encode blob for verification
	encodedBlob, err := e.client.GetCodec().EncodeBlob(payload)
	if err != nil {
		return fmt.Errorf("EigenDA client failed to re-encode blob: %w", err)
	}

	// verify kzg data commitment
	err = e.verifier.VerifyCommitment(cert.BlobHeader.GetCommitment(), encodedBlob)
	if err != nil {
		return fmt.Errorf("failed to verify commitment: %w", err)
	}

	// verify DA certificate against EigenDA's batch metadata that's bridged to Ethereum
	err = e.verifier.VerifyCert(ctx, &cert)
	if errors.Is(err, verify.ErrBatchMetadataHashMismatch) {
		// This error might have been caused by an L1 reorg.
		// See https://github.com/Layr-Labs/eigenda-proxy/blob/main/docs/troubleshooting_v1.md#batch-hash-mismatch-error
		// We would want to update the cert's confirmation block number here,
		// but it's currently not possible because the request_id to be able to query the GetBlobStatus endpoint
		// is not stored in the v1 cert.... :(
		//
		// V2 solves this by using the blob_header_hash as the request_id,
		// but also having moved verification to the client.
		// V1 feels doomed..... maybe just accept this sad state of affairs
		// and try to get people to migrate to V2 asap.
		return fmt.Errorf("failed to verify DA cert: %w", err)
	}
	return err
}
