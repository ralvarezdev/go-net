package auth

import (
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtnethttp "github.com/ralvarezdev/go-jwt/net/http"
	gojwtnethttpctx "github.com/ralvarezdev/go-jwt/net/http/context"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwtvalidator "github.com/ralvarezdev/go-jwt/token/validator"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
	"strings"
)

// Middleware struct
type Middleware struct {
	validator gojwtvalidator.Validator
	handler   gonethttphandler.Handler
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(
	validator gojwtvalidator.Validator,
	handler gonethttphandler.Handler,
) (*Middleware, error) {
	// Check if either the validator, response handler or validator handler is nil
	if validator == nil {
		return nil, gojwtvalidator.ErrNilValidator
	}
	if handler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}

	return &Middleware{
		validator,
		handler,
	}, nil
}

// Authenticate return the middleware function that authenticates the request
func (m *Middleware) Authenticate(
	failHandler func(
		w http.ResponseWriter,
		err string,
		httpStatus int,
		errorCode *string,
	),
	token gojwttoken.Token,
	rawToken string,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Validate the token and get the validated claims
				claims, err := m.validator.ValidateClaims(
					rawToken,
					token,
				)
				if err != nil {
					failHandler(
						w,
						err.Error(),
						http.StatusUnauthorized,
						ErrCodeInvalidTokenClaims,
					)
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

// AuthenticateFromHeader return the middleware function that authenticates the request from the header
func (m *Middleware) AuthenticateFromHeader(
	token gojwttoken.Token,
) func(next http.Handler) http.Handler {
	// Create the fail handler function
	failHandler := func(
		w http.ResponseWriter,
		err string,
		httpStatus int,
		errorCode *string,
	) {
		m.handler.HandleError(
			w,
			gonethttpresponse.NewHeaderError(
				gojwtnethttp.AuthorizationHeaderKey,
				err,
				httpStatus,
				errorCode,
			),
		)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Get the authorization from the header
				authorization := r.Header.Get(gojwtnethttp.AuthorizationHeaderKey)

				// Check if the authorization is a bearer token
				parts := strings.Split(authorization, " ")

				// Return an error if the authorization is missing or invalid
				if len(parts) < 2 || parts[0] != gojwt.BearerPrefix {
					failHandler(
						w,
						ErrInvalidAuthorizationHeader.Error(),
						http.StatusUnauthorized,
						ErrCodeInvalidAuthorizationHeader,
					)
					return
				}

				// Get the raw token from the header
				rawToken := parts[1]

				// Call the Authenticate function
				m.Authenticate(failHandler, token, rawToken)(next).ServeHTTP(
					w,
					r,
				)
			},
		)
	}
}

// AuthenticateFromCookie return the middleware function that authenticates the request from the cookie
func (m *Middleware) AuthenticateFromCookie(
	token gojwttoken.Token,
	cookieName string,
) func(next http.Handler) http.Handler {
	// Create the fail handler function
	failHandler := func(
		w http.ResponseWriter,
		err string,
		httpStatus int,
		errorCode *string,
	) {
		m.handler.HandleError(
			w,
			gonethttpresponse.NewCookieError(
				cookieName,
				err,
				httpStatus,
				errorCode,
			),
		)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Get the cookie
				cookie, err := r.Cookie(cookieName)

				// Return an error if the cookie is missing
				if err != nil {
					failHandler(
						w,
						gonethttp.ErrCookieNotFound.Error(),
						http.StatusUnauthorized,
						gonethttp.ErrCodeCookieNotFound,
					)
					return
				}

				// Get the raw token from the cookie
				rawToken := cookie.Value

				// Call the Authenticate function
				m.Authenticate(failHandler, token, rawToken)(next).ServeHTTP(
					w,
					r,
				)
			},
		)
	}
}
