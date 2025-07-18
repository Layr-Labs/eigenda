package coretypes

import (
	"encoding/json"
	"fmt"

	_ "github.com/Layr-Labs/eigenda/api/clients/codecs"
)

// DerivationErrorStatusCode is an enum for the different error status codes
// that can be returned during EigenDA "derivation" of a payload from a DA cert.
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
		panic(fmt.Sprintf("DerivationError: invalid status code %d, must be between 1 and 4", e.StatusCode))
	}
	// The Msg field should ideally be a human-readable string that explains the error,
	// but we don't enforce it.
}

// These errors can be used as sentinels to indicate specific derivation errors,
// but they should be used with the `WithMessage` method to add context.
// Also see the constructors below for creating these errors with context.
//
// Note: we purposefully don't use StatusCode 0 here, to prevent default value bugs in case people
// create a DerivationError by hand without using the constructors or sentinel errors defined here.
var (
	// Signifies that the input can't be parsed into a versioned cert.
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
