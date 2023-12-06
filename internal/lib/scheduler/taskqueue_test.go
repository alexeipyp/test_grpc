package scheduler_test

import (
	"errors"
	"testing"

	"github.com/alexeipyp/test_grpc/internal/lib/scheduler"
	"go.uber.org/zap"
)

func TestTaskQueueCreation(t *testing.T) {
	tests := []struct {
		name            string
		queueSize       int
		isErrorExpected bool
	}{
		{
			name:            "Success",
			queueSize:       5,
			isErrorExpected: false,
		},
		{
			name:            "Zero queue size value",
			queueSize:       0,
			isErrorExpected: true,
		},
		{
			name:            "Negative queue size value",
			queueSize:       -1,
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Act
			_, err := scheduler.NewTaskQueue(test.queueSize, zap.NewNop(), zap.NewNop())

			//Assert
			if !test.isErrorExpected {
				if err != nil {
					t.Errorf("NewTaskQueue(%v) got an error: %v", test.queueSize, err)
				}
			} else {
				if err == nil {
					t.Errorf("NewTaskQueue(%v) - expected an error", test.queueSize)
				}
			}
		})
	}
}

func TestTaskQueueEnqueue_Empty(t *testing.T) {
	//Arrange
	testOper := func() {}
	queue, err := scheduler.NewTaskQueue(2, zap.NewNop(), zap.NewNop())
	if err != nil {
		t.Fatalf("scheduler.NewTaskQueue %v", err)
	}

	//Act
	err = queue.Enqueue(testOper)

	//Accert
	if err != nil {
		t.Errorf("Failed to enqueue task to empty queue")
	}
}

func TestTaskQueueEnqueue_Full(t *testing.T) {
	//Arrange
	testOper := func() {}
	queue, err := scheduler.NewTaskQueue(2, zap.NewNop(), zap.NewNop())
	if err != nil {
		t.Fatalf("scheduler.NewTaskQueue %v", err)
	}
	err1 := queue.Enqueue(testOper)
	err2 := queue.Enqueue(testOper)
	if err1 != nil || err2 != nil {
		t.Fatalf("scheduler.Enqueue %v", errors.Join(err1, err2))
	}

	//Act
	err3 := queue.Enqueue(testOper)

	//Accert
	if err3 == nil {
		t.Errorf("Enqueue - expected an error")
	}
}
