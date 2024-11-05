package helpers

import "fmt"

// CustomError struct that implements the error interface
type CustomError struct {
	Message string
	Code    int
}

// Error method to implement the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func NewError(code int, message string) *CustomError {
	return &CustomError{
		Message: message,
		Code:    code,
	}
}

var (
	ErrInternal = NewError(500, "internal error")
)
