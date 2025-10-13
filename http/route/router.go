package route

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// Router is the route group struct
	Router struct {
		middlewares  []func(http.Handler) http.Handler
		firstHandler http.Handler
		mux          *http.ServeMux
		pattern      string
		relativePath string
		fullPath     string
		method       string
		mode         *goflagsmode.Flag
		logger       *slog.Logger
	}
)

// SplitPattern returns the method and the path from the pattern
//
// Parameters:
//
//   - pattern: The pattern to split
//
// Returns:
//
//   - string: The method
//   - string: The path
//   - error: The error if any
func SplitPattern(pattern string) (string, string, error) {
	// Trim the pattern
	strings.Trim(pattern, " ")

	// Check if the pattern is empty
	if pattern == "" {
		return "", "", ErrEmptyPattern
	}

	// Iterate over the pattern
	var method string
	path := pattern
	for i, char := range pattern {
		// Split the pattern by the first space
		if char == ' ' {
			// Get the method and the path
			method = pattern[:i]
			path = pattern[i+1:]
			break
		}
	}

	// Trim the path
	strings.Trim(path, " ")

	return method, path, nil
}

// NewRouter creates a new router
//
// Parameters:
//
//   - pattern: The pattern of the router
//   - mode: The flag mode
//   - logger: The logger
//   - middlewares: The middlewares to apply to the router
//
// Returns:
//
//   - RouterWrapper: The router
func NewRouter(
	pattern string,
	mode *goflagsmode.Flag,
	logger *slog.Logger,
	middlewares ...func(next http.Handler) http.Handler,
) (RouterWrapper, error) {
	// Split the method and path from the pattern
	method, path, err := SplitPattern(pattern)
	if err != nil {
		return nil, err
	}

	// Initialize the multiplexer
	mux := http.NewServeMux()

	// Check if there is a nil middleware
	for i, middleware := range middlewares {
		if middleware == nil {
			return nil, fmt.Errorf(ErrNilMiddleware, path, i)
		}
	}

	// Chain the handlers
	firstHandler := ChainHandlers(mux, middlewares...)

	if logger != nil {
		logger = logger.With(
			slog.String("component", "router"),
			slog.String("path", path),
		)
	}

	return &Router{
		middlewares,
		firstHandler,
		mux,
		pattern,
		path,
		path,
		method,
		mode,
		logger,
	}, nil
}

// NewBaseRouter creates a new base router
//
// Parameters:
//
//   - mode: The flag mode
//   - logger: The logger
//   - middlewares: The middlewares to apply to the router
//
// Returns:
//
//   - RouterWrapper: The router
//   - error: The error if any
func NewBaseRouter(
	mode *goflagsmode.Flag,
	logger *slog.Logger,
	middlewares ...func(next http.Handler) http.Handler,
) (RouterWrapper, error) {
	return NewRouter("/", mode, logger, middlewares...)
}

// NewGroup creates a new router group
//
// Parameters:
//
//   - baseRouter: The base router
//   - pattern: The pattern of the group
//   - middlewares: The middlewares to apply to the group
//
// Returns:
//
//   - RouterWrapper: The router group
//   - error: The error if any
func NewGroup(
	baseRouter RouterWrapper,
	pattern string,
	middlewares ...func(next http.Handler) http.Handler,
) (RouterWrapper, error) {
	// Check if the base router is nil
	if baseRouter == nil {
		return nil, ErrNilRouter
	}

	// Split the method and path from the pattern
	method, relativePath, err := SplitPattern(pattern)
	if err != nil {
		return nil, err
	}

	// Check the base router path
	var fullPath string
	if relativePath == "/" && baseRouter.FullPath() == "/" {
		fullPath = "/"
	} else if baseRouter.FullPath()[len(baseRouter.FullPath())-1] == '/' {
		fullPath = baseRouter.FullPath() + relativePath[1:]
	} else {
		fullPath = baseRouter.FullPath() + relativePath
	}

	// Initialize the multiplexer
	mux := http.NewServeMux()

	// Check if there is a nil middleware
	for i, middleware := range middlewares {
		if middleware == nil {
			return nil, fmt.Errorf(ErrNilMiddleware, fullPath, i)
		}
	}

	// Chain the handlers
	firstHandler := ChainHandlers(mux, middlewares...)

	// Create a new router
	instance := &Router{
		middlewares:  middlewares,
		firstHandler: firstHandler,
		mux:          mux,
		logger:       baseRouter.Logger(),
		pattern:      pattern,
		relativePath: relativePath,
		fullPath:     fullPath,
		method:       method,
		mode:         baseRouter.Mode(),
	}

	// Register the group
	baseRouter.RegisterGroup(instance)

	return instance, nil
}

