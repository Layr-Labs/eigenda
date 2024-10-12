package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/mocks"
	"github.com/ethereum/go-ethereum/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	genericPrefix = "\x00"

	// [alt-da, da layer, cert version]
	opKeccakPrefix     = "\x00"
	opGenericPrefixStr = "\x01\x00\x00"

	testCommitStr = "9a7d4f1c3e5b8a09d1c0fa4b3f8e1d7c6b29f1e6d8c4a7b3c2d4e5f6a7b8c9d0"
)

func TestGetHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouter := mocks.NewMockIRouter(ctrl)

	m := metrics.NewMetrics("default")
	server := NewServer("localhost", 8080, mockRouter, log.New(), m)

	tests := []struct {
		name                   string
		url                    string
		mockBehavior           func()
		expectedCode           int
		expectedBody           string
		expectError            bool
		expectedCommitmentMeta commitments.CommitmentMeta
	}{
		{
			name: "Failure - Op Mode InvalidCommitmentKey",
			url:  "/get/0x",
			mockBehavior: func() {
				// Error is triggered before calling the router
			},
			expectedCode:           http.StatusBadRequest,
			expectedBody:           "",
			expectError:            true,
			expectedCommitmentMeta: commitments.CommitmentMeta{},
		},
		{
			name: "Failure - Op Mode InvalidCommitmentKey",
			url:  "/get/0x1",
			mockBehavior: func() {
				// Error is triggered before calling the router
			},
			expectedCode:           http.StatusBadRequest,
			expectedBody:           "",
			expectError:            true,
			expectedCommitmentMeta: commitments.CommitmentMeta{},
		},
		{
			name: "Failure - Op Mode InvalidCommitmentKey",
			url:  "/get/0x999",
			mockBehavior: func() {
				// Error is triggered before calling the router
			},
			expectedCode:           http.StatusBadRequest,
			expectedBody:           "",
			expectError:            true,
			expectedCommitmentMeta: commitments.CommitmentMeta{},
		},
		{
			name: "Failure - OP Keccak256 Internal Server Error",
			url:  fmt.Sprintf("/get/0x00%s", testCommitStr),
			mockBehavior: func() {
				mockRouter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("internal error"))
			},
			expectedCode:           http.StatusInternalServerError,
			expectedBody:           "",
			expectError:            true,
			expectedCommitmentMeta: commitments.CommitmentMeta{},
		},
		{
			name: "Success - OP Keccak256",
			url:  fmt.Sprintf("/get/0x00%s", testCommitStr),
			mockBehavior: func() {
				mockRouter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode:           http.StatusOK,
			expectedBody:           testCommitStr,
			expectError:            false,
			expectedCommitmentMeta: commitments.CommitmentMeta{Mode: commitments.OptimismKeccak, CertVersion: 0},
		},
		{
			name: "Failure - OP Alt-DA Internal Server Error",
			url:  fmt.Sprintf("/get/0x010000%s", testCommitStr),
			mockBehavior: func() {
				mockRouter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("internal error"))
			},
			expectedCode:           http.StatusInternalServerError,
			expectedBody:           "",
			expectError:            true,
			expectedCommitmentMeta: commitments.CommitmentMeta{},
		},
		{
			name: "Success - OP Alt-DA",
			url:  fmt.Sprintf("/get/0x010000%s", testCommitStr),
			mockBehavior: func() {
				mockRouter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode:           http.StatusOK,
			expectedBody:           testCommitStr,
			expectError:            false,
			expectedCommitmentMeta: commitments.CommitmentMeta{Mode: commitments.OptimismGeneric, CertVersion: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			meta, err := server.HandleGet(rec, req)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedCode, rec.Code)
			require.Equal(t, tt.expectedCommitmentMeta, meta)
			require.Equal(t, tt.expectedBody, rec.Body.String())

		})
	}
	for _, tt := range tests {
		t.Run(tt.name+"/CheckMiddlewaresNoPanic", func(_ *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			// we also run the request through the middlewares, to make sure no panic occurs
			// could happen if there's a problem with the metrics. For eg, in the past we saw
			// panic: inconsistent label cardinality: expected 3 label values but got 1 in []string{"GET"}
			handler := WithLogging(WithMetrics(server.HandleGet, server.m), server.log)
			handler(rec, req)
		})
	}
}

func TestPutHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRouter := mocks.NewMockIRouter(ctrl)
	server := NewServer("localhost", 8080, mockRouter, log.New(), metrics.NoopMetrics)

	tests := []struct {
		name                   string
		url                    string
		body                   []byte
		mockBehavior           func()
		expectedCode           int
		expectedBody           string
		expectError            bool
		expectedCommitmentMeta commitments.CommitmentMeta
	}{
		{
			name: "Failure OP Keccak256 - TooShortCommitmentKey",
			url:  "/put/0x",
			body: []byte("some data"),
			mockBehavior: func() {
				// Error is triggered before calling the router
			},
			expectedCode:           http.StatusBadRequest,
			expectedBody:           "",
			expectError:            true,
			expectedCommitmentMeta: commitments.CommitmentMeta{},
		},
		{
			name: "Failure OP Keccak256 - TooShortCommitmentKey",
			url:  "/put/0x1",
			body: []byte("some data"),
			mockBehavior: func() {
				// Error is triggered before calling the router
			},
			expectedCode:           http.StatusBadRequest,
			expectedBody:           "",
			expectError:            true,
			expectedCommitmentMeta: commitments.CommitmentMeta{},
		},
		{
			name: "Failure OP Keccak256 - InvalidCommitmentPrefixBytes",
			url:  fmt.Sprintf("/put/0x999%s", testCommitStr),
			body: []byte("some data"),
			mockBehavior: func() {
				// Error is triggered before calling the router
			},
			expectedCode:           http.StatusBadRequest,
			expectedBody:           "",
			expectError:            true,
			expectedCommitmentMeta: commitments.CommitmentMeta{},
		},
		{
			name: "Failure OP Mode Alt-DA - InternalServerError",
			url:  "/put",
			body: []byte("some data that will trigger an internal error"),
			mockBehavior: func() {
				mockRouter.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("internal error"))
			},
			expectedCode:           http.StatusInternalServerError,
			expectedBody:           "",
			expectError:            true,
			expectedCommitmentMeta: commitments.CommitmentMeta{},
		},
		{
			name: "Success OP Mode Alt-DA",
			url:  "/put",
			body: []byte("some data that will successfully be written to EigenDA"),
			mockBehavior: func() {
				mockRouter.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode:           http.StatusOK,
			expectedBody:           opGenericPrefixStr + testCommitStr,
			expectError:            false,
			expectedCommitmentMeta: commitments.CommitmentMeta{Mode: commitments.OptimismGeneric, CertVersion: 0},
		},
		{
			name: "Success OP Mode Keccak256",
			url:  fmt.Sprintf("/put/0x00%s", testCommitStr),
			body: []byte("some data that will successfully be written to EigenDA"),
			mockBehavior: func() {
				mockRouter.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode:           http.StatusOK,
			expectedBody:           "",
			expectError:            false,
			expectedCommitmentMeta: commitments.CommitmentMeta{Mode: commitments.OptimismKeccak, CertVersion: 0},
		},
		{
			name: "Success Simple Commitment Mode",
			url:  "/put?commitment_mode=simple",
			body: []byte("some data that will successfully be written to EigenDA"),
			mockBehavior: func() {
				mockRouter.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(testCommitStr), nil)
			},
			expectedCode:           http.StatusOK,
			expectedBody:           genericPrefix + testCommitStr,
			expectError:            false,
			expectedCommitmentMeta: commitments.CommitmentMeta{Mode: commitments.SimpleCommitmentMode, CertVersion: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodPut, tt.url, bytes.NewReader(tt.body))
			rec := httptest.NewRecorder()

			meta, err := server.HandlePut(rec, req)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedCode, rec.Code)
			if !tt.expectError && tt.expectedBody != "" {
				require.Equal(t, []byte(tt.expectedBody), rec.Body.Bytes())
			}

			if !tt.expectError && tt.expectedBody == "" {
				require.Equal(t, []byte(nil), rec.Body.Bytes())
			}
			require.Equal(t, tt.expectedCommitmentMeta, meta)
		})
	}

	for _, tt := range tests {
		t.Run(tt.name+"/CheckMiddlewaresNoPanic", func(_ *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			// we also run the request through the middlewares, to make sure no panic occurs
			// could happen if there's a problem with the metrics. For eg, in the past we saw
			// panic: inconsistent label cardinality: expected 3 label values but got 1 in []string{"GET"}
			handler := WithLogging(WithMetrics(server.HandlePut, server.m), server.log)
			handler(rec, req)
		})
	}
}
