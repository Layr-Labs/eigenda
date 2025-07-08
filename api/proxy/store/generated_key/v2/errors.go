package eigenda

import (
	"encoding/json"
	"fmt"
)

// DerivationErrorStatusCode is an enum for the different error status codes
// that can be returned during EigenDA "derivation" of a payload from a DA cert.
type DerivationErrorStatusCode uint8

const (
	InvalidCertDerivationErrorStatusCode DerivationErrorStatusCode = iota + 1
	RecencyCheckFailedDerivationErrorStatusCode
	BlobDecodingFailedDerivationErrorStatusCode
)

type TeapotError interface {
	// Should generate a [derivationErrorTeapotBody] and marshal it to JSON.
	// Panics if the marshaling fails.
	MarshalBody() string
}

var _ TeapotError = (*InvalidCertDerivationError)(nil)
var _ TeapotError = (*RecencyCheckFailedDerivationError)(nil)

// Marshalled to JSON and returned as an HTTP 418 body
// to indicate that the cert should be discarded from rollups' derivation pipeline.
type derivationErrorTeapotBody struct {
	StatusCode uint8
	Msg        string
}

// InvalidCertDerivationError is returned when a cert is invalid and should be discarded
// from the rollups' derivation pipeline.
// The failure can either be a parsing error, or an onchain CertVerifier error status code.
// See the docs in https://github.com/Layr-Labs/hokulea/tree/master/docs#deriving-op-channel-frames-from-da-cert
// for the different ways certs can be invalid.
type InvalidCertDerivationError struct {
	Msg string
}

func (e *InvalidCertDerivationError) Error() string {
	return fmt.Sprintf("invalid cert: %s", e.Msg)
}

func (e *InvalidCertDerivationError) MarshalBody() string {
	body := derivationErrorTeapotBody{
		StatusCode: uint8(InvalidCertDerivationErrorStatusCode),
		Msg:        e.Msg,
	}
	bodyJSON, err := json.Marshal(&body)
	if err != nil {
		panic(fmt.Errorf("failed to marshal invalid cert derivation error: %w", err))
	}
	return string(bodyJSON)
}

// Signifies that the cert is invalid due to a parsing failure,
// meaning that a versioned cert could not be parsed from the serialized hex string.
// For example a CertV3 failed to get rlp.decoded from the hex string.
// It returns an [InvalidCertDerivationError] which gets converted to a TEAPOT 418 HTTP error.
func NewCertParsingFailedError(serializedCertHex string, err string) *InvalidCertDerivationError {
	return &InvalidCertDerivationError{
		Msg: fmt.Sprintf(
			"cert parsing failed for cert %s: %v",
			serializedCertHex, err,
		),
	}
}

// Signifies that the cert is invalid due to a recency check failure,
// meaning that `cert.L1InclusionBlock > batch.RBN + rbnRecencyWindowSize`.
type RecencyCheckFailedDerivationError struct {
	CertRBN              uint64
	CertL1InclusionBlock uint64
	RBNRecencyWindowSize uint64
}

func NewRBNRecencyCheckFailedError(
	certRBN, certL1InclusionBlock, rbnRecencyWindowSize uint64,
) *RecencyCheckFailedDerivationError {
	return &RecencyCheckFailedDerivationError{
		CertRBN:              certRBN,
		CertL1InclusionBlock: certL1InclusionBlock,
		RBNRecencyWindowSize: rbnRecencyWindowSize,
	}
}

func (e *RecencyCheckFailedDerivationError) Error() string {
	return fmt.Sprintf(
		"RBN recency check failed: certL1InclusionBlockNumber (%d) > cert.RBN (%d) + RBNRecencyWindowSize (%d)",
		e.CertL1InclusionBlock, e.CertRBN, e.RBNRecencyWindowSize,
	)
}

func (e *RecencyCheckFailedDerivationError) MarshalBody() string {
	body := &derivationErrorTeapotBody{
		StatusCode: uint8(RecencyCheckFailedDerivationErrorStatusCode),
		Msg:        e.Error(),
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		panic(fmt.Errorf("failed to marshal recency check failed error: %w", err))
	}
	return string(bodyJSON)
}

// TODO: add blob decoding failed error
