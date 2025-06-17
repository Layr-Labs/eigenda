package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/common/types/commitments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/stretchr/testify/require"
)

// Make sue that SetCertVersion/getCertVersion are working correctly,
// by using a mock metrics that makes sure the metrics middleware calls
// recordDur with the correct cert version.
// TODO: we prob should also test the logging middleware, but that's a
// brittle test and logger will probably change soon so inclined to skip it for now.
func TestRequestContext_CertVersionCanBeReadFromMetricsMiddleware(t *testing.T) {
	const testCertVersion = "v42"

	// Handler sets the cert version and echoes it back in JSON
	handler := func(w http.ResponseWriter, r *http.Request) error {
		SetCertVersion(r, testCertVersion)
		return nil
	}
	mockMetrics := &MockMetricer{}
	testLogger := logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{})
	// Compose the middleware chain
	mw := WithCertMiddlewares(
		handler,
		testLogger,
		mockMetrics,
		commitments.OptimismGenericCommitmentMode,
	)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	mw(rec, req)

	require.Equal(t, mockMetrics.recordDurCertVersion, testCertVersion,
		"The cert version should be captured in the metrics middleware")
}

// Mock implementation of the Metricer interface.
// Only used to make sure that the call to recordDur(strconv.Itoa(scw.status), string(mode), certVersion)
// in the metrics middleware contains the correct cert version.
type MockMetricer struct {
	recordDurCertVersion string
}

func (m *MockMetricer) RecordInfo(version string) {}
func (m *MockMetricer) RecordUp()                 {}
func (m *MockMetricer) RecordRPCServerRequest(method string) func(status string, mode string, ver string) {
	return func(status string, mode string, ver string) {
		if m.recordDurCertVersion != "" {
			panic("recordDurCertVersion should only be set once")
		}
		m.recordDurCertVersion = ver // Capture the cert version
	}
}
func (m *MockMetricer) RecordSecondaryRequest(bt string, method string) func(status string) {
	return func(status string) {}
}
func (m *MockMetricer) Document() []opmetrics.DocumentedMetric {
	return nil
}
