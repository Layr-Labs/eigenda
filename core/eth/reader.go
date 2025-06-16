package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	avsdir "github.com/Layr-Labs/eigenda/contracts/bindings/AVSDirectory"
	blsapkreg "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	delegationmgr "github.com/Layr-Labs/eigenda/contracts/bindings/DelegationManager"
	disperserreg "github.com/Layr-Labs/eigenda/contracts/bindings/DisperserRegistry"
	relayreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARelayRegistry"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	thresholdreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAThresholdRegistry"
	ejectionmg "github.com/Layr-Labs/eigenda/contracts/bindings/EjectionManager"
	indexreg "github.com/Layr-Labs/eigenda/contracts/bindings/IIndexRegistry"
	opstateretriever "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	paymentvault "github.com/Layr-Labs/eigenda/contracts/bindings/PaymentVault"
	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/RegistryCoordinator"
	socketreg "github.com/Layr-Labs/eigenda/contracts/bindings/SocketRegistry"
	stakereg "github.com/Layr-Labs/eigenda/contracts/bindings/StakeRegistry"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pingcap/errors"

	blssigner "github.com/Layr-Labs/eigensdk-go/signer/bls"
)

type ContractBindings struct {
	RegCoordinatorAddr    gethcommon.Address
	ServiceManagerAddr    gethcommon.Address
	RelayRegistryAddress  gethcommon.Address
	DelegationManager     *delegationmgr.ContractDelegationManager
	OpStateRetriever      *opstateretriever.ContractOperatorStateRetriever
	BLSApkRegistry        *blsapkreg.ContractBLSApkRegistry
	IndexRegistry         *indexreg.ContractIIndexRegistry
	RegistryCoordinator   *regcoordinator.ContractRegistryCoordinator
	StakeRegistry         *stakereg.ContractStakeRegistry
	EigenDAServiceManager *eigendasrvmg.ContractEigenDAServiceManager
	EjectionManager       *ejectionmg.ContractEjectionManager
	AVSDirectory          *avsdir.ContractAVSDirectory
	SocketRegistry        *socketreg.ContractSocketRegistry
	PaymentVault          *paymentvault.ContractPaymentVault
	RelayRegistry         *relayreg.ContractEigenDARelayRegistry
	ThresholdRegistry     *thresholdreg.ContractEigenDAThresholdRegistry
	DisperserRegistry     *disperserreg.ContractDisperserRegistry
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

	var contractSocketRegistry *socketreg.ContractSocketRegistry
	socketRegistryAddr, err := contractIRegistryCoordinator.SocketRegistry(&bind.CallOpts{})
	if err != nil {
		t.logger.Warn("Failed to fetch SocketRegistry address", "err", err)
		// TODO: don't panic until there is socket registry deployment
		// return err
	} else {
		contractSocketRegistry, err = socketreg.NewContractSocketRegistry(socketRegistryAddr, t.ethClient)
		if err != nil {
			t.logger.Error("Failed to fetch SocketRegistry contract", "err", err)
			return err
		}
	}

	relayRegistryAddress, err := contractEigenDAServiceManager.EigenDARelayRegistry(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch IEigenDARelayRegistry contract", "err", err)
		// TODO(ian-shim): return err when the contract is deployed
	}

	var contractThresholdRegistry *thresholdreg.ContractEigenDAThresholdRegistry
	thresholdRegistryAddr, err := contractEigenDAServiceManager.EigenDAThresholdRegistry(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch EigenDAThresholdRegistry contract", "err", err)
		// TODO(ian-shim): return err when the contract is deployed
	} else {
		contractThresholdRegistry, err = thresholdreg.NewContractEigenDAThresholdRegistry(thresholdRegistryAddr, t.ethClient)
		if err != nil {
			t.logger.Error("Failed to fetch EigenDAThresholdRegistry contract", "err", err)
		}
	}

	var contractPaymentVault *paymentvault.ContractPaymentVault
	paymentVaultAddr, err := contractEigenDAServiceManager.PaymentVault(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch PaymentVault address", "err", err)
		//TODO(hopeyen): return err when the contract is deployed
		// return err
	} else {
		contractPaymentVault, err = paymentvault.NewContractPaymentVault(paymentVaultAddr, t.ethClient)
		if err != nil {
			t.logger.Error("Failed to fetch PaymentVault contract", "err", err)
			return err
		}
	}

	var contractDisperserRegistry *disperserreg.ContractDisperserRegistry
	disperserRegistryAddr, err := contractEigenDAServiceManager.DisperserRegistry(&bind.CallOpts{})
	if err != nil {
		t.logger.Error("Failed to fetch DisperserRegistry address", "err", err)
		// TODO(cody-littley): return err when the contract is deployed
		// return err
	} else {
		contractDisperserRegistry, err =
			disperserreg.NewContractDisperserRegistry(disperserRegistryAddr, t.ethClient)
		if err != nil {
			t.logger.Error("Failed to fetch DisperserRegistry contract", "err", err)
			return err
		}
	}

	t.bindings = &ContractBindings{
		ServiceManagerAddr:    eigenDAServiceManagerAddr,
		RegCoordinatorAddr:    registryCoordinatorAddr,
		RelayRegistryAddress:  relayRegistryAddress,
		AVSDirectory:          contractAVSDirectory,
		SocketRegistry:        contractSocketRegistry,
		OpStateRetriever:      contractBLSOpStateRetr,
		BLSApkRegistry:        contractBLSPubkeyReg,
		IndexRegistry:         contractIIndexReg,
		RegistryCoordinator:   contractIRegistryCoordinator,
		EjectionManager:       contractEjectionManager,
		StakeRegistry:         contractStakeRegistry,
		EigenDAServiceManager: contractEigenDAServiceManager,
		DelegationManager:     contractDelegationManager,
		PaymentVault:          contractPaymentVault,
		ThresholdRegistry:     contractThresholdRegistry,
		DisperserRegistry:     contractDisperserRegistry,
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
	blssigner blssigner.Signer,
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
	sigBytes, err := blssigner.SignG1(ctx, msgToSignG1.Serialize())
	if err != nil {
		return nil, nil, err
	}
	sig := new(core.Signature)
	g, err := sig.Deserialize(sigBytes)
	if err != nil {
		return nil, nil, err
	}
	signature := &core.Signature{
		G1Point: g,
	}

	signedMessageHashParam := regcoordinator.BN254G1Point{
		X: signature.X.BigInt(big.NewInt(0)),
		Y: signature.Y.BigInt(big.NewInt(0)),
	}

	g1KeyHex := blssigner.GetPublicKeyG1()
	g1KeyBytes, err := hex.DecodeString(g1KeyHex)
	if err != nil {
		return nil, nil, err
	}
	g1point := new(core.G1Point)
	g1point, err = g1point.Deserialize(g1KeyBytes)
	if err != nil {
		return nil, nil, err
	}
	g1Point_ := pubKeyG1ToBN254G1Point(g1point)
	g1Point := regcoordinator.BN254G1Point{
		X: g1Point_.X,
		Y: g1Point_.Y,
	}

	g2KeyHex := blssigner.GetPublicKeyG2()
	g2KeyBytes, err := hex.DecodeString(g2KeyHex)
	if err != nil {
		return nil, nil, err
	}
	g2point := new(core.G2Point)
	g2point, err = g2point.Deserialize(g2KeyBytes)
	if err != nil {
		return nil, nil, err
	}
	g2Point_ := pubKeyG2ToBN254G2Point(g2point)
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
		t.logger.Errorf("Failed to fetch operator state: %s", err)
		return nil, fmt.Errorf("failed to fetch operator state: %w", err)
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

// GetOperatorStakesForQuorums returns the stakes of all operators within the supplied quorums. The returned stakes are for the block number supplied.
// The indices of the operators within each quorum are also returned.
func (t *Reader) GetOperatorStakesWithSocketForQuorums(ctx context.Context, quorums []core.QuorumID, blockNumber uint32) (core.OperatorStakesWithSocket, error) {
	quorumBytes := make([]byte, len(quorums))
	for ind, quorum := range quorums {
		quorumBytes[ind] = byte(uint8(quorum))
	}

	// result is a struct{Operators [][]opstateretriever.OperatorStateRetrieverOperator; Sockets [][]string}
	// Operators is a [][]*opstateretriever.OperatorStake with the same length and order as quorumBytes, and then indexed by operator index
	// Sockets is a [][]string with the same length and order as quorumBytes, and then indexed by operator index
	// By contract definition, Operators and Sockets are parallel arrays
	result, err := t.bindings.OpStateRetriever.GetOperatorStateWithSocket(&bind.CallOpts{
		Context: ctx,
	}, t.bindings.RegCoordinatorAddr, quorumBytes, blockNumber)
	if err != nil {
		t.logger.Errorf("Failed to fetch operator state: %s", err)
		return nil, fmt.Errorf("failed to fetch operator state: %w", err)
	}

	state := make(core.OperatorStakesWithSocket, len(result.Operators))
	for i := range result.Operators {
		quorumID := quorums[i]
		state[quorumID] = make(map[core.OperatorIndex]core.OperatorStakeWithSocket, len(result.Operators[i]))
		for j, op := range result.Operators[i] {
			operatorIndex := core.OperatorIndex(j)
			state[quorumID][operatorIndex] = core.OperatorStakeWithSocket{
				Stake:      op.Stake,
				OperatorID: op.OperatorId,
				Socket:     core.OperatorSocket(result.Sockets[i][j]),
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

func (t *Reader) SocketRegistry(ctx context.Context) (gethcommon.Address, error) {
	return t.bindings.RegistryCoordinator.SocketRegistry(&bind.CallOpts{
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

func (t *Reader) BatchOperatorAddressToID(ctx context.Context, addresses []gethcommon.Address) ([]core.OperatorID, error) {
	ids, err := t.bindings.OpStateRetriever.GetBatchOperatorId(&bind.CallOpts{
		Context: ctx,
	}, t.bindings.RegCoordinatorAddr, addresses)
	if err != nil {
		t.logger.Error("Failed to get operator IDs in batch", "err", err)
		return nil, err
	}
	operatorIds := make([]core.OperatorID, len(ids))
	for i, id := range ids {
		operatorIds[i] = core.OperatorID(id)
	}
	return operatorIds, nil
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

func (t *Reader) GetNumBlobVersions(ctx context.Context) (uint16, error) {
	if t.bindings.ThresholdRegistry == nil {
		return 0, errors.New("threshold registry not deployed")
	}

	return t.bindings.ThresholdRegistry.NextBlobVersion(&bind.CallOpts{
		Context: ctx,
	})
}

func (t *Reader) GetVersionedBlobParams(ctx context.Context, blobVersion uint16) (*core.BlobVersionParameters, error) {
	params, err := t.bindings.EigenDAServiceManager.GetBlobParams(&bind.CallOpts{
		Context: ctx,
	}, uint16(blobVersion))
	if err != nil {
		return nil, err
	}
	return &core.BlobVersionParameters{
		CodingRate:      uint32(params.CodingRate),
		NumChunks:       params.NumChunks,
		MaxNumOperators: params.MaxNumOperators,
	}, nil
}

func (t *Reader) GetAllVersionedBlobParams(ctx context.Context) (map[uint16]*core.BlobVersionParameters, error) {
	if t.bindings.ThresholdRegistry == nil {
		return nil, errors.New("threshold registry not deployed")
	}

	numBlobVersions, err := t.GetNumBlobVersions(ctx)
	if err != nil {
		return nil, err
	}

	res := make(map[uint16]*core.BlobVersionParameters)
	for version := uint16(0); version < uint16(numBlobVersions); version++ {
		params, err := t.GetVersionedBlobParams(ctx, version)
		if err != nil && strings.Contains(err.Error(), "execution reverted") {
			break
		} else if err != nil {
			return nil, err
		}

		res[version] = params
	}

	if len(res) == 0 {
		return nil, errors.New("no blob version parameters found")
	}

	return res, nil
}

// TODO: later add quorumIds specifications
func (t *Reader) GetReservedPayments(ctx context.Context, accountIDs []gethcommon.Address) (map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment, error) {
	if t.bindings.PaymentVault == nil {
		return nil, errors.New("payment vault not deployed")
	}

	reservationsMap := make(map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment)

	// In the new API, we need to query reservations per quorum and account
	// For now, there's only one reservation per account, applied to multiple quorums
	// TODO: update this to handle multiple quorums
	for _, accountID := range accountIDs {
		//TODO: add quorumId to the call
		reservation, err := t.bindings.PaymentVault.GetReservation(&bind.CallOpts{
			Context: ctx,
		}, accountID)

		if err != nil {
			t.logger.Warn("failed to get reservation", "account", accountID, "err", err)
			continue
		}

		res, err := ConvertToReservedPayments(reservation)
		if err != nil {
			t.logger.Warn("failed to convert reservation", "account", accountID, "err", err)
			continue
		}

		reservationsMap[accountID] = res
	}

	return reservationsMap, nil
}

// TODO: this function should take in a list of quorumIds
func (t *Reader) GetReservedPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (map[core.QuorumID]*core.ReservedPayment, error) {
	if t.bindings.PaymentVault == nil {
		return nil, errors.New("payment vault not deployed")
	}

	reservation, err := t.bindings.PaymentVault.GetReservation(&bind.CallOpts{
		Context: ctx,
	}, accountID)

	if err != nil {
		return nil, err
	}

	return ConvertToReservedPayments(reservation)
}

func (t *Reader) GetOnDemandPayments(ctx context.Context, accountIDs []gethcommon.Address) (map[gethcommon.Address]*core.OnDemandPayment, error) {
	if t.bindings.PaymentVault == nil {
		return nil, errors.New("payment vault not deployed")
	}

	paymentsMap := make(map[gethcommon.Address]*core.OnDemandPayment)

	// TODO: update to use default of quorum 0 for deposits
	// quorumId := uint64(0)
	for _, accountID := range accountIDs {
		onDemandDeposit, err := t.bindings.PaymentVault.GetOnDemandTotalDeposit(&bind.CallOpts{
			Context: ctx,
		}, accountID)

		if err != nil {
			t.logger.Warn("failed to get on-demand deposit", "account", accountID, "err", err)
			continue
		}

		if onDemandDeposit.Cmp(big.NewInt(0)) == 0 {
			t.logger.Debug("on-demand deposit is zero for account", "account", accountID)
			continue
		}

		paymentsMap[accountID] = &core.OnDemandPayment{
			CumulativePayment: onDemandDeposit,
		}
	}

	return paymentsMap, nil
}

func (t *Reader) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	if t.bindings.PaymentVault == nil {
		return nil, errors.New("payment vault not deployed")
	}

	// In the new API, we need to specify a quorumId (only 0 enabled for deposits)
	onDemandDeposit, err := t.bindings.PaymentVault.GetOnDemandTotalDeposit(&bind.CallOpts{
		Context: ctx,
	}, accountID)

	if err != nil {
		return nil, err
	}

	if onDemandDeposit.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("on-demand deposit does not exist for given account")
	}

	return &core.OnDemandPayment{
		CumulativePayment: onDemandDeposit,
	}, nil
}

func (t *Reader) GetOnDemandGlobalSymbolsPerSecond(ctx context.Context, blockNumber uint32) (uint64, error) {
	if t.bindings.PaymentVault == nil {
		return 0, errors.New("payment vault not deployed")
	}
	globalSymbolsPerSecond, err := t.bindings.PaymentVault.GlobalSymbolsPerPeriod(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return globalSymbolsPerSecond, nil
}

func (t *Reader) GetOnDemandGlobalRatePeriodInterval(ctx context.Context, blockNumber uint32) (uint64, error) {
	if t.bindings.PaymentVault == nil {
		return 0, errors.New("payment vault not deployed")
	}
	globalRateBinInterval, err := t.bindings.PaymentVault.GlobalRatePeriodInterval(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return globalRateBinInterval, nil
}

func (t *Reader) GetMinNumSymbols(ctx context.Context, blockNumber uint32) (uint64, error) {
	if t.bindings.PaymentVault == nil {
		return 0, errors.New("payment vault not deployed")
	}
	minNumSymbols, err := t.bindings.PaymentVault.MinNumSymbols(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return minNumSymbols, nil
}

func (t *Reader) GetPricePerSymbol(ctx context.Context, blockNumber uint32) (uint64, error) {
	if t.bindings.PaymentVault == nil {
		return 0, errors.New("payment vault not deployed")
	}
	pricePerSymbol, err := t.bindings.PaymentVault.PricePerSymbol(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return pricePerSymbol, nil
}

func (t *Reader) GetReservationWindow(ctx context.Context, blockNumber uint32) (uint64, error) {
	if t.bindings.PaymentVault == nil {
		return 0, errors.New("payment vault not deployed")
	}
	reservationWindow, err := t.bindings.PaymentVault.ReservationPeriodInterval(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return 0, err
	}
	return reservationWindow, nil
}

func (t *Reader) GetOperatorSocket(ctx context.Context, operatorId core.OperatorID) (string, error) {
	if t.bindings.SocketRegistry == nil {
		return "", errors.New("socket registry not enabled")
	}
	socket, err := t.bindings.SocketRegistry.GetOperatorSocket(&bind.CallOpts{
		Context: ctx}, [32]byte(operatorId))
	if err != nil {
		return "", err
	}
	if socket == "" {
		return "", errors.New("operator socket string is empty, check operator with id: " + operatorId.Hex())
	}
	return socket, nil
}

func (t *Reader) GetDisperserAddress(ctx context.Context, disperserID uint32) (gethcommon.Address, error) {
	registry := t.bindings.DisperserRegistry
	if registry == nil {
		return gethcommon.Address{}, errors.New("disperser registry not deployed")
	}

	address, err := registry.GetDisperserAddress(
		&bind.CallOpts{
			Context: ctx,
		},
		disperserID)

	var defaultAddress gethcommon.Address
	if err != nil {
		return defaultAddress, fmt.Errorf("failed to get disperser address: %w", err)
	}
	if address == defaultAddress {
		return defaultAddress, fmt.Errorf("disperser with id %d not found", disperserID)
	}

	return address, nil
}

func (t *Reader) GetRelayRegistryAddress() gethcommon.Address {
	return t.bindings.RelayRegistryAddress
}

func (t *Reader) GetRegistryCoordinatorAddress() gethcommon.Address {
	return t.bindings.RegCoordinatorAddr
}
