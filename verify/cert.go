package verify

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"

	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/exp/slices"
)

var ErrBatchMetadataHashNotFound = errors.New("BatchMetadataHash not found for BatchId")

// CertVerifier verifies the DA certificate against on-chain EigenDA contracts
// to ensure disperser returned fields haven't been tampered with
type CertVerifier struct {
	ethConfirmationDepth uint64
	manager              *binding.ContractEigenDAServiceManagerCaller
	finalizedBlockClient *FinalizedBlockClient
	ethClient            *ethclient.Client
}

func NewCertVerifier(cfg *Config, l log.Logger) (*CertVerifier, error) {
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ETH RPC node: %s", err.Error())
	}

	// construct binding
	m, err := binding.NewContractEigenDAServiceManagerCaller(common.HexToAddress(cfg.SvcManagerAddr), client)
	if err != nil {
		return nil, err
	}

	return &CertVerifier{
		manager:              m,
		finalizedBlockClient: NewFinalizedBlockClient(client.Client()),
		ethConfirmationDepth: cfg.EthConfirmationDepth,
		ethClient:            client,
	}, nil
}

func (cv *CertVerifier) VerifyBatch(header *binding.IEigenDAServiceManagerBatchHeader,
	id uint32, recordHash [32]byte, blockNum uint32) error {
	// 0 - Determine block context number
	blockNumber, err := cv.getContextBlock()
	if err != nil {
		return err
	}

	// 1 - Verify batch hash

	// 1.a - ensure that a batch hash can be looked up for a batch ID
	expectedHash, err := cv.manager.BatchIdToBatchMetadataHash(&bind.CallOpts{BlockNumber: blockNumber}, id)
	if err != nil {
		return err
	}
	if bytes.Equal(expectedHash[:], make([]byte, 32)) {
		return ErrBatchMetadataHashNotFound
	}

	// 1.b - ensure that hash generated from local cert matches one stored on-chain
	actualHash, err := HashBatchMetadata(header, recordHash, blockNum)

	if err != nil {
		return err
	}

	equal := slices.Equal(expectedHash[:], actualHash[:])
	if !equal {
		return fmt.Errorf("batch hash mismatch, expected: %x, got: %x", expectedHash, actualHash)
	}

	return nil
}

// VerifyMerkleProof
func (cv *CertVerifier) VerifyMerkleProof(inclusionProof []byte, root []byte, blobIndex uint32, blobHeader BlobHeader) error {
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

// 3 - (TODO) verify blob security params
func (cv *CertVerifier) VerifyBlobParams(inclusionProof []byte, rootHash []byte, leafHash []byte, index uint64) error {
	return nil
}

func (cv *CertVerifier) getContextBlock() (*big.Int, error) {
	var blockNumber *big.Int
	if cv.ethConfirmationDepth == 0 {
		// Get the latest finalized block
		blockHeader, err := cv.finalizedBlockClient.GetBlock(context.Background(), "finalized", false)
		if err != nil {
			return nil, err
		}
		blockNumber = blockHeader.Number()
	} else {
		blockHeader, err := cv.ethClient.BlockByNumber(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		blockNumber = new(big.Int)
		blockNumber.Sub(blockHeader.Number(), big.NewInt(int64(cv.ethConfirmationDepth-1)))
	}
	return blockNumber, nil
}
