package auth

import (
	"errors"
	"net/http"
	"strings"

	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtnethttp "github.com/ralvarezdev/go-jwt/net/http"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwtvalidator "github.com/ralvarezdev/go-jwt/token/validator"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
)

type (
	// Middleware struct is the authentication middleware
	Middleware struct {
		validator        gojwtvalidator.Validator
		responsesHandler gonethttphandler.ResponsesHandler
		options          *Options
	}

	// Options is the options for the authentication middleware
	Options struct {
		CookieRefreshTokenName *string
		CookieAccessTokenName  *string
		RefreshTokenFn         RefreshTokenFn
	}
)

// NewOptions creates a new Options struct
//
// Parameters:
//
//   - cookieRefreshTokenName: The name of the cookie that contains the refresh token
//   - cookieAccessTokenName: The name of the cookie that contains the access token
//   - refreshTokenFn: The function to refresh the access token using the refresh token
//
// Returns:
//
//   - *Options: The options for the authentication middleware
func NewOptions(
	cookieRefreshTokenName,
	cookieAccessTokenName *string,
	refreshTokenFn RefreshTokenFn,
) *Options {
	return &Options{
		cookieRefreshTokenName,
		cookieAccessTokenName,
		refreshTokenFn,
	}
}

// NewMiddleware creates a new authentication middleware
//
// Parameters:
//
//   - responsesHandler: The HTTP handler to handle errors
//   - validator: The JWT validator service (if nil, no validation will be done, can be used for gRPC gateways)
//   - options: The options for the authentication middleware (can be nil)
//
// Returns:
//
//   - *Middleware: The authentication middleware
func NewMiddleware(
	responsesHandler gonethttphandler.ResponsesHandler,
	validator gojwtvalidator.Validator,
	options *Options,
) (*Middleware, error) {
	// Check if either the responses handler is nil or the validator is nil
	if responsesHandler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}
	if validator == nil {
		return nil, gojwtvalidator.ErrNilValidator
	}

	return &Middleware{
		validator,
		responsesHandler,
		options,
	}, nil
}

