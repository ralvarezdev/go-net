package sizelimiter

import (
	"net/http"
)

type (
	// Middleware struct is the size limiter middleware
	Middleware struct{}
)

// NewMiddleware returns a new instance of the size limiter middleware
//
// Returns:
//
//   - *Middleware: The middleware instance
func NewMiddleware() *Middleware {
	return &Middleware{}
}

// Limit is the size limiter middleware
//
// Parameters:
//
//   - bytesSizeLimit: The maximum size of the request body in bytes
//
// Returns:
//
//   - func(next http.Handler) http.Handler: The middleware function
func (m Middleware) Limit(bytesSizeLimit int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				r.Body = http.MaxBytesReader(w, r.Body, bytesSizeLimit)

				// Call the next handler
				next.ServeHTTP(w, r)
			},
		)
	}
}
