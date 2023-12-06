package employeerepoerrors

import "fmt"

type FailedToCacheError struct {
	InternalError error
}

func (e *FailedToCacheError) Error() string {
	return fmt.Sprintf("failed to cache response: %s", e.InternalError.Error())
}
func (e *FailedToCacheError) Unwrap() error {
	return e.InternalError
}

type HTTPServerUnavailableError struct {
	InternalError error
}

func (e *HTTPServerUnavailableError) Error() string {
	return fmt.Sprintf("http request finished with error: %s", e.InternalError.Error())
}
func (e *HTTPServerUnavailableError) Unwrap() error {
	return e.InternalError
}

type HTTPResponseBadStatusError struct {
	Status string
}

func (e *HTTPResponseBadStatusError) Error() string {
	return fmt.Sprintf("bad status recevied: %s", e.Status)
}

type HTTPResponseNoDataError struct{}

func (e *HTTPResponseNoDataError) Error() string {
	return "no data reseived from server"
}

type HTTPResponseBadHTTPStatusError struct {
	ResponseBody string
	Status       int
}

func (e *HTTPResponseBadHTTPStatusError) Error() string {
	return fmt.Sprintf("http server responded with statuscode %d: %s", e.Status, e.ResponseBody)
}

type HTTPResponseFailedToParseBodyError struct {
	InternalError error
}

func (e *HTTPResponseFailedToParseBodyError) Error() string {
	return fmt.Sprintf("unnable to parse http response body: %s", e.InternalError.Error())
}
func (e *HTTPResponseFailedToParseBodyError) Unwrap() error {
	return e.InternalError
}
