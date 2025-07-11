package eigenda

import (
	"encoding/json"
	"fmt"
)

// DerivationErrorStatusCode is an enum for the different error status codes
// that can be returned during EigenDA "derivation" of a payload from a DA cert.
// This error is meant to be marshalled to JSON and returned as an HTTP 418 body
// to indicate that the cert should be discarded from rollups' derivation pipeline.
//
// See https://github.com/Layr-Labs/optimism/pull/45 for how this is
// used in optimism's derivation pipeline.
type DerivationError struct {
	StatusCode uint8
	Msg        string
}

func (e DerivationError) Error() string {
	return fmt.Sprintf("derivation error: status code %d, message: %s", e.StatusCode, e.Msg)
}

// Marshalled to JSON and returned as an HTTP 418 body
// to indicate that the cert should be discarded from rollups' derivation pipeline.
// We panic if marshalling fails, since the caller won't be able to handle the derivation error
// properly, so they'll receive a 500 and must retry.
func (e DerivationError) MarshalToTeapotBody() string {
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

var (
	// A cert can be invalid due to parsing error or an onchain CertVerifier error status code.
	// See https://github.com/Layr-Labs/hokulea/tree/master/docs#deriving-op-channel-frames-from-da-cert
	// for full details.
	ErrInvalidCertDerivationError = DerivationError{StatusCode: 1}
	// Signifies that the cert is invalid due to a recency check failure,
	// meaning that `cert.L1InclusionBlock > batch.RBN + rbnRecencyWindowSize`.
	ErrRecencyCheckFailedDerivationError = DerivationError{StatusCode: 2}
	ErrBlobDecodingFailedDerivationError = DerivationError{StatusCode: 3}
)

// Signifies that the cert is invalid due to a parsing failure,
// meaning that a versioned cert could not be parsed from the serialized hex string.
// For example a CertV3 failed to get rlp.decoded from the hex string.
func NewCertParsingFailedError(serializedCertHex string, err string) DerivationError {
	return ErrInvalidCertDerivationError.WithMessage(
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
