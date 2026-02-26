package middleware

import (
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// WithLogging is a middleware that logs information related to each request.
// It does not write anything to the response, that is the job of the handlers.
// Currently we cannot log the status code because go's default ResponseWriter interface does not expose it.
// TODO: implement a ResponseWriter wrapper that saves the status code: see https://github.com/golang/go/issues/18997
func withLogging(
	handleFn func(http.ResponseWriter, *http.Request) error,
	log logging.Logger,
	mode commitments.CommitmentMode,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		scw := newStatusCaptureWriter(w)
		err := handleFn(scw, r)

		args := []any{
			"method", r.Method, "url", r.URL,
			"commitment_mode", mode, "cert_version", getCertVersion(r),
			"status", scw.status, "duration", time.Since(start),
		}

		// Success-path logging is handled by the handlers (e.g. "Processed request").
		// Keeping middleware logging only for error paths avoids duplicate log lines.
		if err == nil {
			return
		}

		args = append(args, "error", err.Error())
		if scw.status >= 400 && scw.status < 500 {
			log.Warn("request completed with 4xx error", args...)
		} else {
			log.Error("request completed with error", args...)
		}
	}
}
