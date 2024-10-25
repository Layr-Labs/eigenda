package eth

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	avsdir "github.com/Layr-Labs/eigenda/contracts/bindings/AVSDirectory"
	blsapkreg "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	delegationmgr "github.com/Layr-Labs/eigenda/contracts/bindings/DelegationManager"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	ejectionmg "github.com/Layr-Labs/eigenda/contracts/bindings/EjectionManager"
	indexreg "github.com/Layr-Labs/eigenda/contracts/bindings/IIndexRegistry"
	opstateretriever "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/RegistryCoordinator"
	stakereg "github.com/Layr-Labs/eigenda/contracts/bindings/StakeRegistry"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pingcap/errors"
)

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
	EjectionManager       *ejectionmg.ContractEjectionManager
	AVSDirectory          *avsdir.ContractAVSDirectory
}

type Reader struct {
	ethClient common.EthClient
	logger    logging.Logger
	bindings  *ContractBindings
}

var _ core.Reader = (*Reader)(nil)

func NewReader(
	logger logging.Logger,
	client common.EthClient,
	blsOperatorStateRetrieverHexAddr string,
	eigenDAServiceManagerHexAddr string) (*Reader, error) {

	e := &Reader{
		ethClient: client,
		logger:    logger.With("component", "Reader"),
	}

	blsOperatorStateRetrieverAddr := gethcommon.HexToAddress(blsOperatorStateRetrieverHexAddr)
	eigenDAServiceManagerAddr := gethcommon.HexToAddress(eigenDAServiceManagerHexAddr)
	err := e.updateContractBindings(blsOperatorStateRetrieverAddr, eigenDAServiceManagerAddr)

	return e, err
}

func (t *Reader) updateContractBindings(blsOperatorStateRetrieverAddr, eigenDAServiceManagerAddr gethcommon.Address) error {

	contractEigenDAServiceManager, err := eigendasrvmg.NewContractEigenDAServiceManager(eigenDAServiceManagerAddr, t.ethClient)
	if err != nil {
		t.logger.Error("Failed to fetch IEigenDAServiceManager contract", "err", err)
		return err
	}

	delegationManagerAddr, err := contractEigenDAServiceManager.Delegation(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch DelegationManager address", "err", err)
		return err
	}

	avsDirectoryAddr, err := contractEigenDAServiceManager.AvsDirectory(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch AVSDirectory address", "err", err)
		return err
	}

	contractAVSDirectory, err := avsdir.NewContractAVSDirectory(avsDirectoryAddr, t.ethClient)
	if err != nil {
		t.logger.Error("Failed to fetch AVSDirectory contract", "err", err)
		return err
	}

	contractDelegationManager, err := delegationmgr.NewContractDelegationManager(delegationManagerAddr, t.ethClient)
	if err != nil {
		t.logger.Error("Failed to fetch DelegationManager contract", "err", err)
		return err
	}

	registryCoordinatorAddr, err := contractEigenDAServiceManager.RegistryCoordinator(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch RegistryCoordinator address", "err", err)
		return err
	}

	contractIRegistryCoordinator, err := regcoordinator.NewContractRegistryCoordinator(registryCoordinatorAddr, t.ethClient)
	if err != nil {
		t.logger.Error("Failed to fetch IBLSRegistryCoordinatorWithIndices contract", "err", err)
		return err
	}

	contractEjectionManagerAddr, err := contractIRegistryCoordinator.Ejector(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch EjectionManager address", "err", err)
		return err
	}
	contractEjectionManager, err := ejectionmg.NewContractEjectionManager(contractEjectionManagerAddr, t.ethClient)
	if err != nil {
		t.logger.Error("Failed to fetch EjectionManager contract", "err", err)
		return err
	}

	contractBLSOpStateRetr, err := opstateretriever.NewContractOperatorStateRetriever(blsOperatorStateRetrieverAddr, t.ethClient)
	if err != nil {
		t.logger.Error("Failed to fetch BLSOperatorStateRetriever contract", "err", err)
		return err
	}

	blsPubkeyRegistryAddr, err := contractIRegistryCoordinator.BlsApkRegistry(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch BlsPubkeyRegistry address", "err", err)
		return err
	}

	t.logger.Debug("Addresses", "blsOperatorStateRetrieverAddr", blsOperatorStateRetrieverAddr.Hex(), "eigenDAServiceManagerAddr", eigenDAServiceManagerAddr.Hex(), "registryCoordinatorAddr", registryCoordinatorAddr.Hex(), "blsPubkeyRegistryAddr", blsPubkeyRegistryAddr.Hex())

	contractBLSPubkeyReg, err := blsapkreg.NewContractBLSApkRegistry(blsPubkeyRegistryAddr, t.ethClient)
	if err != nil {
		t.logger.Error("Failed to fetch IBLSApkRegistry contract", "err", err)
		return err
	}

	indexRegistryAddr, err := contractIRegistryCoordinator.IndexRegistry(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch IndexRegistry address", "err", err)
		return err
	}

	contractIIndexReg, err := indexreg.NewContractIIndexRegistry(indexRegistryAddr, t.ethClient)
	if err != nil {
		t.logger.Error("Failed to fetch IIndexRegistry contract", "err", err)
		return err
	}

	stakeRegistryAddr, err := contractIRegistryCoordinator.StakeRegistry(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch StakeRegistry address", "err", err)
		return err
	}

	contractStakeRegistry, err := stakereg.NewContractStakeRegistry(stakeRegistryAddr, t.ethClient)
	if err != nil {
		t.logger.Error("Failed to fetch StakeRegistry contract", "err", err)
		return err
	}

	t.bindings = &ContractBindings{
		ServiceManagerAddr:    eigenDAServiceManagerAddr,
		RegCoordinatorAddr:    registryCoordinatorAddr,
		AVSDirectory:          contractAVSDirectory,
		OpStateRetriever:      contractBLSOpStateRetr,
		BLSApkRegistry:        contractBLSPubkeyReg,
		IndexRegistry:         contractIIndexReg,
		RegistryCoordinator:   contractIRegistryCoordinator,
		EjectionManager:       contractEjectionManager,
		StakeRegistry:         contractStakeRegistry,
		EigenDAServiceManager: contractEigenDAServiceManager,
		DelegationManager:     contractDelegationManager,
	}
	return nil
}

