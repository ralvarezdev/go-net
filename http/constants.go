package http

import (
	"net/http"
)

const (
	// XForwardedFor is the header key for the X-Forwarded-For header
	XForwardedFor = "X-Forwarded-For"

	// Authorization is the header key for the Authorization header
	Authorization = "Authorization"
)

var (
	BadRequest          = http.StatusText(http.StatusBadRequest)
	InternalServerError = http.StatusText(http.StatusInternalServerError)
	ServiceUnavailable  = http.StatusText(http.StatusServiceUnavailable)
	Unauthorized        = http.StatusText(http.StatusUnauthorized)
	NotImplemented      = http.StatusText(http.StatusNotImplemented)
	TooManyRequests     = http.StatusText(http.StatusTooManyRequests)
	RequestTimeout      = http.StatusText(http.StatusRequestTimeout)
	NotFound            = http.StatusText(http.StatusNotFound)
)
