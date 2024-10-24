// Code generated by MockGen. DO NOT EDIT.
// Source: notice.go

// Package rpc is a generated GoMock package.
package rpc

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	auth "github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	notice "github.com/morning-night-guild/platform-app/internal/domain/model/notice"
)

// MockNotice is a mock of Notice interface.
type MockNotice struct {
	ctrl     *gomock.Controller
	recorder *MockNoticeMockRecorder
}

// MockNoticeMockRecorder is the mock recorder for MockNotice.
type MockNoticeMockRecorder struct {
	mock *MockNotice
}

// NewMockNotice creates a new mock instance.
func NewMockNotice(ctrl *gomock.Controller) *MockNotice {
	mock := &MockNotice{ctrl: ctrl}
	mock.recorder = &MockNoticeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotice) EXPECT() *MockNoticeMockRecorder {
	return m.recorder
}

// Notify mocks base method.
func (m *MockNotice) Notify(arg0 context.Context, arg1 auth.Email, arg2 notice.Subject, arg3 notice.Message) (notice.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Notify", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(notice.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Notify indicates an expected call of Notify.
func (mr *MockNoticeMockRecorder) Notify(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Notify", reflect.TypeOf((*MockNotice)(nil).Notify), arg0, arg1, arg2, arg3)
}
