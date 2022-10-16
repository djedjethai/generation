// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/djedjethai/generation0/pkg/deleter (interfaces: Deleter)

// Package deleter is a generated GoMock package.
package deleter

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDeleter is a mock of Deleter interface.
type MockDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockDeleterMockRecorder
}

// MockDeleterMockRecorder is the mock recorder for MockDeleter.
type MockDeleterMockRecorder struct {
	mock *MockDeleter
}

// NewMockDeleter creates a new mock instance.
func NewMockDeleter(ctrl *gomock.Controller) *MockDeleter {
	mock := &MockDeleter{ctrl: ctrl}
	mock.recorder = &MockDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleter) EXPECT() *MockDeleterMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDeleter) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDeleterMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDeleter)(nil).Delete), arg0)
}
