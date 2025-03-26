package errs

import "errors"

var (
	ErrUploadingFile = errors.New("error uploading file")
	ErrGettingFile   = errors.New("error getting file")
)
