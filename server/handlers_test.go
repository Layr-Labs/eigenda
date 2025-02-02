package server

// The tests in this file test not only the handlers but also the middlewares,
// because server.registerRoutes(r) registers the handlers wrapped with middlewares.

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/mocks"
	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	testLogger = logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{})
)

const (
	stdCommitmentPrefix = "\x00"

	// [alt-da, da layer, cert version]
	opGenericPrefixStr = "\x01\x00\x00"

	testCommitStr = "9a7d4f1c3e5b8a09d1c0fa4b3f8e1d7c6b29f1e6d8c4a7b3c2d4e5f6a7b8c9d0"
)

func TestHandlerGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorageMgr := mocks.NewMockIManager(ctrl)

	tests := []struct {
		name         string
		url          string
		mockBehavior func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "Failure - OP Keccak256 Internal Server Error",
			url:  fmt.Sprintf("/get/0x00%s", testCommitStr),
			mockBehavior: func() {
				mockStorageMgr.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name: "Success - OP Keccak256",
			url:  fmt.Sprintf("/get/0x00%s", testCommitStr),
			mockBehavior: func() {
				mockStorageMgr.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: testCommitStr,
		},
		{
			name: "Failure - OP Alt-DA Internal Server Error",
			url:  fmt.Sprintf("/get/0x010000%s", testCommitStr),
			mockBehavior: func() {
				mockStorageMgr.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name: "Success - OP Alt-DA",
			url:  fmt.Sprintf("/get/0x010000%s", testCommitStr),
			mockBehavior: func() {
				mockStorageMgr.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: testCommitStr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			// To add the vars to the context,
			// we need to create a router through which we can pass the request.
			r := mux.NewRouter()
			// enable this logger to help debug tests
			server := NewServer("localhost", 0, mockStorageMgr, testLogger, metrics.NoopMetrics)
			server.registerRoutes(r)
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedCode, rec.Code)
			// We only test for bodies for 200s because error messages contain a lot of information
			// that isn't very important to test (plus its annoying to always change if error msg changes slightly).
			if tt.expectedCode == http.StatusOK {
				require.Equal(t, tt.expectedBody, rec.Body.String())
			}

		})
	}
}

func TestHandlerPutSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorageMgr := mocks.NewMockIManager(ctrl)

	tests := []struct {
		name         string
		url          string
		body         []byte
		mockBehavior func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "Success OP Mode Alt-DA",
			url:  "/put",
			body: []byte("some data that will successfully be written to EigenDA"),
			mockBehavior: func() {
				mockStorageMgr.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: opGenericPrefixStr + testCommitStr,
		},
		{
			name: "Success OP Mode Keccak256",
			url:  fmt.Sprintf("/put/0x00%s", testCommitStr),
			body: []byte("some data that will successfully be written to EigenDA"),
			mockBehavior: func() {
				mockStorageMgr.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name: "Success Standard Commitment Mode",
			url:  "/put?commitment_mode=standard",
			body: []byte("some data that will successfully be written to EigenDA"),
			mockBehavior: func() {
				mockStorageMgr.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: stdCommitmentPrefix + testCommitStr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewReader(tt.body))
			rec := httptest.NewRecorder()

			// To add the vars to the context,
			// we need to create a router through which we can pass the request.
			r := mux.NewRouter()
			// enable this logger to help debug tests
			server := NewServer("localhost", 0, mockStorageMgr, testLogger, metrics.NoopMetrics)
			server.registerRoutes(r)
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedCode, rec.Code)
			// We only test for bodies for 200s because error messages contain a lot of information
			// that isn't very important to test (plus its annoying to always change if error msg changes slightly).
			if tt.expectedCode == http.StatusOK {
				require.Equal(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestHandlerPutErrors(t *testing.T) {
	// Each test is run against all 3 different modes.
	modes := []struct {
		name string
		url  string
	}{
		{
			name: "OP Mode Alt-DA",
			url:  "/put",
		},
		{
			name: "OP Mode Keccak256",
			url:  fmt.Sprintf("/put/0x00%s", testCommitStr),
		},
		{
			name: "Standard Commitment Mode",
			url:  "/put?commitment_mode=standard",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorageMgr := mocks.NewMockIManager(ctrl)

	tests := []struct {
		name                         string
		mockStorageMgrPutReturnedErr error
		expectedHTTPCode             int
	}{
		{
			// we only test OK status here. Returned commitment is checked in TestHandlerPut
			name:                         "Success - 200",
			mockStorageMgrPutReturnedErr: nil,
			expectedHTTPCode:             http.StatusOK,
		},
		{
			name:                         "Failure - InternalServerError 500",
			mockStorageMgrPutReturnedErr: fmt.Errorf("internal error"),
			expectedHTTPCode:             http.StatusInternalServerError,
		},
		{
			// if /put results in ErrorFailover (returned by eigenda-client), we should return 503
			name:                         "Failure - Failover 503",
			mockStorageMgrPutReturnedErr: &api.ErrorFailover{},
			expectedHTTPCode:             http.StatusServiceUnavailable,
		},
		{
			name:                         "Failure - TooManyRequests 429",
			mockStorageMgrPutReturnedErr: status.Errorf(codes.ResourceExhausted, "too many requests"),
			expectedHTTPCode:             http.StatusTooManyRequests,
		},
		{
			// only 400s are due to oversized blobs right now
			name:                         "Failure - BadRequest 400",
			mockStorageMgrPutReturnedErr: common.ErrProxyOversizedBlob,
			expectedHTTPCode:             http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		for _, mode := range modes {
			t.Run(tt.name+" / "+mode.name, func(t *testing.T) {
				mockStorageMgr.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, tt.mockStorageMgrPutReturnedErr)

				req := httptest.NewRequest(http.MethodPost, mode.url, strings.NewReader("optional body to be sent to eigenda"))
				rec := httptest.NewRecorder()

				// To add the vars to the context,
				// we need to create a router through which we can pass the request.
				r := mux.NewRouter()
				// enable this logger to help debug tests
				server := NewServer("localhost", 0, mockStorageMgr, testLogger, metrics.NoopMetrics)
				server.registerRoutes(r)
				r.ServeHTTP(rec, req)

				require.Equal(t, tt.expectedHTTPCode, rec.Code)
			})
		}
	}

}
