package error_handler

import (
	"fmt"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	"net/http"
)

// Middleware struct
type Middleware struct {
	handler gonethttphandler.Handler
}

// NewMiddleware creates a new error handler middleware
func NewMiddleware(handler gonethttphandler.Handler) (*Middleware, error) {
	// Check if the handler is nil
	if handler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}
	return &Middleware{
		handler,
	}, nil
}

// HandleError handles the error
func (m *Middleware) HandleError(next http.Handler) http.Handler {
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
					m.handler.HandleError(w, err)
				}
			}()

			// Call the next handler
			next.ServeHTTP(w, r)
		},
	)
}
