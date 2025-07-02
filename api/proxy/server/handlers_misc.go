// handlers_misc.go contains miscellaneous handlers that do not fit into the main request flow.
// These are all health, debug, and testing endpoints.
//
// These handlers SHOULD NOT be wrapped in middlewares, as the middlewares are currently
// hardcoded to log and emit cert related information (we will ideally eventually fix this).
// Handlers in this file thus need to do their own logging and error handling.
//
// DO NOT FORGET to add `http.WriteHeader(http.StatusCodes)` on every error path!
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/proxyerrors"
)

const (
	// HTTP headers
	headerContentType = "Content-Type"

	// Content types
	contentTypeJSON = "application/json"
)

func (svr *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (svr *Server) logDispersalGetError(w http.ResponseWriter, _ *http.Request) {
	svr.log.Warn(`GET method invoked on /put/ endpoint.
		This can occur due to 303 redirects when using incorrect slash ticks.`)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

type EigenDADispersalBackendJSON struct {
	EigenDADispersalBackend string `json:"eigenDADispersalBackend"`
}

// handleGetEigenDADispersalBackend handles the GET request to check the current EigenDA backend used for dispersal.
// This endpoint returns which EigenDA backend version (v1 or v2) is currently being used for blob dispersal.
func (svr *Server) handleGetEigenDADispersalBackend(w http.ResponseWriter, r *http.Request) {
	backend := svr.sm.GetDispersalBackend()
	backendString := common.EigenDABackendToString(backend)

	response := EigenDADispersalBackendJSON{EigenDADispersalBackend: backendString}
	svr.writeJSON(w, r, response)
}

// handleSetEigenDADispersalBackend handles the PUT request to set the EigenDA backend used for dispersal.
// This endpoint configures which EigenDA backend version (v1 or v2) will be used for blob dispersal.
func (svr *Server) handleSetEigenDADispersalBackend(w http.ResponseWriter, r *http.Request) {
	// Read request body to get the new value
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1024)) // Small limit since we only expect a string
	if err != nil {
		svr.log.Error("failed to read request body", "method", r.Method, "path", r.URL.Path, "error", err)
		http.Error(w, proxyerrors.NewReadRequestBodyError(err, 1024).Error(), http.StatusBadRequest)
		return
	}

	// Parse the backend string value
	var eigenDADispersalBackendToSet EigenDADispersalBackendJSON
	if err := json.Unmarshal(body, &eigenDADispersalBackendToSet); err != nil {
		err := proxyerrors.NewUnmarshalJSONError(fmt.Errorf("parsing eigenDADispersalBackend"))
		svr.log.Error("failed to unmarshal body", "method", r.Method, "path", r.URL.Path, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the string to EigenDABackend enum
	backend, err := common.StringToEigenDABackend(eigenDADispersalBackendToSet.EigenDADispersalBackend)
	if err != nil {
		// already a structured error that error middleware knows how to handle
		svr.log.Error(
			"failed to convert string to EigenDABackend",
			"method",
			r.Method,
			"path",
			r.URL.Path,
			"error",
			err,
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	svr.SetDispersalBackend(backend)

	// We return a 200 OK response because the backend was successfully set.
	// Note that writeJSON below can fail to write the response,
	// but we still want to return a 200 OK here to indicate the backend was set.
	// WriteHeader can only be written once, so even if marshalling fails,
	// the WriteHeader(http.StatusInternalServerError) will not overwrite the 200.
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusOK)

	// Exact same logic as GET handler.
	newBackend := svr.sm.GetDispersalBackend()
	backendString := common.EigenDABackendToString(newBackend)

	response := EigenDADispersalBackendJSON{EigenDADispersalBackend: backendString}
	svr.writeJSON(w, r, response)
}

func (svr *Server) writeJSON(w http.ResponseWriter, r *http.Request, response interface{}) {
	jsonData, err := json.Marshal(response)
	if err != nil {
		svr.log.Error("failed to marshal response to json", "method", r.Method, "path", r.URL.Path, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "failed to marshal response to json: %v", err)
		return
	}

	w.Header().Set(contentTypeJSON, "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		svr.log.Error("failed to write response", "method", r.Method, "path", r.URL.Path, "error", err)
	}
}
