// Code generated by MockGen. DO NOT EDIT.
// Source: employee.go
//
// Generated by this command:
//
//	mockgen -source=employee.go -destination=mocks/mock.go
//
// Package mock_employeeservice is a generated GoMock package.
package mock_employeeservice

import (
	context "context"
	reflect "reflect"

	models "github.com/alexeipyp/test_grpc/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockEmployeeRepository is a mock of EmployeeRepository interface.
type MockEmployeeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEmployeeRepositoryMockRecorder
}

// MockEmployeeRepositoryMockRecorder is the mock recorder for MockEmployeeRepository.
type MockEmployeeRepositoryMockRecorder struct {
	mock *MockEmployeeRepository
}

// NewMockEmployeeRepository creates a new mock instance.
func NewMockEmployeeRepository(ctrl *gomock.Controller) *MockEmployeeRepository {
	mock := &MockEmployeeRepository{ctrl: ctrl}
	mock.recorder = &MockEmployeeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmployeeRepository) EXPECT() *MockEmployeeRepositoryMockRecorder {
	return m.recorder
}

// GetEmployeeAbsenceStatus mocks base method.
func (m *MockEmployeeRepository) GetEmployeeAbsenceStatus(ctx context.Context, personId int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmployeeAbsenceStatus", ctx, personId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmployeeAbsenceStatus indicates an expected call of GetEmployeeAbsenceStatus.
func (mr *MockEmployeeRepositoryMockRecorder) GetEmployeeAbsenceStatus(ctx, personId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmployeeAbsenceStatus", reflect.TypeOf((*MockEmployeeRepository)(nil).GetEmployeeAbsenceStatus), ctx, personId)
}

// GetEmployeeInfo mocks base method.
func (m *MockEmployeeRepository) GetEmployeeInfo(ctx context.Context, email string) (*models.Employee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmployeeInfo", ctx, email)
	ret0, _ := ret[0].(*models.Employee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmployeeInfo indicates an expected call of GetEmployeeInfo.
func (mr *MockEmployeeRepositoryMockRecorder) GetEmployeeInfo(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmployeeInfo", reflect.TypeOf((*MockEmployeeRepository)(nil).GetEmployeeInfo), ctx, email)
}
