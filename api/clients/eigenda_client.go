package clients

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	grpcdisperser "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/ethereum/go-ethereum/log"
)

// IEigenDAClient is a wrapper around the DisperserClient interface which
// encodes blobs before dispersing them, and decodes them after retrieving them.
type IEigenDAClient interface {
	GetBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error)
	PutBlob(ctx context.Context, txData []byte) (*grpcdisperser.BlobInfo, error)
	GetCodec() codecs.BlobCodec
	Close() error
}

// See the NewEigenDAClient constructor's documentation for details and usage examples.
type EigenDAClient struct {
	// TODO: all of these should be private, to prevent users from using them directly,
	// which breaks encapsulation and makes it hard for us to do refactors or changes
	Config EigenDAClientConfig
	Log    log.Logger
	Client DisperserClient
	Codec  codecs.BlobCodec
}

var _ IEigenDAClient = &EigenDAClient{}

// EigenDAClient is a wrapper around the DisperserClient which
// encodes blobs before dispersing them, and decodes them after retrieving them.
// It also turns the disperser's async polling-based API (disperseBlob + poll GetBlobStatus)
// into a sync API where PutBlob will poll for the blob to be confirmed or finalized.
//
// DisperserClient is safe to be used concurrently by multiple goroutines.
// Don't forget to call Close() on the client when you're done with it, to close the
// underlying grpc connection maintained by the DiserserClient.
//
// Example usage:
//
//	client, err := NewEigenDAClient(log, EigenDAClientConfig{...})
//	if err != nil {
//	  return err
//	}
//	defer client.Close()
//
//	blobData := []byte("hello world")
//	blobInfo, err := client.PutBlob(ctx, blobData)
//	if err != nil {
//	  return err
//	}
//
//	retrievedData, err := client.GetBlob(ctx, blobInfo.BatchMetadata.BatchHeaderHash, blobInfo.BlobIndex)
//	if err != nil {
//	  return err
//	}
func NewEigenDAClient(log log.Logger, config EigenDAClientConfig) (*EigenDAClient, error) {
	err := config.CheckAndSetDefaults()
	if err != nil {
		return nil, err
	}

	host, port, err := net.SplitHostPort(config.RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EigenDA RPC: %w", err)
	}

	var signer core.BlobRequestSigner
	if len(config.SignerPrivateKeyHex) == 64 {
		signer = auth.NewLocalBlobRequestSigner(config.SignerPrivateKeyHex)
	} else if len(config.SignerPrivateKeyHex) == 0 {
		// noop signer is used when we need a read-only eigenda client
		signer = auth.NewLocalNoopSigner()
	} else {
		return nil, fmt.Errorf("invalid length for signer private key")
	}

	disperserConfig := NewConfig(host, port, config.ResponseTimeout, !config.DisableTLS)
	disperserClient := NewDisperserClient(disperserConfig, signer)

	lowLevelCodec, err := codecs.BlobEncodingVersionToCodec(config.PutBlobEncodingVersion)
	if err != nil {
		return nil, fmt.Errorf("error initializing EigenDA client: %w", err)
	}

	var codec codecs.BlobCodec
	if config.DisablePointVerificationMode {
		codec = codecs.NewNoIFFTCodec(lowLevelCodec)
	} else {
		codec = codecs.NewIFFTCodec(lowLevelCodec)
	}

	return &EigenDAClient{
		Log:    log,
		Config: config,
		Client: disperserClient,
		Codec:  codec,
	}, nil
}

// Deprecated: do not rely on this function. Do not use m.Codec directly either.
// These will eventually be removed and not exposed.
func (m *EigenDAClient) GetCodec() codecs.BlobCodec {
	return m.Codec
}

// GetBlob retrieves a blob from the EigenDA service using the provided context,
// batch header hash, and blob index.  If decode is set to true, the function
// decodes the retrieved blob data. If set to false it returns the encoded blob
// data, which is necessary for generating KZG proofs for data's correctness.
// The function handles potential errors during blob retrieval, data length
// checks, and decoding processes.
func (m *EigenDAClient) GetBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error) {
	data, err := m.Client.RetrieveBlob(ctx, batchHeaderHash, blobIndex)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve blob: %w", err)
	}

	if len(data) == 0 {
		// This should never happen, because empty blobs are rejected from even entering the system:
		// https://github.com/Layr-Labs/eigenda/blob/master/disperser/apiserver/server.go#L930
		return nil, fmt.Errorf("blob has length zero - this should not be possible")
	}

	decodedData, err := m.Codec.DecodeBlob(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding blob: %w", err)
	}

	return decodedData, nil
}

