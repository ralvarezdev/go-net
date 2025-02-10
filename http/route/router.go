package route

import (
	"fmt"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	"net/http"
	"strings"
)

type (
	// RouterWrapper is the interface for the routes
	RouterWrapper interface {
		Handler() http.Handler
		Mux() *http.ServeMux
		GetMiddlewares() *[]func(http.Handler) http.Handler
		HandleFunc(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		ExactHandleFunc(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterRoute(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterExactRoute(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterHandler(pattern string, handler http.Handler)
		NewGroup(
			pattern string,
			middlewares ...func(next http.Handler) http.Handler,
		) RouterWrapper
		RegisterGroup(router RouterWrapper)
		Pattern() string
		RelativePath() string
		FullPath() string
		Method() string
		ServeStaticFiles(pattern, path string)
		Logger() *Logger
		Mode() *goflagsmode.Flag
	}

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
		logger       *Logger
	}
)

// SplitPattern returns the method and the path from the pattern
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
func NewRouter(
	pattern string,
	mode *goflagsmode.Flag,
	logger *Logger,
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
func NewBaseRouter(
	mode *goflagsmode.Flag,
	logger *Logger,
	middlewares ...func(next http.Handler) http.Handler,
) (RouterWrapper, error) {
	return NewRouter("/", mode, logger, middlewares...)
}

// NewGroup creates a new router group
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
func (r *Router) Handler() http.Handler {
	return r.firstHandler
}

// Mux returns the multiplexer
func (r *Router) Mux() *http.ServeMux {
	return r.mux
}

// GetMiddlewares returns the middlewares
func (r *Router) GetMiddlewares() *[]func(http.Handler) http.Handler {
	return &r.middlewares
}

// HandleFunc registers a new route with a path, the handler function and the middlewares
func (r *Router) HandleFunc(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	// Chain the handlers
	firstHandler := ChainHandlers(handler, middlewares...)

	// Register the route
	r.mux.HandleFunc(pattern, firstHandler.ServeHTTP)

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRoute(r.fullPath, pattern)
	}
}

// ExactHandleFunc registers a new route with a path, the handler function and the middlewares
func (r *Router) ExactHandleFunc(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	// Split the method and path from the pattern
	method, path, err := SplitPattern(pattern)
	if err != nil {
		panic(err)
	}

	// Chain the handlers
	firstHandler := ChainHandlers(handler, middlewares...)

	// Add the '$' wildcard to the end of the path to match the exact path
	if path[len(path)-1] == '/' {
		path += "{$}"
	}

	// Register the route
	r.mux.HandleFunc(method+" "+path, firstHandler.ServeHTTP)

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRoute(r.fullPath, pattern)
	}
}

// RegisterRoute registers a new route with a path, the handler function and the middlewares.
// This does not match the exact path
func (r *Router) RegisterRoute(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	r.HandleFunc(pattern, handler, middlewares...)
}

// RegisterExactRoute registers a new route with a path, the handler function and the middlewares.
// This matches the exact path
func (r *Router) RegisterExactRoute(
	pattern string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	r.ExactHandleFunc(pattern, handler, middlewares...)
}

// RegisterHandler registers a new route group with a path and a handler function
func (r *Router) RegisterHandler(pattern string, handler http.Handler) {
	// Check if the pattern contains a trailing slash and remove it
	if len(pattern) > 0 && pattern[len(pattern)-1] == '/' {
		pattern = pattern[:len(pattern)-1]
	}

	// Register the route group
	r.mux.Handle(pattern+"/", http.StripPrefix(pattern, handler))

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRouteGroup(r.fullPath, pattern)
	}
}

// RegisterGroup registers a new router group with a path and a router
func (r *Router) RegisterGroup(router RouterWrapper) {
	r.RegisterHandler(router.Pattern(), router.Handler())
}

// NewGroup creates a new router group with a path
func (r *Router) NewGroup(
	pattern string,
	middlewares ...func(next http.Handler) http.Handler,
) RouterWrapper {
	// Create a new group
	newGroup, _ := NewGroup(r, pattern, middlewares...)
	return newGroup
}

// Pattern returns the pattern
func (r *Router) Pattern() string {
	return r.pattern
}

// RelativePath returns the relative path
func (r *Router) RelativePath() string {
	return r.relativePath
}

// FullPath returns the full path
func (r *Router) FullPath() string {
	return r.fullPath
}

// Method returns the method
func (r *Router) Method() string {
	return r.method
}

// ServeStaticFiles serves the static files
func (r *Router) ServeStaticFiles(
	pattern,
	path string,
) {
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
func (r *Router) Logger() *Logger {
	return r.logger
}

// Mode returns the mode
func (r *Router) Mode() *goflagsmode.Flag {
	return r.mode
}
