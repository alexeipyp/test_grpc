package employeegrpc_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	employeegrpc "github.com/alexeipyp/test_grpc/internal/transport/grpc/employee"
	employeegrpcerrors "github.com/alexeipyp/test_grpc/internal/transport/grpc/employee/errors"
	mock_employeegrpc "github.com/alexeipyp/test_grpc/internal/transport/grpc/employee/mocks"
	employee_v1 "github.com/alexeipyp/test_grpc/pkg/employee/v1"
	"github.com/enescakir/emoji"
	"go.uber.org/mock/gomock"
)

func TestPopulateWithAbsenceStatus_Success(t *testing.T) {
	//Arrange
	c := gomock.NewController(t)
	defer c.Finish()

	expectedEmoji := emoji.House.String()

	empChecker := mock_employeegrpc.NewMockEmployeeAbsenceStatusChecker(c)

	schedulerClient := mock_employeegrpc.NewMockSchedulerClient(c)
	schedulerClient.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedEmoji, nil)

	serverAPI := employeegrpc.New(empChecker, schedulerClient)
	employeeInfo := employee_v1.EmployeeInfo{
		DisplayName: "Иван",
		Email:       "some@email.org",
		MobilePhone: "88005553535",
		WorkPhone:   "1111",
	}
	req := employee_v1.PopulateRequest{Info: &employeeInfo}
	expectedDisplayName := fmt.Sprintf("%s %s", req.GetInfo().GetDisplayName(), expectedEmoji)

	//Act
	resp, err := serverAPI.PopulateWithAbsenceStatus(context.Background(), &req)

	//Assert
	if err != nil {
		t.Errorf("PopulateWithAbsenceStatus got an error: %v", err)
	}
	if resp.GetInfo().GetDisplayName() != expectedDisplayName {
		t.Errorf("PopulateWithAbsenceStatus - expected DisplayName: %s, but got: %s",
			expectedDisplayName,
			resp.GetInfo().GetDisplayName(),
		)
	}
}

func TestPopulateWithAbsenceStatus_ErrorFromSchedulerClient(t *testing.T) {
	//Arrange
	c := gomock.NewController(t)
	defer c.Finish()

	empChecker := mock_employeegrpc.NewMockEmployeeAbsenceStatusChecker(c)

	schedulerClient := mock_employeegrpc.NewMockSchedulerClient(c)
	schedulerClient.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error"))
	serverAPI := employeegrpc.New(empChecker, schedulerClient)
	employeeInfo := employee_v1.EmployeeInfo{
		DisplayName: "Иван",
		Email:       "some@email.org",
		MobilePhone: "88005553535",
		WorkPhone:   "1111",
	}
	req := employee_v1.PopulateRequest{Info: &employeeInfo}

	//Act
	_, err := serverAPI.PopulateWithAbsenceStatus(context.Background(), &req)

	//Assert
	if err == nil {
		t.Errorf("PopulateWithAbsenceStatus - error is expected")
	}
}

func TestValidatePopulateRequest(t *testing.T) {
	tests := []struct {
		name          string
		input         employee_v1.PopulateRequest
		expectedError error
	}{
		{
			name:          "DisplayName is empty string",
			expectedError: &employeegrpcerrors.ValidationError{},
			input: employee_v1.PopulateRequest{
				Info: &employee_v1.EmployeeInfo{
					DisplayName: "",
					Email:       "some@mail.org",
					MobilePhone: "88005553535",
					WorkPhone:   "1111",
				},
			},
		},
		{
			name:          "Email is invalid",
			expectedError: &employeegrpcerrors.ValidationError{},
			input: employee_v1.PopulateRequest{
				Info: &employee_v1.EmployeeInfo{
					DisplayName: "Иван",
					MobilePhone: "88005553535",
					Email:       "notemail",
					WorkPhone:   "1111",
				},
			},
		},
		{
			name:          "Email and DisplayName are invalid",
			expectedError: &employeegrpcerrors.ValidationError{},
			input: employee_v1.PopulateRequest{
				Info: &employee_v1.EmployeeInfo{
					DisplayName: "",
					MobilePhone: "88005553535",
					Email:       "notemail",
					WorkPhone:   "1111",
				},
			},
		},
		{
			name:          "Success",
			expectedError: nil,
			input: employee_v1.PopulateRequest{
				Info: &employee_v1.EmployeeInfo{
					DisplayName: "Иван",
					MobilePhone: "88005553535",
					Email:       "some@mail.ru",
					WorkPhone:   "1111",
				},
			},
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			//Arrange
			c := gomock.NewController(t)
			defer c.Finish()

			empChecker := mock_employeegrpc.NewMockEmployeeAbsenceStatusChecker(c)
			schedulerClient := mock_employeegrpc.NewMockSchedulerClient(c)
			serverAPI := employeegrpc.New(empChecker, schedulerClient)

			//Act
			err := serverAPI.ValidatePopulateRequest(&tests[i].input)

			//Arrange
			if reflect.TypeOf(err) != reflect.TypeOf(tests[i].expectedError) {
				t.Errorf("ValidatePopulateRequest - incorrect error, got: %v, want: %v", reflect.TypeOf(err), reflect.TypeOf(tests[i].expectedError))
			}
		})
	}
}
