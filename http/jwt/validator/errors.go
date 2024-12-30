package validator

import (
	"errors"
)

var (
	ErrNilErrorHandler = errors.New("error handler cannot be nil")
)
