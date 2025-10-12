package auth

import (
	"errors"
	"net/http"
	"strings"

	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtnethttp "github.com/ralvarezdev/go-jwt/net/http"
	gojwtnethttpctx "github.com/ralvarezdev/go-jwt/net/http/context"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwtvalidator "github.com/ralvarezdev/go-jwt/token/validator"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// Middleware struct is the authentication middleware
	Middleware struct {
		validator gojwtvalidator.Validator
		handler   gonethttphandler.Handler
	}
)

// NewMiddleware creates a new authentication middleware
//
// Parameters:
//
//   - handler: The HTTP handler to handle errors
//   - validator: The JWT validator service (if nil, no validation will be done, can be used for gRPC gateways)
//
// Returns:
//
//   - *Middleware: The authentication middleware
func NewMiddleware(
	handler gonethttphandler.Handler,
	validator gojwtvalidator.Validator,
) (*Middleware, error) {
	// Check if either the response handler is nil
	if handler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}

	return &Middleware{
		validator,
		handler,
	}, nil
}

// Authenticate return the middleware function that authenticates the request
//
// Parameters:
//
//   - token: The type of token to authenticate (access or refresh)
//   - rawToken: The raw JWT token string
//   - failHandler: The function to handle authentication failures
//
// Returns:
//
//   - func(next http.Handler) http.Handler: The middleware function
func (m Middleware) Authenticate(
	token gojwttoken.Token,
	rawToken string,
	failHandler FailHandlerFn,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Validate the token and get the validated claims
				if m.validator == nil {
					claims, err := m.validator.ValidateClaims(
						rawToken,
						token,
					)
					if err != nil {
						if failHandler == nil {
							panic(err)
						}
						failHandler(
							w,
							err,
							ErrCodeInvalidTokenClaims,
						)
						return
					}

					// Set the token claims to the context
					r = gojwtnethttpctx.SetCtxTokenClaims(r, claims)
				}

				// Set the raw token to the context
				r, _ = gojwtnethttpctx.SetCtxToken(r, rawToken)

				// Call the next handler
				next.ServeHTTP(w, r)
			},
		)
	}
}

// authenticateFromHeaderFailHandler is the default fail handler for AuthenticateFromHeader
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error that occurred
//   - errorCode: The error code to return
func (m Middleware) authenticateFromHeaderFailHandler(
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

// AuthenticateFromHeader return the middleware function that authenticates the request from the header
//
// Parameters:
//
//   - token: The type of token to authenticate (access or refresh)
//
// Returns:
//
//   - func(next http.Handler) http.Handler: The middleware function
func (m Middleware) AuthenticateFromHeader(
	token gojwttoken.Token,
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
					m.authenticateFromHeaderFailHandler(
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
					m.authenticateFromHeaderFailHandler,
				)(next).ServeHTTP(
					w,
					r,
				)
			},
		)
	}
}

// authenticateFromCookieFailHandler is the default fail handler for AuthenticateFromCookie
//
// Parameters:
//
//   - cookieName: The name of the cookie that contains the token
func (m Middleware) authenticateFromCookieFailHandler(
	cookieName string,
) FailHandlerFn {
	return func(
		w http.ResponseWriter,
		err error,
		errorCode *string,
	) {
		{
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
}

// AuthenticateFromCookie return the middleware function that authenticates the request from the cookie
//
// Parameters:
//
//   - token: The type of token to authenticate (access or refresh)
//   - cookieRefreshTokenName: The name of the cookie that contains the refresh token
//   - cookieAccessTokenName: The name of the cookie that contains the access token
//   - refreshTokenFn: The function to refresh the access token using the refresh token
//
// Returns:
//
//   - func(next http.Handler) http.Handler: The middleware function
func (m Middleware) AuthenticateFromCookie(
	token gojwttoken.Token,
	cookieRefreshTokenName,
	cookieAccessTokenName string,
	refreshTokenFn RefreshTokenFn,
) func(next http.Handler) http.Handler {
	var cookieName string
	if token == gojwttoken.AccessToken {
		cookieName = cookieAccessTokenName
	} else if token == gojwttoken.RefreshToken {
		cookieName = cookieRefreshTokenName
	}

	// Create the fail handler function
	failHandler := m.authenticateFromCookieFailHandler(cookieName)

	// Create the authenticate function
	var authenticateFn func(rawTokens map[gojwttoken.Token]string) func(next http.Handler) http.Handler
	authenticateFn = func(rawTokens map[gojwttoken.Token]string) func(next http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					var (
						rawToken string
						cookie   *http.Cookie
						err      error
						ok       bool
					)

					// Get the cookie
					if rawTokens != nil {
						// Get the raw token from the map
						rawToken, ok = rawTokens[token]

						// Return an error if the token is missing
						if !ok {
							failHandler(
								w,
								gonethttp.ErrCookieNotFound,
								gonethttp.ErrCodeCookieNotFound,
							)
							return
						}
					} else {
						// Get the cookie from the request
						cookie, err = r.Cookie(cookieName)

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
									m.authenticateFromCookieFailHandler(cookieRefreshTokenName)(
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
						failHandler(
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
						failHandler,
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
