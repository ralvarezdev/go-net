package response

import (
	"errors"
)

var (
	ErrNilResponse           = errors.New("response cannot be nil")
	ErrInvalidFieldValueType = "invalid field value type, expected: '%s'"
)
