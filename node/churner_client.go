package node

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	churnerpb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/operators/churner"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ChurnerClient interface {
	// Churn sends a churn request to the churner service
	// The quorumIDs cannot be empty, but may contain quorums that the operator is already registered in.
	// If the operator is already registered in a quorum, the churner will ignore it and continue with the other quorums.
	Churn(ctx context.Context, operatorAddress string, keyPair *core.KeyPair, quorumIDs []core.QuorumID) (*churnerpb.ChurnReply, error)
}

type churnerClient struct {
	churnerURL    string
	useSecureGrpc bool
	timeout       time.Duration
	logger        logging.Logger
}

func NewChurnerClient(churnerURL string, useSecureGrpc bool, timeout time.Duration, logger logging.Logger) ChurnerClient {
	return &churnerClient{
		churnerURL:    churnerURL,
		useSecureGrpc: useSecureGrpc,
		timeout:       timeout,
		logger:        logger.With("component", "ChurnerClient"),
	}
}

func (c *churnerClient) Churn(ctx context.Context, operatorAddress string, keyPair *core.KeyPair, quorumIDs []core.QuorumID) (*churnerpb.ChurnReply, error) {
	if len(quorumIDs) == 0 {
		return nil, errors.New("quorumIDs cannot be empty")
	}
	// generate salt
	privateKeyBytes := []byte(keyPair.PrivKey.String())
	salt := crypto.Keccak256([]byte("churn"), []byte(time.Now().String()), quorumIDs[:], privateKeyBytes)

	churnRequest := &churner.ChurnRequest{
		OperatorAddress:            gethcommon.HexToAddress(operatorAddress),
		OperatorToRegisterPubkeyG1: keyPair.PubKey,
		OperatorToRegisterPubkeyG2: keyPair.GetPubKeyG2(),
		OperatorRequestSignature:   &core.Signature{},
		QuorumIDs:                  quorumIDs,
	}

	copy(churnRequest.Salt[:], salt)

	// sign the request
	churnRequest.OperatorRequestSignature = keyPair.SignMessage(churner.CalculateRequestHash(churnRequest))

	// convert to protobuf
	churnRequestPb := &churnerpb.ChurnRequest{
		OperatorToRegisterPubkeyG1: churnRequest.OperatorToRegisterPubkeyG1.Serialize(),
		OperatorToRegisterPubkeyG2: churnRequest.OperatorToRegisterPubkeyG2.Serialize(),
		OperatorRequestSignature:   churnRequest.OperatorRequestSignature.Serialize(),
		Salt:                       salt[:],
		OperatorAddress:            operatorAddress,
	}

	churnRequestPb.QuorumIds = make([]uint32, len(quorumIDs))
	for i, quorumID := range quorumIDs {
		churnRequestPb.QuorumIds[i] = uint32(quorumID)
	}
	credential := insecure.NewCredentials()
	if c.useSecureGrpc {
		config := &tls.Config{}
		credential = credentials.NewTLS(config)
	}

	conn, err := grpc.NewClient(
		c.churnerURL,
		grpc.WithTransportCredentials(credential),
	)
	if err != nil {
		c.logger.Error("Node cannot connect to churner", "err", err)
		return nil, err
	}
	defer conn.Close()

	gc := churnerpb.NewChurnerClient(conn)
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	opt := grpc.MaxCallSendMsgSize(1024 * 1024 * 300)

	return gc.Churn(ctx, churnRequestPb, opt)
}
