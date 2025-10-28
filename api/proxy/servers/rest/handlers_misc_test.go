package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/test/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestInfoEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)
	mockKeccakManager := mocks.NewMockIKeccakManager(ctrl)

	t.Run("Success - Returns All PublicInfo Fields", func(t *testing.T) {
		// Setup test config with known values
		testPublicInfo := PubliclyExposedInfo{
			Version:             "1.2.3",
			ChainID:             "11155111",
			DirectoryAddress:    "0x1234567890abcdef",
			CertVerifierAddress: "0xfedcba0987654321",
			MaxBlobSizeBytes:    16777216, // 16 MiB
			RecencyWindowSize:   100,
			// DispersalBackend will be set dynamically
		}

		cfg := Config{
			Host: "localhost",
			Port: 0,
			APIsEnabled: &enablement.RestApisEnabled{
				Admin:               true,
				OpGenericCommitment: true,
				OpKeccakCommitment:  true,
				StandardCommitment:  true,
			},
			PublicInfo: testPublicInfo,
		}

		mockEigenDAManager.EXPECT().GetDispersalBackend().Return(common.V1EigenDABackend)

		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		rec := httptest.NewRecorder()

		r := mux.NewRouter()
		server := NewServer(cfg, mockEigenDAManager, mockKeccakManager, testLogger, metrics.NoopMetrics)
		server.RegisterRoutes(r)
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var response PubliclyExposedInfo
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify all fields
		require.Equal(t, testPublicInfo.Version, response.Version)
		require.Equal(t, testPublicInfo.ChainID, response.ChainID)
		require.Equal(t, testPublicInfo.DirectoryAddress, response.DirectoryAddress)
		require.Equal(t, testPublicInfo.CertVerifierAddress, response.CertVerifierAddress)
		require.Equal(t, testPublicInfo.MaxBlobSizeBytes, response.MaxBlobSizeBytes)
		require.Equal(t, testPublicInfo.RecencyWindowSize, response.RecencyWindowSize)
		require.Equal(t, common.EigenDABackendToString(common.V1EigenDABackend), response.DispersalBackend)
	})

	t.Run("Success - Dynamically Updates Dispersal Backend V1", func(t *testing.T) {
		mockEigenDAManager.EXPECT().GetDispersalBackend().Return(common.V1EigenDABackend)

		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		rec := httptest.NewRecorder()

		r := mux.NewRouter()
		server := NewServer(testCfg, mockEigenDAManager, mockKeccakManager, testLogger, metrics.NoopMetrics)
		server.RegisterRoutes(r)
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response PubliclyExposedInfo
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "V1", response.DispersalBackend)
	})

	t.Run("Success - Dynamically Updates Dispersal Backend V2", func(t *testing.T) {
		mockEigenDAManager.EXPECT().GetDispersalBackend().Return(common.V2EigenDABackend)

		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		rec := httptest.NewRecorder()

		r := mux.NewRouter()
		server := NewServer(testCfg, mockEigenDAManager, mockKeccakManager, testLogger, metrics.NoopMetrics)
		server.RegisterRoutes(r)
		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response PubliclyExposedInfo
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "V2", response.DispersalBackend)
	})

	t.Run("Success - Info Endpoint Always Available", func(t *testing.T) {
		// Unlike admin endpoints, /info should always be available
		adminDisabledCfg := Config{
			Host: "localhost",
			Port: 0,
			APIsEnabled: &enablement.RestApisEnabled{
				Admin:               false,
				OpGenericCommitment: false,
				OpKeccakCommitment:  false,
				StandardCommitment:  false,
			},
			PublicInfo: PubliclyExposedInfo{
				Version: "test-version",
			},
		}

		mockEigenDAManager.EXPECT().GetDispersalBackend().Return(common.V1EigenDABackend)

		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		rec := httptest.NewRecorder()

		r := mux.NewRouter()
		server := NewServer(adminDisabledCfg, mockEigenDAManager, mockKeccakManager, testLogger, metrics.NoopMetrics)
		server.RegisterRoutes(r)
		r.ServeHTTP(rec, req)

		// Should succeed even with admin endpoints disabled
		require.Equal(t, http.StatusOK, rec.Code)

		var response PubliclyExposedInfo
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "test-version", response.Version)
	})
}

func TestEigenDADispersalBackendEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)
	mockKeccakManager := mocks.NewMockIKeccakManager(ctrl)

	// Test with admin endpoints disabled - they should not be accessible
	t.Run("Admin Endpoints Disabled", func(t *testing.T) {
		// Create server config with admin endpoints disabled

		adminDisabledCfg := Config{
			Host: "localhost",
			Port: 0,
			APIsEnabled: &enablement.RestApisEnabled{
				Admin:               false,
				OpGenericCommitment: true,
				OpKeccakCommitment:  true,
				StandardCommitment:  true,
			},
		}

		// Test GET endpoint with admin disabled
		req := httptest.NewRequest(http.MethodGet, "/admin/eigenda-dispersal-backend", nil)
		rec := httptest.NewRecorder()

		r := mux.NewRouter()
		server := NewServer(adminDisabledCfg, mockEigenDAManager, mockKeccakManager, testLogger, metrics.NoopMetrics)
		server.RegisterRoutes(r)
		r.ServeHTTP(rec, req)

		// Should get 404 because the endpoint isn't registered
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	// Test with admin endpoints enabled
	t.Run("Admin Endpoints Enabled", func(t *testing.T) {
		// Initial state is false
		mockEigenDAManager.EXPECT().GetDispersalBackend().Return(common.V1EigenDABackend)

		// Test GET endpoint first to verify initial state
		t.Run("Get EigenDA Dispersal Backend", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/admin/eigenda-dispersal-backend", nil)
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			server := NewServer(testCfg, mockEigenDAManager, mockKeccakManager, testLogger, metrics.NoopMetrics)
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
			server := NewServer(testCfg, mockEigenDAManager, mockKeccakManager, testLogger, metrics.NoopMetrics)
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

			mockEigenDAManager.EXPECT().SetDispersalBackend(common.V2EigenDABackend)
			mockEigenDAManager.EXPECT().GetDispersalBackend().Return(common.V2EigenDABackend)

			req := httptest.NewRequest(http.MethodPut, "/admin/eigenda-dispersal-backend", bytes.NewReader(jsonBody))
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			server := NewServer(testCfg, mockEigenDAManager, mockKeccakManager, testLogger, metrics.NoopMetrics)
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
