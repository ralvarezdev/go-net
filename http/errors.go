package http

import (
	"errors"
)

var (
	ErrInvalidRequestBody         = "invalid request body: %v"
	ErrNilRequestBody             = errors.New("request body cannot be nil")
	ErrInDevelopment              = errors.New("in development")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
)
