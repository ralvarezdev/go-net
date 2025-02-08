package auth

import (
	"errors"
)

var (
	ErrCodeInvalidAuthorizationHeader *string
	ErrCodeInvalidTokenClaims         *string
	ErrCodeFailedToRefreshToken       *string
)

var (
	ErrNilAuthenticator           = errors.New("authenticator cannot be nil")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
)
