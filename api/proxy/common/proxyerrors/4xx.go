package proxyerrors

import (
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	_ "github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/v2"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Is400(err error) bool {
	var parsingError ParsingError
	var certHexDecodingError CertHexDecodingError
	var invalidBackendErr common.InvalidBackendError
	var unmarshalJSONErr UnmarshalJSONError
	var l1InclusionBlockNumberParsingError L1InclusionBlockNumberParsingError
	var readRequestBodyErr ReadRequestBodyError
	var s3KeccakKeyValueMismatchErr s3.Keccak256KeyValueMismatchError
	return errors.Is(err, ErrProxyOversizedBlob) ||
		errors.As(err, &parsingError) ||
		errors.As(err, &certHexDecodingError) ||
		errors.As(err, &invalidBackendErr) ||
		errors.As(err, &unmarshalJSONErr) ||
		errors.As(err, &l1InclusionBlockNumberParsingError) ||
		errors.As(err, &readRequestBodyErr) ||
		errors.As(err, &s3KeccakKeyValueMismatchErr) ||
		errors.Is(err, s3.ErrKeccakKeyNotFound)
}

// We return a 418 TEAPOT error for any cert validation error.
// Rollup derivation pipeline should drop any certs that return this error.
// See https://github.com/Layr-Labs/optimism/pull/45 for how this is
// used in optimism's derivation pipeline.
func Is418(err error) bool {
	var invalidCertErr *verification.CertVerificationFailedError
	return errors.As(err, &invalidCertErr)
}

// 429 TOO_MANY_REQUESTS is returned to the client to inform them that they are getting rate-limited
// on the EigenDA disperser. The disperser returns a grpc RESOURCE_EXHAUSTED error, which we convert
// to an HTTP error. It doesn't have any meaning other than to request the client to retry later,
// and/or slow down their rate of requests.
func Is429(err error) bool {
	st, isGRPCError := status.FromError(err)
	return isGRPCError && st.Code() == codes.ResourceExhausted
}

var (
	ErrProxyOversizedBlob = fmt.Errorf("encoded blob is larger than max blob size")
)

type CertHexDecodingError struct {
	serializedCertHex string
	err               error
}

func NewCertHexDecodingError(serializedCertHex string, err error) CertHexDecodingError {
	return CertHexDecodingError{
		serializedCertHex: serializedCertHex,
		err:               err,
	}
}
func (me CertHexDecodingError) Error() string {
	return fmt.Sprintf("decoding cert from hex string: %s, error: %s",
		me.serializedCertHex,
		me.err.Error())
}

// l1_inclusion_block_number is a query param that is used to specify the L1 block number
// at which a cert was included in the batcher inbox. It is used to perform the rbn recency check.
// It is optional, but if it is provided and invalid, we return a 400 error
// to let the client know that they probably have a bug.
type L1InclusionBlockNumberParsingError struct {
	l1BlockNumStr string
	err           error
}

func NewL1InclusionBlockNumberParsingError(l1BlockNumStr string, err error) L1InclusionBlockNumberParsingError {
	return L1InclusionBlockNumberParsingError{
		l1BlockNumStr: l1BlockNumStr,
		err:           err,
	}
}

func (me L1InclusionBlockNumberParsingError) Error() string {
	return fmt.Sprintf("invalid l1_inclusion_block_number %s: %s",
		me.l1BlockNumStr,
		me.err.Error())
}

// ReadRequestBodyError is used to wrap errors that occur when reading the request body.
// This typically happens when we fail to read a payload from a POST request body.
// Reading from body payload should always be limited to a certain size, using
// https://pkg.go.dev/net/http#MaxBytesReader. Unfortunately, MaxBytesReader
// returns an error that doesn't include the limit, so we wrap it in this custom error.
// See https://cs.opensource.google/go/go/+/refs/tags/go1.24.3:src/net/http/request.go;l=1200
// for the dumb error http returns.
type ReadRequestBodyError struct {
	bodyLimit int64
	err       error
}

func NewReadRequestBodyError(err error, bodyLimit int64) ReadRequestBodyError {
	return ReadRequestBodyError{
		bodyLimit: bodyLimit,
		err:       err,
	}
}
func (me ReadRequestBodyError) Error() string {
	return fmt.Sprintf("reading at most %d bytes from body: %s", me.bodyLimit, me.err.Error())
}

type UnmarshalJSONError struct {
	err error
}

func NewUnmarshalJSONError(err error) UnmarshalJSONError {
	return UnmarshalJSONError{
		err: err,
	}
}

func (me UnmarshalJSONError) Error() string {
	return fmt.Sprintf("unmarshalling JSON: %s", me.err.Error())
}

// ParsingError is a very coarse-grained error that's used as a catch-all for any parsing errors
// like parsing a hex string, or parsing a version byte from the request path, reading a query param, etc.
// TODO: should all of these be returned as [eigenda.StatusCertParsingFailed] errors instead,
// to return TEAPOTs instead of 400s?
type ParsingError struct {
	err error
}

func NewParsingError(err error) ParsingError {
	return ParsingError{
		err: err,
	}
}
func (me ParsingError) Error() string {
	return fmt.Sprintf("parsing error: %s", me.err.Error())
}
