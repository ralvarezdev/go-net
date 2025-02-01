package http

import (
	"errors"
)

var (
	ErrCodeCookieNotFound *string
)

var (
	ErrCookieNotFound     = errors.New("cookie not found")
	ErrNilRequestBody     = errors.New("request body cannot be nil")
	ErrInvalidRequestBody = "invalid request body type, expected: %v"
	ErrInDevelopment      = errors.New("in development")
	ErrNilSubmodule       = "%s: submodule at index %d is nil"
)
