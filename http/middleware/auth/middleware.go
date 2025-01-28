package auth

import (
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtnethttp "github.com/ralvarezdev/go-jwt/net/http"
	gojwtnethttpctx "github.com/ralvarezdev/go-jwt/net/http/context"
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
	gojwtvalidator "github.com/ralvarezdev/go-jwt/token/validator"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttpjwtvalidator "github.com/ralvarezdev/go-net/http/jwt/validator"
	"net/http"
	"strings"
)

// Middleware struct
type Middleware struct {
	validator               gojwtvalidator.Validator
	handler                 gonethttphandler.Handler
	jwtValidatorFailHandler gonethttpjwtvalidator.FailHandler
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(
	validator gojwtvalidator.Validator,
	handler gonethttphandler.Handler,
	jwtValidatorFailHandler gonethttpjwtvalidator.FailHandler,
) (*Middleware, error) {
	// Check if either the validator, response handler or validator handler is nil
	if validator == nil {
		return nil, gojwtvalidator.ErrNilValidator
	}
	if handler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}
	if jwtValidatorFailHandler == nil {
		return nil, gonethttpjwtvalidator.ErrNilFailHandler
	}

	return &Middleware{
		validator,
		handler,
		jwtValidatorFailHandler,
	}, nil
}

// Authenticate return the middleware function that authenticates the request
func (m *Middleware) Authenticate(
	interception gojwtinterception.Interception,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Get the authorization from the header
				authorization := r.Header.Get(gojwtnethttp.AuthorizationHeaderKey)

				// Check if the authorization is a bearer token
				parts := strings.Split(authorization, " ")

				// Return an error if the authorization is missing or invalid
				if len(parts) < 2 || parts[0] != gojwt.BearerPrefix {
					m.jwtValidatorFailHandler(
						w,
						ErrInvalidAuthorizationHeader,
					)
					return
				}

				// Get the raw token from the header
				rawToken := parts[1]

				// Validate the token and get the validated claims
				claims, err := m.validator.GetValidatedClaims(
					rawToken,
					interception,
				)
				if err != nil {
					m.jwtValidatorFailHandler(w, err)
					return
				}

				// Set the token claims to the context
				r = gojwtnethttpctx.SetCtxTokenClaims(r, claims)

				// Call the next handler
				next.ServeHTTP(w, r)
			},
		)
	}
}
