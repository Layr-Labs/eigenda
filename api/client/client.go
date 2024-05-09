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
	client, err := NewDisperserClient(config.RPC)
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

func NewDisperserClient(rpc string) (disperser.DisperserClient, error) {
	config := &tls.Config{}
	credential := credentials.NewTLS(config)
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(credential)}
	conn, err := grpc.Dial(rpc, dialOptions...)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("error initializing DisperseBlobAuthenticated() client: %w", err)
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
		return nil, fmt.Errorf("failed sending initial disperse blob authenticated request: %w", err)
	}

	reply, err := disperseBlobAuthClient.Recv()
	if err != nil {
		return nil, fmt.Errorf("failed receiving challenge parameter for disperse blob authenticated request: %w", err)
	}

	authHeaderReply, ok := reply.Payload.(*disperser.AuthenticatedReply_BlobAuthHeader)
	if !ok {
		return nil, fmt.Errorf("expected blob auth header message in response to initial disperse blob authenticated request: %w", err)
	}

	authHeader := core.BlobAuthHeader{
		BlobCommitments: encoding.BlobCommitments{},
		AccountID:       "",
		Nonce:           authHeaderReply.BlobAuthHeader.ChallengeParameter,
	}

	authData, err := m.signer.SignBlobRequest(authHeader)
	if err != nil {
		return nil, fmt.Errorf("error signing challenge parameter while performing disperse blob authenticated request: %w", err)
	}

	// Process challenge and send back challenge_reply
	err = disperseBlobAuthClient.Send(&disperser.AuthenticatedRequest{Payload: &disperser.AuthenticatedRequest_AuthenticationData{
		AuthenticationData: &disperser.AuthenticationData{
			AuthenticationData: authData,
		},
	}})
	if err != nil {
		return nil, fmt.Errorf("error writing signed challenge paramter in disperse blob authenticated request: %w", err)
	}

	reply, err = disperseBlobAuthClient.Recv()
	if err != nil {
		return nil, fmt.Errorf("error receiving signed challenge paramter in disperse blob authenticated request: %w", err)
	}

	disperseResWrapper, ok := reply.Payload.(*disperser.AuthenticatedReply_DisperseReply)
	if !ok {
		return nil, fmt.Errorf("expected disperser reply message in response to signed challenge parameter submission in disperse blob authenticated request: %w", err)
	}

	disperseRes := disperseResWrapper.DisperseReply

	// process response
	if disperseRes.Result == disperser.BlobStatus_UNKNOWN ||
		disperseRes.Result == disperser.BlobStatus_FAILED {
		m.Log.Error("Unable to disperse blob to EigenDA, aborting", "err", err)
		return nil, fmt.Errorf("reply status is %d", disperseRes.Result)
	}

	base64RequestID := base64.StdEncoding.EncodeToString(disperseRes.RequestId)

	m.Log.Info("Blob disepersed to EigenDA, now waiting for confirmation", "requestID", base64RequestID)

	timeoutTime := time.Now().Add(m.StatusQueryTimeout)
	ticker := time.NewTicker(m.StatusQueryRetryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			if time.Now().After(timeoutTime) {
				return nil, fmt.Errorf("timed out waiting for EigenDA blob to confirm blob with request id: %s", base64RequestID)
			}
			statusRes, err := m.client.GetBlobStatus(ctx, &disperser.BlobStatusRequest{
				RequestId: disperseRes.RequestId,
			})
			if err != nil {
				m.Log.Warn("Unable to retrieve blob dispersal status, will retry", "requestID", base64RequestID, "err", err)
				continue
			}

			switch statusRes.Status {
			case disperser.BlobStatus_PROCESSING:
				m.Log.Info("Waiting for confirmation from EigenDA", "requestID", base64RequestID)
			case disperser.BlobStatus_FAILED:
				m.Log.Error("EigenDA blob dispersal failed in processing", "requestID", base64RequestID, "err", err)
			case disperser.BlobStatus_INSUFFICIENT_SIGNATURES:
				m.Log.Error("EigenDA blob dispersal failed in processing with insufficient signatures", "requestID", base64RequestID, "err", err)
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
				return cert, nil
			default:
				return nil, fmt.Errorf("EigenDA blob dispersal failed in processing with reply status %d", statusRes.Status)
			}
		}
	}
}
