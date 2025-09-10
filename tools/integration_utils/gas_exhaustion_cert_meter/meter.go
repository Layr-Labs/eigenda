package gas_exhaustion_cert_meter

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/altdacommitment_parser"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/urfave/cli"

	gnarkbn254 "github.com/consensys/gnark-crypto/ecc/bn254"

	certVerifierBinding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	certTypesBinding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
)

func RunMeterer(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// Read and decode the certificate file
	prefix, versionedCert, err := altdacommitment_parser.ParseAltDACommitmentFromHex(config.CertHexString)
	if err != nil {
		return fmt.Errorf("failed to parse cert hex string: %w", err)
	}

	altdacommitment_parser.DisplayPrefixInfo(prefix)

	cert, err := altdacommitment_parser.ParseCertificateData(versionedCert)
	if err != nil {
		return fmt.Errorf("failed to parse versioned cert: %w", err)
	}

	if err = EstimateGas(config, *cert); err != nil {
		return fmt.Errorf("gas estimation failed: %w", err)
	}

	return nil
}

// EstimateGas calculates the worst-case gas cost for verifying an EigenDA V3 certificate.
// It simulates a scenario where all operators are non-signers, requiring maximum verification work.
// Returns the list of all operator IDs and any error encountered during estimation.
func EstimateGas(
	config *Config,
	certV3 coretypes.EigenDACertV3,
) error {

	blockNumber := certV3.BatchHeader.ReferenceBlockNumber
	quorumBytes := certV3.SignedQuorumNumbers

	allOperatorIDs, err := GetAllOperatorID(config, quorumBytes, blockNumber)
	if err != nil {
		return fmt.Errorf(
			"failed to get all operatorID at block %v for quorumBytes %v: %w",
			blockNumber, quorumBytes, err)
	}

	// Sort operator IDs to match on-chain verification order
	// Reference: https://github.com/Layr-Labs/eigenlayer-middleware/blob/m2-mainnet/src/BLSSignatureChecker.sol#L99
	// Reference: EigenDA core/aggregation.go#L391
	sort.Slice(allOperatorIDs, func(i, j int) bool {
		return bytes.Compare(allOperatorIDs[i][:], allOperatorIDs[j][:]) < 0
	})

	checkSigIndices, err := config.OpStateRetrCaller.GetCheckSignaturesIndices(
		&bind.CallOpts{Context: config.Ctx, BlockNumber: big.NewInt(int64(blockNumber))},
		config.RegistryCoordinatorAddr, blockNumber, quorumBytes, allOperatorIDs)
	if err != nil {
		return fmt.Errorf("eth call failed checkSigIndices: %w", err)
	}

	nonSignerPubKeys := make([]certTypesBinding.BN254G1Point, 0)

	for _, operatorID := range allOperatorIDs {
		operatorAddr, err := config.BLSApkRegistryCaller.PubkeyHashToOperator(&bind.CallOpts{Context: config.Ctx}, operatorID)
		if err != nil {
			return fmt.Errorf("eth-call PubkeyHashToOperator failed: %w", err)
		}
		operatorG1, err := config.BLSApkRegistryCaller.OperatorToPubkey(&bind.CallOpts{Context: config.Ctx}, operatorAddr)
		if err != nil {
			return fmt.Errorf("eth-call OperatorToPubkey failed: %w", err)
		}
		nonSignerPubKeys = append(nonSignerPubKeys, operatorG1)
	}

	// G1 point at infinity
	var sigmaBn254 gnarkbn254.G1Affine
	sigmaBn254.SetInfinity()
	// convert into EigenDA type
	sigma := certTypesBinding.BN254G1Point{
		X: sigmaBn254.X.BigInt(new(big.Int)),
		Y: sigmaBn254.Y.BigInt(new(big.Int)),
	}

	// G2 point at infinity
	var apkG2Bn254 gnarkbn254.G2Affine
	apkG2Bn254.SetInfinity()
	// convert into EigenDA type
	apkG2 := certTypesBinding.BN254G2Point{
		X: [2]*big.Int{apkG2Bn254.X.A1.BigInt(new(big.Int)), apkG2Bn254.X.A0.BigInt(new(big.Int))},
		Y: [2]*big.Int{apkG2Bn254.Y.A1.BigInt(new(big.Int)), apkG2Bn254.Y.A0.BigInt(new(big.Int))},
	}

	// Create worst-case scenario with all operators as non-signers
	worstCaseSignature := certTypesBinding.EigenDATypesV1NonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: checkSigIndices.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:             nonSignerPubKeys,
		QuorumApks:                   certV3.NonSignerStakesAndSignature.QuorumApks,
		ApkG2:                        apkG2, // Set to infinity (worst case)
		Sigma:                        sigma, // Set to infinity (worst case)
		QuorumApkIndices:             checkSigIndices.QuorumApkIndices,
		TotalStakeIndices:            checkSigIndices.TotalStakeIndices,
		NonSignerStakeIndices:        checkSigIndices.NonSignerStakeIndices,
	}

	certV3.NonSignerStakesAndSignature = worstCaseSignature

	certBytes, err := certV3.Serialize(coretypes.CertSerializationABI)
	if err != nil {
		return fmt.Errorf("serialize cert %w", err)
	}

	input, err := BuildCallInput(certBytes)
	if err != nil {
		return fmt.Errorf("BuildCallInput %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &config.CertVerifierAddr,
		Data: input,
	}

	estimate, err := config.EthClient.EstimateGas(config.Ctx, msg)
	if err != nil {
		return fmt.Errorf("EstimateGas %w", err)
	}
	config.Logger.Info("Gas estimation complete", "gasEstimate", estimate, "numOperators", len(allOperatorIDs))

	return nil
}

// BuildCallInput constructs the ABI-encoded input data for calling the checkDACert function.
func BuildCallInput(certBytes []byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(certVerifierBinding.ContractEigenDACertVerifierABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	data, err := a.Pack("checkDACert", certBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI data: %w", err)
	}
	return data, nil
}

// GetAllOperatorID retrieves all operator IDs at a block number for quorums encoded in quorumBytes,
// where each byte encodes a quorumID (uint8). This is similar to retrieving all stakes for operators.
// Reference: https://github.com/Layr-Labs/eigenda/blob/8d1bfff8fecfd0e4bc6c6b8319296a58f76845d5/core/eth/reader.go#L471
func GetAllOperatorID(config *Config, quorumBytes []byte, blockNumber uint32) ([][32]byte, error) {
	// Retrieve operator state for all quorums at the specified block number
	state_, err := config.OpStateRetrCaller.GetOperatorState(&bind.CallOpts{
		Context: context.Background(),
	}, config.RegistryCoordinatorAddr, quorumBytes, blockNumber)

	if err != nil {
		return nil, fmt.Errorf("eth call failed GetOperatorState: %w", err)
	}

	// Collect all unique operator IDs across quorums
	allOperatorIDs := make([][32]byte, 0)
	allOperatorMap := make(map[core.OperatorID]bool)
	for quorum_i := range state_ {
		for _, op := range state_[quorum_i] {
			// An operator may be registered in multiple quorums, so deduplicate
			if !allOperatorMap[op.OperatorId] {
				allOperatorMap[op.OperatorId] = true
				allOperatorIDs = append(allOperatorIDs, op.OperatorId)
			}
		}
	}
	return allOperatorIDs, nil
}
