package verification

// VerificationStatusCode represents the status codes that are returned by
// EigenDACertVerifier.checkDACert contract calls. The enum values below should match exactly
// the status codes defined in the contract:
// https://github.com/Layr-Labs/eigenda/blob/1091f460ba762b84019389cbb82d9b04bb2c2bdb/contracts/src/integrations/cert/libraries/EigenDACertVerificationLib.sol#L48-L54
type VerificationStatusCode uint8

const (
	// NULL_ERROR Unused status code. If this is returned, there is a bug in the code.
	StatusNullError VerificationStatusCode = iota
	// SUCCESS Verification succeeded
	StatusSuccess
	// INVALID_INCLUSION_PROOF Merkle inclusion proof is invalid
	StatusInvalidInclusionProof
	// SECURITY_ASSUMPTIONS_NOT_MET Security assumptions not met
	StatusSecurityAssumptionsNotMet
	// BLOB_QUORUMS_NOT_SUBSET Blob quorums not a subset of confirmed quorums
	StatusBlobQuorumsNotSubset
	// REQUIRED_QUORUMS_NOT_SUBSET Required quorums not a subset of blob quorums
	StatusRequiredQuorumsNotSubset
)

// String returns a human-readable representation of the StatusCode.
func (s VerificationStatusCode) String() string {
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
	default:
		return "Unknown status code"
	}
}
