package verification

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
)

// CertVerifierInputError represents a 4xx-like error (invalid input, misuse, unsupported, etc.)
type CertVerifierInputError struct {
	Msg string
}

func (e *CertVerifierInputError) Error() string {
	return e.Msg
}

// CertVerifierInternalError represents a 5xx-like error (unexpected, internal, infra, etc.)
//
// Our recommendation is to always retry Internal errors in case they were due to a temporary (network) issue.
// TODO: we would want to distinguish temporary vs permanent errors here, to inform the client
// as to whether its worth retrying the request. However, this is currently not possible because the
// underlying geth binding library does not provide this information. For example, a Call()
// does a bunch of things before the actual call, like abi encoding the inputs, which can fail,
// and geth itself does not provide temporary/retryable semantics on its returned errors.
// See https://github.com/ethereum/go-ethereum/blob/a9523b6428238a762e1a1e55e46ead47630c3a23/accounts/abi/bind/base.go#L169
// It seems incredibly difficult in golang (maybe due to golang's laxity with error handling) to distinguish temporary errors.
// See the net package for exampple, which deprecated its Temporary() method: https://pkg.go.dev/net#Error
type CertVerifierInternalError struct {
	Msg string
	// Err is optional and only present if an underlying error is available.
	// Note that we only provide this as a convenience for logging and debugging.
	// Error is NOT part of our public API, so don't match on internal errors,
	// as these errors may change in the future.
	Err error
}

func (e *CertVerifierInternalError) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

// CertVerifierInvalidCertError is returned when cert verification fails:
// [coretypes.VerificationStatusCode] != (StatusSuccess or StatusNullError).
// StatusNullError returns a [CertVerifierInternalError] instead as it is a contract bug
// that should never happen.
type CertVerifierInvalidCertError struct {
	StatusCode coretypes.VerificationStatusCode
	Msg        string
}

func (e *CertVerifierInvalidCertError) Error() string {
	return fmt.Sprintf("invalid cert: call to CertVerifier failed with status code %d: %s", e.StatusCode, e.Msg)
}
