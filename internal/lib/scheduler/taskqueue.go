package scheduler

import (
	"fmt"

	"go.uber.org/zap"
)

type TaskQueue struct {
	logger      *zap.Logger
	debugLogger *zap.Logger
	queue       chan func()
}

func NewTaskQueue(size int, logger *zap.Logger, debugLogger *zap.Logger) (*TaskQueue, error) {
	debugLogger.Debug("try to create task queue with args", zap.Int("size", size))
	if size <= 0 {
		return nil, fmt.Errorf("failed to init queue with size=%d - should be more or equal 1", size)
	}
	return &TaskQueue{
		queue:       make(chan func(), size),
		logger:      logger,
		debugLogger: debugLogger,
	}, nil
}

func (q *TaskQueue) Enqueue(task func()) error {
	q.debugLogger.Debug("try to enqueue task", zap.Bool("is task non-nil", task != nil))
	select {
	case q.queue <- task:
		q.logger.Info("task entered queue")
		return nil
	default:
		q.logger.Warn("attempt to add a task to a full queue",
			zap.Int("queue size", len(q.queue)),
			zap.Int("queue capacity", cap(q.queue)),
		)
		return fmt.Errorf("failed to enqueue task: queue overload")
	}
}
