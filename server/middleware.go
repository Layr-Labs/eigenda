package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/ethereum/go-ethereum/log"
)

// withMetrics is a middleware that records metrics for the route path.
func withMetrics(
	handleFn func(http.ResponseWriter, *http.Request) error,
	m metrics.Metricer,
	mode commitments.CommitmentMode,
) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		recordDur := m.RecordRPCServerRequest(r.Method)

		err := handleFn(w, r)
		if err != nil {
			var metaErr MetaError
			if errors.As(err, &metaErr) {
				recordDur(w.Header().Get("status"), string(metaErr.Meta.Mode), string(metaErr.Meta.CertVersion))
			} else {
				recordDur(w.Header().Get("status"), string("unknown"), string("unknown"))
			}
			return err
		}
		// we assume that every route will set the status header
		versionByte, err := parseVersionByte(w, r)
		if err != nil {
			recordDur(w.Header().Get("status"), string(mode), string(versionByte))
			return fmt.Errorf("metrics middleware: error parsing version byte: %w", err)
		}
		recordDur(w.Header().Get("status"), string(mode), string(versionByte))
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
		err := handleFn(w, r)
		var metaErr MetaError
		//nolint:gocritic // ifElseChain is not a good replacement with errors.As
		if errors.As(err, &metaErr) {
			log.Info("request", "method", r.Method, "url", r.URL, "duration", time.Since(start),
				"err", err, "status", w.Header().Get("status"),
				"commitment_mode", metaErr.Meta.Mode, "cert_version", metaErr.Meta.CertVersion)
		} else if err != nil {
			log.Info("request", "method", r.Method, "url", r.URL, "duration", time.Since(start), "err", err)
		} else {
			log.Info("request", "method", r.Method, "url", r.URL, "duration", time.Since(start))
		}
	}
}
