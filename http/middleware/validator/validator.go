package validator

import (
	"net/http"
)

// Validator interface
type Validator interface {
	Validate(
		body,
		createValidateFn interface{},
	) func(next http.Handler) http.Handler
}
