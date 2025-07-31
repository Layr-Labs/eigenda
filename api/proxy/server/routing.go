//nolint:lll // long lines are expected in this file
package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/api/proxy/common/proxyerrors"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/server/middleware"
	"github.com/gorilla/mux"
)

const (
	routingVarNameKeccakCommitmentHex = "keccak_commitment_hex"
	routingVarNamePayloadHex          = "payload_hex"
	routingVarNameVersionByteHex      = "version_byte_hex"
	routingVarNameCommitTypeByteHex   = "commit_type_byte_hex"
)

func (svr *Server) RegisterRoutes(r *mux.Router) {
	subrouterGET := r.Methods("GET").PathPrefix("/get").Subrouter()
	// std commitments (for nitro)
	subrouterGET.HandleFunc("/"+
		"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
		"{"+routingVarNameVersionByteHex+":[0-9a-fA-F]{2}}"+ // should always be 0x00 for now but we let others through to return a 404
		"{"+routingVarNamePayloadHex+":[0-9a-fA-F]*}",
		middleware.WithCertMiddlewares(svr.handleGetStdCommitment, svr.log, svr.m, commitments.StandardCommitmentMode),
	).Queries("commitment_mode", "standard")
	// op keccak256 commitments (write to S3)
	subrouterGET.HandleFunc(
		"/"+
			"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
			"{"+routingVarNameCommitTypeByteHex+":00}"+ // 00 for keccak256 commitments
			"{"+routingVarNameKeccakCommitmentHex+"}",
		middleware.WithCertMiddlewares(
			svr.handleGetOPKeccakCommitment,
			svr.log,
			svr.m,
			commitments.OptimismKeccakCommitmentMode,
		),
	)
	// op generic commitments (write to EigenDA)
	subrouterGET.HandleFunc(
		"/"+
			"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
			"{"+routingVarNameCommitTypeByteHex+":01}"+ // 01 for generic commitments
			"{da_layer_byte:[0-9a-fA-F]{2}}"+ // should always be 0x00 for eigenDA but we let others through to return a 404
			"{"+routingVarNameVersionByteHex+":[0-9a-fA-F]{2}}"+ // should always be 0x00 for now but we let others through to return a 404
			"{"+routingVarNamePayloadHex+"}",
		middleware.WithCertMiddlewares(
			svr.handleGetOPGenericCommitment,
			svr.log,
			svr.m,
			commitments.OptimismGenericCommitmentMode,
		),
	)
	// unrecognized op commitment type (not 00 or 01)
	subrouterGET.HandleFunc("/"+
		"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
		"{"+routingVarNameCommitTypeByteHex+":[0-9a-fA-F]{2}}",
		func(w http.ResponseWriter, r *http.Request) {
			svr.log.Info(
				"unsupported commitment type",
				routingVarNameCommitTypeByteHex,
				mux.Vars(r)[routingVarNameCommitTypeByteHex],
			)
			commitType := mux.Vars(r)[routingVarNameCommitTypeByteHex]
			http.Error(w, fmt.Sprintf("unsupported commitment type %s", commitType), http.StatusBadRequest)
		},
	).MatcherFunc(notCommitmentModeStandard)

	subrouterPOST := r.Methods("POST").PathPrefix("/put").Subrouter()
	// std commitments (for nitro)
	subrouterPOST.HandleFunc("", // commitment is calculated by the server using the body data
		middleware.WithCertMiddlewares(svr.handlePostStdCommitment, svr.log, svr.m, commitments.StandardCommitmentMode),
	).Queries("commitment_mode", "standard")
	// op keccak256 commitments (write to S3)
	subrouterPOST.HandleFunc(
		"/"+
			"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
			"{"+routingVarNameCommitTypeByteHex+":00}"+ // 00 for keccak256 commitments
			"{"+routingVarNameKeccakCommitmentHex+"}",
		middleware.WithCertMiddlewares(
			svr.handlePostOPKeccakCommitment,
			svr.log,
			svr.m,
			commitments.OptimismKeccakCommitmentMode,
		),
	)
	// op generic commitments (write to EigenDA)
	subrouterPOST.HandleFunc(
		"", // commitment is calculated by the server using the body data
		middleware.WithCertMiddlewares(
			svr.handlePostOPGenericCommitment,
			svr.log,
			svr.m,
			commitments.OptimismGenericCommitmentMode,
		),
	)
	subrouterPOST.HandleFunc(
		"/", // commitment is calculated by the server using the body data
		middleware.WithCertMiddlewares(
			svr.handlePostOPGenericCommitment,
			svr.log,
			svr.m,
			commitments.OptimismGenericCommitmentMode,
		),
	)

	// TODO: should prob setup metrics middlewares to also work for the below routes...
	// right now they only work for the main GET/POST routes.
	r.HandleFunc("/health", svr.handleHealth).Methods("GET")

	// this is done to explicitly log capture potential redirect errors
	r.HandleFunc("/put", svr.logDispersalGetError).Methods("GET")

	// Only register admin endpoints if explicitly enabled in configuration
	//
	// Note: A common pattern for admin endpoints is to generate a random API key on startup for authentication.
	// Since the proxy isn't meant to be exposed publicly, we haven't implemented this here, but it's something
	// that might be done in the future.
	if svr.config.IsAPIEnabled(AdminAPIType) {
		svr.log.Warn("Admin API endpoints are enabled")
		// Admin endpoints to check and set EigenDA backend used for dispersal
		r.HandleFunc("/admin/eigenda-dispersal-backend", svr.handleGetEigenDADispersalBackend).Methods("GET")
		r.HandleFunc("/admin/eigenda-dispersal-backend", svr.handleSetEigenDADispersalBackend).Methods("PUT")
	}
}

func notCommitmentModeStandard(r *http.Request, _ *mux.RouteMatch) bool {
	commitmentMode := r.URL.Query().Get("commitment_mode")
	return commitmentMode == "" || commitmentMode != "standard"
}

// ================== QUERY PARAMS PARSING FUNCTION ==================================================
// These query params don't affect routing, but we keep them here so that everything related to query URLs is in one place,
// and its easy to deduct what kind of queries are supported by the proxy server by just looking at this file.
// The below 2 functions are used in both standard and optimism routes (see handlers_cert.go).

// Parses the l1_inclusion_block_number query param from the request.
// Happy path:
//   - if the l1_inclusion_block_number is provided, it returns the parsed value.
//
// Unhappy paths:
//   - if the l1_inclusion_block_number is not provided, it returns 0 (whose meaning is to skip the check).
//   - if the l1_inclusion_block_number is provided but isn't a valid integer, it returns a [proxyerrors.L1InclusionBlockNumberParsingError].
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

// Parses the return_encoded_payload query parameter from the request (use the first value if multiple are provided).
// Returns true for: ?return_encoded_payload, ?return_encoded_payload=true, ?return_encoded_payload=1
// Anything else returns false, including if the parameter is not present.
func parseReturnEncodedPayloadQueryParam(r *http.Request) bool {
	returnEncodedPayloadValues, exists := r.URL.Query()["return_encoded_payload"]
	if !exists || len(returnEncodedPayloadValues) == 0 {
		return false
	}
	fmt.Println("returnEncodedPayloadValues:", returnEncodedPayloadValues)
	returnEncodedPayload := strings.ToLower(returnEncodedPayloadValues[0])
	if returnEncodedPayload == "" || returnEncodedPayload == "true" || returnEncodedPayload == "1" {
		return true
	}
	return false
}
