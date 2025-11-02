package validator

import (
	"net/http"
)

type (
	// Validator interface
	Validator interface {
		CreateValidateFn(
			bodyExample any,
			cache bool,
			auxiliaryValidatorFns ...any,
		) (func(next http.Handler) http.Handler, error)
		Validate(
			bodyExample any,
			auxiliaryValidatorFns ...any,
		) func(next http.Handler) http.Handler
	}
)
