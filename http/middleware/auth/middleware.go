package auth

import (
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtnethttp "github.com/ralvarezdev/go-jwt/net/http"
	gojwtnethttpctx "github.com/ralvarezdev/go-jwt/net/http/context"
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
	gojwtvalidator "github.com/ralvarezdev/go-jwt/token/validator"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttpjwtvalidator "github.com/ralvarezdev/go-net/http/jwt/validator"
	"github.com/ralvarezdev/go-net/http/response"
	"net/http"
	"strings"
)

// Middleware struct
type Middleware struct {
	validator                gojwtvalidator.Validator
	handler                  gonethttphandler.Handler
	jwtValidatorErrorHandler gonethttpjwtvalidator.ErrorHandler
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(
	validator gojwtvalidator.Validator,
	handler gonethttphandler.Handler,
	jwtValidatorErrorHandler gonethttpjwtvalidator.ErrorHandler,
) (*Middleware, error) {
	// Check if either the validator, response handler or validator handler is nil
	if validator == nil {
		return nil, gojwtvalidator.ErrNilValidator
	}
	if handler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}
	if jwtValidatorErrorHandler == nil {
		return nil, gonethttpjwtvalidator.ErrNilErrorHandler
	}

	return &Middleware{
		validator:                validator,
		handler:                  handler,
		jwtValidatorErrorHandler: jwtValidatorErrorHandler,
	}, nil
}

// Authenticate return the middleware function that authenticates the request
func (m *Middleware) Authenticate(
	interception gojwtinterception.Interception,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Get the context
			ctx := r.Context()

			// Get the authorization from the header
			authorization := ctx.Value(gojwtnethttp.AuthorizationHeaderKey)

			// Parse the authorization to a string
			authorizationStr, ok := authorization.(string)
			if !ok {
				m.handler.HandleResponse(
					w,
					response.NewErrorResponseWithCode(
						gonethttp.ErrInvalidAuthorizationHeader,
						http.StatusUnauthorized,
					),
				)
				return
			}

			// Check if the authorization is a bearer token
			parts := strings.Split(authorizationStr, " ")

			// Return an error if the authorization is missing or invalid
			if len(parts) < 2 || parts[0] != gojwt.BearerPrefix {
				m.handler.HandleResponse(
					w,
					response.NewErrorResponseWithCode(
						gonethttp.ErrInvalidAuthorizationHeader,
						http.StatusUnauthorized,
					),
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
				m.jwtValidatorErrorHandler(w, err)
				return
			}

			// Set the raw token and token claims to the context
			gojwtnethttpctx.SetCtxRawToken(r, &rawToken)
			gojwtnethttpctx.SetCtxTokenClaims(r, claims)

			// Continue
			next.ServeHTTP(w, r)
		},
	)
}
