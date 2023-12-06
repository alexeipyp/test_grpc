package scheduler

import (
	"sync"

	"go.uber.org/zap"
)

type Worker struct {
	taskQueue   *TaskQueue
	quitCh      chan bool
	logger      *zap.Logger
	debugLogger *zap.Logger
	id          int
}

func newWorker(id int, taskQueue *TaskQueue, logger *zap.Logger, debugLogger *zap.Logger) *Worker {
	logger = logger.With(zap.Int("worker id", id))
	debugLogger = debugLogger.With(zap.Int("worker id", id))
	return &Worker{
		taskQueue:   taskQueue,
		quitCh:      make(chan bool),
		logger:      logger,
		debugLogger: debugLogger,
		id:          id,
	}
}

func (w *Worker) init(wg *sync.WaitGroup) {
	w.logger.Info("worker started")
	wg.Add(1)
	for {
		select {
		case task := <-w.taskQueue.queue:
			w.debugLogger.Debug("worker processing task", zap.Bool("is task non-nil value", task != nil))
			task()
		case <-w.quitCh:
			w.logger.Info("worker stopped")
			wg.Done()
			return
		}
	}
}

func (w *Worker) stop() {
	w.quitCh <- true
}
