package auth

import (
	"errors"
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

// Middleware struct is the authentication middleware
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
	token gojwttoken.Token,
	rawToken string,
	failHandler func(
		w http.ResponseWriter,
		err error,
		errorCode *string,
	),
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
						err,
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
		err error,
		errorCode *string,
	) {
		m.handler.HandleError(
			w,
			gonethttpresponse.NewFailResponseError(
				gojwtnethttp.AuthorizationHeaderKey,
				err.Error(),
				errorCode,
				http.StatusUnauthorized,
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
						ErrInvalidAuthorizationHeader,
						ErrCodeInvalidAuthorizationHeader,
					)
					return
				}

				// Get the raw token from the header
				rawToken := parts[1]

				// Call the Authenticate function
				m.Authenticate(
					token,
					rawToken,
					failHandler,
				)(next).ServeHTTP(
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
	cookieRefreshTokenName,
	cookieAccessTokenName string,
	refreshTokenFn func(
		w http.ResponseWriter,
		r *http.Request,
	) (*map[gojwttoken.Token]string, error),
) func(next http.Handler) http.Handler {
	var cookieName string
	if token == gojwttoken.AccessToken {
		cookieName = cookieAccessTokenName
	} else if token == gojwttoken.RefreshToken {
		cookieName = cookieRefreshTokenName
	}

	// Create the fail handler function
	failHandler := func(
		cookieName string,
	) func(
		w http.ResponseWriter,
		err error,
		errorCode *string,
	) {
		return func(
			w http.ResponseWriter,
			err error,
			errorCode *string,
		) {
			m.handler.HandleError(
				w,
				gonethttpresponse.NewFailResponseError(
					cookieName,
					err.Error(),
					errorCode,
					http.StatusUnauthorized,
				),
			)
		}
	}

	// Create the authenticate function
	var authenticateFn func(*map[gojwttoken.Token]string) func(next http.Handler) http.Handler
	authenticateFn = func(rawTokens *map[gojwttoken.Token]string) func(next http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					var rawToken string
					var cookie *http.Cookie
					var err error
					var ok bool

					// Get the cookie
					if rawTokens != nil {
						// Get the raw token from the map
						rawToken, ok = (*rawTokens)[token]

						// Return an error if the token is missing
						if !ok {
							failHandler(cookieName)(
								w,
								gonethttp.ErrCookieNotFound,
								gonethttp.ErrCodeCookieNotFound,
							)
							return
						}
					} else {
						// Get the cookie from the request
						cookie, err = r.Cookie(cookieAccessTokenName)

						// Check if there was an error getting the cookie
						if err == nil {
							// Get the raw token from the cookie
							rawToken = cookie.Value
						} else if errors.Is(err, http.ErrNoCookie) {
							// Check if the token can be refreshed
							if token == gojwttoken.AccessToken && refreshTokenFn != nil {
								// Refresh the token
								rawTokens, err = refreshTokenFn(w, r)
								if err != nil {
									failHandler(cookieRefreshTokenName)(
										w,
										err,
										ErrCodeFailedToRefreshToken,
									)
									return
								}

								// Authenticate again
								authenticateFn(rawTokens)(next).ServeHTTP(
									w,
									r,
								)
								return
							}
						}
					}

					// Check if the raw token is empty
					if rawToken == "" {
						failHandler(cookieAccessTokenName)(
							w,
							gonethttp.ErrCookieNotFound,
							gonethttp.ErrCodeCookieNotFound,
						)
						return
					}

					// Call the Authenticate function
					m.Authenticate(
						token,
						rawToken,
						failHandler(cookieName),
					)(next).ServeHTTP(
						w,
						r,
					)
				},
			)
		}
	}

	return authenticateFn(nil)
}
