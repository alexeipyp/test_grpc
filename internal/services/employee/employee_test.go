package employeeservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexeipyp/test_grpc/internal/models"
	employeeservice "github.com/alexeipyp/test_grpc/internal/services/employee"
	mock_employeeservice "github.com/alexeipyp/test_grpc/internal/services/employee/mocks"
	"github.com/enescakir/emoji"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCheck(t *testing.T) {
	type mockBehavior func(r *mock_employeeservice.MockEmployeeRepository, email string, reasonId int)
	withNoErrorsBehavior := func(r *mock_employeeservice.MockEmployeeRepository, email string, reasonId int) {
		id := 1
		r.EXPECT().GetEmployeeInfo(gomock.Any(), email).Return(&models.Employee{Id: id}, nil)
		r.EXPECT().GetEmployeeAbsenceStatus(gomock.Any(), id).Return(reasonId, nil)
	}

	inputEmail := "sample@example.org"

	tests := []struct {
		name                   string
		mockBehavior           mockBehavior
		reasonIdFromRepository int
		expectedEmoji          string
		isErrorExpected        bool
	}{
		{
			name:                   "Repository returned ReasonId=1",
			reasonIdFromRepository: 1,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.House.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=2",
			reasonIdFromRepository: 2,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          "",
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=3",
			reasonIdFromRepository: 3,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.Airplane.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=4",
			reasonIdFromRepository: 4,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.Airplane.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=5",
			reasonIdFromRepository: 5,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.Thermometer.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=6",
			reasonIdFromRepository: 6,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.Thermometer.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=7",
			reasonIdFromRepository: 7,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          "",
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=8",
			reasonIdFromRepository: 8,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          "",
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=9",
			reasonIdFromRepository: 9,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.GraduationCap.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=10",
			reasonIdFromRepository: 10,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.House.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=11",
			reasonIdFromRepository: 11,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.Sun.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=12",
			reasonIdFromRepository: 12,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.Sun.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned ReasonId=13",
			reasonIdFromRepository: 13,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          emoji.Sun.String(),
			isErrorExpected:        false,
		},
		{
			name:                   "Repository returned unknown ReasonId",
			reasonIdFromRepository: 999,
			mockBehavior:           withNoErrorsBehavior,
			expectedEmoji:          "",
			isErrorExpected:        false,
		},
		{
			name: "Repository error at GetEmployeeInfo",
			mockBehavior: func(r *mock_employeeservice.MockEmployeeRepository, email string, reasonId int) {
				r.EXPECT().GetEmployeeInfo(gomock.Any(), email).Return(nil, errors.New("some error"))
			},
			isErrorExpected: true,
		},
		{
			name: "Repository error at GetEmployeeAbsenceStatus",
			mockBehavior: func(r *mock_employeeservice.MockEmployeeRepository, email string, reasonId int) {
				id := 1
				r.EXPECT().GetEmployeeInfo(gomock.Any(), email).Return(&models.Employee{Id: id}, nil)
				r.EXPECT().GetEmployeeAbsenceStatus(gomock.Any(), id).Return(0, errors.New("some error"))
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Arrange
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_employeeservice.NewMockEmployeeRepository(c)
			test.mockBehavior(repo, inputEmail, test.reasonIdFromRepository)

			empService := employeeservice.New(repo, zap.NewNop(), zap.NewNop())

			//Act
			emoji, err := empService.Check(context.Background(), inputEmail)

			//Assert
			if !test.isErrorExpected {
				if emoji != test.expectedEmoji {
					t.Errorf("Check was incorrect, got: %v, want: %v.", emoji, test.expectedEmoji)
				}
				if err != nil {
					t.Errorf("Check got an error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Check - expected an error")
				}
			}
		})
	}
}
