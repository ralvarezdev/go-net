package size_limiter

import (
	"net/http"
)

// SizeLimiter is the interface for the size limiter middleware
type SizeLimiter interface {
	Limit(bytesSizeLimit int64) func(next http.Handler) http.Handler
}