// Handler returns the first handler
//
// Returns:
//
//   - http.Handler: The first handler
func (r *Router) Handler() http.Handler {
	if r == nil {
		return nil
	}
	return r.firstHandler
}

// Mux returns the multiplexer
//
// Returns:
//
//   - *http.ServeMux: The multiplexer
func (r *Router) Mux() *http.ServeMux {
	if r == nil {
		return nil
	}
	return r.mux
}

// GetMiddlewares returns the middlewares
//
// Returns:
//
//   - []func(http.Handler) http.Handler: The middlewares
func (r *Router) GetMiddlewares() []func(http.Handler) http.Handler {
	if r == nil {
		return nil
	}
	return r.middlewares
}

// HandleFunc registers a new route with a path, the handler function and the middlewares
//
// Parameters:
//
//   - pattern: The pattern of the route
//   - handler: The handler function
//   - middlewares: The middlewares to apply to the route
func (r *Router) HandleFunc(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	if r == nil {
		return
	}

	// Split the method and path from the pattern
	method, path, err := SplitPattern(pattern)
	if err != nil {
		panic(err)
	}

	// Get the wildcards from the path
	parsedPath, wildcards := GetWildcards(path)

	// Add the SetCtxWildcardsMiddleware to the beginning of the middlewares
	AddHandlersToStart(
		&middlewares,
		SetCtxWildcardsMiddleware(wildcards),
		SetCtxQueryParametersMiddleware,
	)

	// Chain the handlers
	firstHandler := ChainHandlers(handler, middlewares...)

	// Create the final pattern
	pattern = method + " " + parsedPath

	// Register the route
	r.mux.HandleFunc(pattern, firstHandler.ServeHTTP)

	if r.mode != nil && r.mode.IsDebug() {
		RegisterRoute(r.fullPath, pattern, r.logger)
	}
}

// ExactHandleFunc registers a new route with a path, the handler function and the middlewares
//
// Parameters:
//
//   - pattern: The pattern of the route
//   - handler: The handler function
//   - middlewares: The middlewares to apply to the route
func (r *Router) ExactHandleFunc(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	if r == nil {
		return
	}

	// Split the method and path from the pattern
	method, path, err := SplitPattern(pattern)
	if err != nil {
		panic(err)
	}

	// Get the wildcards from the path
	parsedPath, wildcards := GetWildcards(path)

	// Add the SetCtxWildcardsMiddleware to the beginning of the middlewares
	AddHandlersToStart(
		&middlewares,
		SetCtxWildcardsMiddleware(wildcards),
		SetCtxQueryParametersMiddleware,
	)

	// Chain the handlers
	firstHandler := ChainHandlers(handler, middlewares...)

	// Add the '$' wildcard to the end of the path to match the exact path
	if parsedPath[len(parsedPath)-1] == '/' {
		parsedPath += "{$}"
	}

	// Create the final pattern
	pattern = method + " " + parsedPath

	// Register the route
	r.mux.HandleFunc(pattern, firstHandler.ServeHTTP)

	if r.mode != nil && r.mode.IsDebug() {
		RegisterRoute(r.fullPath, pattern, r.logger)
	}
}

