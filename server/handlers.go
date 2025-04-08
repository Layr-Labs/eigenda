package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/common"
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

// =================================================================================================
// GET ROUTES
// =================================================================================================

// handleGetStdCommitment handles the GET request for std commitments.
func (svr *Server) handleGetStdCommitment(w http.ResponseWriter, r *http.Request) error {
	versionByte, err := parseVersionByte(w, r)
	if err != nil {
		return fmt.Errorf("error parsing version byte: %w", err)
	}
	commitmentMeta := commitments.CommitmentMeta{
		Mode:    commitments.Standard,
		Version: commitments.EigenDACommitmentType(versionByte),
	}

	rawCommitmentHex, ok := mux.Vars(r)[routingVarNamePayloadHex]
	if !ok {
		return fmt.Errorf("commitment not found in path: %s", r.URL.Path)
	}
	commitment, err := hex.DecodeString(rawCommitmentHex)
	if err != nil {
		return fmt.Errorf("failed to decode commitment %s: %w", rawCommitmentHex, err)
	}

	return svr.handleGetShared(r.Context(), w, commitment, commitmentMeta)
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
	commitmentMeta := commitments.CommitmentMeta{
		Mode:    commitments.OptimismKeccak,
		Version: commitments.CertV0,
	}

	rawCommitmentHex, ok := mux.Vars(r)[routingVarNamePayloadHex]
	if !ok {
		return fmt.Errorf("commitment not found in path: %s", r.URL.Path)
	}
	commitment, err := hex.DecodeString(rawCommitmentHex)
	if err != nil {
		return fmt.Errorf("failed to decode commitment %s: %w", rawCommitmentHex, err)
	}

	return svr.handleGetShared(r.Context(), w, commitment, commitmentMeta)
}

// handleGetOPGenericCommitment handles the GET request for optimism generic commitments.
func (svr *Server) handleGetOPGenericCommitment(w http.ResponseWriter, r *http.Request) error {
	versionByte, err := parseVersionByte(w, r)
	if err != nil {
		return fmt.Errorf("error parsing version byte: %w", err)
	}
	commitmentMeta := commitments.CommitmentMeta{
		Mode:    commitments.OptimismGeneric,
		Version: commitments.EigenDACommitmentType(versionByte),
	}

	rawCommitmentHex, ok := mux.Vars(r)[routingVarNamePayloadHex]
	if !ok {
		return fmt.Errorf("commitment not found in path: %s", r.URL.Path)
	}
	commitment, err := hex.DecodeString(rawCommitmentHex)
	if err != nil {
		return fmt.Errorf("failed to decode commitment %s: %w", rawCommitmentHex, err)
	}

	return svr.handleGetShared(r.Context(), w, commitment, commitmentMeta)
}

func (svr *Server) handleGetShared(
	ctx context.Context,
	w http.ResponseWriter,
	comm []byte,
	meta commitments.CommitmentMeta,
) error {
	commitmentHex := hex.EncodeToString(comm)
	svr.log.Info("Processing GET request", "commitment", commitmentHex, "commitmentMeta", meta)
	input, err := svr.sm.Get(ctx, comm, meta)
	if err != nil {
		err = MetaError{
			Err:  fmt.Errorf("get request failed with commitment %v: %w", commitmentHex, err),
			Meta: meta,
		}
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

// disperseToV2ToEigenDABackend converts the boolean disperseToV2 flag to the corresponding EigenDABackend enum
func disperseToV2ToEigenDABackend(disperseToV2 bool) common.EigenDABackend {
	if disperseToV2 {
		return common.V2EigenDABackend
	}
	return common.V1EigenDABackend
}

// eigenDABackendToDisperseToV2 converts an EigenDABackend enum to the corresponding boolean flag
func eigenDABackendToDisperseToV2(backend common.EigenDABackend) bool {
	return backend == common.V2EigenDABackend
}

// handleGetEigenDADispersalBackend handles the GET request to check the current EigenDA backend used for dispersal.
// This endpoint returns which EigenDA backend version (v1 or v2) is currently being used for blob dispersal.
func (svr *Server) handleGetEigenDADispersalBackend(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusOK)

	backend := disperseToV2ToEigenDABackend(svr.sm.DisperseToV2())
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
	commitmentMeta := commitments.CommitmentMeta{
		Mode:    commitments.Standard,
		Version: commitments.CertV0,
	}

	if svr.sm.DisperseToV2() {
		commitmentMeta.Version = commitments.CertV1
	}

	return svr.handlePostShared(w, r, nil, commitmentMeta)
}

// handlePostOPKeccakCommitment handles the POST request for optimism keccak commitments.
func (svr *Server) handlePostOPKeccakCommitment(w http.ResponseWriter, r *http.Request) error {
	// TODO: do we use a version byte in OPKeccak commitments? README seems to say so, but server_test didn't
	// versionByte, err := parseVersionByte(r)
	// if err != nil {
	// 	err = fmt.Errorf("error parsing version byte: %w", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return err
	// }
	commitmentMeta := commitments.CommitmentMeta{
		Mode:    commitments.OptimismKeccak,
		Version: commitments.CertV0,
	}

	rawCommitmentHex, ok := mux.Vars(r)[routingVarNamePayloadHex]
	if !ok {
		return fmt.Errorf("commitment not found in path: %s", r.URL.Path)
	}
	commitment, err := hex.DecodeString(rawCommitmentHex)
	if err != nil {
		return fmt.Errorf("failed to decode commitment %s: %w", rawCommitmentHex, err)
	}

	return svr.handlePostShared(w, r, commitment, commitmentMeta)
}

// handlePostOPGenericCommitment handles the POST request for optimism generic commitments.
func (svr *Server) handlePostOPGenericCommitment(w http.ResponseWriter, r *http.Request) error {
	commitmentMeta := commitments.CommitmentMeta{
		Mode:    commitments.OptimismGeneric,
		Version: commitments.CertV0,
	}

	if svr.sm.DisperseToV2() {
		commitmentMeta.Version = commitments.CertV1
	}

	return svr.handlePostShared(w, r, nil, commitmentMeta)
}

func (svr *Server) handlePostShared(
	w http.ResponseWriter,
	r *http.Request,
	comm []byte,
	meta commitments.CommitmentMeta,
) error {
	svr.log.Info("Processing POST request", "commitment", hex.EncodeToString(comm), "meta", meta)
	input, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxRequestBodySize))
	if err != nil {
		err = MetaError{
			Err:  fmt.Errorf("failed to read request body: %w", err),
			Meta: meta,
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	commitment, err := svr.sm.Put(r.Context(), meta.Mode, comm, input)
	if err != nil {
		err = MetaError{
			Err:  fmt.Errorf("put request failed with commitment %v (commitment mode %v): %w", comm, meta.Mode, err),
			Meta: meta,
		}
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

	responseCommit, err := commitments.EncodeCommitment(commitment, meta.Mode, meta.Version)
	if err != nil {
		err = MetaError{
			Err:  fmt.Errorf("failed to encode commitment %v (commitment mode %v): %w", commitment, meta.Mode, err),
			Meta: meta,
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	svr.log.Info(fmt.Sprintf("response commitment: %x\n", responseCommit))
	// write commitment to resp body if not in OptimismKeccak mode
	if meta.Mode != commitments.OptimismKeccak {
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

	disperseToV2 := eigenDABackendToDisperseToV2(backend)
	svr.SetDisperseToV2(disperseToV2)

	// Return the current value in the response
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusOK)

	currentBackend := disperseToV2ToEigenDABackend(svr.sm.DisperseToV2())
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
