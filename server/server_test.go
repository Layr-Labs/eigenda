package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/mocks"
	"github.com/ethereum/go-ethereum/log"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

const (
	simpleCommitmentPrefix = "\x00"

	// [alt-da, da layer, cert version]
	opGenericPrefixStr = "\x01\x00\x00"

	testCommitStr = "9a7d4f1c3e5b8a09d1c0fa4b3f8e1d7c6b29f1e6d8c4a7b3c2d4e5f6a7b8c9d0"
)

func TestHandleOPCommitments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouter := mocks.NewMockIRouter(ctrl)

	m := metrics.NewMetrics("default")
	server := NewServer("localhost", 8080, mockRouter, log.New(), m)

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
				mockRouter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name: "Success - OP Keccak256",
			url:  fmt.Sprintf("/get/0x00%s", testCommitStr),
			mockBehavior: func() {
				mockRouter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: testCommitStr,
		},
		{
			name: "Failure - OP Alt-DA Internal Server Error",
			url:  fmt.Sprintf("/get/0x010000%s", testCommitStr),
			mockBehavior: func() {
				mockRouter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name: "Success - OP Alt-DA",
			url:  fmt.Sprintf("/get/0x010000%s", testCommitStr),
			mockBehavior: func() {
				mockRouter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
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
			server.registerRoutes(r)
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedCode, rec.Code)
			// We don't test for bodies because it's a specific error message
			// that contains a lot of information
			// require.Equal(t, tt.expectedBody, rec.Body.String())

		})
	}
}

func TestHandlerPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouter := mocks.NewMockIRouter(ctrl)
	server := NewServer("localhost", 8080, mockRouter, log.New(), metrics.NoopMetrics)

	tests := []struct {
		name         string
		url          string
		body         []byte
		mockBehavior func()
		expectedCode int
		expectedBody string
		expectError  bool
	}{
		{
			name: "Failure OP Mode Alt-DA - InternalServerError",
			url:  "/put",
			body: []byte("some data that will trigger an internal error"),
			mockBehavior: func() {
				mockRouter.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
			expectError:  true,
		},
		{
			name: "Success OP Mode Alt-DA",
			url:  "/put",
			body: []byte("some data that will successfully be written to EigenDA"),
			mockBehavior: func() {
				mockRouter.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: opGenericPrefixStr + testCommitStr,
			expectError:  false,
		},
		{
			name: "Success OP Mode Keccak256",
			url:  fmt.Sprintf("/put/0x00%s", testCommitStr),
			body: []byte("some data that will successfully be written to EigenDA"),
			mockBehavior: func() {
				mockRouter.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "",
			expectError:  false,
		},
		{
			name: "Success Simple Commitment Mode",
			url:  "/put?commitment_mode=simple",
			body: []byte("some data that will successfully be written to EigenDA"),
			mockBehavior: func() {
				mockRouter.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: simpleCommitmentPrefix + testCommitStr,
			expectError:  false,
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
			server.registerRoutes(r)
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedCode, rec.Code)
			if !tt.expectError && tt.expectedBody != "" {
				require.Equal(t, []byte(tt.expectedBody), rec.Body.Bytes())
			}

			if !tt.expectError && tt.expectedBody == "" {
				require.Equal(t, []byte(nil), rec.Body.Bytes())
			}
		})
	}
}
