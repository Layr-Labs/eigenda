package verification

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common/geth"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ICertVerifier is the interface representing a CertVerifier
//
// This interface exists in order to allow verification mocking in unit tests.
type ICertVerifier interface {
	VerifyCertV2(
		ctx context.Context,
		eigenDACert *EigenDACert,
	) error

	GetNonSignerStakesAndSignature(
		ctx context.Context,
		signedBatch *disperser.SignedBatch,
	) (*verifierBindings.NonSignerStakesAndSignature, error)
}

// CertVerifier is responsible for making eth calls against the CertVerifier contract to ensure cryptographic and
// structural integrity of V2 certificates
//
// The cert verifier contract is located at https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDABlobVerifier.sol
type CertVerifier struct {
	// go binding around the EigenDACertVerifier ethereum contract
	certVerifierCaller *verifierBindings.ContractEigenDABlobVerifierCaller
}

var _ ICertVerifier = &CertVerifier{}

// NewCertVerifier constructs a CertVerifier
func NewCertVerifier(
	ethClient geth.EthClient, // the eth client, which should already be set up
	certVerifierAddress string, // the hex address of the EigenDACertVerifier contract
) (*CertVerifier, error) {

	verifierCaller, err := verifierBindings.NewContractEigenDABlobVerifierCaller(
		gethcommon.HexToAddress(certVerifierAddress),
		ethClient)

	if err != nil {
		return nil, fmt.Errorf("bind to verifier contract at %s: %w", certVerifierAddress, err)
	}

	return &CertVerifier{
		certVerifierCaller: verifierCaller,
	}, nil
}

// VerifyCertV2FromSignedBatch calls the verifyCertV2FromSignedBatch view function on the EigenDACertVerifier contract
//
// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
func (cv *CertVerifier) VerifyCertV2FromSignedBatch(
	ctx context.Context,
	// The signed batch that contains the blob whose cert is being verified. This is obtained from the disperser, and
	// is used to verify that the described blob actually exists in a valid batch.
	signedBatch *disperser.SignedBatch,
	// Contains all necessary information about the blob, so that the cert can be verified.
	blobInclusionInfo *disperser.BlobInclusionInfo,
) error {
	convertedSignedBatch, err := SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return fmt.Errorf("convert signed batch: %w", err)
	}

	convertedBlobInclusionInfo, err := InclusionInfoProtoToBinding(blobInclusionInfo)
	if err != nil {
		return fmt.Errorf("convert blob inclusion info: %w", err)
	}

	err = cv.certVerifierCaller.VerifyBlobV2FromSignedBatch(
		&bind.CallOpts{Context: ctx},
		*convertedSignedBatch,
		*convertedBlobInclusionInfo)

	if err != nil {
		return fmt.Errorf("verify cert v2 from signed batch: %w", err)
	}

	return nil
}

// VerifyCertV2 calls the VerifyCertV2 view function on the EigenDACertVerifier contract
//
// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
func (cv *CertVerifier) VerifyCertV2(
	ctx context.Context,
	eigenDACert *EigenDACert,
) error {
	err := cv.certVerifierCaller.VerifyBlobV2(
		&bind.CallOpts{Context: ctx},
		eigenDACert.BatchHeader,
		eigenDACert.BlobInclusionInfo,
		eigenDACert.NonSignerStakesAndSignature)

	if err != nil {
		return fmt.Errorf("verify cert v2: %w", err)
	}

	return nil
}

// GetNonSignerStakesAndSignature calls the getNonSignerStakesAndSignature view function on the EigenDACertVerifier
// contract, and returns the resulting NonSignerStakesAndSignature object.
func (cv *CertVerifier) GetNonSignerStakesAndSignature(
	ctx context.Context,
	signedBatch *disperser.SignedBatch,
) (*verifierBindings.NonSignerStakesAndSignature, error) {

	signedBatchBinding, err := SignedBatchProtoToBinding(signedBatch)
	if err != nil {
		return nil, fmt.Errorf("convert signed batch: %w", err)
	}

	nonSignerStakesAndSignature, err := cv.certVerifierCaller.GetNonSignerStakesAndSignature(
		&bind.CallOpts{Context: ctx},
		*signedBatchBinding)

	if err != nil {
		return nil, fmt.Errorf("get non signer stakes and signature: %w", err)
	}

	return &nonSignerStakesAndSignature, nil
}
