package error_handler

import (
	"errors"
)

var (
	ErrNilErrorHandler = errors.New("nil error handler")
)
