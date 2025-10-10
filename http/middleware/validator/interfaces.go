package validator

import (
	"net/http"
)

type (
	// Validator interface
	Validator interface {
		CreateValidateFn(
			body interface{},
			cache bool,
			auxiliaryValidatorFns ...interface{},
		) (func(next http.Handler) http.Handler, error)
		Validate(
			body interface{},
			auxiliaryValidatorFns ...interface{},
		) func(next http.Handler) http.Handler
	}
)
