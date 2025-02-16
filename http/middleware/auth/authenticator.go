package auth

import (
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	"net/http"
)

// Authenticator interface
type Authenticator interface {
	Authenticate(
		token gojwttoken.Token,
		rawToken string,
		failHandler func(
			w http.ResponseWriter,
			err error,
			errorCode *string,
		),
	) func(next http.Handler) http.Handler
	AuthenticateFromHeader(
		token gojwttoken.Token,
	) func(next http.Handler) http.Handler
	AuthenticateFromCookie(
		token gojwttoken.Token,
		cookieRefreshTokenName,
		cookieAccessTokenName string,
		refreshTokenFn func(
			w http.ResponseWriter,
			r *http.Request,
		) (*map[gojwttoken.Token]string, error),
	) func(next http.Handler) http.Handler
}
