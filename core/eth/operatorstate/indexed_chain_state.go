package operatorstate

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	contractBLSApkRegistry "github.com/Layr-Labs/eigenda/contracts/bindings/BLSApkRegistry"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// IndexedChainState implements core.IndexedChainState using direct on-chain contract queries
type IndexedChainState struct {
	contractClient   *ContractClient
	chainState       core.ChainState // Delegate core.ChainState methods to existing implementation
	blsApkRegistry   *contractBLSApkRegistry.ContractBLSApkRegistry
	contractDirectory *directory.ContractDirectory
}

// NewIndexedChainState creates a new IndexedChainState that queries contracts directly
func NewIndexedChainState(
	ethRpcUrl string,
	registryCoordinatorAddress gethcommon.Address,
	operatorStateRetrieverAddress gethcommon.Address,
	chainState core.ChainState,
	contractDirectory *directory.ContractDirectory,
) (*IndexedChainState, error) {
	client, err := NewContractClient(ethRpcUrl, registryCoordinatorAddress, operatorStateRetrieverAddress)
	if err != nil {
		return nil, err
	}

	// Get BLS APK Registry address from contract directory
	blsApkRegistryAddress, err := contractDirectory.GetContractAddress(context.Background(), directory.BLSApkRegistry)
	if err != nil {
		return nil, fmt.Errorf("failed to get BLS APK Registry address: %w", err)
	}

	// Create ethereum client
	ethClient, err := ethclient.Dial(ethRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	// Create BLS APK Registry contract binding
	blsApkRegistry, err := contractBLSApkRegistry.NewContractBLSApkRegistry(blsApkRegistryAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create BLS APK Registry contract binding: %w", err)
	}

	return &IndexedChainState{
		contractClient:    client,
		chainState:        chainState,
		blsApkRegistry:    blsApkRegistry,
		contractDirectory: contractDirectory,
	}, nil
}

// Start implements core.IndexedChainState.Start
func (s *IndexedChainState) Start(ctx context.Context) error {
	// No background processes needed for direct contract queries
	return nil
}

// GetIndexedOperatorState implements core.IndexedChainState.GetIndexedOperatorState
func (s *IndexedChainState) GetIndexedOperatorState(
	ctx context.Context, blockNumber uint, quorums []core.QuorumID,
) (*core.IndexedOperatorState, error) {
	// Convert core.QuorumID to byte slice
	quorumBytes := make([]byte, len(quorums))
	copy(quorumBytes, quorums)

	// Query the operator state using contract client
	result, err := s.contractClient.GetOperatorStateWithSocket(ctx, quorumBytes, uint64(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("failed to get operator state: %w", err)
	}

	// Convert contract result to core.IndexedOperatorState
	indexedOperators := make(map[core.OperatorID]*core.IndexedOperatorInfo)

	// Create the underlying OperatorState
	operatorState := &core.OperatorState{
		Operators:   make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo),
		Totals:      make(map[core.QuorumID]*core.OperatorInfo),
		BlockNumber: blockNumber,
	}

	// Process each quorum
	for quorumIndex, quorumNumber := range quorumBytes {
		quorumID := quorumNumber
		operatorState.Operators[quorumID] = make(map[core.OperatorID]*core.OperatorInfo)

		totalStake := new(big.Int)
		operatorCount := 0

		// Check if we have operators for this quorum index
		if quorumIndex < len(result.Operators) {
			operators := result.Operators[quorumIndex]
			var sockets []string
			if quorumIndex < len(result.Sockets) {
				sockets = result.Sockets[quorumIndex]
			}

			for opIndex, op := range operators {
				operatorID := core.OperatorID(op.OperatorId)

				// Get socket information
				var socket string
				if opIndex < len(sockets) && sockets[opIndex] != "" {
					socket = sockets[opIndex]
				}

				// Add to operator state
				operatorState.Operators[quorumID][operatorID] = &core.OperatorInfo{
					Stake:  op.Stake,
					Index:  core.OperatorIndex(opIndex),
					Socket: core.OperatorSocket(socket),
				}

				// Create IndexedOperatorInfo if it doesn't exist
				if _, exists := indexedOperators[operatorID]; !exists {
					// Resolve actual pubkeys from BLS registry
					pubkeyG1, pubkeyG2, err := s.getOperatorPubkeys(ctx, op.Operator, blockNumber)
					if err != nil {
						return nil, fmt.Errorf("failed to get pubkeys for operator %s: %w", op.Operator.Hex(), err)
					}

					indexedOperators[operatorID] = &core.IndexedOperatorInfo{
						PubkeyG1: pubkeyG1,
						PubkeyG2: pubkeyG2,
						Socket:   socket,
					}
				}

				totalStake.Add(totalStake, op.Stake)
				operatorCount++
			}
		}

		// Set totals for this quorum
		operatorState.Totals[quorumID] = &core.OperatorInfo{
			Stake:  totalStake,
			Index:  core.OperatorIndex(operatorCount),
			Socket: "", // Empty socket for totals
		}
	}

	// Create aggregate pubkeys (placeholders in this basic version)
	aggregatePubKeys := make(map[core.QuorumID]*core.G1Point)
	for _, quorumNumber := range quorumBytes {
		quorumID := quorumNumber
		aggregatePubKeys[quorumID] = &core.G1Point{G1Affine: &bn254.G1Affine{}} // placeholder
	}

	return &core.IndexedOperatorState{
		OperatorState:    operatorState,
		IndexedOperators: indexedOperators,
		AggKeys:          aggregatePubKeys,
	}, nil
}

// GetIndexedOperators implements core.IndexedChainState.GetIndexedOperators
func (s *IndexedChainState) GetIndexedOperators(
	ctx context.Context, blockNumber uint,
) (map[core.OperatorID]*core.IndexedOperatorInfo, error) {
	// This would require querying all quorums, which is expensive
	// For now, return an error indicating it's not supported
	return nil, fmt.Errorf("GetIndexedOperators is not supported; use GetIndexedOperatorState with specific quorums")
}

// Close closes the underlying contract client
func (s *IndexedChainState) Close() {
	s.contractClient.Close()
}

// Delegate core.ChainState interface methods to the underlying chainState implementation

// GetCurrentBlockNumber implements core.ChainState.GetCurrentBlockNumber
func (s *IndexedChainState) GetCurrentBlockNumber(ctx context.Context) (uint, error) {
	blockNumber, err := s.chainState.GetCurrentBlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get current block number: %w", err)
	}
	return blockNumber, nil
}

