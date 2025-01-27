package validator

import (
	"errors"
)

var (
	ErrNilValidateFn     = errors.New("validate function is nil")
	ErrInvalidValidateFn = errors.New("validate function is invalid")
)
