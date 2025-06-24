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
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/api/proxy/common/proxyerrors"
	eigendav2store "github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/v2"
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
				return &verification.CertVerificationFailedError{
					StatusCode: 99,
					Msg:        "cert failed",
				}
			},
			expectStatus: http.StatusTeapot,
		},
		{
			name: "418 RBNRecencyCheckFailedError",
			handleFn: func(w http.ResponseWriter, r *http.Request) error {
				return eigendav2store.NewRBNRecencyCheckFailedError(1, 2, 3)
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
		expectVerificationStatusCode coretypes.VerificationStatusCode
	}{
		{
			name: "CertVerificationFailedError",
			err: &verification.CertVerificationFailedError{
				StatusCode: 42, Msg: "cert verification failed"},
			expectHTTPStatus:             http.StatusTeapot,
			expectVerificationStatusCode: 42, // 42 is arbitrarily chosen for this test
		},
		{
			name:                         "RBNRecencyCheckFailedError",
			err:                          eigendav2store.NewRBNRecencyCheckFailedError(1, 2, 3),
			expectHTTPStatus:             http.StatusTeapot,
			expectVerificationStatusCode: eigendav2store.StatusRBNRecencyCheckFailed,
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
				StatusCode coretypes.VerificationStatusCode `json:"StatusCode"`
				Msg        string                           `json:"Msg"`
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
