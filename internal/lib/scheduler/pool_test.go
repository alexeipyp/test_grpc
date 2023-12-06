package scheduler_test

import (
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/alexeipyp/test_grpc/internal/lib/scheduler"
	"go.uber.org/zap"
)

func TestWorkerPoolCreation(t *testing.T) {
	tests := []struct {
		name            string
		workerPoolSize  int
		isErrorExpected bool
	}{
		{
			name:            "Success",
			workerPoolSize:  5,
			isErrorExpected: false,
		},
		{
			name:            "Zero worker pool size value",
			workerPoolSize:  0,
			isErrorExpected: true,
		},
		{
			name:            "Negative worker pool size value",
			workerPoolSize:  -1,
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Arrange
			queue, err := scheduler.NewTaskQueue(2, zap.NewNop(), zap.NewNop())
			if err != nil {
				t.Fatalf("scheduler.NewTaskQueue %v", err)
			}

			//Act
			_, err = scheduler.NewWorkerPool(test.workerPoolSize, queue, zap.NewNop(), zap.NewNop())

			//Assert
			if !test.isErrorExpected {
				if err != nil {
					t.Errorf("NewTaskQueue(%v) got an error: %v", test.workerPoolSize, err)
				}
			} else {
				if err == nil {
					t.Errorf("NewTaskQueue(%v) - expected an error", test.workerPoolSize)
				}
			}
		})
	}
}

func TestWorkerPoolRun(t *testing.T) {
	//Arrange
	queue, err := scheduler.NewTaskQueue(2, zap.NewNop(), zap.NewNop())
	if err != nil {
		t.Fatalf("scheduler.NewTaskQueue %v", err)
	}
	expectedNumOfWorkers := 5
	pool, err := scheduler.NewWorkerPool(expectedNumOfWorkers, queue, zap.NewNop(), zap.NewNop())
	if err != nil {
		t.Fatalf("scheduler.NewWorkerPool %v", err)
	}
	numOfGoroutinesBefore := runtime.NumGoroutine()

	//Act
	pool.Run()
	numOfGoroutinesAfter := runtime.NumGoroutine()

	//Assert
	gotNumOfWorkers := numOfGoroutinesAfter - numOfGoroutinesBefore
	if gotNumOfWorkers != expectedNumOfWorkers {
		t.Errorf("Worker pool ran %v workers but expected %v", gotNumOfWorkers, expectedNumOfWorkers)
	}
}

func TestWorkerPoolStop(t *testing.T) {
	//Arrange
	queue, err := scheduler.NewTaskQueue(2, zap.NewNop(), zap.NewNop())
	if err != nil {
		t.Fatalf("scheduler.NewTaskQueue %v", err)
	}
	numOfWorkers := 5
	pool, err := scheduler.NewWorkerPool(numOfWorkers, queue, zap.NewNop(), zap.NewNop())
	if err != nil {
		t.Fatalf("scheduler.NewWorkerPool %v", err)
	}
	pool.Run()
	numOfGoroutinesBefore := runtime.NumGoroutine()

	//Act
	pool.Stop()
	numOfGoroutinesAfter := runtime.NumGoroutine()

	//Assert
	if numOfGoroutinesBefore-numOfGoroutinesAfter != numOfWorkers {
		t.Errorf("Worker pool failed to stop workers")
	}
}

func TestWorkerPool_TasksAreProcessingByWorkers(t *testing.T) {
	//Arrange
	queueSize := 2
	queue, err := scheduler.NewTaskQueue(queueSize, zap.NewNop(), zap.NewNop())
	if err != nil {
		t.Fatalf("scheduler.NewTaskQueue %v", err)
	}
	numOfWorkers := 2
	pool, err := scheduler.NewWorkerPool(numOfWorkers, queue, zap.NewNop(), zap.NewNop())
	if err != nil {
		t.Fatalf("scheduler.NewWorkerPool %v", err)
	}

	isTask1Processing := make(chan bool)
	isTask2Processing := make(chan bool)
	err1 := queue.Enqueue(func() {
		isTask1Processing <- true
	})
	err2 := queue.Enqueue(func() {
		isTask2Processing <- true
	})
	if err1 != nil || err2 != nil {
		t.Fatalf("scheduler.Enqueue %v", errors.Join(err1, err2))
	}

	//Act
	var isTask1Processed, isTask2Processed bool
	pool.Run()
	for i := 0; i < queueSize; i++ {
		select {
		case isTask1Processed = <-isTask1Processing:
		case isTask2Processed = <-isTask2Processing:
		case <-time.After(time.Second * 1):
		}
	}
	//Assert
	if !(isTask1Processed && isTask2Processed) {
		t.Errorf("Workers failed at processing tasks from task queue")
	}
}
