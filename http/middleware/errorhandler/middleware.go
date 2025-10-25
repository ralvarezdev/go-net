package errorhandler

import (
	"fmt"
	"net/http"

	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
)

type (
	// Middleware struct is the error handler middleware
	Middleware struct {
		responsesHandler gonethttphandler.ResponsesHandler
	}
)

// NewMiddleware creates a new error handler middleware
//
// Parameters:
//
//   - responsesHandler: The HTTP handler to handle errors
//
// Returns:
//
//   - *Middleware: The error handler middleware
//   - error: The error if any
func NewMiddleware(responsesHandler gonethttphandler.ResponsesHandler) (
	*Middleware,
	error,
) {
	// Check if the handler is nil
	if responsesHandler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}
	return &Middleware{
		responsesHandler,
	}, nil
}

// HandleError handles the error
//
// Parameters:
//
//   - next: The next HTTP handler
//
// Returns:
//
//   - http.Handler: The HTTP handler
func (m Middleware) HandleError(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if value := recover(); value != nil {
					// Parse the error
					err, ok := value.(error)
					if !ok {
						err = fmt.Errorf("%v", value)
					}

					// Handle the error
					m.responsesHandler.HandleRawError(w, r, err)
				}
			}()

			// Call the next handler
			next.ServeHTTP(w, r)
		},
	)
}
