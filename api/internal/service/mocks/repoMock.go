// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEmailSubscriptionRepo is a mock of EmailSubscriptionRepo interface.
type MockEmailSubscriptionRepo struct {
	ctrl     *gomock.Controller
	recorder *MockEmailSubscriptionRepoMockRecorder
}

// MockEmailSubscriptionRepoMockRecorder is the mock recorder for MockEmailSubscriptionRepo.
type MockEmailSubscriptionRepoMockRecorder struct {
	mock *MockEmailSubscriptionRepo
}

// NewMockEmailSubscriptionRepo creates a new mock instance.
func NewMockEmailSubscriptionRepo(ctrl *gomock.Controller) *MockEmailSubscriptionRepo {
	mock := &MockEmailSubscriptionRepo{ctrl: ctrl}
	mock.recorder = &MockEmailSubscriptionRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailSubscriptionRepo) EXPECT() *MockEmailSubscriptionRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockEmailSubscriptionRepo) Add(email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", email)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockEmailSubscriptionRepoMockRecorder) Add(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockEmailSubscriptionRepo)(nil).Add), email)
}

// GetAll mocks base method.
func (m *MockEmailSubscriptionRepo) GetAll() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockEmailSubscriptionRepoMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockEmailSubscriptionRepo)(nil).GetAll))
}

// MockEmailSendingRepo is a mock of EmailSendingRepo interface.
type MockEmailSendingRepo struct {
	ctrl     *gomock.Controller
	recorder *MockEmailSendingRepoMockRecorder
}

// MockEmailSendingRepoMockRecorder is the mock recorder for MockEmailSendingRepo.
type MockEmailSendingRepoMockRecorder struct {
	mock *MockEmailSendingRepo
}

// NewMockEmailSendingRepo creates a new mock instance.
func NewMockEmailSendingRepo(ctrl *gomock.Controller) *MockEmailSendingRepo {
	mock := &MockEmailSendingRepo{ctrl: ctrl}
	mock.recorder = &MockEmailSendingRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailSendingRepo) EXPECT() *MockEmailSendingRepoMockRecorder {
	return m.recorder
}

// SendToList mocks base method.
func (m *MockEmailSendingRepo) SendToList(emails []string, message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendToList", emails, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendToList indicates an expected call of SendToList.
func (mr *MockEmailSendingRepoMockRecorder) SendToList(emails, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendToList", reflect.TypeOf((*MockEmailSendingRepo)(nil).SendToList), emails, message)
}