package http

import (
	"errors"
)

var (
	ErrNilRequestBody = errors.New("request body cannot be nil")
	ErrInDevelopment  = errors.New("in development")
)
