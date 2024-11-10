package helpers

import (
	"errors"
	"fmt"
)

// CustomError struct that implements the error interface
type CustomError struct {
	Err  error
	Code int
}

// Error method to implement the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Err.Error())
}

func NewError(code int, err error) *CustomError {
	return &CustomError{
		Err:  err,
		Code: code,
	}
}

var (
	ErrInternal = NewError(500, errors.New("internal error"))
)
