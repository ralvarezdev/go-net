package errorhandler

import (
	"net/http"
)

type (
	// ErrorHandler interface for handling the errors
	ErrorHandler interface {
		HandleError(next http.Handler) http.Handler
	}
)
