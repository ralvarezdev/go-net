package route

import (
	"fmt"
	"log/slog"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"

	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
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
		handler      gonethttphandler.Handler
		mode         *goflagsmode.Flag
		logger       *slog.Logger
	}
)

// NewRouter creates a new router
//
// Parameters:
//
//   - pattern: The pattern of the router
//   - mode: The flag mode
//   - handler: The handler to handle the errors
//   - logger: The logger
//   - middlewares: The middlewares to apply to the router
//
// Returns:
//
//   - RouterWrapper: The router
func NewRouter(
	pattern string,
	mode *goflagsmode.Flag,
	handler gonethttphandler.Handler,
	logger *slog.Logger,
	middlewares ...func(next http.Handler) http.Handler,
) (RouterWrapper, error) {
	// Split the method and path from the pattern
	method, path, err := SplitPattern(pattern)
	if err != nil {
		return nil, err
	}

	// Check if the handler is nil
	if handler == nil {
		return nil, gonethttphandler.ErrNilHandler
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
		handler,
		mode,
		logger,
	}, nil
}

// NewBaseRouter creates a new base router
//
// Parameters:
//
//   - mode: The flag mode
//   - handler: The handler to handle the errors
//   - logger: The logger
//   - middlewares: The middlewares to apply to the router
//
// Returns:
//
//   - RouterWrapper: The router
//   - error: The error if any
func NewBaseRouter(
	mode *goflagsmode.Flag,
	handler gonethttphandler.Handler,
	logger *slog.Logger,
	middlewares ...func(next http.Handler) http.Handler,
) (RouterWrapper, error) {
	return NewRouter("/", mode, handler, logger, middlewares...)
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

// chainMiddlewares chains the middlewares to the handler and adds the SetCtxWildcardsMiddleware and
// SetCtxQueryParametersMiddleware
//
// Parameters:
//
//   - pattern: The pattern of the route
//   - exact: Whether the route should match the exact path
//   - handler: The handler to chain the middlewares to
//   - middlewares: The middlewares to add
//
// Returns:
//
//   - string: The final pattern
//   - http.Handler: The chained handler
func (r *Router) chainMiddlewares(
	pattern string,
	exact bool,
	handler http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) (string, http.Handler) {
	if r == nil {
		return "", nil
	}

	// Split the method and path from the pattern
	method, path, err := SplitPattern(pattern)
	if err != nil {
		panic(err)
	}

	// Get the wildcards from the path
	parsedPath, wildcards, err := GetWildcards(path)
	if err != nil {
		panic(err)
	}

	// If the path is exact, add the '$' wildcard to the end of the path to match the exact path
	if exact && parsedPath[len(parsedPath)-1] == '/' {
		parsedPath += "{$}"
	}

	// Add the SetCtxWildcardsMiddleware to the beginning of the middlewares
	AddHandlersToStart(
		&middlewares,
		SetCtxWildcardsMiddleware(wildcards),
		SetCtxQueryParametersMiddleware,
	)

	// Chain the handlers
	firstHandler := ChainHandlers(handler, middlewares...)

	return method + " " + parsedPath, firstHandler
}

// addHandleFunc registers a new route with a path, the handler function and the middlewares
//
// Parameters:
//
//   - pattern: The pattern of the route
//   - handler: The handler function
//   - exact: Whether the route should match the exact path
//   - middlewares: The middlewares to apply to the route
func (r *Router) addHandleFunc(
	pattern string,
	handler http.HandlerFunc,
	exact bool,
	middlewares ...func(http.Handler) http.Handler,
) {
	if r == nil {
		return
	}

	// Check if the handler is nil
	if handler == nil {
		panic(fmt.Sprintf(ErrNilHandlerFunc, pattern))
	}

	// Chain the middlewares
	pattern, firstHandler := r.chainMiddlewares(
		pattern,
		exact,
		handler,
		middlewares...,
	)

	// Register the route
	r.mux.HandleFunc(pattern, firstHandler.ServeHTTP)

	if r.mode != nil && r.mode.IsDebug() {
		AddRouter(r.fullPath, pattern, r.logger)
	}
}

// AddHandleFunc registers a new route with a path, the handler function and the middlewares
//
// Parameters:
//
//   - pattern: The pattern of the route
//   - handler: The handler function
//   - middlewares: The middlewares to apply to the route
func (r *Router) AddHandleFunc(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	if r == nil {
		return
	}

	// Add the route
	r.addHandleFunc(pattern, handler, false, middlewares...)
}

// AddExactHandleFunc registers a new route with a path, the handler function and the middlewares
//
// Parameters:
//
//   - pattern: The pattern of the route
//   - handler: The handler function
//   - middlewares: The middlewares to apply to the route
func (r *Router) AddExactHandleFunc(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	if r == nil {
		return
	}

	// Add the route
	r.addHandleFunc(pattern, handler, true, middlewares...)
}

// AddEndpointHandler adds a new endpoint with a path, the handler function and the middlewares
//
// Parameters:
//
//   - pattern: The pattern of the endpoint
//   - handler: The handler function
//   - middlewares: The middlewares to apply to the endpoint
func (r *Router) AddEndpointHandler(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request) error,
	middlewares ...func(http.Handler) http.Handler,
) {
	if r == nil {
		return
	}

	// Check if the handler is nil
	if handler == nil {
		panic(fmt.Sprintf(ErrNilEndpointHandler, pattern))
	}

	// Wrap the endpoint handler
	wrappedHandler := http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			if err := handler(w, req); err != nil {
				// Handle the error using the handler's HandleRawError method
				r.handler.HandleRawError(w, req, err)
			}
		},
	)

	// Add the endpoint handler
	r.AddHandleFunc(pattern, wrappedHandler, middlewares...)
}

