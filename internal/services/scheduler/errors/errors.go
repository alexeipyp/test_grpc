package schedulerserviceerrors

import "fmt"

type EnqueueTaskError struct {
	InternalError error
}

func (e *EnqueueTaskError) Error() string {
	return fmt.Sprintf("error occurred with task queue: %s", e.InternalError.Error())
}
func (e *EnqueueTaskError) Unwrap() error {
	return e.InternalError
}
