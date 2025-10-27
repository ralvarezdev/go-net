package redis

import (
	"errors"
	"log/slog"
	"net/http"

	goratelimiterredis "github.com/ralvarezdev/go-rate-limiter/redis"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
)

type (
	// Middleware struct
	Middleware struct {
		responsesHandler gonethttphandler.ResponsesHandler
		rateLimiter      goratelimiterredis.RateLimiter
		logger           *slog.Logger
	}
)

// NewMiddleware creates a new rate limiter middleware
//
// Parameters:
//
// responsesHandler gonethttphandler.ResponsesHandler: the HTTP handler to handle errors
// rateLimiter goratelimiterredis.RateLimiter: the rate limiter
// logger *slog.Logger: the logger (optional)
//
// Returns:
//
// *Middleware: the middleware instance
// error: if the rate limiter is nil
func NewMiddleware(
	responsesHandler gonethttphandler.ResponsesHandler,
	rateLimiter goratelimiterredis.RateLimiter,
	logger *slog.Logger,
) (
	*Middleware,
	error,
) {
	// Check if the handler is nil
	if responsesHandler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}

	// Check if the rate limiter is nil
	if rateLimiter == nil {
		return nil, goratelimiterredis.ErrNilRateLimiter
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "http_middleware_rate_limiter_redis"),
		)
	}

	return &Middleware{
		responsesHandler,
		rateLimiter,
		logger,
	}, nil
}

// Limit limits the number of requests per IP address
//
// Returns:
//
//	func(next http.ResponsesHandler) http.ResponsesHandler: the middleware function
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
							gonethttp.TooManyRequests,
							http.StatusTooManyRequests,
						)
						return
					}

					// Log the error
					if m.logger != nil {
						m.logger.Error(
							"Error limiting requests",
							slog.String("ip", ip),
							slog.String("error", err.Error()),
						)
					}

					// Handle other errors
					m.responsesHandler.HandleRawError(w, r, err, nil)
					return
				}

				// Call the next handler
				next.ServeHTTP(w, r)
			},
		)
	}
}
