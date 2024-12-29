package auth

import (
	"errors"
)

var (
	ErrNilAuthenticator = errors.New("authenticator cannot be nil")
)
