package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"slices"

	"github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"

	avsdir "github.com/Layr-Labs/eigenda/contracts/bindings/AVSDirectory"
	blsapkreg "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	delegationmgr "github.com/Layr-Labs/eigenda/contracts/bindings/DelegationManager"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	indexreg "github.com/Layr-Labs/eigenda/contracts/bindings/IIndexRegistry"
	opstateretriever "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/RegistryCoordinator"
	stakereg "github.com/Layr-Labs/eigenda/contracts/bindings/StakeRegistry"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	maxNumberOfQuorums      = 192
	maxNumWorkerPoolThreads = 8
)

type Transactor struct {
	EthClient common.EthClient
	Logger    logging.Logger
	Bindings  *ContractBindings
}

var _ core.Transactor = (*Transactor)(nil)

type ContractBindings struct {
	RegCoordinatorAddr    gethcommon.Address
	ServiceManagerAddr    gethcommon.Address
	DelegationManager     *delegationmgr.ContractDelegationManager
	OpStateRetriever      *opstateretriever.ContractOperatorStateRetriever
	BLSApkRegistry        *blsapkreg.ContractBLSApkRegistry
	IndexRegistry         *indexreg.ContractIIndexRegistry
	RegistryCoordinator   *regcoordinator.ContractRegistryCoordinator
	StakeRegistry         *stakereg.ContractStakeRegistry
	EigenDAServiceManager *eigendasrvmg.ContractEigenDAServiceManager
	AVSDirectory          *avsdir.ContractAVSDirectory
}

type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

func NewTransactor(
	logger logging.Logger,
	client common.EthClient,
	blsOperatorStateRetrieverHexAddr string,
	eigenDAServiceManagerHexAddr string) (*Transactor, error) {

	e := &Transactor{
		EthClient: client,
		Logger:    logger,
	}

	blsOperatorStateRetrieverAddr := gethcommon.HexToAddress(blsOperatorStateRetrieverHexAddr)
	eigenDAServiceManagerAddr := gethcommon.HexToAddress(eigenDAServiceManagerHexAddr)
	err := e.updateContractBindings(blsOperatorStateRetrieverAddr, eigenDAServiceManagerAddr)

	return e, err
}

// GetRegisteredQuorumIdsForOperator returns the quorum ids that the operator is registered in with the given public key.
func (t *Transactor) GetRegisteredQuorumIdsForOperator(ctx context.Context, operator core.OperatorID) ([]core.QuorumID, error) {
	// TODO: Properly handle the case where the operator is not registered in any quorum. The current behavior of the smart contracts is to revert instead of returning an empty bitmap.
	//  We should probably change this.
	emptyBitmapErr := "execution reverted: BLSRegistryCoordinator.getCurrentQuorumBitmapByOperatorId: no quorum bitmap history for operatorId"
	quorumBitmap, err := t.Bindings.RegistryCoordinator.GetCurrentQuorumBitmap(&bind.CallOpts{
		Context: ctx,
	}, operator)
	if err != nil {
		if err.Error() == emptyBitmapErr {
			return []core.QuorumID{}, nil
		} else {
			t.Logger.Error("Failed to fetch current quorum bitmap", "err", err)
			return nil, err
		}
	}

	quorumIds := BitmapToQuorumIds(quorumBitmap)

	return quorumIds, nil
}

