package memconfig

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

var (
	testLogger = logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{AddSource: true})
)

func setup(config Config) (*mux.Router, *SafeConfig) {
	safeConfig := NewSafeConfig(config)
	r := mux.NewRouter()
	api := NewHandlerHTTP(testLogger, safeConfig)
	api.RegisterMemstoreConfigHandlers(r)

	return r, safeConfig
}

func TestHandlersHTTP_GetConfig(t *testing.T) {
	tests := []struct {
		name         string
		inputConfig  Config
		route        string
		expectedCode int
		expectError  bool
	}{
		{
			name:         "empty config",
			inputConfig:  Config{},
			route:        "/memstore/config",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name: "full config",
			inputConfig: Config{
				MaxBlobSizeBytes:        1024,
				BlobExpiration:          1 * time.Hour,
				PutLatency:              1 * time.Second,
				GetLatency:              2 * time.Second,
				PutReturnsFailoverError: true,
			},
			route:        "/memstore/config",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name: "partially filled config",
			inputConfig: Config{
				BlobExpiration: 1 * time.Hour,
				PutLatency:     1 * time.Second,
			},
			route:        "/memstore/config",
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "invalid route",
			inputConfig:  Config{},
			route:        "/memstore/config/",
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:         "invalid route",
			inputConfig:  Config{},
			route:        "/memstore",
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, safeConfig := setup(tt.inputConfig)

			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectError {
				return
			}

			expectedResp, err := safeConfig.Config().MarshalJSON()
			require.NoError(t, err)
			resp := rec.Body.String()
			require.Equal(t, string(expectedResp)+"\n", resp)
		})
	}
}

func TestHandlersHTTP_PatchConfig(t *testing.T) {
	tests := []struct {
		name            string
		initialConfig   Config
		requestBodyJSON string
		expectedStatus  int
		validate        func(*testing.T, Config, *SafeConfig)
	}{
		{
			name: "update single field",
			initialConfig: Config{
				PutLatency: 2 * time.Second,
				GetLatency: 2 * time.Second,
			},
			requestBodyJSON: `{"PutLatency": "5s"}`,
			expectedStatus:  http.StatusOK,
			validate: func(t *testing.T, inputConfig Config, sc *SafeConfig) {
				outputConfig := sc.Config()
				inputConfig.PutLatency = 5 * time.Second
				require.Equal(t, inputConfig, outputConfig)
			},
		},
		{
			name: "invalid PutLatency value (not string) does not update config",
			initialConfig: Config{
				PutLatency: 1 * time.Second,
			},
			requestBodyJSON: `{"PutLatency": 1000}`,
			expectedStatus:  http.StatusBadRequest,
			validate: func(t *testing.T, inputConfig Config, sc *SafeConfig) {
				outputConfig := sc.Config()
				require.Equal(t, inputConfig, outputConfig)
			},
		},
		{
			name:            "update instructed status code return",
			initialConfig:   Config{},
			requestBodyJSON: `{"PutWithGetReturnsDerivationError": {"StatusCode": 3}}`,
			expectedStatus:  http.StatusOK,
			validate: func(t *testing.T, inputConfig Config, sc *SafeConfig) {
				outputConfig := sc.Config()
				inputConfig.PutWithGetReturnsDerivationError = coretypes.ErrInvalidCertDerivationError
				require.Equal(t, inputConfig, outputConfig)
			},
		},
		{
			name:            "invalid update to derivation error with invalid status code (status code 100 does not exist)",
			initialConfig:   Config{},
			requestBodyJSON: `{"PutWithGetReturnsDerivationError": {"StatusCode": 100}}`,
			expectedStatus:  http.StatusBadRequest,
			validate: func(t *testing.T, inputConfig Config, sc *SafeConfig) {
				outputConfig := sc.Config()
				require.Equal(t, inputConfig, outputConfig)
			},
		},
		{
			name: "update multiple fields",
			initialConfig: Config{
				MaxBlobSizeBytes:        1024,
				BlobExpiration:          1 * time.Hour,
				PutLatency:              1 * time.Nanosecond,
				GetLatency:              1 * time.Nanosecond,
				PutReturnsFailoverError: true,
			},
			requestBodyJSON: `{"PutLatency": "5s", "GetLatency": "10s"}`,
			expectedStatus:  http.StatusOK,
			validate: func(t *testing.T, inputConfig Config, sc *SafeConfig) {
				inputConfig.PutLatency = 5 * time.Second
				inputConfig.GetLatency = 10 * time.Second
				outputConfig := sc.Config()
				require.Equal(t, inputConfig, outputConfig)
			},
		},
		{
			name:          "update all fields",
			initialConfig: Config{},
			requestBodyJSON: `{
				"MaxBlobSizeBytes": 1024,
				"BlobExpiration": "1h",
				"PutLatency": "1s",
				"GetLatency": "2s",
				"PutReturnsFailoverError": true
			}`,
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, inputConfig Config, sc *SafeConfig) {
				outputConfig := sc.Config()
				inputConfig.MaxBlobSizeBytes = 1024
				inputConfig.BlobExpiration = 1 * time.Hour
				inputConfig.PutLatency = 1 * time.Second
				inputConfig.GetLatency = 2 * time.Second
				inputConfig.PutReturnsFailoverError = true
				inputConfig.PutWithGetReturnsDerivationError = nil
				require.Equal(t, inputConfig, outputConfig)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, safeConfig := setup(tt.initialConfig)

			req := httptest.NewRequest(
				http.MethodPatch,
				"/memstore/config",
				bytes.NewReader([]byte(tt.requestBodyJSON)),
			)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.validate != nil {
				tt.validate(t, tt.initialConfig, safeConfig)
			}
		})
	}
}
