package employeerepo

import (
	"github.com/alexeipyp/test_grpc/internal/models"
	employeerepoerrors "github.com/alexeipyp/test_grpc/internal/repositories/employee/errors"
)

type GetEmployeesInfoRequest struct {
	Name      string `json:"name,omitempty"`
	WorkPhone string `json:"workPhone,omitempty"`
	Email     string `json:"email,omitempty"`
	DateFrom  string `json:"dateFrom,omitempty"`
	DateTo    string `json:"dateTo,omitempty"`
	Ids       []int  `json:"ids,omitempty"`
}

type GetEmployeesInfoResponse struct {
	Status string            `json:"status"`
	Data   []models.Employee `json:"data"`
}

func (resp GetEmployeesInfoResponse) validate() error {
	if resp.Status != OkStatus {
		return &employeerepoerrors.HTTPResponseBadStatusError{Status: resp.Status}
	}
	if len(resp.Data) < 1 {
		return &employeerepoerrors.HTTPResponseNoDataError{}
	}
	return nil
}

type GetEmployeesAbsenceStatusRequest struct {
	DateFrom string `json:"dateFrom,omitempty"`
	DateTo   string `json:"dateTo,omitempty"`
	Ids      []int  `json:"personIds"`
}

type GetEmployeesAbsenceStatusResponse struct {
	Status string                `json:"status"`
	Data   []EmployeeAbsenceInfo `json:"data"`
}

func (resp GetEmployeesAbsenceStatusResponse) validate() error {
	if resp.Status != OkStatus {
		return &employeerepoerrors.HTTPResponseBadStatusError{Status: resp.Status}
	}
	if len(resp.Data) < 1 {
		return &employeerepoerrors.HTTPResponseNoDataError{}
	}
	return nil
}

type EmployeeAbsenceInfo struct {
	CreatedDate string `json:"createdDate"`
	DateFrom    string `json:"dateFrom,omitempty"`
	DateTo      string `json:"dateTo,omitempty"`
	Id          int    `json:"id"`
	PersonId    int    `json:"personId"`
	ReasonId    int    `json:"reasonId"`
}
