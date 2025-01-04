package http

import (
	"net/http"
)

var (
	BadRequest          = http.StatusText(http.StatusBadRequest)
	InternalServerError = http.StatusText(http.StatusInternalServerError)
	ServiceUnavailable  = http.StatusText(http.StatusServiceUnavailable)
	Unauthorized        = http.StatusText(http.StatusUnauthorized)
)
