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
	l                    log.Logger
	ethConfirmationDepth uint64
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

// verifies on-chain batch ID for equivalence to certificate batch header fields
func (cv *CertVerifier) VerifyBatch(header *binding.IEigenDAServiceManagerBatchHeader,
	id uint32, recordHash [32]byte, confirmationNumber uint32) error {
	blockNumber, err := cv.getContextBlock()
	if err != nil {
		return err
	}

	// 1. ensure that a batch hash can be looked up for a batch ID for a given block number
	expectedHash, err := cv.manager.BatchIdToBatchMetadataHash(&bind.CallOpts{BlockNumber: blockNumber}, id)
	if err != nil {
		return err
	}
	if bytes.Equal(expectedHash[:], make([]byte, 32)) {
		return ErrBatchMetadataHashNotFound
	}

	// 2. ensure that hash generated from local cert matches one stored on-chain
	actualHash, err := HashBatchMetadata(header, recordHash, confirmationNumber)

	if err != nil {
		return err
	}

	equal := slices.Equal(expectedHash[:], actualHash[:])
	if !equal {
		return fmt.Errorf("batch hash mismatch, expected: %x, got: %x", expectedHash, actualHash)
	}

	return nil
}

// verifies the blob batch inclusion proof against the blob root hash
func (cv *CertVerifier) VerifyMerkleProof(inclusionProof []byte, root []byte,
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
func (cv *CertVerifier) getContextBlock() (*big.Int, error) {
	var blockNumber *big.Int
	blockHeader, err := cv.ethClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	if cv.ethConfirmationDepth == 0 {
		return blockHeader.Number(), nil
	}

	blockNumber = new(big.Int)
	blockNumber.Sub(blockHeader.Number(), big.NewInt(int64(cv.ethConfirmationDepth-1))) // #nosec G115

	return blockNumber, nil
}
