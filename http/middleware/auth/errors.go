package auth

import (
	"errors"
	gojwtnethttp "github.com/ralvarezdev/go-jwt/net/http"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

var (
	ErrCodeInvalidAuthorizationHeader *string
	ErrCodeInvalidTokenClaims         *string
)

var (
	ErrNilAuthenticator           = errors.New("authenticator cannot be nil")
	ErrInvalidAuthorizationHeader = gonethttpresponse.NewHeaderError(
		gojwtnethttp.AuthorizationHeaderKey,
		"invalid authorization header",
		http.StatusUnauthorized,
		ErrCodeInvalidAuthorizationHeader,
	)
)
