// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/djedjethai/generation/pkg/getter (interfaces: Getter)

// Package getter is a generated GoMock package.
package getter

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGetter is a mock of Getter interface.
type MockGetter struct {
	ctrl     *gomock.Controller
	recorder *MockGetterMockRecorder
}

// MockGetterMockRecorder is the mock recorder for MockGetter.
type MockGetterMockRecorder struct {
	mock *MockGetter
}

// NewMockGetter creates a new mock instance.
func NewMockGetter(ctrl *gomock.Controller) *MockGetter {
	mock := &MockGetter{ctrl: ctrl}
	mock.recorder = &MockGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetter) EXPECT() *MockGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockGetter) Get(arg0 context.Context, arg1 string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockGetterMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGetter)(nil).Get), arg0, arg1)
}

// GetKeys mocks base method.
func (m *MockGetter) GetKeys(arg0 context.Context) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeys", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetKeys indicates an expected call of GetKeys.
func (mr *MockGetterMockRecorder) GetKeys(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeys", reflect.TypeOf((*MockGetter)(nil).GetKeys), arg0)
}
