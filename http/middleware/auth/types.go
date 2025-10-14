package auth

import (
	"net/http"

	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// FailHandlerFn defines the function signature for handling authentication failures
	FailHandlerFn func(
		w http.ResponseWriter,
		err error,
		errorCode string,
	)

	// RefreshTokenFn defines the function signature for refreshing tokens
	RefreshTokenFn func(
		w http.ResponseWriter,
		r *http.Request,
	) (map[gojwttoken.Token]string, error)
)
