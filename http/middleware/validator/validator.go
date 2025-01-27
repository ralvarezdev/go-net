package validator

import (
	"net/http"
)

// Validator interface
type Validator interface {
	Validate(
		createValidateFn func() (interface{}, func() (interface{}, error)),
	) func(next http.Handler) http.Handler
}
