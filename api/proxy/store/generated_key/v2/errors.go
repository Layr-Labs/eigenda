package eigenda

import (
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
)

// The error status codes defined in this file are reasons for certs to be considered invalid,
// that do not come from the onchain cert verifier contract logic, unlike the rest of the
// [verification.CertVerificationFailedError] errors.
// See the docs in https://github.com/Layr-Labs/hokulea/tree/master/docs#deriving-op-channel-frames-from-da-cert
// for the different ways certs can be invalid. Each of those need to be represented by a StatusCode here.
// TODO: these error codes need to be added to the [coretypes.VerificationStatusCode] enum
// after we move proxy to the monorepo.

// Signifies that the cert is invalid due to a recency check failure,
// meaning that `cert.L1InclusionBlock > batch.RBN + rbnRecencyWindowSize`.
const StatusRBNRecencyCheckFailed = coretypes.VerificationStatusCode(math.MaxUint8)

// Signifies that the cert is invalid due to a parsing failure,
// meaning that a versioned cert could not be parsed from the serialized hex string.
// For example a CertV3 failed to get rlp.decoded from the hex string,
const StatusCertParsingFailed = coretypes.VerificationStatusCode(math.MaxUint8 - 1)

// Returned when `cert.L1InclusionBlock > batch.RBN + rbnRecencyWindowSize`,
// to indicate that the cert should be discarded from rollups' derivation pipeline.
// This should get converted by the proxy to an HTTP 418 TEAPOT error code.
func NewRBNRecencyCheckFailedError(
	certRBN uint64, certL1IBN uint64, rbnRecencyWindowSize uint64,
) *verification.CertVerificationFailedError {
	return &verification.CertVerificationFailedError{
		StatusCode: StatusRBNRecencyCheckFailed,
		Msg: fmt.Sprintf(
			"RBN recency check failed: certL1InclusionBlockNumber (%d) > cert.RBN (%d) + RBNRecencyWindowSize (%d)",
			certL1IBN, certRBN, rbnRecencyWindowSize,
		),
	}
}

func NewCertParsingFailedError(serializedCertHex string, err string) *verification.CertVerificationFailedError {
	return &verification.CertVerificationFailedError{
		StatusCode: StatusCertParsingFailed,
		Msg: fmt.Sprintf(
			"cert parsing failed for cert %s: %v",
			serializedCertHex, err,
		),
	}
}
