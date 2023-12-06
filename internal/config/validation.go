package config

import (
	"errors"
	"fmt"
	"time"
)

func ValidateEnv(env string) error {
	if !(env == EnvDev || env == EnvProd) {
		return fmt.Errorf("incorrect value of env: %s - unknown environment", env)
	}
	return nil
}

func (c SchedulerConfig) Validate() error {
	errList := []error{
		validateInt(c.QueueSize, "queue size"),
		validateInt(c.WorkerPoolSize, "worker pool size"),
	}
	if isErrSliceContainsNonNilValue(errList) {
		return constructAnError(
			"validation failed at scheduler config",
			": ",
			errList...,
		)
	}
	return nil
}

func (c CacheConfig) Validate() error {
	errList := []error{
		validateDuration(c.Lifetime, "lifetime"),
	}
	if isErrSliceContainsNonNilValue(errList) {
		return constructAnError(
			"validation failed at cache config",
			": ",
			errList...,
		)
	}
	return nil
}

func (c ExternalHTTPConnectionConfig) Validate() error {
	errList := []error{
		validateDuration(c.Timeout, "timeout"),
		validateStr(c.Host, "host"),
		validateInt(c.Port, "port"),
		validateStr(c.Username, "username"),
	}
	if isErrSliceContainsNonNilValue(errList) {
		return constructAnError(
			"validation failed at connection config",
			": ",
			errList...,
		)
	}
	return nil
}

func (c GRPCConfig) Validate() error {
	errList := []error{
		validateStr(c.Host, "host"),
		validateInt(c.Port, "port"),
	}
	if isErrSliceContainsNonNilValue(errList) {
		return constructAnError(
			"validation failed at grpc config",
			": ",
			errList...,
		)
	}
	return nil
}

func (c LogConfig) Validate() error {
	errList := []error{
		validateStr(c.LogFilename, "log filepath"),
		validateStr(c.DebugLogFilename, "debug log filepath"),
		validateStr(c.GRPCTraceLogFilename, "grpc trace log filepath"),
		validateStr(c.HTTPTraceLogFilename, "http trace log filepath"),
	}
	if isErrSliceContainsNonNilValue(errList) {
		return constructAnError(
			"validation failed at log config",
			": ",
			errList...,
		)
	}
	return nil
}

func validateDuration(p ParsedDuration, paramName string) error {
	dur, err := p.TryDuration()
	if err != nil {
		return fmt.Errorf("incorrect value of %s: %w", paramName, err)
	}
	if dur <= time.Duration(0) {
		return fmt.Errorf("incorrect value of %s: %v - should be non-zero", paramName, dur)
	}
	return nil
}

func validateInt(i int, paramName string) error {
	if i <= 0 {
		return fmt.Errorf("incorrect value of %s: %d - should be non-zero positive", paramName, i)
	}
	return nil
}

func validateStr(str string, paramName string) error {
	if str == "" {
		return fmt.Errorf("incorrect value of %s: empty value", paramName)
	}
	return nil
}

func isErrSliceContainsNonNilValue(errs []error) bool {
	for _, err := range errs {
		if err != nil {
			return true
		}
	}
	return false
}

func constructAnError(message string, delimiter string, errs ...error) error {
	resErr := errors.New(message)
	for _, err := range errs {
		if err != nil {
			resErr = fmt.Errorf("%w%s%w", resErr, delimiter, err)
		}
	}
	return resErr
}
