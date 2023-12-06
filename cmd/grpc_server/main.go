package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexeipyp/test_grpc/internal/app"
	"github.com/alexeipyp/test_grpc/internal/config"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()
	if err := cfg.Validate(); err != nil {
		zap.L().Log(zap.PanicLevel, "failed to load config", zap.Error(err))
	}

	logger, err := initLogger(cfg.Environment, cfg.Log)
	if err != nil {
		panic(err)
	}
	debugLogger, err := initDebugLogger(cfg.Environment, cfg.Log)
	if err != nil {
		panic(err)
	}
	gRPCTraceLogger, err := initGRPCTraceLogger(cfg.Log)
	if err != nil {
		panic(err)
	}
	httpTraceLogger, err := initHTTPTraceLogger(cfg.Log)
	if err != nil {
		panic(err)
	}

	application := app.New(
		logger,
		gRPCTraceLogger,
		httpTraceLogger,
		debugLogger,
		cfg.GRPC.Host,
		fmt.Sprint(cfg.GRPC.Port),
		cfg.Scheduler.QueueSize,
		cfg.Scheduler.WorkerPoolSize,
		cfg.Cache.Lifetime.Duration(),
		cfg.ExternalHTTPConnection.Username,
		cfg.ExternalHTTPConnection.Password,
		cfg.ExternalHTTPConnection.Timeout.Duration(),
		cfg.ExternalHTTPConnection.Host,
		fmt.Sprint(cfg.ExternalHTTPConnection.Port),
	)

	go func() {
		application.Scheduler.Run()
		application.GRPCServ.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopSignal := <-stop
	logger.Info("stopping application", zap.String("signal", stopSignal.String()))
	application.Scheduler.Stop()
	application.GRPCServ.Stop()
	application.Cache.Stop()

	logger.Info("application stopped")

	if err := logger.Sync(); err != nil {
		panic(err.Error())
	}
	if err := debugLogger.Sync(); err != nil {
		panic(err.Error())
	}
}

func initDebugLogger(env string, cfg config.LogConfig) (*zap.Logger, error) {
	switch env {
	case config.EnvDev:
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.OutputPaths = []string{
			fmt.Sprintf("%s%s", cfg.LogsDir, cfg.DebugLogFilename),
		}
		return loggerConfig.Build()
	case config.EnvProd:
		return zap.NewNop(), nil
	default:
		return nil, fmt.Errorf("unknown env value")
	}
}

func initLogger(env string, cfg config.LogConfig) (*zap.Logger, error) {
	switch env {
	case config.EnvDev:
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		loggerConfig.OutputPaths = []string{
			fmt.Sprintf("%s%s", cfg.LogsDir, cfg.DebugLogFilename),
			fmt.Sprintf("%s%s", cfg.LogsDir, cfg.LogFilename),
		}
		return loggerConfig.Build()
	case config.EnvProd:
		loggerConfig := zap.NewProductionConfig()
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		loggerConfig.OutputPaths = []string{
			fmt.Sprintf("%s%s", cfg.LogsDir, cfg.LogFilename),
		}
		return loggerConfig.Build()
	default:
		return nil, fmt.Errorf("unknown env value")
	}
}

func initGRPCTraceLogger(cfg config.LogConfig) (*zap.Logger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	loggerConfig.OutputPaths = []string{
		fmt.Sprintf("%s%s", cfg.LogsDir, cfg.GRPCTraceLogFilename),
	}
	return loggerConfig.Build()
}

func initHTTPTraceLogger(cfg config.LogConfig) (*zap.Logger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	loggerConfig.OutputPaths = []string{
		fmt.Sprintf("%s%s", cfg.LogsDir, cfg.HTTPTraceLogFilename),
	}
	return loggerConfig.Build()
}
