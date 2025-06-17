//nolint:lll // long lines are expected in this file
package server

import (
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigenda-proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda-proxy/server/middleware"
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
