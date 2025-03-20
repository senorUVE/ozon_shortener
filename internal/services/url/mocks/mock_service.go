// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/url/url.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateURL mocks base method.
func (m *MockService) CreateURL(ctx context.Context, originalUrls []string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateURL", ctx, originalUrls)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateURL indicates an expected call of CreateURL.
func (mr *MockServiceMockRecorder) CreateURL(ctx, originalUrls interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateURL", reflect.TypeOf((*MockService)(nil).CreateURL), ctx, originalUrls)
}

// GetOriginal mocks base method.
func (m *MockService) GetOriginal(ctx context.Context, shortUrls []string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginal", ctx, shortUrls)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginal indicates an expected call of GetOriginal.
func (mr *MockServiceMockRecorder) GetOriginal(ctx, shortUrls interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginal", reflect.TypeOf((*MockService)(nil).GetOriginal), ctx, shortUrls)
}

// PublicURL mocks base method.
func (m *MockService) PublicURL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicURL")
	ret0, _ := ret[0].(string)
	return ret0
}

// PublicURL indicates an expected call of PublicURL.
func (mr *MockServiceMockRecorder) PublicURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicURL", reflect.TypeOf((*MockService)(nil).PublicURL))
}
