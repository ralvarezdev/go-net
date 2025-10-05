package size_limiter

import (
	"net/http"
)

type (
	// SizeLimiter is the interface for the size limiter middleware
	SizeLimiter interface {
		Limit(bytesSizeLimit int64) func(next http.Handler) http.Handler
	}
)
