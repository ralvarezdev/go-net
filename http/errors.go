package http

import (
	"errors"
)

const (
	ErrInvalidRequestBody = "invalid request body type, expected: %v"
	ErrNilSubmodule       = "%s: submodule at index %d is nil"
)

var (
	ErrCodeCookieNotFound *string
)

var (
	ErrCookieNotFound = errors.New("cookie not found")
	ErrNilRequestBody = errors.New("request body cannot be nil")
	ErrInDevelopment  = errors.New("in development")
	ErrNilModule      = errors.New("module cannot be nil")
)
