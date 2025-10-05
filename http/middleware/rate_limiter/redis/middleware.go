package redis

import (
	"errors"
	"net/http"

	gonethttp "github.com/ralvarezdev/go-net/http"
	goratelimiterredis "github.com/ralvarezdev/go-rate-limiter/redis"
)

type (
	// Middleware struct
	Middleware struct {
		rateLimiter goratelimiterredis.RateLimiter
	}
)

// NewMiddleware creates a new rate limiter middleware
//
// Parameters:
//
// rateLimiter goratelimiterredis.RateLimiter: the rate limiter
//
// Returns:
//
// *Middleware: the middleware instance
// error: if the rate limiter is nil
func NewMiddleware(rateLimiter goratelimiterredis.RateLimiter) (
	*Middleware,
	error,
) {
	// Check if the rate limiter is nil
	if rateLimiter == nil {
		return nil, goratelimiterredis.ErrNilRateLimiter
	}

	return &Middleware{
		rateLimiter,
	}, nil
}

// Limit limits the number of requests per IP address
//
// Returns:
//
//	func(next http.Handler) http.Handler: the middleware function
func (m Middleware) Limit() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Get the client IP address
				ip := gonethttp.GetClientIP(r)

				// Limit the number of requests per IP address
				if err := m.rateLimiter.Limit(ip); err != nil {
					// Check if the rate limit is exceeded
					if errors.Is(err, goratelimiterredis.ErrTooManyRequests) {
						http.Error(
							w,
							http.StatusText(http.StatusTooManyRequests),
							http.StatusTooManyRequests,
						)
						return
					}

					// Handle other errors
					http.Error(
						w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
					return
				}

				// Call the next handler
				next.ServeHTTP(w, r)
			},
		)
	}
}
