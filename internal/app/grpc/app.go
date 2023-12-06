package grpcapp

import (
	"net"

	"github.com/alexeipyp/test_grpc/internal/lib/loggers"
	employeegrpc "github.com/alexeipyp/test_grpc/internal/transport/grpc/employee"
	"github.com/alexeipyp/test_grpc/internal/transport/grpc/middleware/errorhandler"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	logger          *zap.Logger
	gRPCTraceLogger *zap.Logger
	debugLogger     *zap.Logger
	gRPCServer      *grpc.Server
	host            string
	port            string
}

func New(logger *zap.Logger,
	gRPCTraceLogger *zap.Logger,
	debugLogger *zap.Logger,
	host string,
	port string,
	empChecker employeegrpc.EmployeeAbsenceStatusChecker,
	schedulerClient employeegrpc.SchedulerClient,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("Recovered from panic", zap.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(loggers.InterceptorLogger(gRPCTraceLogger), loggingOpts...),
		errorhandler.UnaryServerInterceptor(),
	))

	employeegrpc.Register(gRPCServer, empChecker, schedulerClient)

	return &App{
		logger:          logger,
		gRPCTraceLogger: gRPCTraceLogger,
		debugLogger:     debugLogger,
		gRPCServer:      gRPCServer,
		host:            host,
		port:            port,
	}
}

func (a *App) MustRun() {
	err := a.Run()
	if err != nil {
		a.logger.Panic("failed to run grpc server", zap.Error(err))
		return
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", a.host+":"+a.port)
	if err != nil {
		return err
	}

	a.logger.Info("grpc server starting", zap.String("addr", l.Addr().String()))
	err = a.gRPCServer.Serve(l)

	return err
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
	a.logger.Info("grpc server gracefully stopped", zap.String("port", a.port))
}
