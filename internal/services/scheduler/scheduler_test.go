package schedulerservice_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	schedulerservice "github.com/alexeipyp/test_grpc/internal/services/scheduler"
	schedulerserviceerrors "github.com/alexeipyp/test_grpc/internal/services/scheduler/errors"
	mock_schedulerservice "github.com/alexeipyp/test_grpc/internal/services/scheduler/mocks"
	"go.uber.org/mock/gomock"
)

func TestProcess_IncorrectArgs(t *testing.T) {
	tests := []struct {
		name string
		task interface{}
		arg  interface{}
	}{
		{
			name: "Unimplemented task type",
			task: func() {},
		},
		{
			name: "Task is func(string) (string, error) but arg is not string",
			task: func(str string) (string, error) { return "", nil },
			arg:  1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Arrange
			taskQueue := &mock_schedulerservice.MockTaskQueuer{}
			schedulerService := schedulerservice.New(taskQueue)

			//Act
			_, err := schedulerService.Process(context.Background(), test.task, test.arg)

			//Assert
			if err == nil {
				t.Errorf("Process - error is expected")
			}
		})
	}
}

func TestProcess_CorrectArgs(t *testing.T) {
	tests := []struct {
		name            string
		task            func(context.Context, string) (string, error)
		arg             string
		isErrorExpected bool
		expectedOutput  string
	}{
		{
			name:           "Task completed normally",
			expectedOutput: "test succeeded",
			task: func(context.Context, string) (string, error) {
				return "test succeeded", nil
			},
			arg:             "some arg",
			isErrorExpected: false,
		},
		{
			name: "Task returned an error",
			task: func(context.Context, string) (string, error) {
				return "", errors.New("oops, error")
			},
			arg:             "some arg",
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Arrange
			c := gomock.NewController(t)
			defer c.Finish()

			queue := mock_schedulerservice.NewMockTaskQueuer(c)
			queue.EXPECT().Enqueue(gomock.Any()).DoAndReturn(func(task func()) error {
				go task()
				return nil
			})

			schedulerService := schedulerservice.New(queue)

			//Act
			res, err := schedulerService.Process(context.Background(), test.task, test.arg)

			//Assert
			if test.isErrorExpected {
				if err == nil {
					t.Errorf("Process - error is expected")
				}
			} else {
				if err != nil {
					t.Errorf("Process - returned an error: %v", err)
				}
				if res != test.expectedOutput {
					t.Errorf("Process - got: %v, want: %v", res, test.expectedOutput)
				}
			}
		})
	}
}

func TestProcess_FailedToEnqueue(t *testing.T) {
	//Arrange
	c := gomock.NewController(t)
	defer c.Finish()

	queue := mock_schedulerservice.NewMockTaskQueuer(c)
	queue.EXPECT().Enqueue(gomock.Any()).DoAndReturn(func(task func()) error {
		return errors.New("failed to enqueue task")
	})

	task := func(context.Context, string) (string, error) {
		return "", errors.New("oops, error")
	}
	arg := "some arg"
	expectedError := &schedulerserviceerrors.EnqueueTaskError{}

	schedulerService := schedulerservice.New(queue)

	//Act
	_, err := schedulerService.Process(context.Background(), task, arg)

	//Assert
	if reflect.TypeOf(err) != reflect.TypeOf(expectedError) {
		t.Errorf("Process - incorrect error, got: %v, want: %v", reflect.TypeOf(err), reflect.TypeOf(expectedError))
	}
}
