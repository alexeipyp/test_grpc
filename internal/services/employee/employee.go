package employeeservice

import (
	"context"

	"github.com/alexeipyp/test_grpc/internal/models"
	"github.com/enescakir/emoji"
	"go.uber.org/zap"
)

//go:generate mockgen -source=employee.go -destination=mocks/mock.go
type EmployeeRepository interface {
	GetEmployeeInfo(ctx context.Context, email string) (*models.Employee, error)
	GetEmployeeAbsenceStatus(ctx context.Context, personId int) (int, error)
}

type EmployeeCheckService struct {
	logger      *zap.Logger
	debugLogger *zap.Logger
	empRepo     EmployeeRepository
}

func New(empRepo EmployeeRepository, logger *zap.Logger, debugLogger *zap.Logger) *EmployeeCheckService {
	return &EmployeeCheckService{empRepo: empRepo, logger: logger, debugLogger: debugLogger}
}

func (s *EmployeeCheckService) Check(ctx context.Context, email string) (string, error) {
	debugLogger := s.debugLogger.With(
		zap.String("method", "EmployeeCheckService.Check"),
	)
	debugLogger.Debug("calling with args", zap.String("email", email))
	employee, err := s.empRepo.GetEmployeeInfo(ctx, email)
	if err != nil {
		s.logger.Error("failed to get employee info",
			zap.String("employee email", email),
			zap.Error(err),
		)
		return "", err
	}
	reasonId, err := s.empRepo.GetEmployeeAbsenceStatus(ctx, employee.Id)
	if err != nil {
		s.logger.Error("failed to get employee absence status",
			zap.Int("employee id", employee.Id),
			zap.Error(err),
		)
		return "", err
	}
	return getEmojiStringByReasonId(reasonId), nil
}

func getEmojiStringByReasonId(reasonId int) string {
	switch reasonId {
	case 1, 10:
		return emoji.House.String()
	case 3, 4:
		return emoji.Airplane.String()
	case 5, 6:
		return emoji.Thermometer.String()
	case 9:
		return emoji.GraduationCap.String()
	case 11, 12, 13:
		return emoji.Sun.String()
	}
	return ""
}
