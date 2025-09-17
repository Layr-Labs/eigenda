// handlers_cert.go contains the main HTTP handlers for the Eigenda Proxy server.
// These are the handlers that process POST (payload->commitment) and GET (commitment->payload) requests.
// Handlers in this file SHOULD be wrapped in middlewares.
package rest

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/proxyerrors"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/rest/middleware"
	"github.com/gorilla/mux"
)

const (
	// limit requests to only 32 MiB to mitigate potential DoS attacks
	maxPOSTRequestBodySize int64 = 1024 * 1024 * 32
)

// =================================================================================================
// GET ROUTES
// =================================================================================================

// handleGetOPKeccakCommitment handles GET requests for optimism keccak commitments.
func (svr *Server) handleGetOPKeccakCommitment(w http.ResponseWriter, r *http.Request) error {
	if !svr.config.APIsEnabled.OpKeccakCommitment {
		w.WriteHeader(http.StatusForbidden)
		return fmt.Errorf("op-keccak DA Commitment type detected but `op-keccak` API is not enabled")
	}

	keccakCommitmentHex, ok := mux.Vars(r)[routingVarNameKeccakCommitmentHex]
	if !ok {
		return proxyerrors.NewParsingError(fmt.Errorf("keccak commitment not found in path: %s", r.URL.Path))
	}
	keccakCommitment, err := hex.DecodeString(keccakCommitmentHex)
	if err != nil {
		return proxyerrors.NewParsingError(
			fmt.Errorf("failed to decode hex keccak commitment %s: %w", keccakCommitmentHex, err))
	}
	payload, err := svr.keccakMgr.GetOPKeccakValueFromS3(r.Context(), keccakCommitment)
	if err != nil {
		return fmt.Errorf("GET keccakCommitment %v: %w", keccakCommitmentHex, err)
	}

	svr.log.Info("Processed request", "method", r.Method, "url", r.URL.Path,
		"commitmentMode", commitments.OptimismKeccakCommitmentMode, "commitment", keccakCommitmentHex)

	_, err = w.Write(payload)
	if err != nil {
		// If the write fails, we will already have sent a 200 header. But we still return an error
		// here so that the logging middleware can log it.
		return fmt.Errorf("failed to write response for GET keccakCommitment %v: %w", keccakCommitmentHex, err)
	}
	return nil
}

// handleGetOPGenericCommitment handles the GET request for optimism generic commitments.
func (svr *Server) handleGetOPGenericCommitment(w http.ResponseWriter, r *http.Request) error {
	if !svr.config.APIsEnabled.OpGenericCommitment {
		w.WriteHeader(http.StatusForbidden)
		return fmt.Errorf("op-generic DA Commitment type detected but `op-generic` API is not enabled")
	}

	return svr.handleGetShared(w, r)
}

// handleGetStdCommitment handles the GET request for std commitments.
func (svr *Server) handleGetStdCommitment(w http.ResponseWriter, r *http.Request) error {
	if !svr.config.APIsEnabled.StandardCommitment {
		w.WriteHeader(http.StatusForbidden)
		return fmt.Errorf("standard DA Commitment type detected but `standard` API is not enabled")
	}

	return svr.handleGetShared(w, r)
}

func (svr *Server) handleGetShared(
	w http.ResponseWriter,
	r *http.Request,
) error {
	certVersion, err := parseCertVersion(w, r)
	if err != nil {
		return proxyerrors.NewParsingError(fmt.Errorf("parsing version byte: %w", err))
	}
	// used in the metrics middleware... there's prob a better way to do this
	middleware.SetCertVersion(r, string(certVersion))
	serializedCertHex, ok := mux.Vars(r)[routingVarNamePayloadHex]
	if !ok {
		return proxyerrors.NewParsingError(fmt.Errorf("serializedDACert not found in path: %s", r.URL.Path))
	}
	serializedCert, err := hex.DecodeString(serializedCertHex)
	if err != nil {
		return proxyerrors.NewCertHexDecodingError(serializedCertHex, err)
	}
	versionedCert := certs.NewVersionedCert(serializedCert, certVersion)

	l1InclusionBlockNum, err := parseCommitmentInclusionL1BlockNumQueryParam(r)
	if err != nil {
		return err // doesn't need to be wrapped; already a proxyerrors
	}

	// Check if client requested encoded payload
	// This is currently used by secure integrations (e.g. optimism hokulea), which need
	// to decode the payload themselves inside the fpvm.
	returnEncodedPayload := parseReturnEncodedPayloadQueryParam(r)

	payloadOrEncodedPayload, err := svr.certMgr.Get(
		r.Context(),
		versionedCert,
		common.GETOpts{
			L1InclusionBlockNum:  l1InclusionBlockNum,
			ReturnEncodedPayload: returnEncodedPayload,
		},
	)
	if err != nil {
		return fmt.Errorf("get request failed with serializedCert (version %v) %v: %w",
			versionedCert.Version, serializedCertHex, err)
	}

	svr.log.Info("Processed request", "method", r.Method, "url", r.URL.Path, "returnEncodedPayload", returnEncodedPayload,
		"certVersion", versionedCert.Version, "serializedCert", serializedCertHex)

	_, err = w.Write(payloadOrEncodedPayload)
	if err != nil {
		// If the write fails, we will already have sent a 200 header. But we still return an error
		// here so that the logging middleware can log it.
		return fmt.Errorf("failed to write response for GET serializedCert (version %v) %v: %w",
			versionedCert.Version, serializedCertHex, err)
	}
	return nil
}