// GetRegisteredQuorumIdsForOperator returns the quorum ids that the operator is registered in with the given public key.
func (t *Reader) GetRegisteredQuorumIdsForOperator(ctx context.Context, operator core.OperatorID) ([]core.QuorumID, error) {
	// TODO: Properly handle the case where the operator is not registered in any quorum. The current behavior of the smart contracts is to revert instead of returning an empty bitmap.
	//  We should probably change this.
	emptyBitmapErr := "execution reverted: BLSRegistryCoordinator.getCurrentQuorumBitmapByOperatorId: no quorum bitmap history for operatorId"
	quorumBitmap, err := t.bindings.RegistryCoordinator.GetCurrentQuorumBitmap(&bind.CallOpts{
		Context: ctx,
	}, operator)
	if err != nil {
		if err.Error() == emptyBitmapErr {
			return []core.QuorumID{}, nil
		} else {
			t.logger.Error("Failed to fetch current quorum bitmap", "err", err)
			return nil, err
		}
	}

	quorumIds := BitmapToQuorumIds(quorumBitmap)

	return quorumIds, nil
}

func (t *Reader) getRegistrationParams(
	ctx context.Context,
	keypair *core.KeyPair,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
) (*regcoordinator.IBLSApkRegistryPubkeyRegistrationParams, *regcoordinator.ISignatureUtilsSignatureWithSaltAndExpiry, error) {

	operatorAddress := t.ethClient.GetAccountAddress()

	msgToSignG1_, err := t.bindings.RegistryCoordinator.PubkeyRegistrationMessageHash(&bind.CallOpts{
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
	msgToSign, err := t.bindings.AVSDirectory.CalculateOperatorAVSRegistrationDigestHash(
		&bind.CallOpts{
			Context: ctx,
		}, operatorAddress, t.bindings.ServiceManagerAddr, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry)
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

func (t *Reader) BuildEjectOperatorsTxn(ctx context.Context, operatorsByQuorum [][]core.OperatorID) (*types.Transaction, error) {
	byteIdsByQuorum := make([][][32]byte, len(operatorsByQuorum))
	for i, ids := range operatorsByQuorum {
		for _, id := range ids {
			byteIdsByQuorum[i] = append(byteIdsByQuorum[i], [32]byte(id))
		}
	}
	opts, err := t.ethClient.GetNoSendTransactOpts()
	if err != nil {
		t.logger.Error("Failed to generate transact opts", "err", err)
		return nil, err
	}
	return t.bindings.EjectionManager.EjectOperators(opts, byteIdsByQuorum)
}

// GetOperatorStakes returns the stakes of all operators within the quorums that the operator represented by operatorId
// is registered with. The returned stakes are for the block number supplied. The indices of the operators within each quorum
// are also returned.
func (t *Reader) GetOperatorStakes(ctx context.Context, operator core.OperatorID, blockNumber uint32) (core.OperatorStakes, []core.QuorumID, error) {
	quorumBitmap, state_, err := t.bindings.OpStateRetriever.GetOperatorState0(&bind.CallOpts{
		Context: ctx,
	}, t.bindings.RegCoordinatorAddr, operator, blockNumber)
	if err != nil {
		t.logger.Error("Failed to fetch operator state", "err", err, "blockNumber", blockNumber, "operatorID", operator.Hex())
		return nil, nil, err
	}

	// BitmapToQuorumIds returns an ordered list of quorums in ascending order, which is the same order as the state_ returned by the contract
	quorumIds := BitmapToQuorumIds(quorumBitmap)

	state := make(core.OperatorStakes, len(state_))
	for i := range state_ {
		quorumID := quorumIds[i]
		state[quorumID] = make(map[core.OperatorIndex]core.OperatorStake, len(state_[i]))
		for j, op := range state_[i] {
			operatorIndex := core.OperatorIndex(j)
			state[quorumID][operatorIndex] = core.OperatorStake{
				Stake:      op.Stake,
				OperatorID: op.OperatorId,
			}
		}
	}

	return state, quorumIds, nil
}

func (t *Reader) GetBlockStaleMeasure(ctx context.Context) (uint32, error) {
	blockStaleMeasure, err := t.bindings.EigenDAServiceManager.BLOCKSTALEMEASURE(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		t.logger.Error("Failed to fetch BLOCK_STALE_MEASURE", err)
		return *new(uint32), err
	}
	return blockStaleMeasure, nil
}

func (t *Reader) GetStoreDurationBlocks(ctx context.Context) (uint32, error) {
	blockStaleMeasure, err := t.bindings.EigenDAServiceManager.STOREDURATIONBLOCKS(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		t.logger.Error("Failed to fetch STORE_DURATION_BLOCKS", err)
		return *new(uint32), err
	}
	return blockStaleMeasure, nil
}

// GetOperatorStakesForQuorums returns the stakes of all operators within the supplied quorums. The returned stakes are for the block number supplied.
// The indices of the operators within each quorum are also returned.
func (t *Reader) GetOperatorStakesForQuorums(ctx context.Context, quorums []core.QuorumID, blockNumber uint32) (core.OperatorStakes, error) {
	quorumBytes := make([]byte, len(quorums))
	for ind, quorum := range quorums {
		quorumBytes[ind] = byte(uint8(quorum))
	}

	// state_ is a [][]*opstateretriever.OperatorStake with the same length and order as quorumBytes, and then indexed by operator index
	state_, err := t.bindings.OpStateRetriever.GetOperatorState(&bind.CallOpts{
		Context: ctx,
	}, t.bindings.RegCoordinatorAddr, quorumBytes, blockNumber)
	if err != nil {
		t.logger.Error("Failed to fetch operator state", "err", err)
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

func (t *Reader) StakeRegistry(ctx context.Context) (gethcommon.Address, error) {
	return t.bindings.RegistryCoordinator.StakeRegistry(&bind.CallOpts{
		Context: ctx,
	})
}

func (t *Reader) OperatorIDToAddress(ctx context.Context, operatorId core.OperatorID) (gethcommon.Address, error) {
	return t.bindings.BLSApkRegistry.PubkeyHashToOperator(&bind.CallOpts{
		Context: ctx,
	}, operatorId)
}

func (t *Reader) OperatorAddressToID(ctx context.Context, address gethcommon.Address) (core.OperatorID, error) {
	return t.bindings.BLSApkRegistry.GetOperatorId(&bind.CallOpts{
		Context: ctx,
	}, address)
}

func (t *Reader) BatchOperatorIDToAddress(ctx context.Context, operatorIds []core.OperatorID) ([]gethcommon.Address, error) {
	byteIds := make([][32]byte, len(operatorIds))
	for i, id := range operatorIds {
		byteIds[i] = [32]byte(id)
	}
	addresses, err := t.bindings.OpStateRetriever.GetBatchOperatorFromId(&bind.CallOpts{
		Context: ctx,
	}, t.bindings.RegCoordinatorAddr, byteIds)
	if err != nil {
		t.logger.Error("Failed to get operator address in batch", "err", err)
		return nil, err
	}
	return addresses, nil
}

func (t *Reader) GetCurrentQuorumBitmapByOperatorId(ctx context.Context, operatorId core.OperatorID) (*big.Int, error) {
	return t.bindings.RegistryCoordinator.GetCurrentQuorumBitmap(&bind.CallOpts{
		Context: ctx,
	}, operatorId)
}

func (t *Reader) GetQuorumBitmapForOperatorsAtBlockNumber(ctx context.Context, operatorIds []core.OperatorID, blockNumber uint32) ([]*big.Int, error) {
	if len(operatorIds) == 0 {
		return []*big.Int{}, nil
	}
	// When there is just one operator, we can get result by a single RPC with
	// getQuorumBitmapsAtBlockNumber() in OperatorStateRetrievercontract (v.s. 2
	// RPCs in the general case)
	if len(operatorIds) == 1 {
		byteId := [32]byte(operatorIds[0])
		bitmap, err := t.bindings.OpStateRetriever.GetQuorumBitmapsAtBlockNumber(&bind.CallOpts{
			Context: ctx,
		}, t.bindings.RegCoordinatorAddr, [][32]byte{byteId}, blockNumber)
		if err != nil {
			if err.Error() == "execution reverted: RegistryCoordinator.getQuorumBitmapIndexAtBlockNumber: no bitmap update found for operatorId at block number" {
				return []*big.Int{big.NewInt(0)}, nil
			} else {
				return nil, err
			}
		}
		return bitmap, nil
	}

	quorumCount, err := t.GetQuorumCount(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	quorumNumbers := make([]byte, quorumCount)
	for i := 0; i < len(quorumNumbers); i++ {
		quorumNumbers[i] = byte(uint8(i))
	}
	operatorsByQuorum, err := t.bindings.OpStateRetriever.GetOperatorState(&bind.CallOpts{
		Context: ctx,
	}, t.bindings.RegCoordinatorAddr, quorumNumbers, blockNumber)
	if err != nil {
		return nil, err
	}

	quorumsByOperator := make(map[core.OperatorID]map[uint8]bool)
	for i := range operatorsByQuorum {
		for _, op := range operatorsByQuorum[i] {
			if _, ok := quorumsByOperator[op.OperatorId]; !ok {
				quorumsByOperator[op.OperatorId] = make(map[uint8]bool)
			}
			quorumsByOperator[op.OperatorId][uint8(i)] = true
		}
	}
	bitmaps := make([]*big.Int, len(operatorIds))
	for i, op := range operatorIds {
		if quorums, ok := quorumsByOperator[op]; ok {
			bm := big.NewInt(0)
			for id := range quorums {
				bm.SetBit(bm, int(id), 1)
			}
			bitmaps[i] = bm
		} else {
			bitmaps[i] = big.NewInt(0)
		}
	}
	return bitmaps, nil
}

func (t *Reader) GetOperatorSetParams(ctx context.Context, quorumID core.QuorumID) (*core.OperatorSetParam, error) {

	operatorSetParams, err := t.bindings.RegistryCoordinator.GetOperatorSetParams(&bind.CallOpts{
		Context: ctx,
	}, quorumID)
	if err != nil {
		t.logger.Error("Failed to fetch operator set params", "err", err)
		return nil, err
	}

	return &core.OperatorSetParam{
		MaxOperatorCount:         operatorSetParams.MaxOperatorCount,
		ChurnBIPsOfOperatorStake: operatorSetParams.KickBIPsOfOperatorStake,
		ChurnBIPsOfTotalStake:    operatorSetParams.KickBIPsOfTotalStake,
	}, nil
}

// Returns the number of registered operators for the quorum.
func (t *Reader) GetNumberOfRegisteredOperatorForQuorum(ctx context.Context, quorumID core.QuorumID) (uint32, error) {
	return t.bindings.IndexRegistry.TotalOperatorsForQuorum(&bind.CallOpts{
		Context: ctx,
	}, quorumID)
}

func (t *Reader) WeightOfOperatorForQuorum(ctx context.Context, quorumID core.QuorumID, operator gethcommon.Address) (*big.Int, error) {
	return t.bindings.StakeRegistry.WeightOfOperatorForQuorum(&bind.CallOpts{
		Context: ctx,
	}, quorumID, operator)
}

func (t *Reader) CalculateOperatorChurnApprovalDigestHash(
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
	return t.bindings.RegistryCoordinator.CalculateOperatorChurnApprovalDigestHash(&bind.CallOpts{
		Context: ctx,
	}, operatorAddress, operatorId, opKickParams, salt, expiry)
}

func (t *Reader) GetCurrentBlockNumber(ctx context.Context) (uint32, error) {
	bn, err := t.ethClient.BlockNumber(ctx)
	return uint32(bn), err
}

func (t *Reader) GetQuorumCount(ctx context.Context, blockNumber uint32) (uint8, error) {
	return t.bindings.RegistryCoordinator.QuorumCount(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
}

func (t *Reader) GetQuorumSecurityParams(ctx context.Context, blockNumber uint32) ([]core.SecurityParam, error) {
	adversaryThresholdPercentegesBytes, err := t.bindings.EigenDAServiceManager.QuorumAdversaryThresholdPercentages(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return nil, err
	}

	confirmationThresholdPercentegesBytes, err := t.bindings.EigenDAServiceManager.QuorumConfirmationThresholdPercentages(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return nil, err
	}

	if len(adversaryThresholdPercentegesBytes) != len(confirmationThresholdPercentegesBytes) {
		return nil, errors.New("adversaryThresholdPercentegesBytes and confirmationThresholdPercentegesBytes have different lengths")
	}

	securityParams := make([]core.SecurityParam, len(adversaryThresholdPercentegesBytes))

	for i := range adversaryThresholdPercentegesBytes {
		securityParams[i] = core.SecurityParam{
			QuorumID:              core.QuorumID(i),
			AdversaryThreshold:    adversaryThresholdPercentegesBytes[i],
			ConfirmationThreshold: confirmationThresholdPercentegesBytes[i],
		}
	}

	return securityParams, nil

}

func (t *Reader) GetRequiredQuorumNumbers(ctx context.Context, blockNumber uint32) ([]uint8, error) {
	requiredQuorums, err := t.bindings.EigenDAServiceManager.QuorumNumbersRequired(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return nil, err
	}
	return requiredQuorums, nil
}

func (t *Reader) GetActiveReservations(ctx context.Context, blockNumber uint32, accountIDs []string) (map[string]core.ActiveReservation, error) {
	// contract is not implemented yet
	return map[string]core.ActiveReservation{}, nil
}

func (t *Reader) GetActiveReservationByAccount(ctx context.Context, blockNumber uint32, accountID string) (core.ActiveReservation, error) {
	// contract is not implemented yet
	return core.ActiveReservation{}, nil
}

func (t *Reader) GetOnDemandPayments(ctx context.Context, blockNumber uint32, accountIDs []string) (map[string]core.OnDemandPayment, error) {
	// contract is not implemented yet
	return map[string]core.OnDemandPayment{}, nil
}

func (t *Reader) GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint32, accountID string) (core.OnDemandPayment, error) {
	// contract is not implemented yet
	return core.OnDemandPayment{}, nil
}
