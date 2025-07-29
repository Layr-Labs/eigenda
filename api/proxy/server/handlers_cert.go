// handlers_cert.go contains the main HTTP handlers for the Eigenda Proxy server.
// These are the handlers that process POST (payload->commitment) and GET (commitment->payload) requests.
// Handlers in this file SHOULD be wrapped in middlewares.
package server

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/proxyerrors"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/server/middleware"
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
	keccakCommitmentHex, ok := mux.Vars(r)[routingVarNameKeccakCommitmentHex]
	if !ok {
		return proxyerrors.NewParsingError(fmt.Errorf("keccak commitment not found in path: %s", r.URL.Path))
	}
	keccakCommitment, err := hex.DecodeString(keccakCommitmentHex)
	if err != nil {
		return proxyerrors.NewParsingError(
			fmt.Errorf("failed to decode hex keccak commitment %s: %w", keccakCommitmentHex, err))
	}
	payload, err := svr.sm.GetOPKeccakValueFromS3(r.Context(), keccakCommitment)
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
	return svr.handleGetShared(w, r, commitments.OptimismGenericCommitmentMode)
}

// handleGetStdCommitment handles the GET request for std commitments.
func (svr *Server) handleGetStdCommitment(w http.ResponseWriter, r *http.Request) error {
	return svr.handleGetShared(w, r, commitments.StandardCommitmentMode)
}

func (svr *Server) handleGetShared(
	w http.ResponseWriter,
	r *http.Request,
	mode commitments.CommitmentMode,
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
	// This is currently used by secure integration (e.g. optimism hokulea), which need
	// to decode the payload themselves inside the fpvm.
	returnEncodedPayload := parseReturnEncodedPayloadQueryParam(r)

	maybeEncodedPayload, err := svr.sm.Get(
		r.Context(),
		versionedCert,
		mode,
		common.GETOpts{
			L1InclusionBlockNum:  l1InclusionBlockNum,
			ReturnEncodedPayload: returnEncodedPayload,
		},
	)
	if err != nil {
		return fmt.Errorf("get request failed with serializedCert (version %v) %v: %w",
			versionedCert.Version, serializedCertHex, err)
	}

	svr.log.Info("Processed request", "method", r.Method, "url", r.URL.Path,
		"commitmentMode", mode, "returnEncodedPayload", returnEncodedPayload,
		"certVersion", versionedCert.Version, "serializedCert", serializedCertHex)

	_, err = w.Write(maybeEncodedPayload)
	if err != nil {
		// If the write fails, we will already have sent a 200 header. But we still return an error
		// here so that the logging middleware can log it.
		return fmt.Errorf("failed to write response for GET serializedCert (version %v) %v: %w",
			versionedCert.Version, serializedCertHex, err)
	}
	return nil
}

// Parses the l1_inclusion_block_number query param from the request.
// Happy path:
//   - if the l1_inclusion_block_number is provided, it returns the parsed value.
//
// Unhappy paths:
//   - if the l1_inclusion_block_number is not provided, it returns 0 (whose meaning is to skip the check).
//   - if the l1_inclusion_block_number is provided but isn't a valid integer, it returns an error.
func parseCommitmentInclusionL1BlockNumQueryParam(r *http.Request) (uint64, error) {
	l1BlockNumStr := r.URL.Query().Get("l1_inclusion_block_number")
	if l1BlockNumStr != "" {
		l1BlockNum, err := strconv.ParseUint(l1BlockNumStr, 10, 64)
		if err != nil {
			return 0, proxyerrors.NewL1InclusionBlockNumberParsingError(l1BlockNumStr, err)
		}
		return l1BlockNum, nil
	}
	return 0, nil
}

// Parses the return_encoded_payload query parameter from the request.
// Happy path:
//   - if the return_encoded_payload query parameter is present (with any value), it returns true
//   - this means it accepts ?return_encoded_payload, ?return_encoded_payload=true, ?return_encoded_payload=anything
//
// Unhappy paths:
//   - if the return_encoded_payload query parameter is not provided, it returns false
func parseReturnEncodedPayloadQueryParam(r *http.Request) bool {
	_, exists := r.URL.Query()["return_encoded_payload"]
	return exists
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

	err = svr.sm.PutOPKeccakPairInS3(r.Context(), keccakCommitment, payload)
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
	return svr.handlePostShared(w, r, commitments.StandardCommitmentMode)
}

// handlePostOPGenericCommitment handles the POST request for optimism generic commitments.
func (svr *Server) handlePostOPGenericCommitment(w http.ResponseWriter, r *http.Request) error {
	return svr.handlePostShared(w, r, commitments.OptimismGenericCommitmentMode)
}

// This is a shared function for handling POST requests for
func (svr *Server) handlePostShared(
	w http.ResponseWriter,
	r *http.Request,
	mode commitments.CommitmentMode,
) error {
	payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxPOSTRequestBodySize))
	if err != nil {
		return proxyerrors.NewReadRequestBodyError(err, maxPOSTRequestBodySize)
	}

	serializedCert, err := svr.sm.Put(r.Context(), mode, payload)
	if err != nil {
		return fmt.Errorf("post request failed: %w", err)
	}

	var certVersion certs.VersionByte
	switch svr.sm.GetDispersalBackend() {
	case common.V1EigenDABackend:
		certVersion = certs.V0VersionByte
	case common.V2EigenDABackend:
		certVersion = certs.V2VersionByte
	default:
		return fmt.Errorf("unknown dispersal backend: %v", svr.sm.GetDispersalBackend())
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
