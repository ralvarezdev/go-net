package auth

import (
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
	"net/http"
)

// Authenticator interface
type Authenticator interface {
	Authenticate(
		interception gojwtinterception.Interception,
	) func(next http.Handler) http.Handler
}
