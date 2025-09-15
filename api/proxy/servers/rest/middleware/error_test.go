package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/proxy/common/proxyerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWithErrorHandling_HTTPStatusCodes(t *testing.T) {
	type testCase struct {
		name         string
		handleFn     func(http.ResponseWriter, *http.Request) error
		expectStatus int
	}

	testErr := errors.New("test error")

	tests := []testCase{
		{
			name: "400 Bad Request",
			handleFn: func(w http.ResponseWriter, r *http.Request) error {
				// Use a proxyerrors.ParsingError which triggers Is400
				return proxyerrors.NewParsingError(testErr)
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "418 CertVerificationFailedError",
			handleFn: func(w http.ResponseWriter, r *http.Request) error {
				return coretypes.ErrInvalidCertDerivationError
			},
			expectStatus: http.StatusTeapot,
		},
		{
			name: "418 RBNRecencyCheckFailedError",
			handleFn: func(w http.ResponseWriter, r *http.Request) error {
				return coretypes.NewRBNRecencyCheckFailedError(1, 2, 3)
			},
			expectStatus: http.StatusTeapot,
		},
		{
			name: "429 Too Many Requests",
			handleFn: func(w http.ResponseWriter, r *http.Request) error {
				// Simulate a gRPC ResourceExhausted error
				return status.Error(codes.ResourceExhausted, "rate limited")
			},
			expectStatus: http.StatusTooManyRequests,
		},
		{
			name: "503 Service Unavailable",
			handleFn: func(w http.ResponseWriter, r *http.Request) error {
				// Simulate a proxyerrors.Is503 error
				return &api.ErrorFailover{}
			},
			expectStatus: http.StatusServiceUnavailable,
		},
		{
			name: "500 Internal Server Error",
			handleFn: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("unexpected error")
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := withErrorHandling(tc.handleFn)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			err := handler(rr, req)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if rr.Code != tc.expectStatus {
				t.Errorf("expected status %d, got %d", tc.expectStatus, rr.Code)
			}
		})
	}
}

// This one tests that the json body of 418 TEAPOT errors for cert verification failures
// contains the StatusCode, which is used by rollup derivation pipelines.
func TestWithErrorHandling_418TeapotErrors(t *testing.T) {
	tests := []struct {
		name                         string
		err                          error
		expectHTTPStatus             int
		expectVerificationStatusCode uint8
	}{
		{
			name:                         "CertParsingFailedDerivationError",
			err:                          coretypes.ErrCertParsingFailedDerivationError.WithMessage("some arbitrary msg"),
			expectHTTPStatus:             http.StatusTeapot,
			expectVerificationStatusCode: coretypes.ErrCertParsingFailedDerivationError.StatusCode,
		},
		{
			name:                         "RBNRecencyCheckFailedError",
			err:                          coretypes.NewRBNRecencyCheckFailedError(1, 2, 3),
			expectHTTPStatus:             http.StatusTeapot,
			expectVerificationStatusCode: coretypes.ErrRecencyCheckFailedDerivationError.StatusCode,
		},
		{
			name:                         "InvalidCertDerivationError",
			err:                          coretypes.ErrInvalidCertDerivationError.WithMessage("some arbitrary msg"),
			expectHTTPStatus:             http.StatusTeapot,
			expectVerificationStatusCode: coretypes.ErrInvalidCertDerivationError.StatusCode,
		},
		{
			name:                         "BlobDecodingFailedDerivationError",
			err:                          coretypes.ErrBlobDecodingFailedDerivationError.WithMessage("some arbitrary msg"),
			expectHTTPStatus:             http.StatusTeapot,
			expectVerificationStatusCode: coretypes.ErrBlobDecodingFailedDerivationError.StatusCode,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := withErrorHandling(func(w http.ResponseWriter, r *http.Request) error {
				return tc.err
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			err := handler(rr, req)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if rr.Code != tc.expectHTTPStatus {
				t.Errorf("expected status %d, got %d", tc.expectHTTPStatus, rr.Code)
			}
			var resp struct {
				StatusCode uint8  `json:"StatusCode"`
				Msg        string `json:"Msg"`
			}
			dec := json.NewDecoder(strings.NewReader(rr.Body.String()))
			if err := dec.Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp.StatusCode != tc.expectVerificationStatusCode {
				t.Errorf("expected StatusCode %d, got %d", tc.expectVerificationStatusCode, resp.StatusCode)
			}
		})
	}
}
