package employeerepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alexeipyp/test_grpc/internal/models"
	employeerepoerrors "github.com/alexeipyp/test_grpc/internal/repositories/employee/errors"
	"go.uber.org/zap"
)

const (
	OkStatus       = "OK"
	AbsencesRoute  = "Portal/springApi/api/absences"
	EmployeesRoute = "Portal/springApi/api/employees"
	HTTPS          = "https"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, payload []byte) error
}

//go:generate mockgen -source=employee.go -destination=mocks/mock.go
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPEmployeeRepository struct {
	httpClient  HTTPClient
	logger      *zap.Logger
	debugLogger *zap.Logger
	cache       Cache
	host        string
	port        string
}

func New(
	httpClient HTTPClient,
	cache Cache,
	host string,
	port string,
	logger *zap.Logger,
	debugLogger *zap.Logger,
) *HTTPEmployeeRepository {
	return &HTTPEmployeeRepository{
		httpClient:  httpClient,
		cache:       cache,
		host:        host,
		port:        port,
		logger:      logger,
		debugLogger: debugLogger,
	}
}

func (r *HTTPEmployeeRepository) GetEmployeeInfo(ctx context.Context, email string) (*models.Employee, error) {
	debugLogger := r.debugLogger.With(
		zap.String("method", "HTTPEmployeeRepository.GetEmployeeInfo"),
	)
	debugLogger.Debug("calling with args", zap.String("email", email))

	var respBody []byte
	if storedRespBody, err := r.cache.Get(email); err != nil {
		debugLogger.Debug("no cached response got",
			zap.String("cachekey", email),
			zap.Error(err),
		)
		dateFrom, dateTo := getDateTimeRangeInRFC3339()
		reqData := &GetEmployeesInfoRequest{
			Email:    email,
			DateFrom: dateFrom,
			DateTo:   dateTo,
		}
		reqBody, err := json.Marshal(reqData)
		if err != nil {
			return nil, err
		}

		respBody, err = r.makeAndDoPOSTRequest(ctx, reqBody, HTTPS, EmployeesRoute)
		if err != nil {
			return nil, err
		}

		if err := r.cache.Set(email, respBody); err != nil {
			r.logger.Error("Failed to cache http response", zap.String("key", email))
		}
	} else {
		debugLogger.Debug("response fetched from cache",
			zap.String("cachekey", email),
		)
		respBody = storedRespBody
	}

	var parsedResp GetEmployeesInfoResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		return nil, &employeerepoerrors.HTTPResponseFailedToParseBodyError{InternalError: err}
	}
	if err := parsedResp.validate(); err != nil {
		return nil, err
	}

	employee := parsedResp.Data[0]

	return &employee, nil
}

func (r *HTTPEmployeeRepository) GetEmployeeAbsenceStatus(ctx context.Context, personId int) (int, error) {
	debugLogger := r.debugLogger.With(
		zap.String("method", "HTTPEmployeeRepository.GetEmployeeAbsenceStatus"),
	)
	debugLogger.Debug("calling with args", zap.Int("personId", personId))

	var respBody []byte
	if storedRespBody, err := r.cache.Get(fmt.Sprint(personId)); err != nil {
		debugLogger.Debug("no cached response got",
			zap.String("cachekey", fmt.Sprint(personId)),
			zap.Error(err),
		)
		dateFrom, dateTo := getDateTimeRangeInRFC3339()
		reqData := &GetEmployeesAbsenceStatusRequest{
			Ids:      []int{personId},
			DateFrom: dateFrom,
			DateTo:   dateTo,
		}
		reqBody, err := json.Marshal(reqData)
		if err != nil {
			return 0, err
		}

		respBody, err = r.makeAndDoPOSTRequest(ctx, reqBody, HTTPS, AbsencesRoute)
		if err != nil {
			return 0, err
		}

		if err := r.cache.Set(fmt.Sprint(personId), respBody); err != nil {
			r.logger.Error("Failed to cache http response", zap.String("key", fmt.Sprint(personId)))
		}
	} else {
		debugLogger.Debug("response fetched from cache",
			zap.String("cachekey", fmt.Sprint(personId)),
		)
		respBody = storedRespBody
	}

	var parsedResp GetEmployeesAbsenceStatusResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		return 0, &employeerepoerrors.HTTPResponseFailedToParseBodyError{InternalError: err}
	}
	if err := parsedResp.validate(); err != nil {
		return 0, err
	}

	reasonId := parsedResp.Data[0].ReasonId

	return reasonId, nil
}

func (r *HTTPEmployeeRepository) makeAndDoPOSTRequest(ctx context.Context, reqBody []byte, protocol string, route string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s://%s:%s/%s", protocol, r.host, r.port, route),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, &employeerepoerrors.HTTPServerUnavailableError{InternalError: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &employeerepoerrors.HTTPResponseBadHTTPStatusError{Status: resp.StatusCode, ResponseBody: string(respBody)}
	}

	return respBody, nil
}

func getDateTimeRangeInRFC3339() (string, string) {
	currUTCTime := time.Now().UTC()
	currYear := currUTCTime.Year()
	currMonth := currUTCTime.Month()
	currDay := currUTCTime.Day()

	dateFrom := time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(currYear, currMonth, currDay, 23, 59, 59, 0, time.UTC)
	return dateFrom.Format(time.RFC3339), dateTo.Format(time.RFC3339)
}
