package http

import (
	"errors"
	gonethttp "github.com/ralvarezdev/go-net/http"
)

var (
	BadRequest                    = errors.New(gonethttp.BadRequest)
	InternalServerError           = errors.New(gonethttp.InternalServerError)
	ServiceUnavailable            = errors.New(gonethttp.ServiceUnavailable)
	Unauthorized                  = errors.New(gonethttp.Unauthorized)
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	Unauthenticated               = errors.New("missing or invalid bearer token on authentication header")
	InDevelopment                 = errors.New("in development")
)