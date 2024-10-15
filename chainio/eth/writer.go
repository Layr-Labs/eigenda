package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/chainio"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/pingcap/errors"

	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/RegistryCoordinator"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Writer struct {
	*Reader

	EthClient common.EthClient
	Logger    logging.Logger
}

var _ chainio.Writer = (*Writer)(nil)

func NewWriter(
	logger logging.Logger,
	client common.EthClient,
	blsOperatorStateRetrieverHexAddr string,
	eigenDAServiceManagerHexAddr string) (*Writer, error) {

	r := &Reader{
		EthClient: client,
		Logger:    logger.With("component", "Reader"),
	}

	e := &Writer{
		EthClient: client,
		Logger:    logger.With("component", "Writer"),
		Reader:    r,
	}

	blsOperatorStateRetrieverAddr := gethcommon.HexToAddress(blsOperatorStateRetrieverHexAddr)
	eigenDAServiceManagerAddr := gethcommon.HexToAddress(eigenDAServiceManagerHexAddr)
	err := e.updateContractBindings(blsOperatorStateRetrieverAddr, eigenDAServiceManagerAddr)

	return e, err
}

// RegisterOperator registers a new operator with the given public key and socket with the provided quorum ids.
// If the operator is already registered with a given quorum id, the transaction will fail (noop) and an error
// will be returned.
func (t *Writer) RegisterOperator(
	ctx context.Context,
	keypair *bn254.KeyPair,
	socket string,
	quorumIds []chainio.QuorumID,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
) error {

	params, operatorSignature, err := t.getRegistrationParams(ctx, keypair, operatorEcdsaPrivateKey, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry)
	if err != nil {
		t.Logger.Error("Failed to get registration params", "err", err)
		return err
	}

	quorumNumbers := quorumIDsToQuorumNumbers(quorumIds)
	opts, err := t.EthClient.GetNoSendTransactOpts()
	if err != nil {
		t.Logger.Error("Failed to generate transact opts", "err", err)
		return err
	}

	tx, err := t.Bindings.RegistryCoordinator.RegisterOperator(opts, quorumNumbers, socket, *params, *operatorSignature)

	if err != nil {
		t.Logger.Error("Failed to register operator", "err", err)
		return err
	}

	_, err = t.EthClient.EstimateGasPriceAndLimitAndSendTx(context.Background(), tx, "RegisterOperatorWithCoordinator1", nil)
	if err != nil {
		t.Logger.Error("Failed to estimate gas price and limit", "err", err)
		return err
	}
	return nil
}

// RegisterOperatorWithChurn registers a new operator with the given public key and socket with the provided quorum ids
// with the provided signature from the churner
func (t *Writer) RegisterOperatorWithChurn(
	ctx context.Context,
	keypair *bn254.KeyPair,
	socket string,
	quorumIds []chainio.QuorumID,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
	churnReply *churner.ChurnReply,
) error {

	params, operatorSignature, err := t.getRegistrationParams(ctx, keypair, operatorEcdsaPrivateKey, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry)
	if err != nil {
		t.Logger.Error("Failed to get registration params", "err", err)
		return err
	}

	quorumNumbers := quorumIDsToQuorumNumbers(quorumIds)

	operatorsToChurn := make([]regcoordinator.IRegistryCoordinatorOperatorKickParam, len(churnReply.OperatorsToChurn))
	for i := range churnReply.OperatorsToChurn {
		if churnReply.OperatorsToChurn[i].QuorumId >= chainio.MaxQuorumID {
			return errors.New("quorum id is out of range")
		}

		operatorsToChurn[i] = regcoordinator.IRegistryCoordinatorOperatorKickParam{
			QuorumNumber: uint8(churnReply.OperatorsToChurn[i].QuorumId),
			Operator:     gethcommon.BytesToAddress(churnReply.OperatorsToChurn[i].Operator),
		}
	}

	var salt [32]byte
	copy(salt[:], churnReply.SignatureWithSaltAndExpiry.Salt[:])
	churnApproverSignature := regcoordinator.ISignatureUtilsSignatureWithSaltAndExpiry{
		Signature: churnReply.SignatureWithSaltAndExpiry.Signature,
		Salt:      salt,
		Expiry:    new(big.Int).SetInt64(churnReply.SignatureWithSaltAndExpiry.Expiry),
	}

	opts, err := t.EthClient.GetNoSendTransactOpts()
	if err != nil {
		t.Logger.Error("Failed to generate transact opts", "err", err)
		return err
	}

	tx, err := t.Bindings.RegistryCoordinator.RegisterOperatorWithChurn(
		opts,
		quorumNumbers,
		socket,
		*params,
		operatorsToChurn,
		churnApproverSignature,
		*operatorSignature,
	)

	if err != nil {
		t.Logger.Error("Failed to register operator with churn", "err", err)
		return err
	}

	_, err = t.EthClient.EstimateGasPriceAndLimitAndSendTx(context.Background(), tx, "RegisterOperatorWithCoordinatorWithChurn", nil)
	if err != nil {
		t.Logger.Error("Failed to estimate gas price and limit", "err", err)
		return err
	}
	return nil
}

