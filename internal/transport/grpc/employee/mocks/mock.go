// Code generated by MockGen. DO NOT EDIT.
// Source: server.go
//
// Generated by this command:
//
//	mockgen -source=server.go -destination=mocks/mock.go
//
// Package mock_employeegrpc is a generated GoMock package.
package mock_employeegrpc

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEmployeeAbsenceStatusChecker is a mock of EmployeeAbsenceStatusChecker interface.
type MockEmployeeAbsenceStatusChecker struct {
	ctrl     *gomock.Controller
	recorder *MockEmployeeAbsenceStatusCheckerMockRecorder
}

// MockEmployeeAbsenceStatusCheckerMockRecorder is the mock recorder for MockEmployeeAbsenceStatusChecker.
type MockEmployeeAbsenceStatusCheckerMockRecorder struct {
	mock *MockEmployeeAbsenceStatusChecker
}

// NewMockEmployeeAbsenceStatusChecker creates a new mock instance.
func NewMockEmployeeAbsenceStatusChecker(ctrl *gomock.Controller) *MockEmployeeAbsenceStatusChecker {
	mock := &MockEmployeeAbsenceStatusChecker{ctrl: ctrl}
	mock.recorder = &MockEmployeeAbsenceStatusCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmployeeAbsenceStatusChecker) EXPECT() *MockEmployeeAbsenceStatusCheckerMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockEmployeeAbsenceStatusChecker) Check(ctx context.Context, email string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, email)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockEmployeeAbsenceStatusCheckerMockRecorder) Check(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockEmployeeAbsenceStatusChecker)(nil).Check), ctx, email)
}

// MockSchedulerClient is a mock of SchedulerClient interface.
type MockSchedulerClient struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerClientMockRecorder
}

// MockSchedulerClientMockRecorder is the mock recorder for MockSchedulerClient.
type MockSchedulerClientMockRecorder struct {
	mock *MockSchedulerClient
}

// NewMockSchedulerClient creates a new mock instance.
func NewMockSchedulerClient(ctrl *gomock.Controller) *MockSchedulerClient {
	mock := &MockSchedulerClient{ctrl: ctrl}
	mock.recorder = &MockSchedulerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchedulerClient) EXPECT() *MockSchedulerClientMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockSchedulerClient) Process(ctx context.Context, task, arg any) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", ctx, task, arg)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process.
func (mr *MockSchedulerClientMockRecorder) Process(ctx, task, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockSchedulerClient)(nil).Process), ctx, task, arg)
}