// PutBlob encodes and writes a blob to EigenDA, waiting for a desired blob status
// to be reached (guarded by WaitForFinalization config param) before returning.
// This function is resilient to transient failures and timeouts.
//
// PutBlob returned errors all implement the ErrorAPI interface, which allows the caller
// to determine whether the error is a client or server fault, as well as get its status code.
// A 503 error returned is used to signify that eigenda is temporarily unavailable,
// and suggest to the caller (most likely some rollup batcher via the eigenda-proxy)
// to fallback to ethda for some amount of time.
// See https://github.com/ethereum-optimism/specs/issues/434 for more details.
//
// Seriously considered literally returning the ErrorAPI interface instead of error,
// but this seems like a very uncommon pattern in Go, so I'm not sure if it's a good idea.
// https://www.reddit.com/r/golang/comments/18wnmhx/returning_a_type_representing_an_error_rather/
// has some pros and cons, although their main con is when the return type is a struct pointer,
// which wouldn't apply in our case.
func (m *EigenDAClient) PutBlob(ctx context.Context, data []byte) (*grpcdisperser.BlobInfo, error) {
	resultChan, errorChan := m.PutBlobAsync(ctx, data)
	select { // no timeout here because we depend on the configured timeout in PutBlobAsync
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (m *EigenDAClient) PutBlobAsync(ctx context.Context, data []byte) (resultChan chan *grpcdisperser.BlobInfo, errChan chan api.ErrorAPI) {
	resultChan = make(chan *grpcdisperser.BlobInfo, 1)
	errChan = make(chan api.ErrorAPI, 1)
	go m.putBlob(ctx, data, resultChan, errChan)
	return
}

func (m *EigenDAClient) putBlob(ctx context.Context, rawData []byte, resultChan chan *grpcdisperser.BlobInfo, errChan chan api.ErrorAPI) {
	m.Log.Info("Attempting to disperse blob to EigenDA")

	// encode blob
	if m.Codec == nil {
		errChan <- api.NewErrorAPIGeneric(http.StatusBadRequest, fmt.Errorf("codec not initialized"))
		return
	}

	data, err := m.Codec.EncodeBlob(rawData)
	if err != nil {
		errChan <- api.NewErrorAPIGeneric(http.StatusBadRequest, fmt.Errorf("error encoding blob: %w", err))
		return
	}

	customQuorumNumbers := make([]uint8, len(m.Config.CustomQuorumIDs))
	for i, e := range m.Config.CustomQuorumIDs {
		customQuorumNumbers[i] = uint8(e)
	}
	// disperse blob
	// TODO: would be nice to add a trace-id key to the context, to be able to follow requests from batcher->proxy->eigenda
	blobStatus, requestID, err := m.Client.DisperseBlobAuthenticated(ctx, data, customQuorumNumbers)
	if err != nil {
		errChan <- &api.ErrorAPIGeneric{
			Err: fmt.Errorf("error submitting authenticated blob to disperser: %w", err),
			// We set to unknown fault b/c disperser client returns a mix of 400s and 500s currently.
			// TODO: update disperser client to also return ErrorAPIGeneric errors
			Code:  0,
			Fault: api.ErrorFaultUnknown,
		}
		return
	}

	// process response
	if *blobStatus == disperser.Failed {
		// Don't think this state should be reachable. DisperseBlobAuthenticated should return an error instead of a failed status.
		// TODO: we should document that endpoint in the proto file to mention this. If it is the case, then we can remove this code.
		errChan <- api.NewErrorAPIGeneric(http.StatusInternalServerError, fmt.Errorf("dispersed blob immediately in failed status, something bad happened: %w", err))
		return
	}

	base64RequestID := base64.StdEncoding.EncodeToString(requestID)
	m.Log.Info("Blob accepted by EigenDA disperser, now polling for status updates", "requestID", base64RequestID)

	ticker := time.NewTicker(m.Config.StatusQueryRetryInterval)
	defer ticker.Stop()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.Config.StatusQueryTimeout)
	defer cancel()

	alreadyWaitingForDispersal := false
	alreadyWaitingForFinalization := false
	var latestBlobStatus grpcdisperser.BlobStatus
	for {
		select {
		case <-ctx.Done():
			// We can only land here if blob is still in
			// 1. processing or dispersing status: waiting to land onchain
			// 2. or confirmed status: landed onchain, waiting for finalization
			// because all other statuses return immediately below.
			//
			// Assuming that the timeout is correctly set (long enough to both land onchain + finalize),
			// 1. means that there is a problem with EigenDA, so we return 503 to let the batcher failover to ethda
			// 2. means that there is a problem with Ethereum, so we return 500.
			//    batcher would most likely resubmit another blob, which is not ideal but there isn't much to be done...
			//    eigenDA v2 will have idempotency so one can just resubmit the same blob safely.
			if latestBlobStatus == grpcdisperser.BlobStatus_PROCESSING || latestBlobStatus == grpcdisperser.BlobStatus_DISPERSING {
				errChan <- api.NewErrorAPIGeneric(http.StatusServiceUnavailable,
					fmt.Errorf("eigenda might be down. timed out waiting for blob to land onchain (request id=%s): %w", base64RequestID, ctx.Err()))
			}
			// but not else (otherwise it might be a problem with ethereum, so fallbacking to ethda wouldnt help)
			errChan <- api.NewErrorAPIGeneric(http.StatusInternalServerError, fmt.Errorf("timed out waiting for blob that landed onchain to finalize (request id=%s): %w. Either timeout not long enough, or ethereum might be experiencing difficulties.", base64RequestID, ctx.Err()))
			return
		case <-ticker.C:
			statusRes, err := m.Client.GetBlobStatus(ctx, requestID)
			if err != nil {
				m.Log.Warn("Unable to retrieve blob dispersal status, will retry", "requestID", base64RequestID, "err", err)
				continue
			}
			latestBlobStatus = statusRes.Status
			switch statusRes.Status {
			case grpcdisperser.BlobStatus_PROCESSING, grpcdisperser.BlobStatus_DISPERSING:
				// to prevent log clutter, we only log at info level once
				if alreadyWaitingForDispersal {
					m.Log.Debug("Blob is being processed by the EigenDA network", "requestID", base64RequestID)
				} else {
					m.Log.Info("Blob is being processed by the EigenDA network", "requestID", base64RequestID)
					alreadyWaitingForDispersal = true
				}
			case grpcdisperser.BlobStatus_FAILED:
				// TODO: when exactly does this happen? I think only happens if ethereum reorged and the blob was lost
				errChan <- api.NewErrorAPIGeneric(http.StatusInternalServerError, fmt.Errorf("blob dispersal (requestID=%s) reached failed status. please resubmit the blob.", base64RequestID))
				return
			case grpcdisperser.BlobStatus_INSUFFICIENT_SIGNATURES:
				// this might be a temporary condition where some eigenda nodes were temporarily offline, so we should retry
				errChan <- api.NewErrorAPIGeneric(http.StatusInternalServerError, fmt.Errorf("blob dispersal (requestID=%s) failed with insufficient signatures. please resubmit the blob.", base64RequestID))
				return
			case grpcdisperser.BlobStatus_CONFIRMED:
				if m.Config.WaitForFinalization {
					// to prevent log clutter, we only log at info level once
					if alreadyWaitingForFinalization {
						m.Log.Debug("EigenDA blob confirmed, waiting for finalization", "requestID", base64RequestID)
					} else {
						m.Log.Info("EigenDA blob confirmed, waiting for finalization", "requestID", base64RequestID)
						alreadyWaitingForFinalization = true
					}
				} else {
					m.Log.Info("EigenDA blob confirmed", "requestID", base64RequestID)
					resultChan <- statusRes.Info
					return
				}
			case grpcdisperser.BlobStatus_FINALIZED:
				batchHeaderHashHex := fmt.Sprintf("0x%s", hex.EncodeToString(statusRes.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash))
				m.Log.Info("EigenDA blob finalized", "requestID", base64RequestID, "batchHeaderHash", batchHeaderHashHex)
				resultChan <- statusRes.Info
				return
			default:
				// this should never happen. If it does, the blob is in a heisenberg state... it could either eventually get confirmed or fail
				errChan <- api.NewErrorAPIGeneric(http.StatusInternalServerError, fmt.Errorf("unknown reply status %d. ask for assistance from EigenDA team, using requestID %s", statusRes.Status, base64RequestID))
				return
			}
		}
	}
}

// Close simply calls Close() on the wrapped disperserClient, to close the grpc connection to the disperser server.
// It is thread safe and can be called multiple times.
func (c *EigenDAClient) Close() error {
	return c.Client.Close()
}
