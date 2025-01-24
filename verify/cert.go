package verify

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common/consts"
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
	l log.Logger
	// ethConfirmationDepth is used to verify that a blob's batch commitment has been bridged to the EigenDAServiceManager contract at least
	// this many blocks in the past. To do so we make an eth_call to the contract at the current block_number - ethConfirmationDepth.
	// Hence in order to not require an archive node, this value should be kept low. We force it to be < 64 (consts.EthHappyPathFinalizationDepthBlocks).
	// waitForFinalization should be used instead of ethConfirmationDepth if the user wants to wait for finality (typically 64 blocks in happy case).
	ethConfirmationDepth uint64
	waitForFinalization  bool
	manager              *binding.ContractEigenDAServiceManagerCaller
	ethClient            *ethclient.Client
	// The two fields below are fetched from the EigenDAServiceManager contract in the constructor.
	// They are used to verify the quorums in the received certificates.
	// See getQuorumParametersAtLatestBlock for more details.
	quorumsRequired           []uint8
	quorumAdversaryThresholds map[uint8]uint8
}

func NewCertVerifier(cfg *Config, l log.Logger) (*CertVerifier, error) {
	if cfg.EthConfirmationDepth >= uint64(consts.EthHappyPathFinalizationDepthBlocks) {
		// We keep this low (<128) to avoid requiring an archive node.
		return nil, fmt.Errorf("confirmation depth must be less than 64; consider using cfg.WaitForFinalization=true instead")
	}
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

	quorumsRequired, quorumAdversaryThresholds, err := getQuorumParametersAtLatestBlock(m)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quorum parameters from EigenDAServiceManager: %w", err)
	}

	return &CertVerifier{
		l:                         l,
		manager:                   m,
		ethConfirmationDepth:      cfg.EthConfirmationDepth,
		ethClient:                 client,
		quorumsRequired:           quorumsRequired,
		quorumAdversaryThresholds: quorumAdversaryThresholds,
	}, nil
}

// verifyBatchConfirmedOnChain verifies that batchMetadata (typically part of a received cert)
// matches the batch metadata hash stored on-chain
func (cv *CertVerifier) verifyBatchConfirmedOnChain(
	ctx context.Context, batchID uint32, batchMetadata *disperser.BatchMetadata,
) error {
	// 1. Verify that the confirmation status has been reached.
	// The eigenda-client already checks for this, but it is possible for either
	//  1. a reorg to happen, causing the batch to be confirmed by fewer number of blocks than required
	//  2. proxy's node is behind the eigenda_client's node that deemed the batch confirmed, or
	//     even if we use the same url, that the connection drops and we get load-balanced to a different eth node.
	// We retry up to 60 seconds (allowing for reorgs up to 5 blocks deep), but we only wait 3 seconds between each retry,
	// in case (2) is the case and the node simply needs to resync, which could happen fast.
	//
	// Note that we don't verify that the batch is actually onchain at the batchMetadata's state confirmedBlockNumber, because that would require an archive node.
	// This is super unlikely if the disperser is honest, but it could technically happen that a confirmed batch's block gets reorged out,
	// yet the tx is included in an earlier or later block, making the batchMetadata received from the disperser
	// no longer valid. The eigenda batcher does check for these reorgs and updates the batch's confirmation block number:
	// https://github.com/Layr-Labs/eigenda/blob/bee55ed9207f16153c3fd8ebf73c219e68685def/disperser/batcher/finalizer.go#L198
	// confirmedBlockNum                               currentBlock-confirmationDepth        currentBlock
	// | (don't verify here, need archive node)                    | (verify here)               |
	// +-----------------------------------------------------------+-----------------------------+
	onchainHash, err := retry.Do(ctx, 20, retry.Fixed(3*time.Second), func() ([32]byte, error) {
		blockNumber, err := cv.getConfDeepBlockNumber(ctx)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to get context block: %w", err)
		}
		return cv.retrieveBatchMetadataHash(ctx, batchID, blockNumber)
	})
	if err != nil {
		return fmt.Errorf("retrieving batch that was confirmed at block %v: %w", batchMetadata.GetConfirmationBlockNumber(), err)
	}

	// 2. Compute the hash of the batch metadata received as argument.
	header := &binding.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:       [32]byte(batchMetadata.GetBatchHeader().GetBatchRoot()),
		QuorumNumbers:         batchMetadata.GetBatchHeader().GetQuorumNumbers(),
		ReferenceBlockNumber:  batchMetadata.GetBatchHeader().GetReferenceBlockNumber(),
		SignedStakeForQuorums: batchMetadata.GetBatchHeader().GetQuorumSignedPercentages(),
	}
	recordHash := [32]byte(batchMetadata.GetSignatoryRecordHash())
	computedHash, err := HashBatchMetadata(header, recordHash, batchMetadata.GetConfirmationBlockNumber())
	if err != nil {
		return fmt.Errorf("failed to hash batch metadata: %w", err)
	}

	// 3. Ensure that hash generated from local cert matches one stored on-chain.
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
// returns an error if some problem calling the contract happens, or the hash is not found.
// We make an eth_call to the EigenDAServiceManager at the given blockNumber to retrieve the hash.
// Therefore, make sure that blockNumber is <128 blocks behind the latest block, to avoid requiring an archive node.
// This is currently enforced by having EthConfirmationDepth be <64.
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

// getQuorumParametersAtLatestBlock fetches the required quorums and quorum adversary thresholds
// from the EigenDAServiceManager contract at the latest block.
// We then cache these parameters and use them in the Verifier to verify the certificates.
//
// Note: this strategy (fetching once and caching) only works because these parameters are immutable.
// They might be different in different environments (e.g. on a devnet or testnet), but they are fixed on a given network.
// We used to allow these parameters to change (via a setter function on the contract), but that then forced us here in the proxy
// to query for these parameters on every request, at the batch's reference block number (RBN).
// This in turn required rollup validators running this proxy to have an archive node, in case the RBN was >128 blocks in the past,
// which was not ideal. So we decided to make these parameters immutable, and cache them here.
func getQuorumParametersAtLatestBlock(
	manager *binding.ContractEigenDAServiceManagerCaller,
) ([]uint8, map[uint8]uint8, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	requiredQuorums, err := manager.QuorumNumbersRequired(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch QuorumNumbersRequired from EigenDAServiceManager: %w", err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	thresholds, err := manager.QuorumAdversaryThresholdPercentages(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch QuorumAdversaryThresholdPercentages from EigenDAServiceManager: %w", err)
	}
	if len(thresholds) > math.MaxUint8 {
		return nil, nil, fmt.Errorf("thresholds received from EigenDAServiceManager contains %d > 256 quorums, which isn't possible", len(thresholds))
	}
	var quorumAdversaryThresholds = make(map[uint8]uint8)
	for quorumNum, threshold := range thresholds {
		quorumAdversaryThresholds[uint8(quorumNum)] = threshold //nolint:gosec // disable G115 // We checked the length of thresholds above
	}
	// Sanity check: ensure that the required quorums are a subset of the quorums for which we have adversary thresholds
	for _, quorum := range requiredQuorums {
		if _, ok := quorumAdversaryThresholds[quorum]; !ok {
			return nil, nil, fmt.Errorf("required quorum %d does not have an adversary threshold. Was the EigenDAServiceManager properly deployed?", quorum)
		}
	}
	return requiredQuorums, quorumAdversaryThresholds, nil
}
