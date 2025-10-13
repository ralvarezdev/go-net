package grpc

import (
	"net/http"
)

type (
	// Authenticator interface
	Authenticator interface {
		AuthenticateFromHeader(
			rpcMethod string,
		) func(next http.Handler) http.Handler
		AuthenticateFromCookie(
			rpcMethod string,
		) func(next http.Handler) http.Handler
	}
)
