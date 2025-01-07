package verification

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// BlobVerifier is responsible for making eth calls to verify blobs that have been received by the client
//
// Blob verification is not threadsafe.
type BlobVerifier struct {
	// go binding around the verifyBlobV2FromSignedBatch ethereum contract
	blobVerifierCaller *verifierBindings.ContractEigenDABlobVerifierCaller
}

// NewBlobVerifier constructs a BlobVerifier
func NewBlobVerifier(
	ethClient *common.EthClient,               // the eth client, which should already be set up
	verifyBlobV2FromSignedBatchAddress string, // the hex address of the verifyBlobV2FromSignedBatch contract
) (*BlobVerifier, error) {

	verifierCaller, err := verifierBindings.NewContractEigenDABlobVerifierCaller(
		gethcommon.HexToAddress(verifyBlobV2FromSignedBatchAddress),
		*ethClient)

	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %s", verifyBlobV2FromSignedBatchAddress, err)
	}

	return &BlobVerifier{
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
