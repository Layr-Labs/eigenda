package cert_gas_meter

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common"
	blsapkregistry "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	opstateretriever "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"

	gnarkbn254 "github.com/consensys/gnark-crypto/ecc/bn254"

	certVerifierBinding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	certTypesBinding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
)

func OperatorIDsScan(
	quorumBytes []byte,
	blockNumber uint32,
	ethClient common.EthClient,
	logger logging.Logger,
	certV3 coretypes.EigenDACertV3,
) ([][32]byte, error) {

	ctx := context.Background()
	operatorStateRetrieverAddr := gethcommon.HexToAddress("0x22478d082E9edaDc2baE8443E4aC9473F6E047Ff")

	contractOpStateRetr, err := opstateretriever.NewContractOperatorStateRetriever(operatorStateRetrieverAddr, ethClient)
	if err != nil {
		logger.Error("Failed to fetch OperatorStateRetriever contract", "err", err)
		return nil, err
	}

	blsApkRegistryAddr := gethcommon.HexToAddress("0xA8fF891E5b8cA255A0e884129bc14977F7A742BC")

	bLSApkRegistry, err := blsapkregistry.NewContractBLSApkRegistry(blsApkRegistryAddr, ethClient)
	if err != nil {
		logger.Error("Failed to fetch NewContractBLSApkRegistry contract", "err", err)
		return nil, err
	}

	registryCoordinatorAddr := gethcommon.HexToAddress("0xAF21d3811B5d23D5466AC83BA7a9c34c261A8D81")

	quorums := make([]core.QuorumID, len(quorumBytes))
	for i := range len(quorumBytes) {
		quorums[i] = uint8(quorumBytes[i])
	}

	// state_ is a [][]*opstateretriever.OperatorStake with the same length and order as quorumBytes, and then indexed by operator index
	state_, err := contractOpStateRetr.GetOperatorState(&bind.CallOpts{
		Context: context.Background(),
	}, registryCoordinatorAddr, quorumBytes, blockNumber)

	allOperatorIDs := make([][32]byte, 0)
	allOperatorMaps := make(map[core.OperatorID]bool)
	//
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
			_, found := allOperatorMaps[op.OperatorId]
			if !found {
				allOperatorMaps[op.OperatorId] = true
				allOperatorIDs = append(allOperatorIDs, op.OperatorId)
			}
		}
	}

	// sort non signer keys according to how it's checked onchain
	// ref: https://github.com/Layr-Labs/eigenlayer-middleware/blob/m2-mainnet/src/BLSSignatureChecker.sol#L99
	sort.Slice(allOperatorIDs, func(i, j int) bool {
		hash1 := allOperatorIDs[i]
		hash2 := allOperatorIDs[j]
		// sort in nonSignerOperatorIDs order
		return bytes.Compare(hash1[:], hash2[:]) == -1
	})

	checkSigIndices, err := contractOpStateRetr.GetCheckSignaturesIndices(&bind.CallOpts{Context: ctx, BlockNumber: big.NewInt(int64(blockNumber))},
		registryCoordinatorAddr, uint32(blockNumber), quorumBytes, allOperatorIDs)
	if err != nil {
		// We log the call parameters for debugging purposes: input them into tenderly to simulate the call and get more context.
		logger.Error("eth-call failed")
		return nil, fmt.Errorf("check sig indices call: %w", err)
	}

	//fmt.Println("checkSigIndices")
	//fmt.Println(checkSigIndices)

	nonSignerPubKeysBN254 := make([]certTypesBinding.BN254G1Point, 0)

	for i := range len(allOperatorIDs) {
		operatorID_ := allOperatorIDs[i]
		operatorAddr, err := bLSApkRegistry.PubkeyHashToOperator(&bind.CallOpts{Context: ctx}, operatorID_)
		if err != nil {
			logger.Error("eth-call PubkeyHashToOperator failed")
			return nil, fmt.Errorf("check sig indices call: %w", err)
		}
		operatorG1, err := bLSApkRegistry.OperatorToPubkey(&bind.CallOpts{Context: ctx}, operatorAddr)
		if err != nil {
			logger.Error("eth-call OperatorToPubkey failed")
			return nil, fmt.Errorf("check sig indices call: %w", err)
		}
		nonSignerPubKeysBN254 = append(nonSignerPubKeysBN254, operatorG1)
	}

	//for i, operatorID := range allOperatorIDs {
	//	fmt.Printf("%v operatorID %v PubKey %v\n", i, core.OperatorID(operatorID).Hex(), nonSignerPubKeysBN254[i])
	//}

	var sigmaBn254 gnarkbn254.G1Affine
	sigmaBn254.SetInfinity()

	sigma := certTypesBinding.BN254G1Point{
		X: sigmaBn254.X.BigInt(new(big.Int)),
		Y: sigmaBn254.Y.BigInt(new(big.Int)),
	}

	var apkG2Bn254 gnarkbn254.G2Affine
	apkG2Bn254.SetInfinity()

	apkG2 := certTypesBinding.BN254G2Point{
		X: [2]*big.Int{apkG2Bn254.X.A1.BigInt(new(big.Int)), apkG2Bn254.X.A0.BigInt(new(big.Int))},
		Y: [2]*big.Int{apkG2Bn254.Y.A1.BigInt(new(big.Int)), apkG2Bn254.Y.A0.BigInt(new(big.Int))},
	}

	getNonSignerStakesAndSignature := certTypesBinding.EigenDATypesV1NonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: checkSigIndices.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:             nonSignerPubKeysBN254,                         // everyone
		QuorumApks:                   certV3.NonSignerStakesAndSignature.QuorumApks, // keep
		ApkG2:                        apkG2,                                         // made up 0
		Sigma:                        sigma,                                         // made up 0
		QuorumApkIndices:             checkSigIndices.QuorumApkIndices,
		TotalStakeIndices:            checkSigIndices.TotalStakeIndices,
		NonSignerStakeIndices:        checkSigIndices.NonSignerStakeIndices,
	}

	fmt.Println("number of operators are", len(allOperatorIDs))

	fmt.Println("checkSigIndices.NonSignerQuorumBitmapIndices")
	fmt.Println(checkSigIndices.NonSignerQuorumBitmapIndices)
	fmt.Println(nonSignerPubKeysBN254)
	fmt.Println(certV3.NonSignerStakesAndSignature.QuorumApks)
	fmt.Println("apkG2", apkG2)
	fmt.Println("sigma", sigma)
	fmt.Println(checkSigIndices.QuorumApkIndices)
	fmt.Println(checkSigIndices.TotalStakeIndices)
	fmt.Println(checkSigIndices.NonSignerStakeIndices)

	fmt.Println("signed quorum bytes", certV3.SignedQuorumNumbers)
	fmt.Println("blob inclusion quprum", certV3.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers)

	V3CertVerifierAddr := gethcommon.HexToAddress("0x58D2B844a894f00b7E6F9F492b9F43aD54Cd4429")

	certVerifierCaller, err := certVerifierBinding.NewContractEigenDACertVerifier(V3CertVerifierAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", err)
	}

	certV3.NonSignerStakesAndSignature = getNonSignerStakesAndSignature
	certV3.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers = make([]byte, 0)

	fmt.Println("before certV3.Serialize")
	certBytes, err := certV3.Serialize(coretypes.CertSerializationABI)
	if err != nil {
		return nil, fmt.Errorf("serialize cert %w", err)
	}

	// TODO: determine if there's any merit in passing call options to impose better determinism and
	// safety on the operation
	result, err := certVerifierCaller.CheckDACert(
		&bind.CallOpts{Context: ctx},
		certBytes,
	)

	fmt.Println("result %v", result)
	fmt.Println("result %v", coretypes.VerificationStatusCode(result))

	input, err := BuildCallInput(certBytes)
	if err != nil {
		return nil, fmt.Errorf("BuildCallInput %w", err)
	}

	msg := ethereum.CallMsg{
		From: gethcommon.HexToAddress("0x41d52a62591282784bF7ACAFc0a321B8A468e07a"),
		To:   &V3CertVerifierAddr,
		Data: input,
	}

	estimate, err := ethClient.EstimateGas(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("EstimateGas %w", err)
	}
	fmt.Println("estimate %v", estimate)

	return allOperatorIDs, nil
}

func BuildCallInput(certBytes []byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(certVerifierBinding.ContractEigenDACertVerifierABI))
	if err != nil {
		return nil, err
	}
	data, err := a.Pack("checkDACert", certBytes)
	if err != nil {
		return nil, err
	}
	return data, nil
}

/*
// GetOperatorStakesForQuorums returns the stakes of all operators within the supplied quorums. The returned stakes are for the block number supplied.
// The indices of the operators within each quorum are also returned.
func GetOperatorStakesForQuorums(ctx context.Context, quorums []core.QuorumID, blockNumber uint32) (core.OperatorStakes, error) {

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
*/
