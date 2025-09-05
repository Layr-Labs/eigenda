package verification

// CheckDACertStatusCode represents the status codes that are returned by
// EigenDACertVerifier.checkDACert contract calls. The enum values below should match exactly
// the status codes defined in the contract:
// https://github.com/Layr-Labs/eigenda/blob/1091f460ba762b84019389cbb82d9b04bb2c2bdb/contracts/src/integrations/cert/libraries/EigenDACertVerificationLib.sol#L48-L54
type CheckDACertStatusCode uint8

// Since v3.1.0 of the CertVerifier, checkDACert calls are classified into: success (200), invalid_cert (400), and internal_error (500).
const (
	// Introduced in CertVerifier v3.0.0.
	// NULL_ERROR Unused status code. If this is returned, there is a bug in the code.
	StatusNullError CheckDACertStatusCode = iota
	// Introduced in CertVerifier v3.0.0.
	// SUCCESS Verification succeeded
	StatusSuccess
	// Introduced in CertVerifier v3.0.0. Deprecated in v3.1.0 (mapped to INVALID_CERT instead)
	// INVALID_INCLUSION_PROOF Merkle inclusion proof is invalid
	StatusInvalidInclusionProof
	// Introduced in CertVerifier v3.0.0. Deprecated in v3.1.0 (mapped to INVALID_CERT instead)
	// SECURITY_ASSUMPTIONS_NOT_MET Security assumptions not met
	StatusSecurityAssumptionsNotMet
	// Introduced in CertVerifier v3.0.0. Deprecated in v3.1.0 (mapped to INVALID_CERT instead)
	// BLOB_QUORUMS_NOT_SUBSET Blob quorums not a subset of confirmed quorums
	StatusBlobQuorumsNotSubset
	// Introduced in CertVerifier v3.0.0. Deprecated in v3.1.0 (mapped to INVALID_CERT instead)
	// REQUIRED_QUORUMS_NOT_SUBSET Required quorums not a subset of blob quorums
	StatusRequiredQuorumsNotSubset
	// Introduced in CertVerifier v3.1.0
	// INVALID_CERT Certificate is invalid due to some revert from the onchain verification library
	StatusInvalidCert
	// Introduced in CertVerifier v3.1.0
	// INTERNAL_ERROR Bug or misconfiguration in the CertVerifier contract itself. This includes solidity panics and evm reverts.
	StatusContractInternalError
)

// String returns a human-readable representation of the StatusCode.
func (s CheckDACertStatusCode) String() string {
	switch s {
	case StatusNullError:
		return "Null Error: Unused status code. If this is returned, there is a bug in the code."
	case StatusSuccess:
		return "Success: Verification succeeded"
	case StatusInvalidInclusionProof:
		return "Invalid inclusion proof detected: Merkle inclusion proof for blob batch is invalid"
	case StatusSecurityAssumptionsNotMet:
		return "Security assumptions not met: BLS signer weight is less than the required threshold"
	case StatusBlobQuorumsNotSubset:
		return "Blob quorums are not a subset of the confirmed quorums"
	case StatusRequiredQuorumsNotSubset:
		return "Required quorums are not a subset of the blob quorums"
	case StatusInvalidCert:
		return "Invalid certificate: Certificate is invalid due to some revert from the verification library"
	case StatusContractInternalError:
		return "Contract Internal error: Bug or misconfiguration in the CertVerifier contract itself. This includes solidity panics and evm reverts."
	default:
		return "Unknown status code"
	}
}
