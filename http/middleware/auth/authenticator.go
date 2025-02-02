package auth

import (
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	"net/http"
)

// Authenticator interface
type Authenticator interface {
	Authenticate(
		failHandler func(
			w http.ResponseWriter,
			err string,
			httpStatus int,
			errorCode *string,
		),
		token gojwttoken.Token,
		rawToken string,
	) func(next http.Handler) http.Handler
	AuthenticateFromHeader(
		token gojwttoken.Token,
	) func(next http.Handler) http.Handler
	AuthenticateFromCookie(
		token gojwttoken.Token,
		cookieName string,
	) func(next http.Handler) http.Handler
}
