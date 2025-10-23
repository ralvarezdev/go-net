package auth

import (
	"errors"
)

var (
	ErrCodeInvalidAuthorizationHeader string
	ErrCodeInvalidTokenClaims         string
	ErrCodeFailedToRefreshToken       string
)

var (
	ErrNilAuthenticator              = errors.New("authenticator cannot be nil")
	ErrInvalidAuthorizationHeader    = errors.New("invalid authorization header")
	ErrNilOptions                    = errors.New("options cannot be nil")
	ErrNilRefreshTokenFn             = errors.New("refresh token function cannot be nil")
	ErrNilCookieRefreshTokenName     = errors.New("cookie refresh token name cannot be nil")
	ErrNilCookieAccessTokenName      = errors.New("cookie access token name cannot be nil")
	ErrCodeFailedToSetTokenInContext = errors.New("failed to set token in context")
)
