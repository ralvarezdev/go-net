package response

import (
	"errors"
)

var (
	ErrNilResponse           = errors.New("response cannot be nil")
	ErrNilResponseHTTPStatus = errors.New("response http status cannot be nil")
)
