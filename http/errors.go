package http

import (
	"errors"
)

var (
	ErrCodeCookieNotFound string
)

const (
	ErrInvalidRequestBody = "invalid request body type, expected: %v"
	ErrNilSubmodule       = "%s: submodule at index %d is nil"
)

var (
	ErrCookieNotFound = errors.New("cookie not found")
	ErrNilRequestBody = errors.New("request body cannot be nil")
	ErrInDevelopment  = errors.New("in development")
	ErrNilModule      = errors.New("module cannot be nil")
)

var (
	ErrBadRequest          = errors.New(BadRequest)
	ErrInternalServerError = errors.New(InternalServerError)
	ErrServiceUnavailable  = errors.New(ServiceUnavailable)
	ErrUnauthorized        = errors.New(Unauthorized)
	ErrNotImplemented      = errors.New(NotImplemented)
	ErrTooManyRequests     = errors.New(TooManyRequests)
	ErrRequestTimeout      = errors.New(RequestTimeout)
	ErrNotFound            = errors.New(NotFound)
	ErrUnauthenticated     = errors.New("missing or invalid bearer token on authentication header")
)
