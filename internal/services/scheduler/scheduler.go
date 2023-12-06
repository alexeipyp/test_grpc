package schedulerservice

import (
	"context"
	"errors"
	"fmt"

	schedulerserviceerrors "github.com/alexeipyp/test_grpc/internal/services/scheduler/errors"
)

//go:generate mockgen -source=scheduler.go -destination=mocks/mock.go
type TaskQueuer interface {
	Enqueue(task func()) error
}

type SchedulerService struct {
	queue TaskQueuer
}

func New(queue TaskQueuer) *SchedulerService {
	return &SchedulerService{queue: queue}
}

func (s *SchedulerService) Process(ctx context.Context, task interface{}, arg interface{}) (interface{}, error) {
	switch t := task.(type) {
	case func(context.Context, string) (string, error):
		argString, ok := arg.(string)
		if !ok {
			return nil, errors.New("incorrect type of argument")
		}
		return process[string, string](ctx, s.queue, task.(func(context.Context, string) (string, error)), argString)
	default:
		return nil, fmt.Errorf("unimplemented to such type of task: %v", t)
	}
}

func process[T any, R any](ctx context.Context, queue TaskQueuer, task func(context.Context, T) (R, error), arg T) (R, error) {
	var errResult R
	resultCh := make(chan R)
	errorCh := make(chan error)
	err := queue.Enqueue(func() {
		res, err := task(ctx, arg)
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- res
	})
	if err != nil {
		return errResult, &schedulerserviceerrors.EnqueueTaskError{InternalError: err}
	}
	select {
	case res := <-resultCh:
		return res, nil
	case err := <-errorCh:
		return errResult, err
	}
}
