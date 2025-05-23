package status

import (
	"net/http"
)

var (
	BadRequest          = http.StatusText(http.StatusBadRequest)
	InternalServerError = http.StatusText(http.StatusInternalServerError)
	ServiceUnavailable  = http.StatusText(http.StatusServiceUnavailable)
	Unauthorized        = http.StatusText(http.StatusUnauthorized)
	NotImplemented      = http.StatusText(http.StatusNotImplemented)
)

const (
	// XForwardedFor is the header key for the X-Forwarded-For header
	XForwardedFor = "X-Forwarded-For"
)
