// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/djedjethai/generation0/pkg/setter (interfaces: Setter)

// Package setter is a generated GoMock package.
package setter

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSetter is a mock of Setter interface.
type MockSetter struct {
	ctrl     *gomock.Controller
	recorder *MockSetterMockRecorder
}

// MockSetterMockRecorder is the mock recorder for MockSetter.
type MockSetterMockRecorder struct {
	mock *MockSetter
}

// NewMockSetter creates a new mock instance.
func NewMockSetter(ctrl *gomock.Controller) *MockSetter {
	mock := &MockSetter{ctrl: ctrl}
	mock.recorder = &MockSetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSetter) EXPECT() *MockSetterMockRecorder {
	return m.recorder
}

// Set mocks base method.
func (m *MockSetter) Set(arg0 string, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockSetterMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockSetter)(nil).Set), arg0, arg1)
}
