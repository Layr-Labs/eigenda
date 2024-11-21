package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/ethereum/go-ethereum/log"
)

// Used to capture the status code of the response, so that we can use it in middlewares.
// See https://github.com/golang/go/issues/18997
// TODO: right now instantiating a separate scw for logging and metrics... is there a better way?
type statusCaptureWriter struct {
	http.ResponseWriter
	status int
}

func (scw *statusCaptureWriter) WriteHeader(status int) {
	scw.status = status
	scw.ResponseWriter.WriteHeader(status)
}

func newStatusCaptureWriter(w http.ResponseWriter) *statusCaptureWriter {
	return &statusCaptureWriter{
		ResponseWriter: w,
		// 200 status code is only added to response by outer layer http framework,
		// since WriteHeader(200) is typically not called by handlers.
		// So we initialize status as 200, and assume that any other status code
		// will be set by the handler.
		status: http.StatusOK,
	}
}

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
		if err != nil {
			commitMode := "unknown"
			certVersion := "unknown"
			var metaErr MetaError
			if errors.As(err, &metaErr) {
				commitMode = string(metaErr.Meta.Mode)
				certVersion = string(metaErr.Meta.CertVersion)
			}
			recordDur(strconv.Itoa(scw.status), commitMode, certVersion)
			return err
		}
		versionByte, err := parseVersionByte(w, r)
		if err != nil {
			recordDur(strconv.Itoa(scw.status), string(mode), "unknown")
			return fmt.Errorf("metrics middleware: parsing version byte: %w", err)
		}
		recordDur(strconv.Itoa(scw.status), string(mode), string(versionByte))
		return nil
	}
}

// withLogging is a middleware that logs information related to each request.
// It does not write anything to the response, that is the job of the handlers.
// Currently we cannot log the status code because go's default ResponseWriter interface does not expose it.
// TODO: implement a ResponseWriter wrapper that saves the status code: see https://github.com/golang/go/issues/18997
func withLogging(
	handleFn func(http.ResponseWriter, *http.Request) error,
	log log.Logger,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		scw := newStatusCaptureWriter(w)
		err := handleFn(scw, r)

		args := []any{
			"method", r.Method, "url", r.URL, "status", scw.status, "duration", time.Since(start),
		}
		if err != nil {
			args = append(args, "err", err)
		}
		var metaErr MetaError
		if errors.As(err, &metaErr) {
			args = append(args, "commitment_mode", metaErr.Meta.Mode, "cert_version", metaErr.Meta.CertVersion)
		}
		log.Info("request", args...)
	}
}
