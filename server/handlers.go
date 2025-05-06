package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda-proxy/common/types/commitments"
	"github.com/gorilla/mux"
)

const (
	// limit requests to only 32 mib to mitigate potential DoS attacks
	maxRequestBodySize int64 = 1024 * 1024 * 32

	// HTTP headers
	headerContentType = "Content-Type"

	// Content types
	contentTypeJSON = "application/json"
)

func (svr *Server) handleHealth(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

func (svr *Server) logDispersalGetError(w http.ResponseWriter, _ *http.Request) error {
	svr.log.Warn(`GET method invoked on /put/ endpoint.
		This can occur due to 303 redirects when using incorrect slash ticks.`)
	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

// =================================================================================================
// GET ROUTES
// =================================================================================================

// handleGetStdCommitment handles the GET request for std commitments.
func (svr *Server) handleGetStdCommitment(w http.ResponseWriter, r *http.Request) error {
	certVersion, err := parseCertVersion(w, r)
	if err != nil {
		return fmt.Errorf("error parsing version byte: %w", err)
	}
	serializedCertHex, ok := mux.Vars(r)[routingVarNamePayloadHex]
	if !ok {
		return fmt.Errorf("serializedDACert not found in path: %s", r.URL.Path)
	}
	serializedCert, err := hex.DecodeString(serializedCertHex)
	if err != nil {
		return fmt.Errorf("failed to decode from hex serializedDACert %s: %w", serializedCertHex, err)
	}
	versionedCert := certs.NewVersionedCert(serializedCert, certVersion)

	return svr.handleGetShared(r.Context(), w, versionedCert, commitments.StandardCommitmentMode)
}

// handleGetOPKeccakCommitment handles the GET request for optimism keccak commitments.
func (svr *Server) handleGetOPKeccakCommitment(w http.ResponseWriter, r *http.Request) error {
	// TODO: do we use a version byte in OPKeccak commitments? README seems to say so, but server_test didn't
	// versionByte, err := parseVersionByte(r)
	// if err != nil {
	// 	err = fmt.Errorf("error parsing version byte: %w", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return err
	// }

	rawCommitmentHex, ok := mux.Vars(r)[routingVarNamePayloadHex]
	if !ok {
		return fmt.Errorf("commitment not found in path: %s", r.URL.Path)
	}
	commitment, err := hex.DecodeString(rawCommitmentHex)
	if err != nil {
		return fmt.Errorf("failed to decode hex commitment %s: %w", rawCommitmentHex, err)
	}
	// We use certV0 arbitrarily here, as it isn't used. Keccak commitments are not versioned.
	// TODO: We should probably create a new route for this which doesn't require a versionedCert.
	versionedCert := certs.NewVersionedCert(commitment, certs.V0VersionByte)

	return svr.handleGetShared(r.Context(), w, versionedCert, commitments.OptimismKeccakCommitmentMode)
}

// handleGetOPGenericCommitment handles the GET request for optimism generic commitments.
func (svr *Server) handleGetOPGenericCommitment(w http.ResponseWriter, r *http.Request) error {
	certVersion, err := parseCertVersion(w, r)
	if err != nil {
		return fmt.Errorf("error parsing version byte: %w", err)
	}
	serializedCertHex, ok := mux.Vars(r)[routingVarNamePayloadHex]
	if !ok {
		return fmt.Errorf("serializedDACert not found in path: %s", r.URL.Path)
	}
	commitment, err := hex.DecodeString(serializedCertHex)
	if err != nil {
		return fmt.Errorf("failed to decode from hex serializedDACert %s: %w", serializedCertHex, err)
	}
	versionedCert := certs.NewVersionedCert(commitment, certVersion)

	return svr.handleGetShared(r.Context(), w, versionedCert, commitments.OptimismGenericCommitmentMode)
}

func (svr *Server) handleGetShared(
	ctx context.Context,
	w http.ResponseWriter,
	versionedCert certs.VersionedCert,
	mode commitments.CommitmentMode,
) error {
	serializedCertHex := hex.EncodeToString(versionedCert.SerializedCert)
	svr.log.Info("Processing GET request", "commitmentMode", mode,
		"certVersion", versionedCert.Version, "serializedCert", serializedCertHex)
	input, err := svr.sm.Get(ctx, versionedCert, mode)
	if err != nil {
		err = NewGETError(
			fmt.Errorf("get request failed with serializedCert %v: %w", serializedCertHex, err),
			versionedCert.Version,
			mode,
		)
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	svr.writeResponse(w, input)
	return nil
}

// handleGetEigenDADispersalBackend handles the GET request to check the current EigenDA backend used for dispersal.
// This endpoint returns which EigenDA backend version (v1 or v2) is currently being used for blob dispersal.
func (svr *Server) handleGetEigenDADispersalBackend(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusOK)

	backend := svr.sm.GetDispersalBackend()
	backendString := common.EigenDABackendToString(backend)

	response := struct {
		EigenDADispersalBackend string `json:"eigenDADispersalBackend"`
	}{
		EigenDADispersalBackend: backendString,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return err
	}

	return nil
}

// =================================================================================================
// POST ROUTES
// =================================================================================================

// handlePostStdCommitment handles the POST request for std commitments.
func (svr *Server) handlePostStdCommitment(w http.ResponseWriter, r *http.Request) error {
	return svr.handlePostShared(w, r, nil, commitments.StandardCommitmentMode)
}

// handlePostOPKeccakCommitment handles the POST request for optimism keccak commitments.
func (svr *Server) handlePostOPKeccakCommitment(w http.ResponseWriter, r *http.Request) error {
	rawCommitmentHex, ok := mux.Vars(r)[routingVarNamePayloadHex]
	if !ok {
		return fmt.Errorf("commitment not found in path: %s", r.URL.Path)
	}
	commitment, err := hex.DecodeString(rawCommitmentHex)
	if err != nil {
		return fmt.Errorf("failed to decode commitment %s: %w", rawCommitmentHex, err)
	}

	return svr.handlePostShared(w, r, commitment, commitments.OptimismKeccakCommitmentMode)
}

// handlePostOPGenericCommitment handles the POST request for optimism generic commitments.
func (svr *Server) handlePostOPGenericCommitment(w http.ResponseWriter, r *http.Request) error {
	return svr.handlePostShared(w, r, nil, commitments.OptimismGenericCommitmentMode)
}

func (svr *Server) handlePostShared(
	w http.ResponseWriter,
	r *http.Request,
	comm []byte, // only non-nil for OPKeccak commitments
	mode commitments.CommitmentMode,
) error {
	svr.log.Info("Processing POST request", "commitment", hex.EncodeToString(comm), "mode", mode)
	input, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxRequestBodySize))
	if err != nil {
		err = NewPOSTError(fmt.Errorf("failed to read request body: %w", err), mode)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	serializedCert, err := svr.sm.Put(r.Context(), mode, comm, input)
	if err != nil {
		err = NewPOSTError(fmt.Errorf("post request failed with commitment %v: %w", comm, err), mode)
		switch {
		case is400(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case is429(err):
			http.Error(w, err.Error(), http.StatusTooManyRequests)
		case is503(err):
			// this tells the caller (batcher) to failover to ethda b/c eigenda is temporarily down
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	var certVersion certs.VersionByte
	switch svr.sm.GetDispersalBackend() {
	case common.V1EigenDABackend:
		certVersion = certs.V0VersionByte
	case common.V2EigenDABackend:
		certVersion = certs.V1VersionByte
	default:
		return fmt.Errorf("unknown dispersal backend: %v", svr.sm.GetDispersalBackend())
	}
	versionedCert := certs.NewVersionedCert(serializedCert, certVersion)

	responseCommit, err := commitments.EncodeCommitment(versionedCert, mode)
	if err != nil {
		err = NewPOSTError(fmt.Errorf("failed to encode serializedCert %v: %w", serializedCert, err), mode)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	svr.log.Info(fmt.Sprintf("response commitment: %x\n", responseCommit))
	// write commitment to resp body if not in OptimismKeccak mode
	if mode != commitments.OptimismKeccakCommitmentMode {
		svr.writeResponse(w, responseCommit)
	}
	return nil
}

// handleSetEigenDADispersalBackend handles the PUT request to set the EigenDA backend used for dispersal.
// This endpoint configures which EigenDA backend version (v1 or v2) will be used for blob dispersal.
func (svr *Server) handleSetEigenDADispersalBackend(w http.ResponseWriter, r *http.Request) error {
	// Read request body to get the new value
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1024)) // Small limit since we only expect a string
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request body: %v", err), http.StatusBadRequest)
		return err
	}

	// Parse the backend string value
	var setRequest struct {
		EigenDADispersalBackend string `json:"eigenDADispersalBackend"`
	}

	if err := json.Unmarshal(body, &setRequest); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse JSON request: %v", err), http.StatusBadRequest)
		return err
	}

	// Convert the string to EigenDABackend enum
	backend, err := common.StringToEigenDABackend(setRequest.EigenDADispersalBackend)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid eigenDADispersalBackend value: %v", err), http.StatusBadRequest)
		return err
	}

	svr.SetDispersalBackend(backend)

	// Return the current value in the response
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusOK)

	currentBackend := svr.sm.GetDispersalBackend()
	backendString := common.EigenDABackendToString(currentBackend)

	response := struct {
		EigenDADispersalBackend string `json:"eigenDADispersalBackend"`
	}{
		EigenDADispersalBackend: backendString,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return err
	}

	return nil
}
