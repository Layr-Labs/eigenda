package verify

import (
	"fmt"

	proxy_common "github.com/Layr-Labs/eigenda-proxy/common"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

// CertVerifier verifies the DA certificate against on-chain EigenDA contracts
// to ensure disperser returned fields haven't been tampered with
type CertVerifier struct {
	manager *binding.ContractEigenDAServiceManagerCaller
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
		manager: m,
	}, nil
}

func (cv *CertVerifier) VerifyBatch(header *binding.IEigenDAServiceManagerBatchHeader,
	id uint32, recordHash [32]byte, blockNum uint32) error {
	// 1 - Verify batch hash

	// 1.a - ensure that a batch hash can be looked up for a batch ID
	expectedHash, err := cv.manager.BatchIdToBatchMetadataHash(nil, id)
	if err != nil {
		return err
	}

	// 1.b - ensure that hash generated from local cert matches one stored on-chain

	actualHash, err := HashBatchMetadata(header, recordHash, blockNum)

	if err != nil {
		return err
	}

	equal := proxy_common.EqualSlices(expectedHash[:], actualHash[:])
	if !equal {
		return fmt.Errorf("batch hash mismatch, expected: %x, got: %x", expectedHash, actualHash)
	}

	return nil
}

// 2 - (TODO) merkle proof verification

func (cv *CertVerifier) VerifyMerkleProof(inclusionProof []byte, rootHash []byte, leafHash []byte, index uint64) error {
	return nil
}

// 3 - (TODO) verify blob security params
func (cv *CertVerifier) Verify(inclusionProof []byte, rootHash []byte, leafHash []byte, index uint64) error {
	return nil
}
