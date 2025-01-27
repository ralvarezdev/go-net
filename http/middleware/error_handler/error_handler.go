package error_handler

import (
	"net/http"
)

// ErrorHandler interface for handling the errors
type ErrorHandler interface {
	HandleError(next http.Handler) http.Handler
}
