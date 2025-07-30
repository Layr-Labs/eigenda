package memconfig

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"
)

// NullableDerivationError is a custom type for managing the PutWithGetReturnsDerivationError configuration.
// It allows users to distinguish between three states:
// 1. Field omitted from JSON: no change to current configuration
// 2. Reset=false with embedded DerivationError: sets the derivation error to the embedded values
// 3. Reset=true: resets the derivation error to nil (disabled)
//
// Usage examples:
// - To set an error: {"PutWithGetReturnsDerivationError": {"StatusCode": 3, "Msg": "test error", "Reset": false}}
// - To reset to nil: {"PutWithGetReturnsDerivationError": {"Reset": true}}
// - To leave unchanged: omit the field entirely from the JSON request
type NullableDerivationError struct {
	// Embed the DerivationError directly. Only used when Reset=false.
	coretypes.DerivationError
	// Reset indicates the user's intent:
	// - true: reset PutWithGetReturnsDerivationError to nil (disabled)
	// - false: set PutWithGetReturnsDerivationError to the embedded DerivationError
	Reset bool `json:"Reset"`
}

// JSON bodies received by the PATCH /memstore/config endpoint are deserialized into this struct,
// which is then used to update the memstore configuration.
type ConfigUpdate struct {
	MaxBlobSizeBytes                 *uint64                  `json:"MaxBlobSizeBytes,omitempty"`
	PutLatency                       *string                  `json:"PutLatency,omitempty"`
	GetLatency                       *string                  `json:"GetLatency,omitempty"`
	PutReturnsFailoverError          *bool                    `json:"PutReturnsFailoverError,omitempty"`
	BlobExpiration                   *string                  `json:"BlobExpiration,omitempty"`
	PutWithGetReturnsDerivationError *NullableDerivationError `json:"PutWithGetReturnsDerivationError,omitempty"`
}

// HandlerHTTP is an admin HandlerHTTP for GETting and PATCHing the memstore configuration.
// It adds routes to the proxy's main router (to be served on same port as the main proxy routes):
// - GET /memstore/config: returns the current memstore configuration
// - PATCH /memstore/config: updates the memstore configuration
type HandlerHTTP struct {
	log        logging.Logger
	safeConfig *SafeConfig
}

func NewHandlerHTTP(log logging.Logger, safeConfig *SafeConfig) HandlerHTTP {
	return HandlerHTTP{
		log:        log,
		safeConfig: safeConfig,
	}
}

func (api HandlerHTTP) RegisterMemstoreConfigHandlers(r *mux.Router) {
	memstore := r.PathPrefix("/memstore").Subrouter()
	memstore.HandleFunc("/config", api.handleGetConfig).Methods("GET")
	memstore.HandleFunc("/config", api.handleUpdateConfig).Methods("PATCH")
}

// Returns the config of the memstore in json format.
// TODO: we prob want to use out custom Duration type instead of time.Duration
// since time.Duration serializes to nanoseconds, which is hard to read.
func (api HandlerHTTP) handleGetConfig(w http.ResponseWriter, _ *http.Request) {
	// Return the current configuration
	err := json.NewEncoder(w).Encode(api.safeConfig.Config())
	if err != nil {
		api.log.Error("failed to encode config", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api HandlerHTTP) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var update ConfigUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		// TODO: wrap this error?
		api.log.Info("received bad update memstore config update", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Only update fields that were included in the request
	if update.PutLatency != nil {
		duration, err := time.ParseDuration(*update.PutLatency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		api.safeConfig.SetLatencyPUTRoute(duration)
	}

	if update.GetLatency != nil {
		duration, err := time.ParseDuration(*update.GetLatency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		api.safeConfig.SetLatencyGETRoute(duration)
	}

	if update.PutReturnsFailoverError != nil {
		api.safeConfig.SetPUTReturnsFailoverError(*update.PutReturnsFailoverError)
	}

	if update.MaxBlobSizeBytes != nil {
		api.safeConfig.SetMaxBlobSizeBytes(*update.MaxBlobSizeBytes)
	}

	if update.BlobExpiration != nil {
		duration, err := time.ParseDuration(*update.BlobExpiration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		api.safeConfig.SetBlobExpiration(duration)
	}

	if update.PutWithGetReturnsDerivationError != nil {
		if update.PutWithGetReturnsDerivationError.Reset {
			// Reset is true means reset to nil
			_ = api.safeConfig.SetPUTWithGetReturnsDerivationError(nil)
		} else {
			// Reset is false means set the provided value
			err := api.safeConfig.SetPUTWithGetReturnsDerivationError(update.PutWithGetReturnsDerivationError.DerivationError)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	// Return the current configuration
	err := json.NewEncoder(w).Encode(api.safeConfig.Config())
	if err != nil {
		api.log.Error("failed to encode config", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
