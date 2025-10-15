package response

import (
	"errors"
)

const (
	ErrInvalidFieldValueType = "invalid field value type, expected: '%s'"
)

var (
	ErrCodeNilFailFieldError string
	ErrCodeNilFailDataError  string
	ErrCodeNilError          string
	ErrCodeNilResponse       string
)

var (
	ErrNilResponse       = errors.New("response cannot be nil")
	ErrNilEncoder        = errors.New("json encoder is nil")
	ErrNilFailFieldError = errors.New("fail field error cannot be nil")
	ErrNilFailDataError  = errors.New("fail data error cannot be nil")
	ErrNilError          = errors.New("error cannot be nil")
)
