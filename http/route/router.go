package route

import (
	"fmt"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	"net/http"
)

type (
	// RouterWrapper is the interface for the routes
	RouterWrapper interface {
		Handler() http.Handler
		Mux() *http.ServeMux
		GetMiddlewares() *[]func(http.Handler) http.Handler
		HandleFunc(
			relativePath string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		ExactHandleFunc(
			relativePath string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterRoute(
			relativePath string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterExactRoute(
			relativePath string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterHandler(relativePath string, handler http.Handler)
		NewGroup(
			relativePath string,
			middlewares ...func(next http.Handler) http.Handler,
		) *Router
		RegisterGroup(router *Router)
		RelativePath() string
		FullPath() string
	}

	// Router is the route group struct
	Router struct {
		middlewares  []func(http.Handler) http.Handler
		firstHandler http.Handler
		mux          *http.ServeMux
		relativePath string
		fullPath     string
		mode         *goflagsmode.Flag
		logger       *Logger
	}
)

// AddSlash adds a slash to the path
func AddSlash(path string) string {
	if path == "" {
		return "/"
	} else if path[0] != '/' {
		return "/" + path
	}
	return path
}

// NewRouter creates a new router
func NewRouter(
	path string,
	mode *goflagsmode.Flag,
	logger *Logger,
	middlewares ...func(next http.Handler) http.Handler,
) (*Router, error) {
	// Add a slash to the path if it does not have it
	path = AddSlash(path)

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
		path,
		path,
		mode,
		logger,
	}, nil
}

// NewBaseRouter creates a new base router
func NewBaseRouter(
	mode *goflagsmode.Flag,
	logger *Logger,
	middlewares ...func(next http.Handler) http.Handler,
) (*Router, error) {
	return NewRouter("", mode, logger, middlewares...)
}

// NewGroup creates a new router group
func NewGroup(
	baseRouter *Router,
	relativePath string,
	middlewares ...func(next http.Handler) http.Handler,
) (*Router, error) {
	// Check if the base router is nil
	if baseRouter == nil {
		return nil, ErrNilRouter
	}

	// Add a slash to the path if it does not have it
	relativePath = AddSlash(relativePath)

	// Check the base router path
	var fullPath string
	if baseRouter.fullPath != "/" {
		fullPath = baseRouter.fullPath + relativePath
	} else {
		fullPath = relativePath
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
		logger:       baseRouter.logger,
		relativePath: relativePath,
		fullPath:     fullPath,
		mode:         baseRouter.mode,
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
	relativePath string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	// Chain the handlers
	firstHandler := ChainHandlers(handler, middlewares...)

	// Register the route
	r.mux.HandleFunc(relativePath, firstHandler.ServeHTTP)

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRoute(r.relativePath, relativePath)
	}
}

// ExactHandleFunc registers a new route with a path, the handler function and the middlewares
func (r *Router) ExactHandleFunc(
	relativePath string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	// Add slash to the path
	relativePath = AddSlash(relativePath)

	// Chain the handlers
	firstHandler := ChainHandlers(handler, middlewares...)

	// Add the '$' wildcard to the end of the path to match the exact path
	if relativePath[len(relativePath)-1] == '/' {
		relativePath += "{$}"
	}

	// Register the route
	r.mux.HandleFunc(relativePath, firstHandler.ServeHTTP)

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRoute(r.relativePath, relativePath)
	}
}

// RegisterRoute registers a new route with a path, the handler function and the middlewares.
// This does not match the exact path
func (r *Router) RegisterRoute(
	relativePath string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	r.HandleFunc(relativePath, handler, middlewares...)
}

// RegisterExactRoute registers a new route with a path, the handler function and the middlewares.
// This matches the exact path
func (r *Router) RegisterExactRoute(
	relativePath string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	r.ExactHandleFunc(relativePath, handler, middlewares...)
}

// RegisterHandler registers a new route group with a path and a handler function
func (r *Router) RegisterHandler(relativePath string, handler http.Handler) {
	// Check if the path contains a trailing slash and remove it
	if len(relativePath) > 1 && relativePath[len(relativePath)-1] == '/' {
		relativePath = relativePath[:len(relativePath)-1]
	}

	// Register the route group
	r.mux.Handle(relativePath+"/", http.StripPrefix(relativePath, handler))

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRouteGroup(r.relativePath, relativePath)
	}
}

// RegisterGroup registers a new router group with a path and a router
func (r *Router) RegisterGroup(router *Router) {
	r.RegisterHandler(router.RelativePath(), router.mux)
}

// NewGroup creates a new router group with a path
func (r *Router) NewGroup(
	relativePath string,
	middlewares ...func(next http.Handler) http.Handler,
) *Router {
	// Create the middlewares slice
	var fullMiddlewares []func(http.Handler) http.Handler

	// Append the base router middlewares
	fullMiddlewares = append(fullMiddlewares, r.middlewares...)

	// Append the new middlewares
	fullMiddlewares = append(fullMiddlewares, middlewares...)

	// Create a new group
	newGroup, _ := NewGroup(r, relativePath, fullMiddlewares...)
	return newGroup
}

// RelativePath returns the relative path
func (r *Router) RelativePath() string {
	return r.relativePath
}

// FullPath returns the full path
func (r *Router) FullPath() string {
	return r.fullPath
}
