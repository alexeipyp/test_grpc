package schedulerapp

import "go.uber.org/zap"

type WorkerPool interface {
	Run()
	Stop()
}

type TaskQueue interface {
	Close()
}

type App struct {
	logger     *zap.Logger
	workerPool WorkerPool
}

func New(logger *zap.Logger, workerPool WorkerPool) *App {
	return &App{
		workerPool: workerPool,
		logger:     logger,
	}
}

func (a *App) Run() {
	a.workerPool.Run()
	a.logger.Info("scheduler successfully ran")
}

func (a *App) Stop() {
	a.workerPool.Stop()
	a.logger.Info("scheduler successfully stopped")
}
