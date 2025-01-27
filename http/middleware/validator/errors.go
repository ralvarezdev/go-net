package validator

import (
	"errors"
)

var (
	ErrNilValidateFn       = errors.New("validate function is nil")
	ErrInvalidValidateFn   = errors.New("validate function is invalid")
	ErrNilCreateValidateFn = errors.New("create validate function is nil")
	ErrNilParametersValues = errors.New("parameters value is nil")
)
