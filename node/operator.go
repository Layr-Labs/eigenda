package node

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"time"

	grpcchurner "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Operator struct {
	Address    string
	Socket     string
	Timeout    time.Duration
	PrivKey    *ecdsa.PrivateKey
	KeyPair    *core.KeyPair
	OperatorId core.OperatorID
	QuorumIDs  []core.QuorumID
}

// RegisterOperator operator registers the operator with the given public key for the given quorum IDs.
func RegisterOperator(ctx context.Context, operator *Operator, transactor core.Transactor, churnerUrl string, useSecureGrpc bool, logger common.Logger) error {
	registeredQuorumIds, err := transactor.GetRegisteredQuorumIdsForOperator(ctx, operator.OperatorId)
	if err != nil {
		return fmt.Errorf("failed to get registered quorum ids for an operator: %w", err)
	}

	logger.Debug("Registered quorum ids", "registeredQuorumIds", registeredQuorumIds)
	if len(registeredQuorumIds) != 0 {
		return nil
	}

	logger.Info("Quorums to register for", "quorums", operator.QuorumIDs)

	if len(operator.QuorumIDs) == 0 {
		return errors.New("an operator should be in at least one quorum to be useful")
	}

	// register for quorums
	shouldCallChurner := false
	// check if one of the quorums to register for is full
	for _, quorumID := range operator.QuorumIDs {
		operatorSetParams, err := transactor.GetOperatorSetParams(ctx, quorumID)
		if err != nil {
			return err
		}

		numberOfRegisteredOperators, err := transactor.GetNumberOfRegisteredOperatorForQuorum(ctx, quorumID)
		if err != nil {
			return err
		}

		// if the quorum is full, we need to call the churner
		if operatorSetParams.MaxOperatorCount == numberOfRegisteredOperators {
			shouldCallChurner = true
			break
		}
	}

	logger.Info("Should call churner", "shouldCallChurner", shouldCallChurner)

	// Generate salt and expiry

	privateKeyBytes := []byte(operator.KeyPair.PrivKey.String())
	salt := [32]byte{}
	copy(salt[:], crypto.Keccak256([]byte("churn"), []byte(time.Now().String()), operator.QuorumIDs[:], privateKeyBytes))

	// Get the current block number
	expiry := big.NewInt((time.Now().Add(10 * time.Minute)).Unix())

	// if we should call the churner, call it
	if shouldCallChurner {
		churnReply, err := requestChurnApproval(ctx, operator, churnerUrl, useSecureGrpc, logger)
		if err != nil {
			return fmt.Errorf("failed to request churn approval: %w", err)
		}

		return transactor.RegisterOperatorWithChurn(ctx, operator.KeyPair, operator.Socket, operator.QuorumIDs, operator.PrivKey, salt, expiry, churnReply)
	} else {
		// other wise just register normally
		return transactor.RegisterOperator(ctx, operator.KeyPair, operator.Socket, operator.QuorumIDs, operator.PrivKey, salt, expiry)
	}
}

// DeregisterOperator deregisters the operator with the given public key from the all the quorums that it is registered with at the supplied block number.
func DeregisterOperator(ctx context.Context, KeyPair *core.KeyPair, transactor core.Transactor) error {
	blockNumber, err := transactor.GetCurrentBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	return transactor.DeregisterOperator(ctx, KeyPair.GetPubKeyG1(), blockNumber)
}

// UpdateOperatorQuorums updates the quorums for the given operator
func UpdateOperatorQuorums(
	ctx context.Context,
	operator *Operator,
	transactor core.Transactor,
	churnerUrl string,
	useSecureGrpc bool,
	logger common.Logger,
) error {
	err := DeregisterOperator(ctx, operator.KeyPair, transactor)
	if err != nil {
		return fmt.Errorf("failed to deregister operator: %w", err)
	}
	return RegisterOperator(ctx, operator, transactor, churnerUrl, useSecureGrpc, logger)
}

// UpdateOperatorSocket updates the socket for the given operator
func UpdateOperatorSocket(ctx context.Context, transactor core.Transactor, socket string) error {
	return transactor.UpdateOperatorSocket(ctx, socket)
}

func requestChurnApproval(ctx context.Context, operator *Operator, churnerUrl string, useSecureGrpc bool, logger common.Logger) (*grpcchurner.ChurnReply, error) {
	logger.Info("churner url", "url", churnerUrl)

	credential := insecure.NewCredentials()
	if useSecureGrpc {
		config := &tls.Config{}
		credential = credentials.NewTLS(config)
	}

	conn, err := grpc.Dial(
		churnerUrl,
		grpc.WithTransportCredentials(credential),
	)
	if err != nil {
		logger.Error("Node cannot connect to churner", "err", err)
		return nil, err
	}
	defer conn.Close()

	gc := grpcchurner.NewChurnerClient(conn)
	ctx, cancel := context.WithTimeout(ctx, operator.Timeout)
	defer cancel()

	request := newChurnRequest(operator.Address, operator.KeyPair, operator.QuorumIDs)
	opt := grpc.MaxCallSendMsgSize(1024 * 1024 * 300)

	return gc.Churn(ctx, request, opt)
}

func newChurnRequest(address string, KeyPair *core.KeyPair, QuorumIDs []core.QuorumID) *grpcchurner.ChurnRequest {

	// generate salt
	privateKeyBytes := []byte(KeyPair.PrivKey.String())
	salt := crypto.Keccak256([]byte("churn"), []byte(time.Now().String()), QuorumIDs[:], privateKeyBytes)

	churnRequest := &churner.ChurnRequest{
		OperatorAddress:            gethcommon.HexToAddress(address),
		OperatorToRegisterPubkeyG1: KeyPair.PubKey,
		OperatorToRegisterPubkeyG2: KeyPair.GetPubKeyG2(),
		OperatorRequestSignature:   &core.Signature{},
		QuorumIDs:                  QuorumIDs,
	}

	copy(churnRequest.Salt[:], salt)

	// sign the request
	churnRequest.OperatorRequestSignature = KeyPair.SignMessage(churner.CalculateRequestHash(churnRequest))

	// convert to protobuf
	churnRequestPb := &grpcchurner.ChurnRequest{
		OperatorToRegisterPubkeyG1: churnRequest.OperatorToRegisterPubkeyG1.Serialize(),
		OperatorToRegisterPubkeyG2: churnRequest.OperatorToRegisterPubkeyG2.Serialize(),
		OperatorRequestSignature:   churnRequest.OperatorRequestSignature.Serialize(),
		Salt:                       salt[:],
		OperatorAddress:            address,
	}

	churnRequestPb.QuorumIds = make([]uint32, len(QuorumIDs))
	for i, quorumID := range QuorumIDs {
		churnRequestPb.QuorumIds[i] = uint32(quorumID)
	}

	return churnRequestPb
}
