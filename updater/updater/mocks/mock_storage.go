// Code generated by MockGen. DO NOT EDIT.
// Source: updater/updater (interfaces: StorageInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	feed "updater/updater/model/feed"

	gomock "github.com/golang/mock/gomock"
)

// MockStorageInterface is a mock of StorageInterface interface.
type MockStorageInterface struct {
	ctrl     *gomock.Controller
	recorder *MockStorageInterfaceMockRecorder
}

// MockStorageInterfaceMockRecorder is the mock recorder for MockStorageInterface.
type MockStorageInterfaceMockRecorder struct {
	mock *MockStorageInterface
}

// NewMockStorageInterface creates a new mock instance.
func NewMockStorageInterface(ctrl *gomock.Controller) *MockStorageInterface {
	mock := &MockStorageInterface{ctrl: ctrl}
	mock.recorder = &MockStorageInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageInterface) EXPECT() *MockStorageInterfaceMockRecorder {
	return m.recorder
}

// UpdateHTMLFeed mocks base method.
func (m *MockStorageInterface) UpdateHTMLFeed(arg0 feed.Source, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateHTMLFeed", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateHTMLFeed indicates an expected call of UpdateHTMLFeed.
func (mr *MockStorageInterfaceMockRecorder) UpdateHTMLFeed(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateHTMLFeed", reflect.TypeOf((*MockStorageInterface)(nil).UpdateHTMLFeed), arg0, arg1)
}

// UpdateRSSFeed mocks base method.
func (m *MockStorageInterface) UpdateRSSFeed(arg0 feed.Source, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRSSFeed", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRSSFeed indicates an expected call of UpdateRSSFeed.
func (mr *MockStorageInterfaceMockRecorder) UpdateRSSFeed(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRSSFeed", reflect.TypeOf((*MockStorageInterface)(nil).UpdateRSSFeed), arg0, arg1)
}
