package app

import (
	"time"

	cacheapp "github.com/alexeipyp/test_grpc/internal/app/cache"
	grpcapp "github.com/alexeipyp/test_grpc/internal/app/grpc"
	schedulerapp "github.com/alexeipyp/test_grpc/internal/app/scheduler"
	"github.com/alexeipyp/test_grpc/internal/lib/scheduler"
	employeerepo "github.com/alexeipyp/test_grpc/internal/repositories/employee"
	employeeservice "github.com/alexeipyp/test_grpc/internal/services/employee"
	schedulerservice "github.com/alexeipyp/test_grpc/internal/services/scheduler"
	httpcache "github.com/alexeipyp/test_grpc/internal/storage/cache"
	"github.com/alexeipyp/test_grpc/internal/storage/httpclient"
	"go.uber.org/zap"
)

type App struct {
	GRPCServ  *grpcapp.App
	Scheduler *schedulerapp.App
	Cache     *cacheapp.App
}

func New(logger *zap.Logger,
	gRPCTraceLogger *zap.Logger,
	httpTraceLogger *zap.Logger,
	debugLogger *zap.Logger,
	grpcHost string,
	grpcPort string,
	queueSize int,
	workerPoolSize int,
	cacheLifeTime time.Duration,
	username string,
	password string,
	httpClientTimeout time.Duration,
	httpHost string,
	httpPort string,
) *App {
	taskQueue, err := scheduler.NewTaskQueue(queueSize, logger, debugLogger)
	if err != nil {
		logger.Panic("failed to init task queue", zap.Error(err))
	}
	workerPool, err := scheduler.NewWorkerPool(workerPoolSize, taskQueue, logger, debugLogger)
	if err != nil {
		logger.Panic("failed to init worker pool", zap.Error(err))
	}
	schedulerApp := schedulerapp.New(logger, workerPool)

	cache, err := httpcache.New(cacheLifeTime)
	if err != nil {
		logger.Panic("failed to init cache", zap.Error(err))
	}
	cacheApp := cacheapp.New(cache, logger)

	httpClient := httpclient.New(username, password, httpClientTimeout, httpTraceLogger)
	employeeRepo := employeerepo.New(httpClient, cache, httpHost, httpPort, logger, debugLogger)
	employeeAbsenceStatusService := employeeservice.New(employeeRepo, logger, debugLogger)
	schedulerClient := schedulerservice.New(taskQueue)

	grpcApp := grpcapp.New(logger, gRPCTraceLogger, debugLogger, grpcHost, grpcPort, employeeAbsenceStatusService, schedulerClient)

	return &App{
		GRPCServ:  grpcApp,
		Scheduler: schedulerApp,
		Cache:     cacheApp,
	}
}
