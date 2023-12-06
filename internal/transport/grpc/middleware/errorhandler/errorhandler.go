package errorhandler

import (
	"context"

	employeerepoerrors "github.com/alexeipyp/test_grpc/internal/repositories/employee/errors"
	schedulerserviceerrors "github.com/alexeipyp/test_grpc/internal/services/scheduler/errors"
	employeegrpcerrors "github.com/alexeipyp/test_grpc/internal/transport/grpc/employee/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		m, err := handler(ctx, req)
		err = populateErrorWithGRPCStatusCode(err)
		return m, err
	}
}

func populateErrorWithGRPCStatusCode(err error) error {
	if err != nil {
		switch err.(type) {
		case *schedulerserviceerrors.EnqueueTaskError,
			*employeerepoerrors.HTTPServerUnavailableError:
			return status.Error(codes.Unavailable, err.Error())
		case *employeerepoerrors.HTTPResponseBadHTTPStatusError,
			*employeerepoerrors.HTTPResponseBadStatusError,
			*employeerepoerrors.HTTPResponseNoDataError:
			return status.Error(codes.Unknown, err.Error())
		case *employeegrpcerrors.ValidationError:
			return status.Error(codes.InvalidArgument, err.Error())
		default:
			return status.Error(codes.Internal, "internal error")
		}
	}
	return nil
}
