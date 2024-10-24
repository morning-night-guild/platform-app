// Code generated by MockGen. DO NOT EDIT.
// Source: core_article.go

// Package usecase is a generated GoMock package.
package usecase

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCoreArticle is a mock of CoreArticle interface.
type MockCoreArticle struct {
	ctrl     *gomock.Controller
	recorder *MockCoreArticleMockRecorder
}

// MockCoreArticleMockRecorder is the mock recorder for MockCoreArticle.
type MockCoreArticleMockRecorder struct {
	mock *MockCoreArticle
}

// NewMockCoreArticle creates a new mock instance.
func NewMockCoreArticle(ctrl *gomock.Controller) *MockCoreArticle {
	mock := &MockCoreArticle{ctrl: ctrl}
	mock.recorder = &MockCoreArticleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoreArticle) EXPECT() *MockCoreArticleMockRecorder {
	return m.recorder
}

// AddToUser mocks base method.
func (m *MockCoreArticle) AddToUser(arg0 context.Context, arg1 CoreArticleAddToUserInput) (CoreArticleAddToUserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToUser", arg0, arg1)
	ret0, _ := ret[0].(CoreArticleAddToUserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddToUser indicates an expected call of AddToUser.
func (mr *MockCoreArticleMockRecorder) AddToUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToUser", reflect.TypeOf((*MockCoreArticle)(nil).AddToUser), arg0, arg1)
}

// Delete mocks base method.
func (m *MockCoreArticle) Delete(arg0 context.Context, arg1 CoreArticleDeleteInput) (CoreArticleDeleteOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(CoreArticleDeleteOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockCoreArticleMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCoreArticle)(nil).Delete), arg0, arg1)
}

// List mocks base method.
func (m *MockCoreArticle) List(arg0 context.Context, arg1 CoreArticleListInput) (CoreArticleListOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(CoreArticleListOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCoreArticleMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCoreArticle)(nil).List), arg0, arg1)
}

// ListByUser mocks base method.
func (m *MockCoreArticle) ListByUser(arg0 context.Context, arg1 CoreArticleListByUserInput) (CoreArticleListByUserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByUser", arg0, arg1)
	ret0, _ := ret[0].(CoreArticleListByUserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByUser indicates an expected call of ListByUser.
func (mr *MockCoreArticleMockRecorder) ListByUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByUser", reflect.TypeOf((*MockCoreArticle)(nil).ListByUser), arg0, arg1)
}

// RemoveFromUser mocks base method.
func (m *MockCoreArticle) RemoveFromUser(arg0 context.Context, arg1 CoreArticleRemoveFromUserInput) (CoreArticleRemoveFromUserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromUser", arg0, arg1)
	ret0, _ := ret[0].(CoreArticleRemoveFromUserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveFromUser indicates an expected call of RemoveFromUser.
func (mr *MockCoreArticleMockRecorder) RemoveFromUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromUser", reflect.TypeOf((*MockCoreArticle)(nil).RemoveFromUser), arg0, arg1)
}

// Share mocks base method.
func (m *MockCoreArticle) Share(arg0 context.Context, arg1 CoreArticleShareInput) (CoreArticleShareOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Share", arg0, arg1)
	ret0, _ := ret[0].(CoreArticleShareOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Share indicates an expected call of Share.
func (mr *MockCoreArticleMockRecorder) Share(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Share", reflect.TypeOf((*MockCoreArticle)(nil).Share), arg0, arg1)
}
