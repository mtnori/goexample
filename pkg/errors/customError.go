package errors

import "fmt"

type CustomError struct {
	OriginalError error
	Message       string
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("%s: %v", ce.Message, ce.OriginalError)
}

func WrapError(err error, message string) error {
	return &CustomError{
		OriginalError: err,
		Message:       message,
	}
}

func UnwrapError(err error) error {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.OriginalError
	}
	return err
}
