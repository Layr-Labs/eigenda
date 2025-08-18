package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/test/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestEigenDADispersalBackendEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorageMgr := mocks.NewMockIManager(ctrl)

	// Test with admin endpoints disabled - they should not be accessible
	t.Run("Admin Endpoints Disabled", func(t *testing.T) {
		// Create server config with admin endpoints disabled
		adminDisabledCfg := Config{
			Host:        "localhost",
			Port:        0,
			EnabledAPIs: []string{}, // Empty list means no APIs are enabled
		}

		// Test GET endpoint with admin disabled
		req := httptest.NewRequest(http.MethodGet, "/admin/eigenda-dispersal-backend", nil)
		rec := httptest.NewRecorder()

		r := mux.NewRouter()
		server := NewServer(adminDisabledCfg, mockStorageMgr, testLogger, metrics.NoopMetrics)
		server.RegisterRoutes(r)
		r.ServeHTTP(rec, req)

		// Should get 404 because the endpoint isn't registered
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	// Test with admin endpoints enabled
	t.Run("Admin Endpoints Enabled", func(t *testing.T) {
		// Initial state is false
		mockStorageMgr.EXPECT().GetDispersalBackend().Return(common.V1EigenDABackend)

		// Test GET endpoint first to verify initial state
		t.Run("Get EigenDA Dispersal Backend", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/admin/eigenda-dispersal-backend", nil)
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			server := NewServer(testCfg, mockStorageMgr, testLogger, metrics.NoopMetrics)
			server.RegisterRoutes(r)
			r.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)

			var response struct {
				EigenDADispersalBackend string `json:"eigenDADispersalBackend"`
			}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)
			require.Equal(t, common.EigenDABackendToString(common.V1EigenDABackend), response.EigenDADispersalBackend)
		})

		// Test PUT endpoint with invalid input
		t.Run("Set EigenDA Dispersal Backend With Invalid Value", func(t *testing.T) {
			requestBody := struct {
				EigenDADispersalBackend string `json:"eigenDADispersalBackend"`
			}{
				EigenDADispersalBackend: "invalid",
			}
			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/admin/eigenda-dispersal-backend", bytes.NewReader(jsonBody))
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			server := NewServer(testCfg, mockStorageMgr, testLogger, metrics.NoopMetrics)
			server.RegisterRoutes(r)
			r.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		})

		// Test PUT endpoint to set the EigenDA dispersal backend
		t.Run("Set EigenDA Dispersal Backend", func(t *testing.T) {
			requestBody := struct {
				EigenDADispersalBackend string `json:"eigenDADispersalBackend"`
			}{
				EigenDADispersalBackend: common.EigenDABackendToString(common.V2EigenDABackend),
			}
			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			mockStorageMgr.EXPECT().SetDispersalBackend(common.V2EigenDABackend)
			mockStorageMgr.EXPECT().GetDispersalBackend().Return(common.V2EigenDABackend)

			req := httptest.NewRequest(http.MethodPut, "/admin/eigenda-dispersal-backend", bytes.NewReader(jsonBody))
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			server := NewServer(testCfg, mockStorageMgr, testLogger, metrics.NoopMetrics)
			server.RegisterRoutes(r)
			r.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)

			var response struct {
				EigenDADispersalBackend string `json:"eigenDADispersalBackend"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)
			require.Equal(t, common.EigenDABackendToString(common.V2EigenDABackend), response.EigenDADispersalBackend)
		})
	})
}
