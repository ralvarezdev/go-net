package size_limiter

import (
	"net/http"
)

// Middleware struct is the size limiter middleware
type Middleware struct{}

// NewMiddleware returns a new instance of the size limiter middleware
func NewMiddleware() *Middleware {
	return &Middleware{}
}

// Limit is the size limiter middleware
func (m *Middleware) Limit(bytesSizeLimit int64) func(next http.Handler) http.Handler {
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
