package clients

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
)

// ICertVerifier is an interface for interacting with the EigenDACertVerifier contract.
type ICertVerifier interface {
	// VerifyCertV2 calls the VerifyCertV2 view function on the EigenDACertVerifier contract.
	//
	// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
	VerifyCertV2(ctx context.Context, eigenDACert *coretypes.EigenDACert) error

	// GetNonSignerStakesAndSignature calls the getNonSignerStakesAndSignature view function on the EigenDACertVerifier
	// contract, and returns the resulting NonSignerStakesAndSignature object.
	GetNonSignerStakesAndSignature(
		ctx context.Context,
		signedBatch *disperser.SignedBatch,
	) (*verifierBindings.NonSignerStakesAndSignature, error)

	// GetQuorumNumbersRequired queries the cert verifier contract for the configured set of quorum numbers that must
	// be set in the BlobHeader, and verified in VerifyDACertV2 and verifyDACertV2FromSignedBatch
	GetQuorumNumbersRequired(ctx context.Context) ([]uint8, error)
}
