package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/RegistryCoordinator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pingcap/errors"
)

type Writer struct {
	*Reader

	ethClient common.EthClient
	logger    logging.Logger
}

var _ core.Writer = (*Writer)(nil)

func NewWriter(
	logger logging.Logger,
	client common.EthClient,
	blsOperatorStateRetrieverHexAddr string,
	eigenDAServiceManagerHexAddr string) (*Writer, error) {

	r := &Reader{
		ethClient: client,
		logger:    logger.With("component", "Reader"),
	}

	e := &Writer{
		ethClient: client,
		logger:    logger.With("component", "Writer"),
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
	keypair *core.KeyPair,
	socket string,
	quorumIds []core.QuorumID,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
) error {

	params, operatorSignature, err := t.getRegistrationParams(ctx, keypair, operatorEcdsaPrivateKey, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry)
	if err != nil {
		t.logger.Error("Failed to get registration params", "err", err)
		return err
	}

	quorumNumbers := quorumIDsToQuorumNumbers(quorumIds)
	opts, err := t.ethClient.GetNoSendTransactOpts()
	if err != nil {
		t.logger.Error("Failed to generate transact opts", "err", err)
		return err
	}

	tx, err := t.bindings.RegistryCoordinator.RegisterOperator(opts, quorumNumbers, socket, *params, *operatorSignature)

	if err != nil {
		t.logger.Error("Failed to register operator", "err", err)
		return err
	}

	_, err = t.ethClient.EstimateGasPriceAndLimitAndSendTx(context.Background(), tx, "RegisterOperatorWithCoordinator1", nil)
	if err != nil {
		t.logger.Error("Failed to estimate gas price and limit", "err", err)
		return err
	}
	return nil
}

// RegisterOperatorWithChurn registers a new operator with the given public key and socket with the provided quorum ids
// with the provided signature from the churner
func (t *Writer) RegisterOperatorWithChurn(
	ctx context.Context,
	keypair *core.KeyPair,
	socket string,
	quorumIds []core.QuorumID,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
	churnReply *churner.ChurnReply,
) error {

	params, operatorSignature, err := t.getRegistrationParams(ctx, keypair, operatorEcdsaPrivateKey, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry)
	if err != nil {
		t.logger.Error("Failed to get registration params", "err", err)
		return err
	}

	quorumNumbers := quorumIDsToQuorumNumbers(quorumIds)

	operatorsToChurn := make([]regcoordinator.IRegistryCoordinatorOperatorKickParam, len(churnReply.OperatorsToChurn))
	for i := range churnReply.OperatorsToChurn {
		if churnReply.OperatorsToChurn[i].QuorumId >= core.MaxQuorumID {
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

	opts, err := t.ethClient.GetNoSendTransactOpts()
	if err != nil {
		t.logger.Error("Failed to generate transact opts", "err", err)
		return err
	}

	tx, err := t.bindings.RegistryCoordinator.RegisterOperatorWithChurn(
		opts,
		quorumNumbers,
		socket,
		*params,
		operatorsToChurn,
		churnApproverSignature,
		*operatorSignature,
	)

	if err != nil {
		t.logger.Error("Failed to register operator with churn", "err", err)
		return err
	}

	_, err = t.ethClient.EstimateGasPriceAndLimitAndSendTx(context.Background(), tx, "RegisterOperatorWithCoordinatorWithChurn", nil)
	if err != nil {
		t.logger.Error("Failed to estimate gas price and limit", "err", err)
		return err
	}
	return nil
}

// DeregisterOperator deregisters an operator with the given public key from the specified the quorums that it is
// registered with at the supplied block number. To fully deregister an operator, this function should be called
// with the current block number.
// If the operator isn't registered with any of the specified quorums, this function will return error, and
// no quorum will be deregistered.
func (t *Writer) DeregisterOperator(ctx context.Context, pubkeyG1 *core.G1Point, blockNumber uint32, quorumIds []core.QuorumID) error {
	if len(quorumIds) == 0 {
		return errors.New("no quorum is specified to deregister from")
	}
	// Make sure the operator is registered in all the quorums it tries to deregister.
	operatorId := HashPubKeyG1(pubkeyG1)
	quorumBitmap, _, err := t.bindings.OpStateRetriever.GetOperatorState0(&bind.CallOpts{
		Context: ctx,
	}, t.bindings.RegCoordinatorAddr, operatorId, blockNumber)
	if err != nil {
		t.logger.Error("Failed to fetch operator state", "err", err)
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

	opts, err := t.ethClient.GetNoSendTransactOpts()
	if err != nil {
		t.logger.Error("Failed to generate transact opts", "err", err)
		return err
	}
	tx, err := t.bindings.RegistryCoordinator.DeregisterOperator(
		opts,
		quorumIds,
	)
	if err != nil {
		t.logger.Error("Failed to deregister operator", "err", err)
		return err
	}

	_, err = t.ethClient.EstimateGasPriceAndLimitAndSendTx(context.Background(), tx, "DeregisterOperatorWithCoordinator", nil)
	if err != nil {
		t.logger.Error("Failed to estimate gas price and limit", "err", err)
		return err
	}
	return nil
}

// UpdateOperatorSocket updates the socket of the operator in all the quorums that it is
func (t *Writer) UpdateOperatorSocket(ctx context.Context, socket string) error {
	opts, err := t.ethClient.GetNoSendTransactOpts()
	if err != nil {
		t.logger.Error("Failed to generate transact opts", "err", err)
		return err
	}
	tx, err := t.bindings.RegistryCoordinator.UpdateSocket(opts, socket)
	if err != nil {
		t.logger.Error("Failed to update operator socket", "err", err)
		return err
	}

	_, err = t.ethClient.EstimateGasPriceAndLimitAndSendTx(context.Background(), tx, "UpdateOperatorSocket", nil)
	if err != nil {
		t.logger.Error("Failed to estimate gas price and limit", "err", err)
		return err
	}
	return nil
}

// BuildConfirmBatchTxn builds a transaction to confirm a batch header and signature aggregation. The signature aggregation must satisfy the quorum thresholds
// specified in the batch header. If the signature aggregation does not satisfy the quorum thresholds, the transaction will fail.
// Note that this function returns a transaction without publishing it to the blockchain. The caller is responsible for publishing the transaction.
func (t *Writer) BuildConfirmBatchTxn(ctx context.Context, batchHeader *core.BatchHeader, quorums map[core.QuorumID]*core.QuorumResult, signatureAggregation *core.SignatureAggregation) (*types.Transaction, error) {
	quorumNumbers := quorumParamsToQuorumNumbers(quorums)
	nonSignerOperatorIds := make([][32]byte, len(signatureAggregation.NonSigners))
	for i := range signatureAggregation.NonSigners {
		// TODO: instead of recalculating the operator id, we should just pass it in from the caller
		nonSignerOperatorIds[i] = HashPubKeyG1(signatureAggregation.NonSigners[i])
	}

	checkSignaturesIndices, err := t.bindings.OpStateRetriever.GetCheckSignaturesIndices(
		&bind.CallOpts{
			Context: ctx,
		},
		t.bindings.RegCoordinatorAddr,
		uint32(batchHeader.ReferenceBlockNumber),
		quorumNumbers,
		nonSignerOperatorIds,
	)
	if err != nil {
		t.logger.Error("Failed to fetch checkSignaturesIndices", "err", err)
		return nil, err
	}

	nonSignerPubkeys := make([]eigendasrvmg.BN254G1Point, len(signatureAggregation.NonSigners))
	for i := range signatureAggregation.NonSigners {
		signature := signatureAggregation.NonSigners[i]
		nonSignerPubkeys[i] = pubKeyG1ToBN254G1Point(signature)
	}

	signedStakeForQuorums := serializeSignedStakeForQuorums(quorums)
	batchH := eigendasrvmg.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:       batchHeader.BatchRoot,
		QuorumNumbers:         quorumNumbers,
		SignedStakeForQuorums: signedStakeForQuorums,
		ReferenceBlockNumber:  uint32(batchHeader.ReferenceBlockNumber),
	}
	t.logger.Debug("batch header", "batchHeaderReferenceBlock", batchH.ReferenceBlockNumber, "batchHeaderRoot", gethcommon.Bytes2Hex(batchH.BlobHeadersRoot[:]), "quorumNumbers", gethcommon.Bytes2Hex(batchH.QuorumNumbers), "quorumThresholdPercentages", gethcommon.Bytes2Hex(batchH.SignedStakeForQuorums))

	sigma := signatureToBN254G1Point(signatureAggregation.AggSignature)

	apkG2 := pubKeyG2ToBN254G2Point(signatureAggregation.AggPubKey)

	quorumApks := make([]eigendasrvmg.BN254G1Point, len(signatureAggregation.QuorumAggPubKeys))
	for i := range signatureAggregation.QuorumAggPubKeys {
		quorumApks[i] = pubKeyG1ToBN254G1Point(signatureAggregation.QuorumAggPubKeys[i])
	}

	signatureChecker := eigendasrvmg.IBLSSignatureCheckerNonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: checkSignaturesIndices.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:             nonSignerPubkeys,
		QuorumApks:                   quorumApks,
		ApkG2:                        apkG2,
		Sigma:                        sigma,
		QuorumApkIndices:             checkSignaturesIndices.QuorumApkIndices,
		TotalStakeIndices:            checkSignaturesIndices.TotalStakeIndices,
		NonSignerStakeIndices:        checkSignaturesIndices.NonSignerStakeIndices,
	}
	sigChecker, err := json.Marshal(signatureChecker)
	if err == nil {
		t.logger.Debug("signature checker", "signatureChecker", string(sigChecker))
	}

	opts, err := t.ethClient.GetNoSendTransactOpts()
	if err != nil {
		t.logger.Error("Failed to generate transact opts", "err", err)
		return nil, err
	}
	return t.bindings.EigenDAServiceManager.ConfirmBatch(opts, batchH, signatureChecker)
}

// ConfirmBatch confirms a batch header and signature aggregation. The signature aggregation must satisfy the quorum thresholds
// specified in the batch header. If the signature aggregation does not satisfy the quorum thresholds, the transaction will fail.
func (t *Writer) ConfirmBatch(ctx context.Context, batchHeader *core.BatchHeader, quorums map[core.QuorumID]*core.QuorumResult, signatureAggregation *core.SignatureAggregation) (*types.Receipt, error) {
	tx, err := t.BuildConfirmBatchTxn(ctx, batchHeader, quorums, signatureAggregation)
	if err != nil {
		t.logger.Error("Failed to build a ConfirmBatch txn", "err", err)
		return nil, err
	}

	t.logger.Info("confirming batch onchain")
	receipt, err := t.ethClient.EstimateGasPriceAndLimitAndSendTx(ctx, tx, "ConfirmBatch", nil)
	if err != nil {
		t.logger.Error("Failed to estimate gas price and limit", "err", err)
		return nil, err
	}
	return receipt, nil
}
