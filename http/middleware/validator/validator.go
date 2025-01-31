package validator

import (
	"net/http"
)

// Validator interface
type Validator interface {
	Validate(
		body,
		auxiliaryValidatorFn interface{},
	) func(next http.Handler) http.Handler
}
