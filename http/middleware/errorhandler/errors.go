package errorhandler

import (
	"errors"
)

var (
	ErrNilErrorHandler = errors.New("nil error handler")
)