// =================================================================================================
// POST ROUTES
// =================================================================================================

// handlePostOPKeccakCommitment handles the POST request for optimism keccak commitments.
func (svr *Server) handlePostOPKeccakCommitment(w http.ResponseWriter, r *http.Request) error {
	keccakCommitmentHex, ok := mux.Vars(r)[routingVarNameKeccakCommitmentHex]
	if !ok {
		return proxyerrors.NewParsingError(fmt.Errorf("keccak commitment not found in path: %s", r.URL.Path))
	}
	keccakCommitment, err := hex.DecodeString(keccakCommitmentHex)
	if err != nil {
		return proxyerrors.NewParsingError(
			fmt.Errorf("failed to decode hex keccak commitment %s: %w", keccakCommitmentHex, err))
	}
	payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxPOSTRequestBodySize))
	if err != nil {
		return proxyerrors.NewReadRequestBodyError(err, maxPOSTRequestBodySize)
	}

	err = svr.keccakMgr.PutOPKeccakPairInS3(r.Context(), keccakCommitment, payload)
	if err != nil {
		return fmt.Errorf("keccak POST request failed for commitment %v: %w", keccakCommitmentHex, err)
	}

	svr.log.Info("Processed request", "method", r.Method, "url", r.URL.Path,
		"commitmentMode", commitments.OptimismKeccakCommitmentMode, "commitment", keccakCommitmentHex)
	// No need to return the keccak commitment because it's already known by the client (keccak(payload)).
	return nil
}

// handlePostStdCommitment handles the POST request for std commitments.
func (svr *Server) handlePostStdCommitment(w http.ResponseWriter, r *http.Request) error {
	if !svr.config.APIsEnabled.StandardCommitment {
		w.WriteHeader(http.StatusForbidden)
		return fmt.Errorf("standard DA Commitment type detected but `standard` API is not enabled")
	}

	return svr.handlePostShared(w, r, commitments.StandardCommitmentMode)
}

// handlePostOPGenericCommitment handles the POST request for optimism generic commitments.
func (svr *Server) handlePostOPGenericCommitment(w http.ResponseWriter, r *http.Request) error {
	if !svr.config.APIsEnabled.OpGenericCommitment {
		w.WriteHeader(http.StatusForbidden)
		return fmt.Errorf("op-generic DA Commitment type detected but `op-generic` API is not enabled")
	}

	return svr.handlePostShared(w, r, commitments.OptimismGenericCommitmentMode)
}

// This is a shared function for handling POST requests for
func (svr *Server) handlePostShared(
	w http.ResponseWriter,
	r *http.Request,
	mode commitments.CommitmentMode,
) error {
	if !svr.config.APIsEnabled.StandardCommitment && mode == commitments.StandardCommitmentMode {
		w.WriteHeader(http.StatusForbidden)
		return fmt.Errorf("standard DA Commitment type detected but `standard` API is not enabled")
	}

	payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxPOSTRequestBodySize))
	if err != nil {
		return proxyerrors.NewReadRequestBodyError(err, maxPOSTRequestBodySize)
	}

	serializedCert, err := svr.certMgr.Put(r.Context(), payload)
	if err != nil {
		return fmt.Errorf("post request failed: %w", err)
	}

	var certVersion certs.VersionByte
	switch svr.certMgr.GetDispersalBackend() {
	case common.V1EigenDABackend:
		certVersion = certs.V0VersionByte
	case common.V2EigenDABackend:
		certVersion = certs.V2VersionByte
	default:
		return fmt.Errorf("unknown dispersal backend: %v", svr.certMgr.GetDispersalBackend())
	}
	versionedCert := certs.NewVersionedCert(serializedCert, certVersion)

	responseCommit, err := commitments.EncodeCommitment(versionedCert, mode)
	if err != nil {
		// This error is only possible if we have a bug in the code.
		return fmt.Errorf("failed to encode serializedCert %v: %w", serializedCert, err)
	}

	svr.log.Info("Processed request", "method", r.Method, "url", r.URL.Path, "commitmentMode", mode,
		"certVersion", versionedCert.Version, "cert", hex.EncodeToString(serializedCert))

	// We write the commitment as bytes directly instead of hex encoded.
	// The spec https://specs.optimism.io/experimental/alt-da.html#da-server says it should be hex-encoded,
	// but the client expects it to be raw bytes.
	// See
	// https://github.com/Layr-Labs/optimism/blob/89ac40d0fddba2e06854b253b9f0266f36350af2/op-alt-da/daclient.go#L151
	_, err = w.Write(responseCommit)
	if err != nil {
		// If the write fails, we will already have sent a 200 header. But we still return an error
		// here so that the logging middleware can log it.
		return fmt.Errorf("failed to write response for POST serializedCert (version %v) %x: %w",
			versionedCert.Version, serializedCert, err)
	}
	return nil
}
