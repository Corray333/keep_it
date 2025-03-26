package errs

import (
	"errors"
)

var (
	ErrWrongVerificationCode = errors.New("wrong verification code")
	ErrWrongCodeRequestType  = errors.New("wrong code request type")
	ErrWrongPassword         = errors.New("wrong password")
	ErrWrongPasswordFormat   = errors.New("wrong password format")
)