// GetOperatorState implements core.ChainState.GetOperatorState
func (s *IndexedChainState) GetOperatorState(
	ctx context.Context, blockNumber uint, quorums []core.QuorumID,
) (*core.OperatorState, error) {
	operatorState, err := s.chainState.GetOperatorState(ctx, blockNumber, quorums)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator state: %w", err)
	}
	return operatorState, nil
}

// GetOperatorStateWithSocket implements core.ChainState.GetOperatorStateWithSocket
func (s *IndexedChainState) GetOperatorStateWithSocket(
	ctx context.Context, blockNumber uint, quorums []core.QuorumID,
) (*core.OperatorState, error) {
	operatorState, err := s.chainState.GetOperatorStateWithSocket(ctx, blockNumber, quorums)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator state with socket: %w", err)
	}
	return operatorState, nil
}

// GetOperatorStateByOperator implements core.ChainState.GetOperatorStateByOperator
func (s *IndexedChainState) GetOperatorStateByOperator(
	ctx context.Context, blockNumber uint, operator core.OperatorID,
) (*core.OperatorState, error) {
	operatorState, err := s.chainState.GetOperatorStateByOperator(ctx, blockNumber, operator)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator state by operator: %w", err)
	}
	return operatorState, nil
}

// GetOperatorSocket implements core.ChainState.GetOperatorSocket
func (s *IndexedChainState) GetOperatorSocket(
	ctx context.Context, blockNumber uint, operator core.OperatorID,
) (string, error) {
	socket, err := s.chainState.GetOperatorSocket(ctx, blockNumber, operator)
	if err != nil {
		return "", fmt.Errorf("failed to get operator socket: %w", err)
	}
	return socket, nil
}

// getOperatorPubkeys retrieves the G1 and G2 public keys for an operator from the BLS APK Registry
func (s *IndexedChainState) getOperatorPubkeys(
	ctx context.Context, 
	operatorAddress gethcommon.Address, 
	blockNumber uint,
) (*core.G1Point, *core.G2Point, error) {
	// Prepare call options with specific block number
	callOpts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: new(big.Int).SetUint64(uint64(blockNumber)),
	}

	// Get the registered pubkey from BLS APK Registry
	// This returns (G1Point, pubkeyHash)
	g1Point, _, err := s.blsApkRegistry.GetRegisteredPubkey(callOpts, operatorAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get registered pubkey: %w", err)
	}

	// Convert contract G1Point to core.G1Point
	pubkeyG1 := &core.G1Point{
		G1Affine: &bn254.G1Affine{},
	}
	pubkeyG1.X.SetBytes(g1Point.X.Bytes())
	pubkeyG1.Y.SetBytes(g1Point.Y.Bytes())

	// For G2 pubkey, we need to query the operator's full registration data
	// The BLS APK Registry doesn't expose G2 pubkeys directly in a simple query,
	// so we'll create a placeholder G2 for now. In a full implementation,
	// we would need to listen to NewPubkeyRegistration events or query
	// the operator's original registration transaction.
	pubkeyG2 := &core.G2Point{
		G2Affine: &bn254.G2Affine{},
	}
	// TODO: Implement G2 pubkey resolution from registration events or storage

	return pubkeyG1, pubkeyG2, nil
}
