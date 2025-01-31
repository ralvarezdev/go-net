package errors

import (
	"errors"
	gonethttpstatus "github.com/ralvarezdev/go-net/http/status"
)

var (
	BadRequest          = errors.New(gonethttpstatus.BadRequest)
	InternalServerError = errors.New(gonethttpstatus.InternalServerError)
	ServiceUnavailable  = errors.New(gonethttpstatus.ServiceUnavailable)
	Unauthorized        = errors.New(gonethttpstatus.Unauthorized)
	NotImplemented      = errors.New(gonethttpstatus.NotImplemented)
	Unauthenticated     = errors.New("missing or invalid bearer token on authentication header")
)
