package clients

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	grpcdisperser "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/ethereum/go-ethereum/log"
)

type IEigenDAClient interface {
	GetBlob(ctx context.Context, BatchHeaderHash []byte, BlobIndex uint32) ([]byte, error)
	PutBlob(ctx context.Context, txData []byte) (*grpcdisperser.BlobInfo, error)
}

type EigenDAClient struct {
	Config   EigenDAClientConfig
	Log      log.Logger
	Client   DisperserClient
	PutCodec BlobCodec
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

	codec, err := BlobEncodingVersionToCodec(config.PutBlobEncodingVersion)
	if err != nil {
		return nil, fmt.Errorf("error initializing EigenDA client: %w", err)
	}

	return &EigenDAClient{
		Log:      log,
		Config:   config,
		Client:   llClient,
		PutCodec: codec,
	}, nil
}

func (m EigenDAClient) GetBlob(ctx context.Context, BatchHeaderHash []byte, BlobIndex uint32) ([]byte, error) {
	data, err := m.Client.RetrieveBlob(ctx, BatchHeaderHash, BlobIndex)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("blob has length zero")
	}

	version := BlobEncodingVersion(data[0])
	codec, err := BlobEncodingVersionToCodec(version)
	if err != nil {
		return nil, fmt.Errorf("error getting blob: %w", err)
	}

	rawData, err := codec.DecodeBlob(data)
	if err != nil {
		return nil, fmt.Errorf("error getting blob: %w", err)
	}

	return rawData, nil
}

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
	if m.PutCodec == nil {
		errChan <- fmt.Errorf("PutCodec cannot be nil")
		return
	}
	data := m.PutCodec.EncodeBlob(rawData)

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
				m.Log.Info("EigenDA blob confirmed, waiting for finalization", "requestID", base64RequestID)
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
