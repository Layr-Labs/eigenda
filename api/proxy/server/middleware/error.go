package middleware

import (
	"errors"
	"net/http"

	"github.com/Layr-Labs/eigenda/api/proxy/common/proxyerrors"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/v2"
)

// Error handling middleware (innermost) transforms internal errors to HTTP errors,
func withErrorHandling(
	handleFn func(http.ResponseWriter, *http.Request) error,
) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := handleFn(w, r)
		if err == nil {
			return nil
		}

		// TODO: should we add request specific information like GET vs POST,
		// commitment mode, cert version, etc. to each error?
		// Or maybe we should just add a requestID to the error, and log the request-specific information
		// in the logging middleware, so that we can correlate the error with the request?

		var teapotErr eigenda.TeapotError
		switch {
		case proxyerrors.Is400(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		// 418 TEAPOT errors don't follow the pattern proxyerrors.Is418(err),
		// because we need to marshal the correct json body.
		case errors.As(err, &teapotErr):
			http.Error(w, teapotErr.MarshalBody(), http.StatusTeapot)
		case proxyerrors.Is429(err):
			http.Error(w, err.Error(), http.StatusTooManyRequests)
		case proxyerrors.Is503(err):
			// this tells the caller (batcher) to failover to ethda b/c eigenda is temporarily down
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		default:
			// Default to 500 for unexpected errors.
			// Note that this includes grpc 4xx errors returned from the disperser server.
			// because those are due to formatting bugs in proxy code, e.g. badly
			// IFFT'ing or encoding the blob, so we shouldn't return a 400 to the client.
			// See https://github.com/Layr-Labs/eigenda/blob/bee55ed9207f16153c3fd8ebf73c219e68685def/api/errors.go#L22
			// for the 400s returned by the disperser server (currently only INVALID_ARGUMENT).
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// forward error to the logging middleware (through the metrics middleware)
		// so that the error is logged.
		return err
	}
}
