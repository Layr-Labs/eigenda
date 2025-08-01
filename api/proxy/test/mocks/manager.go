// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Layr-Labs/eigenda/api/proxy/store (interfaces: IManager)
//
// Generated by this command:
//
//	mockgen -package mocks --destination ../test/mocks/manager.go . IManager
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	common "github.com/Layr-Labs/eigenda/api/proxy/common"
	certs "github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	commitments "github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	gomock "go.uber.org/mock/gomock"
)

// MockIManager is a mock of IManager interface.
type MockIManager struct {
	ctrl     *gomock.Controller
	recorder *MockIManagerMockRecorder
	isgomock struct{}
}

// MockIManagerMockRecorder is the mock recorder for MockIManager.
type MockIManagerMockRecorder struct {
	mock *MockIManager
}

// NewMockIManager creates a new mock instance.
func NewMockIManager(ctrl *gomock.Controller) *MockIManager {
	mock := &MockIManager{ctrl: ctrl}
	mock.recorder = &MockIManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIManager) EXPECT() *MockIManagerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockIManager) Get(ctx context.Context, versionedCert certs.VersionedCert, cm commitments.CommitmentMode, opts common.GETOpts) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, versionedCert, cm, opts)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIManagerMockRecorder) Get(ctx, versionedCert, cm, opts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIManager)(nil).Get), ctx, versionedCert, cm, opts)
}

// GetDispersalBackend mocks base method.
func (m *MockIManager) GetDispersalBackend() common.EigenDABackend {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDispersalBackend")
	ret0, _ := ret[0].(common.EigenDABackend)
	return ret0
}

// GetDispersalBackend indicates an expected call of GetDispersalBackend.
func (mr *MockIManagerMockRecorder) GetDispersalBackend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDispersalBackend", reflect.TypeOf((*MockIManager)(nil).GetDispersalBackend))
}

// GetOPKeccakValueFromS3 mocks base method.
func (m *MockIManager) GetOPKeccakValueFromS3(ctx context.Context, key []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOPKeccakValueFromS3", ctx, key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOPKeccakValueFromS3 indicates an expected call of GetOPKeccakValueFromS3.
func (mr *MockIManagerMockRecorder) GetOPKeccakValueFromS3(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOPKeccakValueFromS3", reflect.TypeOf((*MockIManager)(nil).GetOPKeccakValueFromS3), ctx, key)
}

// Put mocks base method.
func (m *MockIManager) Put(ctx context.Context, cm commitments.CommitmentMode, value []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, cm, value)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put.
func (mr *MockIManagerMockRecorder) Put(ctx, cm, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockIManager)(nil).Put), ctx, cm, value)
}

// PutOPKeccakPairInS3 mocks base method.
func (m *MockIManager) PutOPKeccakPairInS3(ctx context.Context, key, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutOPKeccakPairInS3", ctx, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutOPKeccakPairInS3 indicates an expected call of PutOPKeccakPairInS3.
func (mr *MockIManagerMockRecorder) PutOPKeccakPairInS3(ctx, key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutOPKeccakPairInS3", reflect.TypeOf((*MockIManager)(nil).PutOPKeccakPairInS3), ctx, key, value)
}

// SetDispersalBackend mocks base method.
func (m *MockIManager) SetDispersalBackend(backend common.EigenDABackend) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDispersalBackend", backend)
}

// SetDispersalBackend indicates an expected call of SetDispersalBackend.
func (mr *MockIManagerMockRecorder) SetDispersalBackend(backend any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDispersalBackend", reflect.TypeOf((*MockIManager)(nil).SetDispersalBackend), backend)
}
