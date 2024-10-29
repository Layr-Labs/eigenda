package verify

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"

	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/exp/slices"
)

// CertVerifier verifies the DA certificate against on-chain EigenDA contracts
// to ensure disperser returned fields haven't been tampered with
type CertVerifier struct {
	l                    log.Logger
	ethConfirmationDepth uint64
	waitForFinalization  bool
	manager              *binding.ContractEigenDAServiceManagerCaller
	ethClient            *ethclient.Client
}

func NewCertVerifier(cfg *Config, l log.Logger) (*CertVerifier, error) {
	log.Info("Enabling certificate verification", "confirmation_depth", cfg.EthConfirmationDepth)

	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ETH RPC node: %s", err.Error())
	}

	// construct caller binding
	m, err := binding.NewContractEigenDAServiceManagerCaller(common.HexToAddress(cfg.SvcManagerAddr), client)
	if err != nil {
		return nil, err
	}

	return &CertVerifier{
		l:                    l,
		manager:              m,
		ethConfirmationDepth: cfg.EthConfirmationDepth,
		ethClient:            client,
	}, nil
}

// verifyBatchConfirmedOnChain verifies that batchMetadata (typically part of a received cert)
// matches the batch metadata hash stored on-chain
func (cv *CertVerifier) verifyBatchConfirmedOnChain(
	ctx context.Context, batchID uint32, batchMetadata *disperser.BatchMetadata,
) error {
	// 1. Verify batch is actually onchain at the batchMetadata's state confirmedBlockNumber.
	// This is super unlikely if the disperser is honest, but it could technically happen that a confirmed batch's block gets reorged out,
	// yet the tx is included in an earlier or later block, making the batchMetadata received from the disperser
	// no longer valid. The eigenda batcher does check for these reorgs and updates the batch's confirmation block number:
	// https://github.com/Layr-Labs/eigenda/blob/bee55ed9207f16153c3fd8ebf73c219e68685def/disperser/batcher/finalizer.go#L198
	// TODO: We could require the disperser for the new batch, or try to reconstruct it ourselves by querying the chain,
	// but for now we opt to simply fail the verification, which will force the batcher to resubmit the batch to eigenda.
	confirmationBlockNumber := batchMetadata.GetConfirmationBlockNumber()
	confirmationBlockNumberBigInt := big.NewInt(0).SetInt64(int64(confirmationBlockNumber))
	_, err := cv.retrieveBatchMetadataHash(ctx, batchID, confirmationBlockNumberBigInt)
	if err != nil {
		return fmt.Errorf("batch not found onchain at supposedly confirmed block %d: %w", confirmationBlockNumber, err)
	}

	// 2. Verify that the confirmation status has been reached.
	// The eigenda-client already checks for this, but it is possible for either
	//  1. a reorg to happen, causing the batch to be confirmed by fewer number of blocks than required
	//  2. proxy's node is behind the eigenda_client's node that deemed the batch confirmed, or
	//     even if we use the same url, that the connection drops and we get load-balanced to a different eth node.
	// We retry up to 60 seconds (allowing for reorgs up to 5 blocks deep), but we only wait 3 seconds between each retry,
	// in case (2) is the case and the node simply needs to resync, which could happen fast.
	onchainHash, err := retry.Do(ctx, 20, retry.Fixed(3*time.Second), func() ([32]byte, error) {
		blockNumber, err := cv.getConfDeepBlockNumber(ctx)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to get context block: %w", err)
		}
		return cv.retrieveBatchMetadataHash(ctx, batchID, blockNumber)
	})
	if err != nil {
		return fmt.Errorf("retrieving batch that was confirmed at block %v: %w", confirmationBlockNumber, err)
	}

	// 3. Compute the hash of the batch metadata received as argument.
	header := &binding.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:       [32]byte(batchMetadata.GetBatchHeader().GetBatchRoot()),
		QuorumNumbers:         batchMetadata.GetBatchHeader().GetQuorumNumbers(),
		ReferenceBlockNumber:  batchMetadata.GetBatchHeader().GetReferenceBlockNumber(),
		SignedStakeForQuorums: batchMetadata.GetBatchHeader().GetQuorumSignedPercentages(),
	}
	recordHash := [32]byte(batchMetadata.GetSignatoryRecordHash())
	computedHash, err := HashBatchMetadata(header, recordHash, confirmationBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to hash batch metadata: %w", err)
	}

	// 4. Ensure that hash generated from local cert matches one stored on-chain.
	equal := slices.Equal(onchainHash[:], computedHash[:])
	if !equal {
		return fmt.Errorf("batch hash mismatch, onchain: %x, computed: %x", onchainHash, computedHash)
	}

	return nil
}

// verifies the blob batch inclusion proof against the blob root hash
func (cv *CertVerifier) verifyMerkleProof(inclusionProof []byte, root []byte,
	blobIndex uint32, blobHeader BlobHeader) error {
	leafHash, err := HashEncodeBlobHeader(blobHeader)
	if err != nil {
		return err
	}

	generatedRoot, err := ProcessInclusionProof(inclusionProof, leafHash, uint64(blobIndex))
	if err != nil {
		return err
	}

	equal := slices.Equal(root, generatedRoot.Bytes())
	if !equal {
		return fmt.Errorf("root hash mismatch, expected: %x, got: %x", root, generatedRoot)
	}

	return nil
}

// fetches a block number provided a subtraction of a user defined conf depth from latest block
func (cv *CertVerifier) getConfDeepBlockNumber(ctx context.Context) (*big.Int, error) {
	if cv.waitForFinalization {
		var header = types.Header{}
		// We ask for the latest finalized block. The second parameter "hydrated txs" is set to false because we don't need full txs.
		// See https://github.com/ethereum/execution-apis/blob/4140e528360fea53c34a766d86a000c6c039100e/src/eth/block.yaml#L61
		// This is equivalent to `cast block finalized`, as opposed to `cast block finalized --full`.
		err := cv.ethClient.Client().CallContext(ctx, &header, "eth_getBlockByNumber", "finalized", false)
		if err != nil {
			return nil, fmt.Errorf("failed to get finalized block: %w", err)
		}
		return header.Number, nil
	}
	blockNumber, err := cv.ethClient.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %w", err)
	}
	if blockNumber < cv.ethConfirmationDepth {
		return big.NewInt(0), nil
	}
	return new(big.Int).SetUint64(blockNumber - cv.ethConfirmationDepth), nil
}

// retrieveBatchMetadataHash retrieves the batch metadata hash stored on-chain at a specific blockNumber for a given batchID
// returns an error if some problem calling the contract happens, or the hash is not found
func (cv *CertVerifier) retrieveBatchMetadataHash(ctx context.Context, batchID uint32, blockNumber *big.Int) ([32]byte, error) {
	onchainHash, err := cv.manager.BatchIdToBatchMetadataHash(&bind.CallOpts{Context: ctx, BlockNumber: blockNumber}, batchID)
	if err != nil {
		return [32]byte{}, fmt.Errorf("calling EigenDAServiceManager.BatchIdToBatchMetadataHash: %w", err)
	}
	if bytes.Equal(onchainHash[:], make([]byte, 32)) {
		return [32]byte{}, fmt.Errorf("BatchMetadataHash not found for BatchId %d at block %d", batchID, blockNumber.Uint64())
	}
	return onchainHash, nil
}
