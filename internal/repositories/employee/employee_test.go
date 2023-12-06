package employeerepo_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/alexeipyp/test_grpc/internal/models"
	employeerepo "github.com/alexeipyp/test_grpc/internal/repositories/employee"
	employeerepoerrors "github.com/alexeipyp/test_grpc/internal/repositories/employee/errors"
	mock_employeerepo "github.com/alexeipyp/test_grpc/internal/repositories/employee/mocks"
	"github.com/jarcoal/httpmock"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestGetEmployeeInfo(t *testing.T) {
	type mockBehavior func(c *mock_employeerepo.MockHTTPClient)
	sampleEmployee := models.Employee{
		Name:      "sample name",
		Email:     "sample@example.org",
		WorkPhone: "1234",
		Id:        1234,
	}

	tests := []struct {
		name           string
		mockBehavior   mockBehavior
		expectedError  error
		expectedOutput *models.Employee
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedOutput: &sampleEmployee,
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				respJson := &employeerepo.GetEmployeesInfoResponse{
					Status: employeerepo.OkStatus,
					Data: []models.Employee{
						sampleEmployee,
					},
				}

				resp, err := httpmock.NewJsonResponse(200, respJson)
				if err != nil {
					t.Fatalf("Error in httpmock: %v", err)
				}
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
		{
			name:           "External HTTP server unavailable",
			expectedOutput: nil,
			expectedError:  &employeerepoerrors.HTTPServerUnavailableError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("server unavailable"))
			},
		},
		{
			name:           "External HTTP server responded with non-OK (not 200) status",
			expectedOutput: nil,
			expectedError:  &employeerepoerrors.HTTPResponseBadHTTPStatusError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				resp := httpmock.NewStringResponse(500, "internal error")
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
		{
			name:           "Bad status in HTTP OK server response body",
			expectedOutput: nil,
			expectedError:  &employeerepoerrors.HTTPResponseBadStatusError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				respJson := &employeerepo.GetEmployeesInfoResponse{
					Status: "Some bad status",
				}
				resp, err := httpmock.NewJsonResponse(200, respJson)
				if err != nil {
					t.Fatalf("Error in httpmock: %v", err)
				}
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
		{
			name:           "Empty data field in HTTP OK server response body",
			expectedOutput: nil,
			expectedError:  &employeerepoerrors.HTTPResponseNoDataError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				respJson := &employeerepo.GetEmployeesInfoResponse{
					Status: employeerepo.OkStatus,
				}
				resp, err := httpmock.NewJsonResponse(200, respJson)
				if err != nil {
					t.Fatalf("Error in httpmock: %v", err)
				}
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
		{
			name:           "Unknown HTTP response body structure",
			expectedOutput: nil,
			expectedError:  &employeerepoerrors.HTTPResponseFailedToParseBodyError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				resp := httpmock.NewStringResponse(200, "some resp")
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Arrange
			c := gomock.NewController(t)
			defer c.Finish()

			httpClient := mock_employeerepo.NewMockHTTPClient(c)
			test.mockBehavior(httpClient)
			cache := mock_employeerepo.NewMockCache(c)
			cache.EXPECT().Get(gomock.Any()).Return(nil, errors.New("unavailable")).AnyTimes()
			cache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(errors.New("unavailable")).AnyTimes()

			empRepo := employeerepo.New(httpClient, cache, "example.org", "5432", zap.NewNop(), zap.NewNop())

			//Act
			res, err := empRepo.GetEmployeeInfo(context.Background(), "sample@example.org")

			//Accert
			if reflect.TypeOf(err) != reflect.TypeOf(test.expectedError) {
				t.Errorf("GetEmployeeInfo - incorrect error, got: %v, want: %v", reflect.TypeOf(err), reflect.TypeOf(test.expectedError))
			}
			if !reflect.DeepEqual(res, test.expectedOutput) {
				t.Errorf("GetEmployeeInfo - incorrect output, got: %+v, want: %+v", res, test.expectedOutput)
			}
		})
	}
}

func TestEmployeeAbsenceStatus(t *testing.T) {
	type mockBehavior func(c *mock_employeerepo.MockHTTPClient)

	tests := []struct {
		name           string
		mockBehavior   mockBehavior
		expectedError  error
		expectedOutput int
	}{
		{
			name:           "Success",
			expectedOutput: 1,
			expectedError:  nil,
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				respJson := &employeerepo.GetEmployeesAbsenceStatusResponse{
					Status: employeerepo.OkStatus,
					Data: []employeerepo.EmployeeAbsenceInfo{
						{
							ReasonId: 1,
						},
					},
				}

				resp, err := httpmock.NewJsonResponse(200, respJson)
				if err != nil {
					t.Fatalf("Error in httpmock: %v", err)
				}
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
		{
			name:           "External HTTP server unavailable",
			expectedOutput: 0,
			expectedError:  &employeerepoerrors.HTTPServerUnavailableError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				c.EXPECT().Do(gomock.Any()).Return(nil, errors.New("server unavailable"))
			},
		},
		{
			name:           "External HTTP server responded with non-OK (not 200) status",
			expectedOutput: 0,
			expectedError:  &employeerepoerrors.HTTPResponseBadHTTPStatusError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				resp := httpmock.NewStringResponse(500, "internal error")
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
		{
			name:           "Bad status in HTTP OK server response body",
			expectedOutput: 0,
			expectedError:  &employeerepoerrors.HTTPResponseBadStatusError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				respJson := &employeerepo.GetEmployeesInfoResponse{
					Status: "Some bad status",
				}
				resp, err := httpmock.NewJsonResponse(200, respJson)
				if err != nil {
					t.Fatalf("Error in httpmock: %v", err)
				}
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
		{
			name:           "Empty data field in HTTP OK server response body",
			expectedOutput: 0,
			expectedError:  &employeerepoerrors.HTTPResponseNoDataError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				respJson := &employeerepo.GetEmployeesInfoResponse{
					Status: employeerepo.OkStatus,
				}
				resp, err := httpmock.NewJsonResponse(200, respJson)
				if err != nil {
					t.Fatalf("Error in httpmock: %v", err)
				}
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
		{
			name:           "Unknown HTTP response body structure",
			expectedOutput: 0,
			expectedError:  &employeerepoerrors.HTTPResponseFailedToParseBodyError{},
			mockBehavior: func(c *mock_employeerepo.MockHTTPClient) {
				resp := httpmock.NewStringResponse(200, "some resp")
				c.EXPECT().Do(gomock.Any()).Return(resp, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Arrange
			c := gomock.NewController(t)
			defer c.Finish()

			httpClient := mock_employeerepo.NewMockHTTPClient(c)
			test.mockBehavior(httpClient)
			cache := mock_employeerepo.NewMockCache(c)
			cache.EXPECT().Get(gomock.Any()).Return(nil, errors.New("unavailable")).AnyTimes()
			cache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(errors.New("unavailable")).AnyTimes()

			empRepo := employeerepo.New(httpClient, cache, "example.org", "5432", zap.NewNop(), zap.NewNop())

			//Act
			res, err := empRepo.GetEmployeeAbsenceStatus(context.Background(), 1)

			//Accert
			if reflect.TypeOf(err) != reflect.TypeOf(test.expectedError) {
				t.Errorf("GetEmployeeInfo - incorrect error, got: %v, want: %v", reflect.TypeOf(err), reflect.TypeOf(test.expectedError))
			}
			if res != test.expectedOutput {
				t.Errorf("GetEmployeeInfo - incorrect reasonId, got: %v, want: %v", res, test.expectedOutput)
			}
		})
	}
}