// RegisterRoute registers a new route with a path, the handler function and the middlewares.
// This does not match the exact path
//
// Parameters:
//
//   - pattern: The pattern of the route
//   - handler: The handler function
//   - middlewares: The middlewares to apply to the route
func (r *Router) RegisterRoute(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	if r == nil {
		return
	}
	r.HandleFunc(pattern, handler, middlewares...)
}

// RegisterExactRoute registers a new route with a path, the handler function and the middlewares.
// This matches the exact path
//
// Parameters:
//
//   - pattern: The pattern of the route
//   - handler: The handler function
//   - middlewares: The middlewares to apply to the route
func (r *Router) RegisterExactRoute(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	if r == nil {
		return
	}
	r.ExactHandleFunc(pattern, handler, middlewares...)
}

// RegisterHandler registers a new route group with a path and a handler function
//
// Parameters:
//
//   - pattern: The pattern of the route group
//   - handler: The handler function
func (r *Router) RegisterHandler(pattern string, handler http.Handler) {
	if r == nil {
		return
	}

	// Check if the pattern contains a trailing slash and remove it
	if len(pattern) > 0 && pattern[len(pattern)-1] == '/' {
		pattern = pattern[:len(pattern)-1]
	}

	// Register the route group
	r.mux.Handle(pattern+"/", http.StripPrefix(pattern, handler))

	if r.logger != nil && r.mode != nil && r.mode.IsDebug() {
		r.logger.Debug(
			"registering route group",
			slog.String("full_path", r.fullPath),
			slog.String("relative_path", r.relativePath),
			slog.String("pattern", pattern),
		)
	}
}

// RegisterGroup registers a new router group with a path and a router
//
// Parameters:
//
//   - router: The router group
func (r *Router) RegisterGroup(router RouterWrapper) {
	if r == nil {
		return
	}
	r.RegisterHandler(router.Pattern(), router.Handler())
}

// NewGroup creates a new router group with a path
//
// Parameters:
//
//   - pattern: The pattern of the group
//   - middlewares: The middlewares to apply to the group
//
// Returns:
//
//   - RouterWrapper: The router group
func (r *Router) NewGroup(
	pattern string,
	middlewares ...func(next http.Handler) http.Handler,
) RouterWrapper {
	if r == nil {
		return nil
	}

	// Create a new group
	newGroup, _ := NewGroup(r, pattern, middlewares...)
	return newGroup
}

// Pattern returns the pattern
//
// Returns:
//
//   - string: The pattern
func (r *Router) Pattern() string {
	if r == nil {
		return ""
	}
	return r.pattern
}

// RelativePath returns the relative path
//
// Returns:
//
//   - string: The relative path
func (r *Router) RelativePath() string {
	if r == nil {
		return ""
	}
	return r.relativePath
}

// FullPath returns the full path
//
// Returns:
//
//   - string: The full path
func (r *Router) FullPath() string {
	if r == nil {
		return ""
	}
	return r.fullPath
}

// Method returns the method
//
// Returns:
//
//   - string: The method
func (r *Router) Method() string {
	if r == nil {
		return ""
	}
	return r.method
}

// ServeStaticFiles serves the static files
//
// Parameters:
//
//   - pattern: The pattern to serve the static files
//   - path: The path to the static files
func (r *Router) ServeStaticFiles(
	pattern,
	path string,
) {
	if r == nil {
		return
	}

	// Check if the pattern contains a trailing slash and add it
	if len(pattern) == 0 || pattern[len(pattern)-1] != '/' {
		pattern = pattern + "/"
	}

	// Serve the static files
	r.mux.HandleFunc(
		pattern,
		http.StripPrefix(pattern, http.FileServer(http.Dir(path))).ServeHTTP,
	)
}

// Logger returns the logger
//
// Returns:
//
//   - *slog.Logger: The logger
func (r *Router) Logger() *slog.Logger {
	if r == nil {
		return nil
	}
	return r.logger
}

// Mode returns the mode
//
// Returns:
//
//   - *goflagsmode.Flag: The mode
func (r *Router) Mode() *goflagsmode.Flag {
	if r == nil {
		return nil
	}
	return r.mode
}
