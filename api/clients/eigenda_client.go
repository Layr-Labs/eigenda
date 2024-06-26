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
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/ethereum/go-ethereum/log"
)

type IEigenDAClient interface {
	GetBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error)
	PutBlob(ctx context.Context, txData []byte) (*grpcdisperser.BlobInfo, error)
	GetCodec() codecs.BlobCodec
}

type EigenDAClient struct {
	Config EigenDAClientConfig
	Log    log.Logger
	Client DisperserClient
	Codec  codecs.BlobCodec
}

var _ IEigenDAClient = EigenDAClient{}

func NewEigenDAClient(log log.Logger, config EigenDAClientConfig) (*EigenDAClient, error) {
	err := config.CheckAndSetDefaults()
	if err != nil {
		return nil, err
	}

	host, port, err := net.SplitHostPort(config.RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EigenDA RPC: %w", err)
	}

	signer := auth.NewLocalBlobRequestSigner(config.SignerPrivateKeyHex)
	llConfig := NewConfig(host, port, config.ResponseTimeout, !config.DisableTLS)
	llClient := NewDisperserClient(llConfig, signer)

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
		Client: llClient,
		Codec:  codec,
	}, nil
}

func (m EigenDAClient) GetCodec() codecs.BlobCodec {
	return m.Codec
}

// GetBlob retrieves a blob from the EigenDA service using the provided context,
// batch header hash, and blob index.  If decode is set to true, the function
// decodes the retrieved blob data. If set to false it returns the encoded blob
// data, which is necessary for generating KZG proofs for data's correctness.
// The function handles potential errors during blob retrieval, data length
// checks, and decoding processes.
func (m EigenDAClient) GetBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error) {
	data, err := m.Client.RetrieveBlob(ctx, batchHeaderHash, blobIndex)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("blob has length zero")
	}

	decodedData, err := m.Codec.DecodeBlob(data)
	if err != nil {
		return nil, fmt.Errorf("error getting blob: %w", err)
	}

	return decodedData, nil
}

// PutBlob encodes and writes a blob to EigenDA, waiting for it to be finalized
// before returning. This function is resiliant to transient failures and
// timeouts.
func (m EigenDAClient) PutBlob(ctx context.Context, data []byte) (*grpcdisperser.BlobInfo, error) {
	resultChan, errorChan := m.PutBlobAsync(ctx, data)
	select { // no timeout here because we depend on the configured timeout in PutBlobAsync
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (m EigenDAClient) PutBlobAsync(ctx context.Context, data []byte) (resultChan chan *grpcdisperser.BlobInfo, errChan chan error) {
	resultChan = make(chan *grpcdisperser.BlobInfo, 1)
	errChan = make(chan error, 1)
	go m.putBlob(ctx, data, resultChan, errChan)
	return
}

func (m EigenDAClient) putBlob(ctx context.Context, rawData []byte, resultChan chan *grpcdisperser.BlobInfo, errChan chan error) {
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
	blobStatus, requestID, err := m.Client.DisperseBlobAuthenticated(ctx, data, customQuorumNumbers)
	if err != nil {
		errChan <- fmt.Errorf("error initializing DisperseBlobAuthenticated() client: %w", err)
		return
	}

	// process response
	if *blobStatus == disperser.Failed {
		m.Log.Error("Unable to disperse blob to EigenDA, aborting", "err", err)
		errChan <- fmt.Errorf("reply status is %d", blobStatus)
		return
	}

	base64RequestID := base64.StdEncoding.EncodeToString(requestID)
	m.Log.Info("Blob dispersed to EigenDA, now waiting for confirmation", "requestID", base64RequestID)

	ticker := time.NewTicker(m.Config.StatusQueryRetryInterval)
	defer ticker.Stop()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.Config.StatusQueryTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			errChan <- fmt.Errorf("timed out waiting for EigenDA blob to confirm blob with request id=%s: %w", base64RequestID, ctx.Err())
			return
		case <-ticker.C:
			statusRes, err := m.Client.GetBlobStatus(ctx, requestID)
			if err != nil {
				m.Log.Error("Unable to retrieve blob dispersal status, will retry", "requestID", base64RequestID, "err", err)
				continue
			}

			switch statusRes.Status {
			case grpcdisperser.BlobStatus_PROCESSING, grpcdisperser.BlobStatus_DISPERSING:
				m.Log.Info("Blob submitted, waiting for dispersal from EigenDA", "requestID", base64RequestID)
			case grpcdisperser.BlobStatus_FAILED:
				m.Log.Error("EigenDA blob dispersal failed in processing", "requestID", base64RequestID, "err", err)
				errChan <- fmt.Errorf("EigenDA blob dispersal failed in processing, requestID=%s: %w", base64RequestID, err)
				return
			case grpcdisperser.BlobStatus_INSUFFICIENT_SIGNATURES:
				m.Log.Error("EigenDA blob dispersal failed in processing with insufficient signatures", "requestID", base64RequestID, "err", err)
				errChan <- fmt.Errorf("EigenDA blob dispersal failed in processing with insufficient signatures, requestID=%s: %w", base64RequestID, err)
				return
			case grpcdisperser.BlobStatus_CONFIRMED:
				if m.Config.WaitForFinalization {
					m.Log.Info("EigenDA blob confirmed, waiting for finalization", "requestID", base64RequestID)
				} else {
					m.Log.Info("EigenDA blob confirmed", "requestID", base64RequestID)
					resultChan <- statusRes.Info
					return
				}
			case grpcdisperser.BlobStatus_FINALIZED:
				batchHeaderHashHex := fmt.Sprintf("0x%s", hex.EncodeToString(statusRes.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash))
				m.Log.Info("Successfully dispersed blob to EigenDA", "requestID", base64RequestID, "batchHeaderHash", batchHeaderHashHex)
				resultChan <- statusRes.Info
				return
			default:
				errChan <- fmt.Errorf("EigenDA blob dispersal failed in processing with reply status %d", statusRes.Status)
				return
			}
		}
	}
}
