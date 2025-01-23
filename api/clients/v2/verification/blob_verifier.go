package verification

import (
	"context"
	"fmt"

	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/Layr-Labs/eigenda/common"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// BlobVerifier is responsible for making eth calls against the BlobVerifier contract to ensure cryptographic and
// structural integrity of V2 certificates
//
// The blob verifier contract is located at https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDABlobVerifier.sol
type BlobVerifier struct {
	// go binding around the EigenDABlobVerifier ethereum contract
	blobVerifierCaller *verifierBindings.ContractEigenDABlobVerifierCaller
}

// NewBlobVerifier constructs a BlobVerifier
func NewBlobVerifier(
	ethClient *common.EthClient, // the eth client, which should already be set up
	blobVerifierAddress string, // the hex address of the EigenDABlobVerifier contract
) (*BlobVerifier, error) {

	verifierCaller, err := verifierBindings.NewContractEigenDABlobVerifierCaller(
		gethcommon.HexToAddress(blobVerifierAddress),
		*ethClient)

	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %s", blobVerifierAddress, err)
	}

	return &BlobVerifier{
		blobVerifierCaller: verifierCaller,
	}, nil
}

// VerifyBlobV2FromSignedBatch calls the verifyBlobV2FromSignedBatch view function on the EigenDABlobVerifier contract
//
// This method returns nil if the blob is successfully verified. Otherwise, it returns an error.
func (v *BlobVerifier) VerifyBlobV2FromSignedBatch(
	ctx context.Context,
	// The signed batch that contains the blob being verified. This is obtained from the disperser, and is used
	// to verify that the described blob actually exists in a valid batch.
	signedBatch *disperser.SignedBatch,
	// Contains all necessary information about the blob, so that it can be verified.
	blobVerificationProof *disperser.BlobInclusionInfo,
) error {
	convertedSignedBatch, err := verifierBindings.ConvertSignedBatch(signedBatch)
	if err != nil {
		return fmt.Errorf("convert signed batch: %s", err)
	}

	convertedBlobVerificationProof, err := verifierBindings.ConvertVerificationProof(blobVerificationProof)
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

// VerifyBlobV2 calls the VerifyBlobV2 view function on the EigenDABlobVerifier contract
//
// This method returns nil if the blob is successfully verified. Otherwise, it returns an error.
func (v *BlobVerifier) VerifyBlobV2(
	ctx context.Context,
	// The header of the batch that the blob is contained in
	batchHeader *commonv2.BatchHeader,
	// Contains data pertaining to the blob's inclusion in the batch
	blobVerificationProof *disperser.BlobInclusionInfo,
	// Contains data that can be used to verify that the blob actually exists in the claimed batch
	nonSignerStakesAndSignature verifierBindings.NonSignerStakesAndSignature,
) error {
	convertedBatchHeader, err := verifierBindings.ConvertBatchHeader(batchHeader)
	if err != nil {
		return fmt.Errorf("convert batch header: %s", err)
	}

	convertedBlobVerificationProof, err := verifierBindings.ConvertVerificationProof(blobVerificationProof)
	if err != nil {
		return fmt.Errorf("convert blob verification proof: %s", err)
	}

	err = v.blobVerifierCaller.VerifyBlobV2(
		&bind.CallOpts{Context: ctx},
		*convertedBatchHeader,
		*convertedBlobVerificationProof,
		nonSignerStakesAndSignature)

	if err != nil {
		return fmt.Errorf("verify blob v2: %s", err)
	}

	return nil
}
