package server

import (
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/gorilla/mux"
)

const (
	routingVarNameRawCommitmentHex  = "raw_commitment_hex"
	routingVarNameVersionByteHex    = "version_byte_hex"
	routingVarNameCommitTypeByteHex = "commit_type_byte_hex"
)

func (svr *Server) registerRoutes(r *mux.Router) {
	subrouterGET := r.Methods("GET").PathPrefix("/get").Subrouter()
	// simple commitments (for nitro)
	subrouterGET.HandleFunc("/"+
		"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
		"{"+routingVarNameVersionByteHex+":[0-9a-fA-F]{2}}"+ // should always be 0x00 for now but we let others through to return a 404
		"{"+routingVarNameRawCommitmentHex+":[0-9a-fA-F]*}",
		withLogging(withMetrics(svr.handleGetSimpleCommitment, svr.m, commitments.SimpleCommitmentMode), svr.log),
	).Queries("commitment_mode", "simple")
	// op keccak256 commitments (write to S3)
	subrouterGET.HandleFunc("/"+
		"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
		"{"+routingVarNameCommitTypeByteHex+":00}"+ // 00 for keccak256 commitments
		// we don't use version_byte for keccak commitments, because not expecting keccak commitments to change,
		// but perhaps we should (in case we want a v2 to use another hash for eg?)
		// "{version_byte_hex:[0-9a-fA-F]{2}}"+ // should always be 0x00 for now but we let others through to return a 404
		"{"+routingVarNameRawCommitmentHex+"}",
		withLogging(withMetrics(svr.handleGetOPKeccakCommitment, svr.m, commitments.OptimismKeccak), svr.log),
	)
	// op generic commitments (write to EigenDA)
	subrouterGET.HandleFunc("/"+
		"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
		"{"+routingVarNameCommitTypeByteHex+":01}"+ // 01 for generic commitments
		"{da_layer_byte:[0-9a-fA-F]{2}}"+ // should always be 0x00 for eigenDA but we let others through to return a 404
		"{"+routingVarNameVersionByteHex+":[0-9a-fA-F]{2}}"+ // should always be 0x00 for now but we let others through to return a 404
		"{"+routingVarNameRawCommitmentHex+"}",
		withLogging(withMetrics(svr.handleGetOPGenericCommitment, svr.m, commitments.OptimismGeneric), svr.log),
	)
	// unrecognized op commitment type (not 00 or 01)
	subrouterGET.HandleFunc("/"+
		"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
		"{"+routingVarNameCommitTypeByteHex+":[0-9a-fA-F]{2}}",
		func(w http.ResponseWriter, r *http.Request) {
			svr.log.Info("unsupported commitment type", routingVarNameCommitTypeByteHex, mux.Vars(r)[routingVarNameCommitTypeByteHex])
			commitType := mux.Vars(r)[routingVarNameCommitTypeByteHex]
			http.Error(w, fmt.Sprintf("unsupported commitment type %s", commitType), http.StatusBadRequest)
		},
	).MatcherFunc(notCommitmentModeSimple)

	subrouterPOST := r.Methods("POST").PathPrefix("/put").Subrouter()
	// simple commitments (for nitro)
	subrouterPOST.HandleFunc("", // commitment is calculated by the server using the body data
		withLogging(withMetrics(svr.handlePostSimpleCommitment, svr.m, commitments.SimpleCommitmentMode), svr.log),
	).Queries("commitment_mode", "simple")
	// op keccak256 commitments (write to S3)
	subrouterPOST.HandleFunc("/"+
		"{optional_prefix:(?:0x)?}"+ // commitments can be prefixed with 0x
		// TODO: do we need this 0x00 byte? keccak commitments are the only ones that have anything after /put/
		"{"+routingVarNameCommitTypeByteHex+":00}"+ // 00 for keccak256 commitments
		// we don't use version_byte for keccak commitments, because not expecting keccak commitments to change,
		// but perhaps we should (in case we want a v2 to use another hash for eg?)
		// "{version_byte_hex:[0-9a-fA-F]{2}}"+ // should always be 0x00 for now but we let others through to return a 404
		"{"+routingVarNameRawCommitmentHex+"}",
		withLogging(withMetrics(svr.handlePostOPKeccakCommitment, svr.m, commitments.OptimismKeccak), svr.log),
	)
	// op generic commitments (write to EigenDA)
	subrouterPOST.HandleFunc("", // commitment is calculated by the server using the body data
		withLogging(withMetrics(svr.handlePostOPGenericCommitment, svr.m, commitments.OptimismGeneric), svr.log),
	)
	subrouterPOST.HandleFunc("/", // commitment is calculated by the server using the body data
		withLogging(withMetrics(svr.handlePostOPGenericCommitment, svr.m, commitments.OptimismGeneric), svr.log),
	)

	r.HandleFunc("/health", withLogging(svr.handleHealth, svr.log)).Methods("GET")
}

func notCommitmentModeSimple(r *http.Request, _ *mux.RouteMatch) bool {
	commitmentMode := r.URL.Query().Get("commitment_mode")
	return commitmentMode == "" || commitmentMode != "simple"
}
