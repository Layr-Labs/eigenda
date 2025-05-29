package verification

import "github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"

// CertVerifierInputError represents a 4xx-like error (invalid input, misuse, unsupported, etc.)
type CertVerifierInputError struct {
	Msg string
}

func (e *CertVerifierInputError) Error() string {
	return e.Msg
}

// CertVerifierInternalError represents a 5xx-like error (unexpected, internal, infra, etc.)
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

// CertVerificationFailedError is returned when cert verification fails:
// [coretypes.VerificationStatusCode] != (StatusSuccess or StatusNullError).
// StatusNullError returns a [CertVerifierInternalError] instead as it is a contract bug
// that should never happen.
type CertVerificationFailedError struct {
	StatusCode coretypes.VerificationStatusCode
	Msg        string
}

func (e *CertVerificationFailedError) Error() string {
	return e.Msg
}
