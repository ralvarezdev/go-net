package response

import (
	"errors"
)

const (
	ErrInvalidFieldValueType = "invalid field value type, expected: '%s'"
)

var (
	ErrNilResponse      = errors.New("response cannot be nil")
	ErrNilEncoder       = errors.New("json encoder is nil")
	ErrNilFailBodyError = errors.New("fail body error cannot be nil")
)
