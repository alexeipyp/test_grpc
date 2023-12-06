package scheduler

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type WorkerPool struct {
	logger      *zap.Logger
	debugLogger *zap.Logger
	taskQueue   *TaskQueue
	workers     []*Worker
	size        int
	wg          sync.WaitGroup
}

func NewWorkerPool(size int, taskQueue *TaskQueue, logger *zap.Logger, debugLogger *zap.Logger) (*WorkerPool, error) {
	debugLogger.Debug("try to init worker pool with args", zap.Int("worker pool size", size))
	if size <= 0 {
		return nil, fmt.Errorf("failed to init worker pool with size=%d - should be more or equal 1", size)
	}
	return &WorkerPool{
		size:        size,
		taskQueue:   taskQueue,
		logger:      logger,
		debugLogger: debugLogger,
	}, nil
}

func (wp *WorkerPool) Run() {
	for i := 1; i <= wp.size; i++ {
		worker := newWorker(i, wp.taskQueue, wp.logger, wp.debugLogger)
		wp.workers = append(wp.workers, worker)
		go worker.init(&wp.wg)
	}
	wp.logger.Info("worker pool ran", zap.Int("num of workers", wp.size))
}

func (wp *WorkerPool) Stop() {
	for i := range wp.workers {
		wp.workers[i].stop()
	}
	wp.wg.Wait()
	wp.logger.Info("worker pool stopped", zap.Int("num of workers", wp.size))
}
