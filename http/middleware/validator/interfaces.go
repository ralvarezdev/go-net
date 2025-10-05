package validator

import (
	"net/http"
)

type (
	// Validator interface
	Validator interface {
		Validate(
			body,
			auxiliaryValidatorFn interface{},
		) func(next http.Handler) http.Handler
	}
)
