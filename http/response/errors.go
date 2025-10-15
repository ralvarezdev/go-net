package response

import (
	"errors"
)

const (
	ErrInvalidFieldValueType = "invalid field value type, expected: '%s'"
)

var (
	ErrCodeNilFailError string
	ErrCodeNilError     string
)

var (
	ErrNilResponse  = errors.New("response cannot be nil")
	ErrNilEncoder   = errors.New("json encoder is nil")
	ErrNilFailError = errors.New("fail error cannot be nil")
	ErrNilError     = errors.New("error cannot be nil")
)
