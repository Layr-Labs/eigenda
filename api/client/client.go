package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	grpcdisperser "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/clients"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/ethereum/go-ethereum/log"
)

type BlobEncodingVersion byte

var NoIFFT BlobEncodingVersion = 0x00

type IEigenDAClient interface {
	GetBlob(ctx context.Context, BatchHeaderHash []byte, BlobIndex uint32) ([]byte, error)
	PutBlob(ctx context.Context, txData []byte) (*Cert, error)
}

type EigenDAClient struct {
	Config
	Log    log.Logger
	client clients.DisperserClient
}

var _ IEigenDAClient = EigenDAClient{}

func NewEigenDAClient(log log.Logger, config Config) (*EigenDAClient, error) {
	err := config.Check()
	if err != nil {
		return nil, err
	}

	host, port, err := net.SplitHostPort(config.RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EigenDA RPC: %w", err)
	}

	signer := auth.NewLocalBlobRequestSigner(config.SignerPrivateKeyHex)
	llConfig := clients.NewConfig(host, port, config.ResponseTimeout, !config.DisableTLS)
	llClient := clients.NewDisperserClient(llConfig, signer)
	return &EigenDAClient{
		Log:    log,
		Config: config,
		client: llClient,
	}, nil
}

func (m EigenDAClient) GetBlob(ctx context.Context, BatchHeaderHash []byte, BlobIndex uint32) ([]byte, error) {
	data, err := m.client.RetrieveBlob(ctx, BatchHeaderHash, BlobIndex)
	if err != nil {
		return nil, err
	}

	// decode modulo bn254
	decodedData := codec.RemoveEmptyByteFromPaddedBytes(data)

	// Return exact data with buffer removed
	reader := bytes.NewReader(decodedData)

	// read version byte, we will not use it for now since there is only one version
	_, err = reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to read version byte")
	}

	// read length uvarint
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to decode length uvarint prefix")
	}

	result := make([]byte, length)
	n, err := reader.Read(result)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to copy unpadded data into final buffer")
	}
	if uint64(n) != length {
		return nil, fmt.Errorf("EigenDA client failed, data length does not match length prefix")
	}

	return result, nil
}

func (m EigenDAClient) PutBlob(ctx context.Context, data []byte) (*Cert, error) {
	resultChan, errorChan := m.PutBlobAsync(ctx, data)
	select { // no timeout here because we depend on the configured timeout in PutBlobAsync
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (m EigenDAClient) PutBlobAsync(ctx context.Context, data []byte) (resultChan chan *Cert, errChan chan error) {
	resultChan = make(chan *Cert, 1)
	errChan = make(chan error, 1)
	go m.putBlob(ctx, data, resultChan, errChan)
	return
}

func (m EigenDAClient) putBlob(ctx context.Context, data []byte, resultChan chan *Cert, errChan chan error) {
	m.Log.Info("Attempting to disperse blob to EigenDA")

	// encode current blob encoding version byte
	data = append([]byte{byte(NoIFFT)}, data...)

	// encode data length
	data = append(ConvertIntToVarUInt(len(data)), data...)

	// encode modulo bn254
	data = codec.ConvertByPaddingEmptyByte(data)

	customQuorumNumbers := make([]uint8, len(m.Config.CustomQuorumIDs))
	for i, e := range m.Config.CustomQuorumIDs {
		customQuorumNumbers[i] = uint8(e)
	}

	// do auth handshake
	blobStatus, requestID, err := m.client.DisperseBlobAuthenticated(ctx, data, customQuorumNumbers)
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
	m.Log.Info("Blob disepersed to EigenDA, now waiting for confirmation", "requestID", base64RequestID)

	ticker := time.NewTicker(m.StatusQueryRetryInterval)
	defer ticker.Stop()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.StatusQueryTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			errChan <- fmt.Errorf("timed out waiting for EigenDA blob to confirm blob with request id=%s: %w", base64RequestID, ctx.Err())
			return
		case <-ticker.C:
			statusRes, err := m.client.GetBlobStatus(ctx, requestID)
			if err != nil {
				m.Log.Error("Unable to retrieve blob dispersal status, will retry", "requestID", base64RequestID, "err", err)
				continue
			}

			switch statusRes.Status {
			case grpcdisperser.BlobStatus_PROCESSING:
				m.Log.Info("Waiting for confirmation from EigenDA", "requestID", base64RequestID)
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
				blobInfo := statusRes.Info
				quorumIDs := make([]uint32, len(blobInfo.BlobHeader.BlobQuorumParams))
				for i := range quorumIDs {
					quorumIDs[i] = blobInfo.BlobHeader.BlobQuorumParams[i].QuorumNumber
				}
				cert := &Cert{
					BatchHeaderHash:      blobInfo.BlobVerificationProof.BatchMetadata.BatchHeaderHash,
					BlobIndex:            blobInfo.BlobVerificationProof.BlobIndex,
					ReferenceBlockNumber: blobInfo.BlobVerificationProof.BatchMetadata.BatchHeader.ReferenceBlockNumber,
					QuorumIDs:            quorumIDs,
					BlobCommitment:       blobInfo.BlobHeader.Commitment,
				}
				resultChan <- cert
				return
			default:
				errChan <- fmt.Errorf("EigenDA blob dispersal failed in processing with reply status %d", statusRes.Status)
				return
			}
		}
	}
}