func (t *Transactor) getRegistrationParams(
	ctx context.Context,
	keypair *core.KeyPair,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
) (*regcoordinator.IBLSApkRegistryPubkeyRegistrationParams, *regcoordinator.ISignatureUtilsSignatureWithSaltAndExpiry, error) {

	operatorAddress := t.EthClient.GetAccountAddress()

	msgToSignG1_, err := t.Bindings.RegistryCoordinator.PubkeyRegistrationMessageHash(&bind.CallOpts{
		Context: ctx,
	}, operatorAddress)
	if err != nil {
		return nil, nil, err
	}

	msgToSignG1 := core.NewG1Point(msgToSignG1_.X, msgToSignG1_.Y)
	signature := keypair.SignHashedToCurveMessage(msgToSignG1)

	signedMessageHashParam := regcoordinator.BN254G1Point{
		X: signature.X.BigInt(big.NewInt(0)),
		Y: signature.Y.BigInt(big.NewInt(0)),
	}

	g1Point_ := pubKeyG1ToBN254G1Point(keypair.GetPubKeyG1())
	g1Point := regcoordinator.BN254G1Point{
		X: g1Point_.X,
		Y: g1Point_.Y,
	}
	g2Point_ := pubKeyG2ToBN254G2Point(keypair.GetPubKeyG2())
	g2Point := regcoordinator.BN254G2Point{
		X: g2Point_.X,
		Y: g2Point_.Y,
	}

	params := regcoordinator.IBLSApkRegistryPubkeyRegistrationParams{
		PubkeyRegistrationSignature: signedMessageHashParam,
		PubkeyG1:                    g1Point,
		PubkeyG2:                    g2Point,
	}

	// params to register operator in delegation manager's operator-avs mapping
	msgToSign, err := t.Bindings.AVSDirectory.CalculateOperatorAVSRegistrationDigestHash(
		&bind.CallOpts{
			Context: ctx,
		}, operatorAddress, t.Bindings.ServiceManagerAddr, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry)
	if err != nil {
		return nil, nil, err
	}
	operatorSignature, err := crypto.Sign(msgToSign[:], operatorEcdsaPrivateKey)
	if err != nil {
		return nil, nil, err
	}
	// this is annoying, and not sure why its needed, but seems like some historical baggage
	// see https://github.com/ethereum/go-ethereum/issues/28757#issuecomment-1874525854
	// and https://twitter.com/pcaversaccio/status/1671488928262529031
	operatorSignature[64] += 27
	operatorSignatureWithSaltAndExpiry := regcoordinator.ISignatureUtilsSignatureWithSaltAndExpiry{
		Signature: operatorSignature,
		Salt:      operatorToAvsRegistrationSigSalt,
		Expiry:    operatorToAvsRegistrationSigExpiry,
	}

	return &params, &operatorSignatureWithSaltAndExpiry, nil

}

