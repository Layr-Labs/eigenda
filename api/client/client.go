package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type BlobEncodingVersion byte

var NoIFFT BlobEncodingVersion = 0x00

type IEigenDAClient interface {
	GetBlob(ctx context.Context, BatchHeaderHash []byte, BlobIndex uint32) ([]byte, error)
	PutBlob(ctx context.Context, txData []byte) (*Cert, error)
}

type EigenDAClient struct {
	Config

	Log log.Logger

	client disperser.DisperserClient

	signer *auth.LocalBlobRequestSigner
}

var _ IEigenDAClient = EigenDAClient{}

func NewEigenDAClient(log log.Logger, config Config) (*EigenDAClient, error) {
	err := config.Check()
	if err != nil {
		return nil, err
	}
	client, err := NewDisperserClient(config.RPC, config.DisableTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to instatiated EigenDA GRPC client: %w", err)
	}
	return &EigenDAClient{
		Log:    log,
		Config: config,
		client: client,
		signer: auth.NewLocalBlobRequestSigner(config.SignerPrivateKeyHex),
	}, nil
}

func NewDisperserClient(rpc string, disableTLS bool) (disperser.DisperserClient, error) {
	var credentialOptions grpc.DialOption
	if !disableTLS {
		config := &tls.Config{}
		credential := credentials.NewTLS(config)
		credentialOptions = grpc.WithTransportCredentials(credential)
	} else {
		credentialOptions = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.Dial(rpc, credentialOptions)
	if err != nil {
		return nil, fmt.Errorf("error dialing EigenDA GRPC connection: %w", err)
	}
	client := disperser.NewDisperserClient(conn)
	return client, nil
}

func (m EigenDAClient) GetBlob(ctx context.Context, BatchHeaderHash []byte, BlobIndex uint32) ([]byte, error) {
	reply, err := m.client.RetrieveBlob(ctx, &disperser.RetrieveBlobRequest{
		BatchHeaderHash: BatchHeaderHash,
		BlobIndex:       BlobIndex,
	})
	if err != nil {
		return nil, err
	}

	// decode modulo bn254
	decodedData := codec.RemoveEmptyByteFromPaddedBytes(reply.Data)

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
	data := make([]byte, length)
	n, err := reader.Read(data)
	if err != nil {
		return nil, fmt.Errorf("EigenDA client failed to copy unpadded data into final buffer")
	}
	if uint64(n) != length {
		return nil, fmt.Errorf("EigenDA client failed, data length does not match length prefix")
	}

	return data, nil
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

	// do auth handshake
	disperseBlobAuthClient, err := m.client.DisperseBlobAuthenticated(ctx)
	if err != nil {
		errChan <- fmt.Errorf("error initializing DisperseBlobAuthenticated() client: %w", err)
		return
	}

	customQuorumNumbers := make([]uint32, len(m.Config.CustomQuorumIDs))
	for i, e := range m.Config.CustomQuorumIDs {
		customQuorumNumbers[i] = uint32(e)
	}
	err = disperseBlobAuthClient.Send(&disperser.AuthenticatedRequest{
		Payload: &disperser.AuthenticatedRequest_DisperseRequest{
			DisperseRequest: &disperser.DisperseBlobRequest{
				Data:                data,
				CustomQuorumNumbers: customQuorumNumbers,
				AccountId:           m.signer.GetAccountID(),
			},
		},
	})
	if err != nil {
		errChan <- fmt.Errorf("failed sending initial disperse blob authenticated request: %w", err)
		return
	}

	reply, err := disperseBlobAuthClient.Recv()
	if err != nil {
		errChan <- fmt.Errorf("failed receiving challenge parameter for disperse blob authenticated request: %w", err)
		return
	}

	authHeaderReply, ok := reply.Payload.(*disperser.AuthenticatedReply_BlobAuthHeader)
	if !ok {
		errChan <- fmt.Errorf("expected blob auth header message in response to initial disperse blob authenticated request: %w", err)
		return
	}

	authHeader := core.BlobAuthHeader{
		BlobCommitments: encoding.BlobCommitments{},
		AccountID:       m.signer.GetAccountID(),
		Nonce:           authHeaderReply.BlobAuthHeader.ChallengeParameter,
	}

	authData, err := m.signer.SignBlobRequest(authHeader)
	if err != nil {
		errChan <- fmt.Errorf("error signing challenge parameter while performing disperse blob authenticated request: %w", err)
		return
	}

	// Process challenge and send back challenge_reply
	err = disperseBlobAuthClient.Send(&disperser.AuthenticatedRequest{Payload: &disperser.AuthenticatedRequest_AuthenticationData{
		AuthenticationData: &disperser.AuthenticationData{
			AuthenticationData: authData,
		},
	}})
	if err != nil {
		errChan <- fmt.Errorf("error writing signed challenge paramter in disperse blob authenticated request: %w", err)
		return
	}

	reply, err = disperseBlobAuthClient.Recv()
	if err != nil {
		errChan <- fmt.Errorf("error receiving signed challenge paramter in disperse blob authenticated request: %w", err)
		return
	}

	disperseResWrapper, ok := reply.Payload.(*disperser.AuthenticatedReply_DisperseReply)
	if !ok {
		errChan <- fmt.Errorf("expected disperser reply message in response to signed challenge parameter submission in disperse blob authenticated request: %w", err)
		return
	}

	disperseRes := disperseResWrapper.DisperseReply

	// process response
	if disperseRes.Result == disperser.BlobStatus_UNKNOWN ||
		disperseRes.Result == disperser.BlobStatus_FAILED {
		m.Log.Error("Unable to disperse blob to EigenDA, aborting", "err", err)
		errChan <- fmt.Errorf("reply status is %d", disperseRes.Result)
		return
	}

	base64RequestID := base64.StdEncoding.EncodeToString(disperseRes.RequestId)

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
			statusRes, err := m.client.GetBlobStatus(ctx, &disperser.BlobStatusRequest{
				RequestId: disperseRes.RequestId,
			})
			if err != nil {
				m.Log.Error("Unable to retrieve blob dispersal status, will retry", "requestID", base64RequestID, "err", err)
				continue
			}

			switch statusRes.Status {
			case disperser.BlobStatus_PROCESSING:
				m.Log.Info("Waiting for confirmation from EigenDA", "requestID", base64RequestID)
			case disperser.BlobStatus_FAILED:
				m.Log.Error("EigenDA blob dispersal failed in processing", "requestID", base64RequestID, "err", err)
				errChan <- fmt.Errorf("EigenDA blob dispersal failed in processing, requestID=%s: %w", base64RequestID, err)
				return
			case disperser.BlobStatus_INSUFFICIENT_SIGNATURES:
				m.Log.Error("EigenDA blob dispersal failed in processing with insufficient signatures", "requestID", base64RequestID, "err", err)
				errChan <- fmt.Errorf("EigenDA blob dispersal failed in processing with insufficient signatures, requestID=%s: %w", base64RequestID, err)
				return
			case disperser.BlobStatus_CONFIRMED:
				m.Log.Info("EigenDA blob confirmed, waiting for finalization", "requestID", base64RequestID)
			case disperser.BlobStatus_FINALIZED:
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
