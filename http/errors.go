package http

import (
	"errors"
	gojwtnethttp "github.com/ralvarezdev/go-jwt/net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

var (
	ErrNilRequestBody             = errors.New("request body cannot be nil")
	ErrInDevelopment              = errors.New("in development")
	ErrInvalidAuthorizationHeader = gonethttpresponse.NewHeaderError(
		gojwtnethttp.AuthorizationHeaderKey,
		"invalid authorization header",
	)
)
