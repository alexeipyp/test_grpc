package employeegrpcerrors

import "fmt"

type ValidationError struct {
	InternalError error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("incorrect input: %v", e.InternalError)
}
func (e *ValidationError) Unwrap() error {
	return e.InternalError
}