// RegisterOperator registers a new operator with the given public key and socket with the provided quorum ids.
// If the operator is already registered with a given quorum id, the transaction will fail (noop) and an error
// will be returned.
func (t *Transactor) RegisterOperator(
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
func (t *Transactor) RegisterOperatorWithChurn(
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
		t.Logger.Error("Failed to get registration params", "err", err)
		return err
	}

	quorumNumbers := quorumIDsToQuorumNumbers(quorumIds)

	operatorsToChurn := make([]regcoordinator.IRegistryCoordinatorOperatorKickParam, len(churnReply.OperatorsToChurn))
	for i := range churnReply.OperatorsToChurn {
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

// DeregisterOperator deregisters an operator with the given public key from the all the quorums that it is
// registered with at the supplied block number. To fully deregister an operator, this function should be called
// with the current block number.
func (t *Transactor) DeregisterOperator(ctx context.Context, pubkeyG1 *core.G1Point, blockNumber uint32) error {
	operatorId := HashPubKeyG1(pubkeyG1)
	quorumBitmap, opStates, err := t.Bindings.OpStateRetriever.GetOperatorState0(&bind.CallOpts{
		Context: ctx,
	}, t.Bindings.RegCoordinatorAddr, operatorId, blockNumber)
	if err != nil {
		t.Logger.Error("Failed to fetch operator state", "err", err)
		return err
	}

	operatorIdsToSwap := make([][32]byte, len(opStates))
	for i := range opStates {
		quorum := opStates[i]
		operatorIdsToSwap[i] = quorum[len(opStates[i])-1].OperatorId
	}

	quorumNumbers := bitmapToBytesArray(quorumBitmap)

	opts, err := t.EthClient.GetNoSendTransactOpts()
	if err != nil {
		t.Logger.Error("Failed to generate transact opts", "err", err)
		return err
	}
	tx, err := t.Bindings.RegistryCoordinator.DeregisterOperator(
		opts,
		quorumNumbers,
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
func (t *Transactor) UpdateOperatorSocket(ctx context.Context, socket string) error {
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

// GetOperatorStakes returns the stakes of all operators within the quorums that the operator represented by operatorId
// is registered with. The returned stakes are for the block number supplied. The indices of the operators within each quorum
// are also returned.
func (t *Transactor) GetOperatorStakes(ctx context.Context, operator core.OperatorID, blockNumber uint32) (core.OperatorStakes, []core.QuorumID, error) {
	quorumBitmap, state_, err := t.Bindings.OpStateRetriever.GetOperatorState0(&bind.CallOpts{
		Context: ctx,
	}, t.Bindings.RegCoordinatorAddr, operator, blockNumber)
	if err != nil {
		t.Logger.Error("Failed to fetch operator state", "err", err)
		return nil, nil, err
	}

	state := make(core.OperatorStakes, len(state_))
	for i := range state_ {
		quorumID := core.QuorumID(i)
		state[quorumID] = make(map[core.OperatorIndex]core.OperatorStake, len(state_[i]))
		for j, op := range state_[i] {
			operatorIndex := core.OperatorIndex(j)
			state[quorumID][operatorIndex] = core.OperatorStake{
				Stake:      op.Stake,
				OperatorID: op.OperatorId,
			}
		}
	}

	quorumIds := BitmapToQuorumIds(quorumBitmap)

	return state, quorumIds, nil
}

func (t *Transactor) GetBlockStaleMeasure(ctx context.Context) (uint32, error) {
	blockStaleMeasure, err := t.Bindings.EigenDAServiceManager.BLOCKSTALEMEASURE(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		t.Logger.Error("Failed to fetch BLOCK_STALE_MEASURE", err)
		return *new(uint32), err
	}
	return blockStaleMeasure, nil
}

func (t *Transactor) GetStoreDurationBlocks(ctx context.Context) (uint32, error) {
	blockStaleMeasure, err := t.Bindings.EigenDAServiceManager.STOREDURATIONBLOCKS(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		t.Logger.Error("Failed to fetch STORE_DURATION_BLOCKS", err)
		return *new(uint32), err
	}
	return blockStaleMeasure, nil
}

// GetOperatorStakesForQuorums returns the stakes of all operators within the supplied quorums. The returned stakes are for the block number supplied.
// The indices of the operators within each quorum are also returned.
func (t *Transactor) GetOperatorStakesForQuorums(ctx context.Context, quorums []core.QuorumID, blockNumber uint32) (core.OperatorStakes, error) {
	quorumBytes := make([]byte, len(quorums))
	for ind, quorum := range quorums {
		quorumBytes[ind] = byte(uint8(quorum))
	}

	// state_ is a [][]*opstateretriever.OperatorStake with the same length and order as quorumBytes, and then indexed by operator index
	state_, err := t.Bindings.OpStateRetriever.GetOperatorState(&bind.CallOpts{
		Context: ctx,
	}, t.Bindings.RegCoordinatorAddr, quorumBytes, blockNumber)
	if err != nil {
		t.Logger.Error("Failed to fetch operator state", err)
		return nil, err
	}

	state := make(core.OperatorStakes, len(state_))
	for i := range state_ {
		quorumID := quorums[i]
		state[quorumID] = make(map[core.OperatorIndex]core.OperatorStake, len(state_[i]))
		for j, op := range state_[i] {
			operatorIndex := core.OperatorIndex(j)
			state[quorumID][operatorIndex] = core.OperatorStake{
				Stake:      op.Stake,
				OperatorID: op.OperatorId,
			}
		}
	}

	return state, nil
}

// BuildConfirmBatchTxn builds a transaction to confirm a batch header and signature aggregation. The signature aggregation must satisfy the quorum thresholds
// specified in the batch header. If the signature aggregation does not satisfy the quorum thresholds, the transaction will fail.
// Note that this function returns a transaction without publishing it to the blockchain. The caller is responsible for publishing the transaction.
func (t *Transactor) BuildConfirmBatchTxn(ctx context.Context, batchHeader *core.BatchHeader, quorums map[core.QuorumID]*core.QuorumResult, signatureAggregation *core.SignatureAggregation) (*types.Transaction, error) {
	quorumNumbers := quorumParamsToQuorumNumbers(quorums)
	nonSignerOperatorIds := make([][32]byte, len(signatureAggregation.NonSigners))
	for i := range signatureAggregation.NonSigners {
		// TODO: instead of recalculating the operator id, we should just pass it in from the caller
		nonSignerOperatorIds[i] = HashPubKeyG1(signatureAggregation.NonSigners[i])
	}
	sigAgg, err := json.Marshal(signatureAggregation)
	if err == nil {
		t.Logger.Debug("[BuildConfirmBatchTxn]", "signatureAggregation", string(sigAgg))
	}

	t.Logger.Debug("[GetCheckSignaturesIndices]", "regCoordinatorAddr", t.Bindings.RegCoordinatorAddr.Hex(), "refBlockNumber", batchHeader.ReferenceBlockNumber, "quorumNumbers", gethcommon.Bytes2Hex(quorumNumbers))
	for _, ns := range nonSignerOperatorIds {
		t.Logger.Debug("[GetCheckSignaturesIndices]", "nonSignerOperatorId", gethcommon.Bytes2Hex(ns[:]))
	}
	checkSignaturesIndices, err := t.Bindings.OpStateRetriever.GetCheckSignaturesIndices(
		&bind.CallOpts{
			Context: ctx,
		},
		t.Bindings.RegCoordinatorAddr,
		uint32(batchHeader.ReferenceBlockNumber),
		quorumNumbers,
		nonSignerOperatorIds,
	)
	if err != nil {
		t.Logger.Error("Failed to fetch checkSignaturesIndices", "err", err)
		return nil, err
	}

	nonSignerPubkeys := make([]eigendasrvmg.BN254G1Point, len(signatureAggregation.NonSigners))
	for i := range signatureAggregation.NonSigners {
		signature := signatureAggregation.NonSigners[i]
		nonSignerPubkeys[i] = pubKeyG1ToBN254G1Point(signature)
	}

	quorumThresholdPercentages := quorumParamsToThresholdPercentages(quorums)
	batchH := eigendasrvmg.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:            batchHeader.BatchRoot,
		QuorumNumbers:              quorumNumbers,
		QuorumThresholdPercentages: quorumThresholdPercentages,
		ReferenceBlockNumber:       uint32(batchHeader.ReferenceBlockNumber),
	}
	t.Logger.Debug("[ConfirmBatch] batch header", "batchHeaderReferenceBlock", batchH.ReferenceBlockNumber, "batchHeaderRoot", gethcommon.Bytes2Hex(batchH.BlobHeadersRoot[:]), "quorumNumbers", gethcommon.Bytes2Hex(batchH.QuorumNumbers), "quorumThresholdPercentages", gethcommon.Bytes2Hex(batchH.QuorumThresholdPercentages))

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
		t.Logger.Debug("[ConfirmBatch] signature checker", "signatureChecker", string(sigChecker))
	}

	opts, err := t.EthClient.GetNoSendTransactOpts()
	if err != nil {
		t.Logger.Error("Failed to generate transact opts", "err", err)
		return nil, err
	}
	return t.Bindings.EigenDAServiceManager.ConfirmBatch(opts, batchH, signatureChecker)
}

// ConfirmBatch confirms a batch header and signature aggregation. The signature aggregation must satisfy the quorum thresholds
// specified in the batch header. If the signature aggregation does not satisfy the quorum thresholds, the transaction will fail.
func (t *Transactor) ConfirmBatch(ctx context.Context, batchHeader *core.BatchHeader, quorums map[core.QuorumID]*core.QuorumResult, signatureAggregation *core.SignatureAggregation) (*types.Receipt, error) {
	tx, err := t.BuildConfirmBatchTxn(ctx, batchHeader, quorums, signatureAggregation)
	if err != nil {
		t.Logger.Error("Failed to build a ConfirmBatch txn", "err", err)
		return nil, err
	}

	t.Logger.Info("confirming batch onchain")
	receipt, err := t.EthClient.EstimateGasPriceAndLimitAndSendTx(ctx, tx, "ConfirmBatch", nil)
	if err != nil {
		t.Logger.Error("Failed to estimate gas price and limit", "err", err)
		return nil, err
	}
	return receipt, nil
}

func (t *Transactor) StakeRegistry(ctx context.Context) (gethcommon.Address, error) {
	return t.Bindings.RegistryCoordinator.StakeRegistry(&bind.CallOpts{
		Context: ctx,
	})
}

func (t *Transactor) OperatorIDToAddress(ctx context.Context, operatorId core.OperatorID) (gethcommon.Address, error) {
	return t.Bindings.BLSApkRegistry.PubkeyHashToOperator(&bind.CallOpts{
		Context: ctx,
	}, operatorId)
}

func (t *Transactor) BatchOperatorIDToAddress(ctx context.Context, operatorIds []core.OperatorID) ([]gethcommon.Address, error) {
	type AddressOrError struct {
		address gethcommon.Address
		index   int
		err     error
	}
	resultChan := make(chan AddressOrError, len(operatorIds))
	pool := workerpool.New(maxNumWorkerPoolThreads)
	for i, operatorId := range operatorIds {
		idx := i
		op := operatorId
		pool.Submit(func() {
			addr, err := t.Bindings.BLSApkRegistry.PubkeyHashToOperator(&bind.CallOpts{
				Context: ctx,
			}, op)
			resultChan <- AddressOrError{address: addr, index: idx, err: err}
		})
	}
	pool.StopWait()
	close(resultChan)

	addresses := make([]gethcommon.Address, len(operatorIds))
	for result := range resultChan {
		if result.err != nil {
			return nil, result.err
		}
		addresses[result.index] = result.address
	}
	return addresses, nil
}

func (t *Transactor) GetCurrentQuorumBitmapByOperatorId(ctx context.Context, operatorId core.OperatorID) (*big.Int, error) {
	return t.Bindings.RegistryCoordinator.GetCurrentQuorumBitmap(&bind.CallOpts{
		Context: ctx,
	}, operatorId)
}

func (t *Transactor) GetQuorumBitmapForOperatorsAtBlockNumber(ctx context.Context, operatorIds []core.OperatorID, blockNumber uint32) ([]*big.Int, error) {
	// Get indices in batch (1 RPC).
	byteOperatorIds := make([][32]byte, len(operatorIds))
	for i := range operatorIds {
		byteOperatorIds[i] = [32]byte(operatorIds[i])
	}
	indices, err := t.Bindings.RegistryCoordinator.GetQuorumBitmapIndicesAtBlockNumber(&bind.CallOpts{
		Context: ctx,
	}, blockNumber, byteOperatorIds)
	if err != nil {
		return nil, err
	}

	// Get bitmaps in N RPCs, but in parallel.
	type BitmapOrError struct {
		bitmap *big.Int
		index  uint32
		err    error
	}
	resultChan := make(chan BitmapOrError, len(indices))
	pool := workerpool.New(maxNumWorkerPoolThreads)
	for i, index := range indices {
		idx := index
		op := operatorIds[i]
		pool.Submit(func() {
			bm, err := t.Bindings.RegistryCoordinator.GetQuorumBitmapAtBlockNumberByIndex(&bind.CallOpts{
				Context: ctx,
			}, op, blockNumber, big.NewInt(int64(idx)))
			resultChan <- BitmapOrError{bitmap: bm, index: idx, err: err}
		})
	}
	pool.StopWait()
	close(resultChan)

	bitmaps := make([]*big.Int, len(indices))
	for result := range resultChan {
		if result.err != nil {
			return nil, result.err
		}
		bitmaps[result.index] = result.bitmap
	}
	return bitmaps, nil
}

func (t *Transactor) GetOperatorSetParams(ctx context.Context, quorumID core.QuorumID) (*core.OperatorSetParam, error) {

	operatorSetParams, err := t.Bindings.RegistryCoordinator.GetOperatorSetParams(&bind.CallOpts{
		Context: ctx,
	}, quorumID)
	if err != nil {
		t.Logger.Error("Failed to fetch operator set params", "err", err)
		return nil, err
	}

	return &core.OperatorSetParam{
		MaxOperatorCount:         operatorSetParams.MaxOperatorCount,
		ChurnBIPsOfOperatorStake: operatorSetParams.KickBIPsOfOperatorStake,
		ChurnBIPsOfTotalStake:    operatorSetParams.KickBIPsOfTotalStake,
	}, nil
}

// Returns the number of registered operators for the quorum.
func (t *Transactor) GetNumberOfRegisteredOperatorForQuorum(ctx context.Context, quorumID core.QuorumID) (uint32, error) {
	return t.Bindings.IndexRegistry.TotalOperatorsForQuorum(&bind.CallOpts{
		Context: ctx,
	}, quorumID)
}

func (t *Transactor) WeightOfOperatorForQuorum(ctx context.Context, quorumID core.QuorumID, operator gethcommon.Address) (*big.Int, error) {
	return t.Bindings.StakeRegistry.WeightOfOperatorForQuorum(&bind.CallOpts{
		Context: ctx,
	}, quorumID, operator)
}

func (t *Transactor) CalculateOperatorChurnApprovalDigestHash(
	ctx context.Context,
	operatorAddress gethcommon.Address,
	operatorId core.OperatorID,
	operatorsToChurn []core.OperatorToChurn,
	salt [32]byte,
	expiry *big.Int,
) ([32]byte, error) {
	opKickParams := make([]regcoordinator.IRegistryCoordinatorOperatorKickParam, len(operatorsToChurn))
	for i := range operatorsToChurn {

		opKickParams[i] = regcoordinator.IRegistryCoordinatorOperatorKickParam{
			QuorumNumber: operatorsToChurn[i].QuorumId,
			Operator:     operatorsToChurn[i].Operator,
		}
	}
	return t.Bindings.RegistryCoordinator.CalculateOperatorChurnApprovalDigestHash(&bind.CallOpts{
		Context: ctx,
	}, operatorAddress, operatorId, opKickParams, salt, expiry)
}

func (t *Transactor) GetCurrentBlockNumber(ctx context.Context) (uint32, error) {
	bn, err := t.EthClient.BlockNumber(ctx)
	return uint32(bn), err
}

func (t *Transactor) GetQuorumCount(ctx context.Context, blockNumber uint32) (uint8, error) {
	return t.Bindings.RegistryCoordinator.QuorumCount(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
}

func (t *Transactor) updateContractBindings(blsOperatorStateRetrieverAddr, eigenDAServiceManagerAddr gethcommon.Address) error {

	contractEigenDAServiceManager, err := eigendasrvmg.NewContractEigenDAServiceManager(eigenDAServiceManagerAddr, t.EthClient)
	if err != nil {
		t.Logger.Error("Failed to fetch IEigenDAServiceManager contract", "err", err)
		return err
	}

	delegationManagerAddr, err := contractEigenDAServiceManager.Delegation(&bind.CallOpts{})
	if err != nil {
		t.Logger.Error("Failed to fetch DelegationManager address", "err", err)
		return err
	}

	avsDirectoryAddr, err := contractEigenDAServiceManager.AvsDirectory(&bind.CallOpts{})
	if err != nil {
		t.Logger.Error("Failed to fetch AVSDirectory address", "err", err)
		return err
	}

	contractAVSDirectory, err := avsdir.NewContractAVSDirectory(avsDirectoryAddr, t.EthClient)
	if err != nil {
		t.Logger.Error("Failed to fetch AVSDirectory contract", "err", err)
		return err
	}

	contractDelegationManager, err := delegationmgr.NewContractDelegationManager(delegationManagerAddr, t.EthClient)
	if err != nil {
		t.Logger.Error("Failed to fetch DelegationManager contract", "err", err)
		return err
	}

	registryCoordinatorAddr, err := contractEigenDAServiceManager.RegistryCoordinator(&bind.CallOpts{})
	if err != nil {
		t.Logger.Error("Failed to fetch RegistryCoordinator address", "err", err)
		return err
	}

	contractIRegistryCoordinator, err := regcoordinator.NewContractRegistryCoordinator(registryCoordinatorAddr, t.EthClient)
	if err != nil {
		t.Logger.Error("Failed to fetch IBLSRegistryCoordinatorWithIndices contract", "err", err)
		return err
	}

	contractBLSOpStateRetr, err := opstateretriever.NewContractOperatorStateRetriever(blsOperatorStateRetrieverAddr, t.EthClient)
	if err != nil {
		t.Logger.Error("Failed to fetch BLSOperatorStateRetriever contract", "err", err)
		return err
	}

	blsPubkeyRegistryAddr, err := contractIRegistryCoordinator.BlsApkRegistry(&bind.CallOpts{})
	if err != nil {
		t.Logger.Error("Failed to fetch BlsPubkeyRegistry address", "err", err)
		return err
	}

	t.Logger.Debug("Addresses", "blsOperatorStateRetrieverAddr", blsOperatorStateRetrieverAddr.Hex(), "eigenDAServiceManagerAddr", eigenDAServiceManagerAddr.Hex(), "registryCoordinatorAddr", registryCoordinatorAddr.Hex(), "blsPubkeyRegistryAddr", blsPubkeyRegistryAddr.Hex())

	contractBLSPubkeyReg, err := blsapkreg.NewContractBLSApkRegistry(blsPubkeyRegistryAddr, t.EthClient)
	if err != nil {
		t.Logger.Error("Failed to fetch IBLSApkRegistry contract", "err", err)
		return err
	}

	indexRegistryAddr, err := contractIRegistryCoordinator.IndexRegistry(&bind.CallOpts{})
	if err != nil {
		t.Logger.Error("Failed to fetch IndexRegistry address", "err", err)
		return err
	}

	contractIIndexReg, err := indexreg.NewContractIIndexRegistry(indexRegistryAddr, t.EthClient)
	if err != nil {
		t.Logger.Error("Failed to fetch IIndexRegistry contract", "err", err)
		return err
	}

	stakeRegistryAddr, err := contractIRegistryCoordinator.StakeRegistry(&bind.CallOpts{})
	if err != nil {
		t.Logger.Error("Failed to fetch StakeRegistry address", "err", err)
		return err
	}

	contractStakeRegistry, err := stakereg.NewContractStakeRegistry(stakeRegistryAddr, t.EthClient)
	if err != nil {
		t.Logger.Error("Failed to fetch StakeRegistry contract", "err", err)
		return err
	}

	t.Bindings = &ContractBindings{
		ServiceManagerAddr:    eigenDAServiceManagerAddr,
		RegCoordinatorAddr:    registryCoordinatorAddr,
		AVSDirectory:          contractAVSDirectory,
		OpStateRetriever:      contractBLSOpStateRetr,
		BLSApkRegistry:        contractBLSPubkeyReg,
		IndexRegistry:         contractIIndexReg,
		RegistryCoordinator:   contractIRegistryCoordinator,
		StakeRegistry:         contractStakeRegistry,
		EigenDAServiceManager: contractEigenDAServiceManager,
		DelegationManager:     contractDelegationManager,
	}
	return nil
}

func signatureToBN254G1Point(s *core.Signature) eigendasrvmg.BN254G1Point {
	return eigendasrvmg.BN254G1Point{
		X: s.X.BigInt(new(big.Int)),
		Y: s.Y.BigInt(new(big.Int)),
	}
}

func pubKeyG1ToBN254G1Point(p *core.G1Point) eigendasrvmg.BN254G1Point {
	return eigendasrvmg.BN254G1Point{
		X: p.X.BigInt(new(big.Int)),
		Y: p.Y.BigInt(new(big.Int)),
	}
}

func pubKeyG2ToBN254G2Point(p *core.G2Point) eigendasrvmg.BN254G2Point {
	return eigendasrvmg.BN254G2Point{
		X: [2]*big.Int{p.X.A1.BigInt(new(big.Int)), p.X.A0.BigInt(new(big.Int))},
		Y: [2]*big.Int{p.Y.A1.BigInt(new(big.Int)), p.Y.A0.BigInt(new(big.Int))},
	}
}

func quorumIDsToQuorumNumbers(quorumIds []core.QuorumID) []byte {
	quorumNumbers := make([]byte, len(quorumIds))
	for i, quorumId := range quorumIds {
		quorumNumbers[i] = byte(quorumId)
	}
	return quorumNumbers
}

func quorumParamsToQuorumNumbers(quorumParams map[core.QuorumID]*core.QuorumResult) []byte {
	quorumNumbers := make([]byte, len(quorumParams))
	quorums := make([]uint8, len(quorumParams))
	i := 0
	for k := range quorumParams {
		quorums[i] = k
		i++
	}
	slices.Sort(quorums)
	i = 0
	for _, quorum := range quorums {
		qp := quorumParams[quorum]
		quorumNumbers[i] = byte(qp.QuorumID)
		i++
	}
	return quorumNumbers
}

func quorumParamsToThresholdPercentages(quorumParams map[core.QuorumID]*core.QuorumResult) []byte {
	thresholdPercentages := make([]byte, len(quorumParams))
	quorums := make([]uint8, len(quorumParams))
	i := 0
	for k := range quorumParams {
		quorums[i] = k
		i++
	}
	slices.Sort(quorums)
	i = 0
	for _, quorum := range quorums {
		qp := quorumParams[quorum]
		thresholdPercentages[i] = byte(qp.PercentSigned)
		i++
	}
	return thresholdPercentages
}

func HashPubKeyG1(pk *core.G1Point) [32]byte {
	gp := pubKeyG1ToBN254G1Point(pk)
	xBytes := make([]byte, 32)
	yBytes := make([]byte, 32)
	gp.X.FillBytes(xBytes)
	gp.Y.FillBytes(yBytes)
	return crypto.Keccak256Hash(append(xBytes, yBytes...))
}

func BitmapToQuorumIds(bitmap *big.Int) []core.QuorumID {
	// loop through each index in the bitmap to construct the array

	quorumIds := make([]core.QuorumID, 0, maxNumberOfQuorums)
	for i := 0; i < maxNumberOfQuorums; i++ {
		if bitmap.Bit(i) == 1 {
			quorumIds = append(quorumIds, core.QuorumID(i))
		}
	}
	return quorumIds
}

func bitmapToBytesArray(bitmap *big.Int) []byte {
	// initialize an empty uint64 to be used as a bitmask inside the loop
	var (
		bytesArray []byte
	)
	// loop through each index in the bitmap to construct the array
	for i := 0; i < maxNumberOfQuorums; i++ {
		// check if the i-th bit is flipped in the bitmap
		if bitmap.Bit(i) == 1 {
			// if the i-th bit is flipped, then add a byte encoding the value 'i' to the `bytesArray`
			bytesArray = append(bytesArray, byte(uint8(i)))
		}
	}
	return bytesArray
}
