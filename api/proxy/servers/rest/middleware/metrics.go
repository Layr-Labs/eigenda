package middleware

import (
	"net/http"
	"strconv"

	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
)

// withMetrics is a middleware that records metrics for the route path.
// It does not write anything to the response, that is the job of the handlers.
func withMetrics(
	handleFn func(http.ResponseWriter, *http.Request) error,
	m metrics.Metricer,
	mode commitments.CommitmentMode,
) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		recordDur := m.RecordRPCServerRequest(r.Method)

		scw := newStatusCaptureWriter(w)
		err := handleFn(scw, r)

		certVersion := getCertVersion(r)
		// Prob should use different metric for POST and GET errors.
		recordDur(strconv.Itoa(scw.status), string(mode), certVersion)

		// Forward error to the logging middleware
		return err
	}
}
