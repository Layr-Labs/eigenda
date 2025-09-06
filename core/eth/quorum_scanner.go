package eth

import (
	"context"
	"fmt"
	"math/big"

	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARegistryCoordinator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// A utility that is capable of producing a list of all registered quorums.
//
// This utility is fully thread safe, but it does not cache results. Calling into this utility multiple times
// with the same reference block number is wasteful.
type QuorumScanner interface {

	// Get all quorums registered at the given reference block number. Quorums are returned
	// sorted from least to greatest.
	GetQuorums(ctx context.Context, referenceBlockNumber uint64) ([]core.QuorumID, error)
}

var _ QuorumScanner = (*quorumScanner)(nil)

type quorumScanner struct {
	// A handle for communicating with the registry coordinator contract.
	registryCoordinator *regcoordinator.ContractEigenDARegistryCoordinator
}

// Create a new QuorumScanner instance.
func NewQuorumScanner(
	contractBackend bind.ContractBackend,
	registryCoordinatorAddress gethcommon.Address,
) (QuorumScanner, error) {

	registryCoordinator, err := regcoordinator.NewContractEigenDARegistryCoordinator(
		registryCoordinatorAddress,
		contractBackend)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry coordinator client: %w", err)
	}

	return &quorumScanner{
		registryCoordinator: registryCoordinator,
	}, nil
}

func (q *quorumScanner) GetQuorums(ctx context.Context, referenceBlockNumber uint64) ([]core.QuorumID, error) {
	// Quorums are assigned starting at 0, and then sequentially without gaps. If we
	// know the number of quorums, we can generate a list of quorum IDs.

	quorumCount, err := q.registryCoordinator.QuorumCount(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: new(big.Int).SetUint64(referenceBlockNumber),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get quorum count: %w", err)
	}

	quorums := make([]core.QuorumID, quorumCount)
	for i := uint8(0); i < quorumCount; i++ {
		quorums[i] = i
	}

	return quorums, nil
}
