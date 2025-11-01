package grpc

import (
	"log/slog"
	"net/http"

	gojwttoken "github.com/ralvarezdev/go-jwt/token"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttpmiddlewareauth "github.com/ralvarezdev/go-net/http/middleware/auth"
)

type (
	// Middleware struct is the authentication middleware for a REST API that is a gateway to a gRPC service.
	Middleware struct {
		responsesHandler gonethttphandler.ResponsesHandler
		interceptions    map[string]*gojwttoken.Token
		authenticator    gonethttpmiddlewareauth.Authenticator
		logger           *slog.Logger
	}
)

// NewMiddleware creates a new instance of the authentication middleware.
//
// Parameters:
//
//   - interceptions: A map of string keys to JWT tokens for interception purposes.
//   - responsesHandler: An instance of a ResponsesHandler to manage HTTP responses.
//   - authenticator: An instance of an Authenticator to handle authentication logic.
//   - logger: A pointer to a slog.Logger for logging purposes.
//
// Returns:
//
//   - *Middleware: A pointer to the newly created Middleware instance.
//   - error: An error if the middleware could not be created.
func NewMiddleware(
	interceptions map[string]*gojwttoken.Token,
	responsesHandler gonethttphandler.ResponsesHandler,
	authenticator gonethttpmiddlewareauth.Authenticator,
	logger *slog.Logger,
) (*Middleware, error) {
	// Check if the interceptions are nil
	if interceptions == nil {
		return nil, ErrNilInterceptions
	}

	// Check if the handler is nil
	if responsesHandler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}

	// Check if the authenticator is nil
	if authenticator == nil {
		return nil, gonethttpmiddlewareauth.ErrNilAuthenticator
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "http_middleware_auth_grpc"),
		)
	}

	return &Middleware{
		interceptions:    interceptions,
		responsesHandler: responsesHandler,
		authenticator:    authenticator,
		logger:           logger,
	}, nil
}

// interceptionNotFoundHandler is a default handler for cases when no interception is found for a given RPC method.
//
// Parameters:
//
//   - rpcMethod: The RPC method for which the interception was not found.
//
// Returns:
//
//   - func(next http.Handler) http.Handler: A middleware function that logs a warning and calls the fail handler.
func (m Middleware) interceptionNotFoundHandler(
	rpcMethod string,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Log a warning if the interception was not found
				if m.logger != nil {
					m.logger.Warn(
						"Interception not found for RPC method",
						slog.String("rpc_method", rpcMethod),
					)
				}

				// Handle the response with an internal server error
				m.responsesHandler.HandleDebugErrorWithCode(
					w,
					r,
					ErrInterceptionNotFound,
					gonethttp.ErrInternalServerError,
					ErrCodeInterceptionNotFound,
					http.StatusInternalServerError,
				)
			},
		)
	}
}

// AuthenticateFromHeader is a middleware function that authenticates requests based on the provided RPC method using
// tokens from the request header.
//
// Parameters:
//
//   - rpcMethod: The RPC method to authenticate against.
//
// Returns:
//
//   - func(next http.Handler) http.Handler: A middleware function that authenticates requests.
func (m Middleware) AuthenticateFromHeader(
	rpcMethod string,
) func(next http.Handler) http.Handler {
	// Try to find the interception for the given RPC method
	token, ok := m.interceptions[rpcMethod]
	if !ok {
		return m.interceptionNotFoundHandler(
			rpcMethod,
		)
	}

	// Check if the authentication is needed
	if token != nil {
		return m.authenticator.AuthenticateFromHeader(
			*token,
		)
	}

	// If no authentication is needed, return a no-op middleware
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			},
		)
	}
}

// AuthenticateFromCookie is a middleware function that authenticates requests based on the provided RPC method using
// tokens from cookies.
//
// Parameters:
//
//   - rpcMethod: The RPC method to authenticate against.
//
// Returns:
//
//   - func(next http.Handler) http.Handler: A middleware function that authenticates requests.
func (m Middleware) AuthenticateFromCookie(
	rpcMethod string,
) func(next http.Handler) http.Handler {
	// Try to find the interception for the given RPC method
	token, ok := m.interceptions[rpcMethod]
	if !ok {
		return m.interceptionNotFoundHandler(
			rpcMethod,
		)
	}

	// Check if the authentication is needed
	if token != nil {
		return m.authenticator.AuthenticateFromCookie(
			*token,
		)
	}

	// If no authentication is needed, return a no-op middleware
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			},
		)
	}
}