// AddExactEndpointHandler adds a new endpoint with a path, the handler function and the middlewares
//
// Parameters:
//
//   - pattern: The pattern of the endpoint
//   - handler: The handler function
//   - middlewares: The middlewares to apply to the endpoint
func (r *Router) AddExactEndpointHandler(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request) error,
	middlewares ...func(next http.Handler) http.Handler,
) {
	if r == nil {
		return
	}

	// Check if the handler is nil
	if handler == nil {
		panic(fmt.Sprintf(ErrNilEndpointHandler, pattern))
	}

	// Wrap the endpoint handler
	wrappedHandler := http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			if err := handler(w, req); err != nil {
				// Handle the error using the handler's HandleRawError method
				r.handler.HandleRawError(w, req, err)
			}
		},
	)

	// Add the endpoint handler
	r.AddExactHandleFunc(pattern, wrappedHandler, middlewares...)
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
	if pattern != "" && pattern[len(pattern)-1] == '/' {
		pattern = pattern[:len(pattern)-1]
	}

	// Register the route group
	r.mux.Handle(pattern+"/", http.StripPrefix(pattern, handler))

	if r.logger != nil && r.mode != nil && r.mode.IsDebug() {
		r.logger.Debug(
			"Registering route group",
			slog.String("full_path", r.fullPath),
			slog.String("relative_path", r.relativePath),
			slog.String("pattern", pattern),
		)
	}
}

// NewRouter creates a new router group with a path
//
// Parameters:
//
//   - pattern: The pattern of the group
//   - middlewares: The middlewares to apply to the group
//
// Returns:
//
//   - RouterWrapper: The router group
//   - error: The error if any
func (r *Router) NewRouter(
	pattern string,
	middlewares ...func(next http.Handler) http.Handler,
) (RouterWrapper, error) {
	if r == nil {
		return nil, ErrNilRouter
	}

	// Split the method and path from the pattern
	method, relativePath, err := SplitPattern(pattern)
	if err != nil {
		return nil, err
	}

	// Check the base router path
	var fullPath string
	switch {
	case relativePath == "/" && r.FullPath() == "/":
		fullPath = "/"
	case r.FullPath()[len(r.FullPath())-1] == '/':
		fullPath = r.FullPath() + relativePath[1:]
	default:
		fullPath = r.FullPath() + relativePath
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
		logger:       r.Logger(),
		pattern:      pattern,
		relativePath: relativePath,
		fullPath:     fullPath,
		method:       method,
		mode:         r.Mode(),
	}

	// Add the new router to the parent router
	r.AddRouter(instance)
	return instance, nil
}

// AddRouter registers a new router group with a path and a router
//
// Parameters:
//
//   - router: The router group
func (r *Router) AddRouter(router RouterWrapper) {
	if r == nil {
		return
	}
	r.RegisterHandler(router.Pattern(), router.Handler())
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
	if pattern == "" || pattern[len(pattern)-1] != '/' {
		pattern += "/"
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
