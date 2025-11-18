package coretypes

import (
	"encoding/json"
	"fmt"

	_ "github.com/Layr-Labs/eigenda/api/clients/codecs"
	_ "github.com/Layr-Labs/eigenda/encoding"
)

// Sentinel [DerivationError] errors that set the correct StatusCode.
// If used directly, extend them using [DerivationError.WithMessage] to add context.
// Otherwise, see the specific constructors below for creating these errors with context.
//
// Note: we purposefully don't use StatusCode 0 here, to prevent default value bugs in case people
// create a DerivationError by hand without using the constructors or sentinel errors defined here.
var (
	// Signifies that the input can't be parsed into a versioned cert,
	// meaning either the cert has an invalid version byte, or failed to get rlp.decoded from the given hex string
	ErrCertParsingFailedDerivationError = DerivationError{StatusCode: 1}
	// Signifies that the cert is invalid due to a recency check failure,
	// meaning that `cert.L1InclusionBlock > batch.RBN + rbnRecencyWindowSize`.
	// See https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#1-rbn-recency-validation
	ErrRecencyCheckFailedDerivationError = DerivationError{StatusCode: 2}
	// Signifies that the CertVerifier.checkDACert eth-call returned an error status code.
	// See https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#2-cert-validation
	ErrInvalidCertDerivationError = DerivationError{StatusCode: 3}
	// Signifies that the blob is incorrectly encoded, and cannot be decoded into a valid payload.
	// See [codecs.PayloadEncodingVersion] for the different supported encodings.
	// See https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation
	ErrBlobDecodingFailedDerivationError = DerivationError{StatusCode: 4}
)

// Sentinel [MaliciousOperatorsError] errors.
// Extend these with [MaliciousOperatorsError.WithBlobKey] to add context.
var (
	// [encoding.BlobCommitments.Length] needs to be a power of 2, and that is checked by the eigenda validators:
	// https://github.com/Layr-Labs/eigenda/blob/cc392dbabef362f2e03a4b35616a407d45fad510/core/v2/assignment.go#L308
	// Therefore, if we ever receive a cert with a blob length that is not a power of 2,
	// it means that the eigenda validators are colluding and doing something fishy.
	ErrCertCommitmentBlobLengthNotPowerOf2MaliciousOperatorsError = MaliciousOperatorsError{
		Msg: "blob length in cert commitment is not a power of 2",
	}
)

// DerivationErrorStatusCode is an enum for the different error status codes
// that can be returned during EigenDA "derivation" of a payload from a DA cert,
// and signify that the cert is invalid and should be dropped.
// For more details, see the EigenDA spec on derivation:
// https://github.com/Layr-Labs/eigenda/blob/f4ef5cd5/docs/spec/src/integration/spec/6-secure-integration.md#derivation-process
//
// This error is meant to be marshalled to JSON and returned as an HTTP 418 body
// to indicate that the cert should be discarded from rollups' derivation pipelines.
//
// See https://github.com/Layr-Labs/optimism/pull/50 for how this is
// used in optimism's derivation pipeline.
type DerivationError struct {
	StatusCode uint8
	Msg        string
}

func (e DerivationError) Error() string {
	return fmt.Sprintf("derivation error: status code %d, message: %s", e.StatusCode, e.Msg)
}

// Marshalled to JSON and returned as an HTTP 418 body
// to indicate that the cert should be discarded from rollups' derivation pipelines.
// We panic if marshalling fails, since the caller won't be able to handle the derivation error
// properly, so they'll receive a 500 and must retry.
func (e DerivationError) MarshalToTeapotBody() string {
	e.Validate()
	bodyJSON, err := json.Marshal(e)
	if err != nil {
		panic(fmt.Errorf("failed to marshal derivation error: %w", err))
	}
	return string(bodyJSON)
}

// Used to add context to the sentinel errors below. For example:
// ErrInvalidCertDerivationError.WithMessage("failed to parse cert")
func (e DerivationError) WithMessage(msg string) DerivationError {
	return DerivationError{
		StatusCode: e.StatusCode,
		Msg:        msg,
	}
}

// Validate that the DerivationError has a valid status code.
// The only valid status codes are 1-4, as defined in the sentinel errors below, eg [ErrCertParsingFailedDerivationError].
func (e DerivationError) Validate() {
	if e.StatusCode < 1 || e.StatusCode > 4 {
		panic(fmt.Errorf("DerivationError: invalid status code %d, must be between 1 and 4", e.StatusCode))
	}
	// The Msg field should ideally be a human-readable string that explains the error,
	// but we don't enforce it.
}

// Signifies that the cert is invalid due to a parsing failure,
// meaning that a versioned cert could not be parsed from the serialized hex string.
// For example a CertV3 failed to get rlp.decoded from the hex string.
func NewCertParsingFailedError(serializedCertHex string, err string) DerivationError {
	return ErrCertParsingFailedDerivationError.WithMessage(
		fmt.Sprintf("cert parsing failed for cert %s: %v", serializedCertHex, err),
	)
}

// Signifies that the cert is invalid due to a recency check failure,
// meaning that `cert.L1InclusionBlock > batch.RBN + rbnRecencyWindowSize`.
func NewRBNRecencyCheckFailedError(
	certRBN, certL1InclusionBlock, rbnRecencyWindowSize uint64,
) DerivationError {
	return ErrRecencyCheckFailedDerivationError.WithMessage(
		fmt.Sprintf(
			"RBN recency check failed: certL1InclusionBlockNumber (%d) > cert.RBN (%d) + RBNRecencyWindowSize (%d)",
			certL1InclusionBlock, certRBN, rbnRecencyWindowSize,
		))
}

// MaliciousOperatorsErrors are kept separate from [DerivationError]s because
// they are triggered by errors that should have been validated by the EigenDA operators.
// This means that certs that trigger these errors should never have been signed.
//
// Although the certs that trigger these errors could also be dropped from rollup derivation
// pipelines the same way that DerivationErrors are, they are more serious errors and
// signify that the eigenda validators are possibly colluding and attempting something fishy.
// These errors should cause the software to crash, stopping the rollup and raising alarms
// to investigate the validators or the issue.
//
// If a bug explaining these errors is not found, then very likely the validators
// should get slashed.
type MaliciousOperatorsError struct {
	// The BlobKey can be used to retrieve the BlobStatus to reconstruct the DACert.
	BlobKey string
	// The Msg field contains a human-readable error message explaining the issue.
	// We don't need a status code for these errors because there is no way to
	// programmatically deal with them.
	Msg string
}

func (e MaliciousOperatorsError) Error() string {
	return fmt.Sprintf("malicious operators error: blob key %s, message: %s", e.BlobKey, e.Msg)
}

func (e MaliciousOperatorsError) WithBlobKey(blobKey string) MaliciousOperatorsError {
	return MaliciousOperatorsError{
		BlobKey: blobKey,
		Msg:     e.Msg,
	}
}
