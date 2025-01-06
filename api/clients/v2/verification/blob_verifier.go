package verification

import (
	"context"
	"fmt"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlobVerifier is responsible for making eth calls to verify blobs that have been received by the client
//
// Blob verification is not threadsafe.
type BlobVerifier struct {
	// the eth client that calls will be made to
	ethClient *ethclient.Client
	// go binding around the verifyBlobV2FromSignedBatch ethereum contract
	blobVerifierCaller *verifierBindings.ContractEigenDABlobVerifierCaller
}

// NewBlobVerifier constructs a BlobVerifier
func NewBlobVerifier(
	// the eth RPC URL that will be connected to
	ethRpcUrl string,
	// the hex address of the verifyBlobV2FromSignedBatch contract
	verifyBlobV2FromSignedBatchAddress string) (*BlobVerifier, error) {

	ethClient, err := ethclient.Dial(ethRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("dial ETH RPC node: %s", err)
	}

	verifierCaller, err := verifierBindings.NewContractEigenDABlobVerifierCaller(
		ethcommon.HexToAddress(verifyBlobV2FromSignedBatchAddress),
		ethClient)

	if err != nil {
		ethClient.Close()

		return nil, fmt.Errorf("bind to verifier contract at %s: %s", verifyBlobV2FromSignedBatchAddress, err)
	}

	return &BlobVerifier{
		ethClient:          ethClient,
		blobVerifierCaller: verifierCaller,
	}, nil
}

// VerifyBlobV2FromSignedBatch makes a call to the verifyBlobV2FromSignedBatch contract
//
// This method returns nil if the blob is successfully verified. Otherwise, it returns an error.
//
// This method is not threadsafe.
func (v *BlobVerifier) VerifyBlobV2FromSignedBatch(
	ctx context.Context,
	// The signed batch that contains the blob being verified. This is obtained from the disperser, and is used
	// to verify that the described blob actually exists in a valid batch.
	signedBatch *disperser.SignedBatch,
	// Contains all necessary information about the blob, so that it can be verified.
	blobVerificationProof *disperser.BlobVerificationInfo,
) error {
	convertedSignedBatch, err := signedBatch.ToBinding()
	if err != nil {
		return fmt.Errorf("convert signed batch: %s", err)
	}

	convertedBlobVerificationProof, err := blobVerificationProof.ToBinding()
	if err != nil {
		return fmt.Errorf("convert blob verification proof: %s", err)
	}

	err = v.blobVerifierCaller.VerifyBlobV2FromSignedBatch(
		&bind.CallOpts{Context: ctx},
		*convertedSignedBatch,
		*convertedBlobVerificationProof)

	if err != nil {
		return fmt.Errorf("verify blob v2 from signed batch: %s", err)
	}

	return nil
}

// Close closes the eth client. This method is threadsafe.
func (v *BlobVerifier) Close() {
	v.ethClient.Close()
}
