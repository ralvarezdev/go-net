package http

import (
	"errors"
	"net/http"
)

var (
	BadRequest                    = http.StatusText(http.StatusBadRequest)
	InternalServerError           = http.StatusText(http.StatusInternalServerError)
	ServiceUnavailable            = http.StatusText(http.StatusServiceUnavailable)
	Unauthorized                  = http.StatusText(http.StatusUnauthorized)
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	Unauthenticated               = errors.New("missing or invalid bearer token on authentication header")
	InDevelopment                 = errors.New("in development")
)
