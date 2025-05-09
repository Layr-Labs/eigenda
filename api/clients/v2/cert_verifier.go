package clients

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV2"
	cert_type_binding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
)

// IV2CertVerifier is an interface for interacting with the EigenDACertVerifier contract.
// NOTE: This is legacy code and should be removed in the future. Currently its ran on a few
//       EigenDA blazzar testnets but should be nuked before mainnet.
type IV2CertVerifier interface {
	// VerifyCertV2 calls the VerifyCertV2 view function on the EigenDACertVerifier contract.
	//
	// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
	VerifyCertV2(ctx context.Context, eigenDACert *coretypes.EigenDACertV2) error

	// GetNonSignerStakesAndSignature calls the getNonSignerStakesAndSignature view function on the EigenDACertVerifier
	// contract, and returns the resulting NonSignerStakesAndSignature object.
	GetNonSignerStakesAndSignature(
		ctx context.Context,
		signedBatch *disperser.SignedBatch,
	) (*verifierBindings.EigenDATypesV1NonSignerStakesAndSignature, error)

	// GetQuorumNumbersRequired queries the cert verifier contract for the configured set of quorum numbers that must
	// be set in the BlobHeader, and verified in VerifyDACertV2 and verifyDACertV2FromSignedBatch
	GetQuorumNumbersRequired(ctx context.Context) ([]uint8, error)

	// GetConfirmationThreshold queries the cert verifier contract for the configured ConfirmationThreshold.
	// The ConfirmationThreshold is an integer value between 0 and 100 (inclusive), where the value represents
	// a percentage of validator stake that needs to have signed for availability, for the blob to be considered
	// "available".
	GetConfirmationThreshold(ctx context.Context, referenceBlockNumber uint64) (uint8, error)
}

// IGenericCertVerifier is an interface for interacting with the updated EigenDACertVerifier contract
// that takes a low-level bytes interface and is compatible with the EigenDACertVerifierRouter contract.
type IGenericCertVerifier interface {
	// CheckDACert calls the CheckDACert view function on the EigenDACertVerifier contract.
	//
	// This method returns nil if the cert is successfully verified. Otherwise, it returns an error.
	CheckDACert(ctx context.Context, rbn uint64, daCert []byte) error

	// GetNonSignerStakesAndSignature calls the getNonSignerStakesAndSignature view function on the EigenDACertVerifier
	// contract, and returns the resulting NonSignerStakesAndSignature object.
	GetNonSignerStakesAndSignature(
		ctx context.Context,
		signedBatch *disperser.SignedBatch,
	) (*cert_type_binding.EigenDATypesV1NonSignerStakesAndSignature, error)

	// GetQuorumNumbersRequired queries the cert verifier contract for the configured set of quorum numbers that must
	// be set in the BlobHeader, and verified in VerifyDACertV2 and verifyDACertV2FromSignedBatch
	GetQuorumNumbersRequired(ctx context.Context, referenceBlockNumber uint64) ([]uint8, error)

	// GetConfirmationThreshold queries the cert verifier contract for the configured ConfirmationThreshold.
	// The ConfirmationThreshold is an integer value between 0 and 100 (inclusive), where the value represents
	// a percentage of validator stake that needs to have signed for availability, for the blob to be considered
	// "available".
	GetConfirmationThreshold(ctx context.Context, referenceBlockNumber uint64) (uint8, error)

	// GetCertVersion queries the cert verifier contract for the configured certificate version.
	// The version is an integer value that represents the version of the certificate and is used
	// for different struct encodings.
	GetCertVersion(ctx context.Context, referenceBlockNumber uint64) (uint8, error)
}
