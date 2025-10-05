package auth

import (
	"net/http"

	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// Authenticator interface
	Authenticator interface {
		Authenticate(
			token gojwttoken.Token,
			rawToken string,
			failHandler FailHandlerFn,
		) func(next http.Handler) http.Handler
		AuthenticateFromHeader(
			token gojwttoken.Token,
		) func(next http.Handler) http.Handler
		AuthenticateFromCookie(
			token gojwttoken.Token,
			cookieRefreshTokenName,
			cookieAccessTokenName string,
			refreshTokenFn RefreshTokenFn,
		) func(next http.Handler) http.Handler
	}
)
