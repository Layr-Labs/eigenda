package clients

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"time"

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
func (m *EigenDAClient) PutBlob(ctx context.Context, data []byte) (*grpcdisperser.BlobInfo, error) {
	resultChan, errorChan := m.PutBlobAsync(ctx, data)
	select { // no timeout here because we depend on the configured timeout in PutBlobAsync
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (m *EigenDAClient) PutBlobAsync(ctx context.Context, data []byte) (resultChan chan *grpcdisperser.BlobInfo, errChan chan error) {
	resultChan = make(chan *grpcdisperser.BlobInfo, 1)
	errChan = make(chan error, 1)
	go m.putBlob(ctx, data, resultChan, errChan)
	return
}

func (m *EigenDAClient) putBlob(ctx context.Context, rawData []byte, resultChan chan *grpcdisperser.BlobInfo, errChan chan error) {
	m.Log.Info("Attempting to disperse blob to EigenDA")

	// encode blob
	if m.Codec == nil {
		errChan <- fmt.Errorf("Codec cannot be nil")
		return
	}

	data, err := m.Codec.EncodeBlob(rawData)
	if err != nil {
		errChan <- fmt.Errorf("error encoding blob: %w", err)
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
		errChan <- fmt.Errorf("error initializing DisperseBlobAuthenticated() client: %w", err)
		return
	}

	// process response
	if *blobStatus == disperser.Failed {
		errChan <- fmt.Errorf("unable to disperse blob to eigenda (reply status %d): %w", blobStatus, err)
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
	for {
		select {
		case <-ctx.Done():
			errChan <- fmt.Errorf("timed out waiting for EigenDA blob to confirm blob with request id=%s: %w", base64RequestID, ctx.Err())
			return
		case <-ticker.C:
			statusRes, err := m.Client.GetBlobStatus(ctx, requestID)
			if err != nil {
				m.Log.Warn("Unable to retrieve blob dispersal status, will retry", "requestID", base64RequestID, "err", err)
				continue
			}

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
				errChan <- fmt.Errorf("EigenDA blob dispersal failed in processing, requestID=%s: %w", base64RequestID, err)
				return
			case grpcdisperser.BlobStatus_INSUFFICIENT_SIGNATURES:
				errChan <- fmt.Errorf("EigenDA blob dispersal failed in processing with insufficient signatures, requestID=%s: %w", base64RequestID, err)
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
				errChan <- fmt.Errorf("unknown reply status %d. ask for assistance from EigenDA team, using requestID %s", statusRes.Status, base64RequestID)
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
