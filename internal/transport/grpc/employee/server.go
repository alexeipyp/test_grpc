package employeegrpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexeipyp/test_grpc/internal/lib/emailvalidator"
	employeegrpcerrors "github.com/alexeipyp/test_grpc/internal/transport/grpc/employee/errors"
	employee_v1 "github.com/alexeipyp/test_grpc/pkg/employee/v1"
	"google.golang.org/grpc"
)

//go:generate mockgen -source=server.go -destination=mocks/mock.go
type EmployeeAbsenceStatusChecker interface {
	Check(ctx context.Context, email string) (string, error)
}

type SchedulerClient interface {
	Process(ctx context.Context, task interface{}, arg interface{}) (interface{}, error)
}

type ServerAPI struct {
	employee_v1.UnimplementedEmployeeServer
	empChecker      EmployeeAbsenceStatusChecker
	schedulerClient SchedulerClient
}

func New(empChecker EmployeeAbsenceStatusChecker, schedulerClient SchedulerClient) *ServerAPI {
	return &ServerAPI{
		empChecker:      empChecker,
		schedulerClient: schedulerClient,
	}
}

func Register(gRPC *grpc.Server,
	empChecker EmployeeAbsenceStatusChecker,
	schedulerClient SchedulerClient) {
	employee_v1.RegisterEmployeeServer(gRPC, New(empChecker, schedulerClient))
}

func (s *ServerAPI) PopulateWithAbsenceStatus(
	ctx context.Context,
	req *employee_v1.PopulateRequest,
) (*employee_v1.PopulateResponse, error) {
	validationError := s.ValidatePopulateRequest(req)
	if validationError != nil {
		return nil, validationError
	}

	emoji, err := s.schedulerClient.Process(ctx, s.empChecker.Check, req.GetInfo().GetEmail())
	if err != nil {
		return nil, err
	}
	employeeInfo := employee_v1.EmployeeInfo{
		DisplayName: fmt.Sprintf("%s %s", req.GetInfo().GetDisplayName(), emoji),
		Email:       req.GetInfo().GetEmail(),
		WorkPhone:   req.GetInfo().GetWorkPhone(),
		MobilePhone: req.GetInfo().GetMobilePhone(),
	}
	return &employee_v1.PopulateResponse{Info: &employeeInfo}, nil
}

func (s *ServerAPI) ValidatePopulateRequest(req *employee_v1.PopulateRequest) error {
	var displayNameValidationError error
	var emailValidationError error
	if req.GetInfo().GetDisplayName() == "" {
		displayNameValidationError = errors.New("[expect DisplayName to be non-empty string]")
	}
	if !emailvalidator.IsEmailValid(req.GetInfo().GetEmail()) {
		emailValidationError = errors.New("[email is invalid]")
	}
	validationError := errors.Join(displayNameValidationError, emailValidationError)
	if validationError != nil {
		return &employeegrpcerrors.ValidationError{InternalError: validationError}
	}
	return nil
}