// authenticate return the middleware function that authenticates the request
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
func (m Middleware) authenticate(
	token gojwttoken.Token,
	rawToken string,
	failHandler FailHandlerFn,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Get the context from the request
				ctx := r.Context()

				// Validate the token claims
				claims, err := m.validator.ValidateClaims(
					ctx,
					rawToken,
					token,
				)
				if err != nil {
					if failHandler == nil {
						panic(err)
					}
					failHandler(
						w,
						r,
						err,
						ErrCodeInvalidTokenClaims,
					)
					return
				}

				// Set the token claims to the context
				r = gojwtnethttp.SetCtxTokenClaims(r, claims)

				// Set the raw token to the context
				r, err = gojwtnethttp.SetCtxToken(r, rawToken)
				if err != nil {
					if failHandler == nil {
						panic(err)
					}
					failHandler(
						w,
						r,
						err,
						ErrCodeFailedToSetTokenInContext.Error(),
					)
					return
				}

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
//   - r: The HTTP request
//   - err: The error that occurred
//   - errorCode: The error code to return
func (m Middleware) authenticateFromHeaderFailHandler(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	errorCode string,
) {
	m.responsesHandler.HandleFailFieldErrorWithCode(
		w,
		r,
		gonethttp.Authorization,
		err,
		errorCode,
		http.StatusUnauthorized,
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
				authorization := r.Header.Get(gonethttp.Authorization)

				// Check if the authorization is a bearer token
				parts := strings.Split(authorization, " ")

				// Return an error if the authorization is missing or invalid
				if len(parts) < 2 || parts[0] != gojwt.BearerPrefix {
					m.authenticateFromHeaderFailHandler(
						w,
						r,
						ErrInvalidAuthorizationHeader,
						ErrCodeInvalidAuthorizationHeader,
					)
					return
				}

				// Get the raw token from the header
				rawToken := parts[1]

				// Call the authenticate function
				m.authenticate(
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
//
// Returns:
//
//   - FailHandlerFn: The fail handler function
func (m Middleware) authenticateFromCookieFailHandler(
	cookieName string,
) FailHandlerFn {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		err error,
		errorCode string,
	) {
		m.responsesHandler.HandleFailFieldErrorWithCode(
			w,
			r,
			cookieName,
			err,
			errorCode,
			http.StatusUnauthorized,
		)
	}
}

// AuthenticateFromCookie return the middleware function that authenticates the request from the cookie
//
// Parameters:
//
//   - token: The type of token to authenticate (access or refresh)
//
// Returns:
//
//   - func(next http.Handler) http.Handler: The middleware function
func (m Middleware) AuthenticateFromCookie(
	token gojwttoken.Token,
) func(next http.Handler) http.Handler {
	// Validate that all the required options are set
	if m.options == nil {
		panic(ErrNilOptions)
	}
	if m.options.CookieAccessTokenName == nil {
		panic(ErrNilCookieAccessTokenName)
	}
	if m.options.CookieRefreshTokenName == nil {
		panic(ErrNilCookieRefreshTokenName)
	}
	if m.options.RefreshTokenFn == nil {
		panic(ErrNilRefreshTokenFn)
	}

	// Determine the cookie name based on the token type
	refreshTokenCookieName := *m.options.CookieRefreshTokenName
	accessTokenCookieName := *m.options.CookieAccessTokenName
	var currentCookieName string
	switch token {
	case gojwttoken.AccessToken:
		currentCookieName = accessTokenCookieName
	case gojwttoken.RefreshToken:
		currentCookieName = refreshTokenCookieName
	}

	// Create the fail handler functions
	accessTokenFailHandler := m.authenticateFromCookieFailHandler(accessTokenCookieName)
	refreshTokenFailHandler := m.authenticateFromCookieFailHandler(refreshTokenCookieName)
	var currentFailHandler FailHandlerFn
	switch token {
	case gojwttoken.AccessToken:
		currentFailHandler = accessTokenFailHandler
	case gojwttoken.RefreshToken:
		currentFailHandler = refreshTokenFailHandler
	}

	//nolint:nestif // This function requires nested ifs for clarity
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var (
					rawToken string
					cookie   *http.Cookie
					err      error
					ok       bool
				)

				// Get the cookie from the request
				cookie, err = r.Cookie(currentCookieName)

				// Check if there was an error getting the cookie
				if err == nil {
					// Get the raw token from the cookie
					rawToken = cookie.Value
				} else if errors.Is(err, http.ErrNoCookie) {
					// Check if there's no option to refresh the token or if the token is a refresh token
					if token == gojwttoken.RefreshToken || m.options.RefreshTokenFn == nil {
						currentFailHandler(
							w,
							r,
							gonethttp.ErrCookieNotFound,
							gonethttp.ErrCodeCookieNotFound,
						)
						return
					}

					// Check if the refresh token cookie is present
					if _, err = r.Cookie(refreshTokenCookieName); err != nil {
						currentFailHandler(
							w,
							r,
							gonethttp.ErrCookieNotFound,
							gonethttp.ErrCodeCookieNotFound,
						)
						return
					}

					// Refresh the token
					rawTokens, refreshErr := m.options.RefreshTokenFn(w, r)
					if refreshErr != nil {
						refreshTokenFailHandler(
							w,
							r,
							refreshErr,
							ErrCodeFailedToRefreshToken,
						)
						return
					}

					// Get the raw token from the map, if not found set it to an empty string
					if rawToken, ok = rawTokens[token]; !ok {
						rawToken = ""
					}
				}

				// Check if the raw token is empty
				if rawToken == "" {
					currentFailHandler(
						w,
						r,
						gonethttp.ErrCookieNotFound,
						gonethttp.ErrCodeCookieNotFound,
					)
					return
				}

				// Call the authenticate function
				m.authenticate(
					token,
					rawToken,
					currentFailHandler,
				)(next).ServeHTTP(
					w,
					r,
				)
			},
		)
	}
}
