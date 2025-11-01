package validator

import (
	"net/http"
)

type (
	// Validator interface
	Validator interface {
		CreateValidateFn(
			body any,
			decode bool,
			cache bool,
			auxiliaryValidatorFns ...any,
		) (func(next http.Handler) http.Handler, error)
		Validate(
			body any,
			auxiliaryValidatorFns ...any,
		) func(next http.Handler) http.Handler
	}
)
