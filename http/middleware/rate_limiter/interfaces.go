package rate_limiter

import (
	"net/http"
)

type (
	// RateLimiter interface
	RateLimiter interface {
		Limit() func(next http.Handler) http.Handler
	}
)
