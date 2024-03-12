// Code generated by MockGen. DO NOT EDIT.
// Source: service/shamir_client.go

// Package testutil is a generated GoMock package.
package testutil

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIShamirClient is a mock of IShamirClient interface.
type MockIShamirClient struct {
	ctrl     *gomock.Controller
	recorder *MockIShamirClientMockRecorder
}

// MockIShamirClientMockRecorder is the mock recorder for MockIShamirClient.
type MockIShamirClientMockRecorder struct {
	mock *MockIShamirClient
}

// NewMockIShamirClient creates a new mock instance.
func NewMockIShamirClient(ctrl *gomock.Controller) *MockIShamirClient {
	mock := &MockIShamirClient{ctrl: ctrl}
	mock.recorder = &MockIShamirClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIShamirClient) EXPECT() *MockIShamirClientMockRecorder {
	return m.recorder
}

// IssueShamirTransaction mocks base method.
func (m *MockIShamirClient) IssueShamirTransaction(amount uint64, address string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IssueShamirTransaction", amount, address)
	ret0, _ := ret[0].(error)
	return ret0
}

// IssueShamirTransaction indicates an expected call of IssueShamirTransaction.
func (mr *MockIShamirClientMockRecorder) IssueShamirTransaction(amount, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IssueShamirTransaction", reflect.TypeOf((*MockIShamirClient)(nil).IssueShamirTransaction), amount, address)
}