// DeregisterOperator deregisters an operator with the given public key from the specified the quorums that it is
// registered with at the supplied block number. To fully deregister an operator, this function should be called
// with the current block number.
// If the operator isn't registered with any of the specified quorums, this function will return error, and
// no quorum will be deregistered.
func (t *Writer) DeregisterOperator(ctx context.Context, pubkeyG1 *bn254.G1Point, blockNumber uint32, quorumIds []chainio.QuorumID) error {
	if len(quorumIds) == 0 {
		return errors.New("no quorum is specified to deregister from")
	}
	// Make sure the operator is registered in all the quorums it tries to deregister.
	operatorId := HashPubKeyG1(pubkeyG1)
	quorumBitmap, _, err := t.Bindings.OpStateRetriever.GetOperatorState0(&bind.CallOpts{
		Context: ctx,
	}, t.Bindings.RegCoordinatorAddr, operatorId, blockNumber)
	if err != nil {
		t.Logger.Error("Failed to fetch operator state", "err", err)
		return err
	}

	quorumNumbers := bitmapToBytesArray(quorumBitmap)
	for _, quorumToDereg := range quorumIds {
		found := false
		for _, currentQuorum := range quorumNumbers {
			if quorumToDereg == currentQuorum {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("operatorId %s is not registered in quorum %d at block %d", hex.EncodeToString(operatorId[:]), quorumToDereg, blockNumber)
		}
	}

	opts, err := t.EthClient.GetNoSendTransactOpts()
	if err != nil {
		t.Logger.Error("Failed to generate transact opts", "err", err)
		return err
	}
	tx, err := t.Bindings.RegistryCoordinator.DeregisterOperator(
		opts,
		quorumIds,
	)
	if err != nil {
		t.Logger.Error("Failed to deregister operator", "err", err)
		return err
	}

	_, err = t.EthClient.EstimateGasPriceAndLimitAndSendTx(context.Background(), tx, "DeregisterOperatorWithCoordinator", nil)
	if err != nil {
		t.Logger.Error("Failed to estimate gas price and limit", "err", err)
		return err
	}
	return nil
}

// UpdateOperatorSocket updates the socket of the operator in all the quorums that it is
func (t *Writer) UpdateOperatorSocket(ctx context.Context, socket string) error {
	opts, err := t.EthClient.GetNoSendTransactOpts()
	if err != nil {
		t.Logger.Error("Failed to generate transact opts", "err", err)
		return err
	}
	tx, err := t.Bindings.RegistryCoordinator.UpdateSocket(opts, socket)
	if err != nil {
		t.Logger.Error("Failed to update operator socket", "err", err)
		return err
	}

	_, err = t.EthClient.EstimateGasPriceAndLimitAndSendTx(context.Background(), tx, "UpdateOperatorSocket", nil)
	if err != nil {
		t.Logger.Error("Failed to estimate gas price and limit", "err", err)
		return err
	}
	return nil
}

func (t *Writer) BuildEjectOperatorsTxn(ctx context.Context, operatorsByQuorum [][][32]byte) (*types.Transaction, error) {
	byteIdsByQuorum := make([][][32]byte, len(operatorsByQuorum))
	for i, ids := range operatorsByQuorum {
		for _, id := range ids {
			byteIdsByQuorum[i] = append(byteIdsByQuorum[i], [32]byte(id))
		}
	}
	opts, err := t.EthClient.GetNoSendTransactOpts()
	if err != nil {
		t.Logger.Error("Failed to generate transact opts", "err", err)
		return nil, err
	}
	return t.Bindings.EjectionManager.EjectOperators(opts, byteIdsByQuorum)
}
