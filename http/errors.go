package http

import (
	"errors"
)

var (
	ErrInvalidRequestBody         = errors.New("invalid request body")
	ErrNilRequestBody             = errors.New("request body cannot be nil")
	ErrInDevelopment              = errors.New("in development")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
)
