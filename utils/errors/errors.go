package errors

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// NewErrorWithStack create a new error with given message and stack info.
func NewErrorWithStack(msg string) error {
	return errors.New(msg)
}

// WrapErrorWithStack wrap an exist error with stack.
// If err is nil, WrapfErrorWithStack returns nil.
func WrapErrorWithStack(err error) error {
	return errors.WithStack(err)
}

// WrapfErrorWithStack wrap an exist error with given format message and stack info.
// If err is nil, WrapfErrorWithStack returns nil.
func WrapfErrorWithStack(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// AttachErrorMessage attach a format message to a given erros.
func AttachErrorMessage(err error, format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args...)
}

// SprintError return the error string with both message and stack.
func SprintError(err error) string {
	if err == nil {
		return ""
	}

	return fmt.Sprintf("%s\nStack trace:\n\t%s",
		GetErrorMessage(err),
		strings.ReplaceAll(GetErrorStack(err), "\n", "\n\t"),
	)
}

// GetErrorMessage get the error message without stack.
func GetErrorMessage(err error) string {
	return err.Error()
}

// GetErrorStack get the stack of error.
func GetErrorStack(err error) string {
	errMsg := fmt.Sprintf("%+v", err)
	return cleanPath(errMsg)
}

// cleanPath remove the parent path of current work directory.
func cleanPath(s string) string {
	return strings.ReplaceAll(s, getCurrentPath()+"/", "")
}

// getCurrentPath get current work directory.
func getCurrentPath() string {
	getwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return strings.Replace(getwd, "\\", "/", -1)
}

func Cause(err error) error {
	return errors.Cause(err)
}
