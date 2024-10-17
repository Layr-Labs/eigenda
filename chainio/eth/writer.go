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

	avsdir "github.com/Layr-Labs/eigenda/contracts/bindings/AVSDirectory"
	blsapkreg "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	delegationmgr "github.com/Layr-Labs/eigenda/contracts/bindings/DelegationManager"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	ejectionmg "github.com/Layr-Labs/eigenda/contracts/bindings/EjectionManager"
	indexreg "github.com/Layr-Labs/eigenda/contracts/bindings/IIndexRegistry"
	opstateretriever "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/RegistryCoordinator"
	stakereg "github.com/Layr-Labs/eigenda/contracts/bindings/StakeRegistry"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Writer struct {
	EthClient common.EthClient
	Logger    logging.Logger
	Bindings  *ContractBindings
}

var _ chainio.Writer = (*Writer)(nil)

func NewWriter(
	logger logging.Logger,
	client common.EthClient,
	blsOperatorStateRetrieverHexAddr string,
	eigenDAServiceManagerHexAddr string) (*Writer, error) {

	e := &Writer{
		EthClient: client,
		Logger:    logger.With("component", "Writer"),
	}

	blsOperatorStateRetrieverAddr := gethcommon.HexToAddress(blsOperatorStateRetrieverHexAddr)
	eigenDAServiceManagerAddr := gethcommon.HexToAddress(eigenDAServiceManagerHexAddr)
	err := e.updateContractBindings(blsOperatorStateRetrieverAddr, eigenDAServiceManagerAddr)

	return e, err
}

func (t *Writer) getRegistrationParams(
	ctx context.Context,
	keypair *bn254.KeyPair,
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

	msgToSignG1 := bn254.NewG1Point(msgToSignG1_.X, msgToSignG1_.Y)
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

func (t *Writer) updateContractBindings(blsOperatorStateRetrieverAddr, eigenDAServiceManagerAddr gethcommon.Address) error {

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

	contractEjectionManagerAddr, err := contractIRegistryCoordinator.Ejector(&bind.CallOpts{})
	if err != nil {
		t.Logger.Error("Failed to fetch EjectionManager address", "err", err)
		return err
	}
	contractEjectionManager, err := ejectionmg.NewContractEjectionManager(contractEjectionManagerAddr, t.EthClient)
	if err != nil {
		t.Logger.Error("Failed to fetch EjectionManager contract", "err", err)
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
		EjectionManager:       contractEjectionManager,
		StakeRegistry:         contractStakeRegistry,
		EigenDAServiceManager: contractEigenDAServiceManager,
		DelegationManager:     contractDelegationManager,
	}
	return nil
}
