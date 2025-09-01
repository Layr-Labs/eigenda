package operatorstate

import (
	"context"
	"fmt"
	"math/big"

	contractOperatorStateRetriever "github.com/Layr-Labs/eigenda/contracts/bindings/OperatorStateRetriever"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ContractClient provides methods to query operator state directly from on-chain contracts
type ContractClient struct {
	ethClient                  *ethclient.Client
	registryCoordinatorAddress gethcommon.Address
	operatorStateRetriever     *contractOperatorStateRetriever.ContractOperatorStateRetriever
}

// NewContractClient creates a new operator state contract client
func NewContractClient(
	ethRpcUrl string,
	registryCoordinatorAddress gethcommon.Address,
	operatorStateRetrieverAddress gethcommon.Address,
) (*ContractClient, error) {
	ethClient, err := ethclient.Dial(ethRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	if operatorStateRetrieverAddress == (gethcommon.Address{}) {
		return nil, fmt.Errorf("operator-state-retriever-address is required")
	}

	operatorStateRetriever, err := contractOperatorStateRetriever.NewContractOperatorStateRetriever(
		operatorStateRetrieverAddress,
		ethClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create operator state retriever contract binding: %w", err)
	}

	return &ContractClient{
		ethClient:                  ethClient,
		registryCoordinatorAddress: registryCoordinatorAddress,
		operatorStateRetriever:     operatorStateRetriever,
	}, nil
}

// OperatorStateResult represents the result of GetOperatorStateWithSocket
type OperatorStateResult struct {
	Operators [][]contractOperatorStateRetriever.OperatorStateRetrieverOperator
	Sockets   [][]string
}

// GetOperatorStateWithSocket retrieves the operator state for specific quorums at a given block
// Returns the raw contract result which includes operators and sockets for each quorum
func (c *ContractClient) GetOperatorStateWithSocket(
	ctx context.Context, quorums []byte, blockNumber uint64,
) (*OperatorStateResult, error) {
	// Determine the actual block number to query
	var actualBlockNumber uint64
	if blockNumber == 0 {
		// Get latest block
		latestBlock, err := c.ethClient.BlockNumber(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest block number: %w", err)
		}
		actualBlockNumber = latestBlock
	} else {
		actualBlockNumber = blockNumber
	}

	// Prepare call options with specific block number
	callOpts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: new(big.Int).SetUint64(actualBlockNumber),
	}

	// Query operator state with sockets for better information
	result, err := c.operatorStateRetriever.GetOperatorStateWithSocket(
		callOpts,
		c.registryCoordinatorAddress,
		quorums,
		uint32(actualBlockNumber),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator state: %w", err)
	}

	return &OperatorStateResult{
		Operators: result.Operators,
		Sockets:   result.Sockets,
	}, nil
}

// OperatorStateByIdResult represents the result of GetOperatorStateWithSocket0
type OperatorStateByIdResult struct {
	QuorumBitmap *big.Int
	Operators    [][]contractOperatorStateRetriever.OperatorStateRetrieverOperator
	Sockets      [][]string
}

// GetOperatorStateByOperatorId retrieves information about a specific operator by ID
func (c *ContractClient) GetOperatorStateByOperatorId(
	ctx context.Context, operatorId [32]byte, blockNumber uint64,
) (*OperatorStateByIdResult, error) {
	// Determine the actual block number to query
	var actualBlockNumber uint64
	if blockNumber == 0 {
		// Get latest block
		latestBlock, err := c.ethClient.BlockNumber(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest block number: %w", err)
		}
		actualBlockNumber = latestBlock
	} else {
		actualBlockNumber = blockNumber
	}

	// Prepare call options with specific block number
	callOpts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: new(big.Int).SetUint64(actualBlockNumber),
	}

	// Query operator state by operator ID
	result, err := c.operatorStateRetriever.GetOperatorStateWithSocket0(
		callOpts,
		c.registryCoordinatorAddress,
		operatorId,
		uint32(actualBlockNumber),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator state by ID: %w", err)
	}

	return &OperatorStateByIdResult{
		QuorumBitmap: result.QuorumBitmap,
		Operators:    result.Operators,
		Sockets:      result.Sockets,
	}, nil
}

// GetCurrentBlockNumber returns the current block number from the Ethereum client
func (c *ContractClient) GetCurrentBlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := c.ethClient.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get current block number: %w", err)
	}
	return blockNumber, nil
}

// Close closes the Ethereum client connection
func (c *ContractClient) Close() {
	if c.ethClient != nil {
		c.ethClient.Close()
	}
}
