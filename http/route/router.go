package route

import (
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
			path string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterRoute(
			path string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterHandler(path string, handler http.Handler)
		NewGroup(
			path string,
			middlewares ...func(next http.Handler) http.Handler,
		) *Router
		RegisterGroup(path string, router *Router)
	}

	// Router is the route group struct
	Router struct {
		middlewares  []func(http.Handler) http.Handler
		firstHandler http.Handler
		mux          *http.ServeMux
		path         string
		mode         *goflagsmode.Flag
		logger       *Logger
	}
)

// NewRouter creates a new router
func NewRouter(
	path string,
	mode *goflagsmode.Flag,
	logger *Logger,
	middlewares ...func(next http.Handler) http.Handler,
) *Router {
	// Check if the path is empty
	if path == "" {
		path = "/"
	}

	// Initialize the multiplexer
	mux := http.NewServeMux()

	// Chain the handlers
	firstHandler := ChainHandlers(mux, middlewares...)

	return &Router{
		middlewares,
		firstHandler,
		mux,
		path,
		mode,
		logger,
	}
}

// NewBaseRouter creates a new base router
func NewBaseRouter(
	mode *goflagsmode.Flag,
	logger *Logger,
	middlewares ...func(next http.Handler) http.Handler,
) *Router {
	return NewRouter("", mode, logger, middlewares...)
}

// NewGroup creates a new router group
func NewGroup(
	baseRouter *Router,
	path string,
	middlewares ...func(next http.Handler) http.Handler,
) (*Router, error) {
	// Check if the base router is nil
	if baseRouter == nil {
		return nil, ErrNilRouter
	}

	// Check the base router path
	routerPath := path
	if baseRouter.path != "/" {
		routerPath = baseRouter.path + path
	}

	// Initialize the multiplexer
	mux := http.NewServeMux()

	// Chain the handlers
	firstHandler := ChainHandlers(mux, middlewares...)

	// Create a new router
	instance := &Router{
		middlewares:  middlewares,
		firstHandler: firstHandler,
		mux:          mux,
		logger:       baseRouter.logger,
		path:         routerPath,
		mode:         baseRouter.mode,
	}

	// Register the group
	baseRouter.RegisterGroup(path, instance)

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
	path string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	// Chain the handlers
	firstHandler := ChainHandlers(handler, middlewares...)

	// Register the route
	r.mux.HandleFunc(path, firstHandler.ServeHTTP)

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRoute(r.path, path)
	}
}

// RegisterRoute registers a new route with a path, the handler function and the middlewares
func (r *Router) RegisterRoute(
	path string,
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) {
	r.HandleFunc(path, handler, middlewares...)
}

// RegisterHandler registers a new route group with a path and a handler function
func (r *Router) RegisterHandler(path string, handler http.Handler) {
	// Check if the path contains a trailing slash and remove it
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Register the route group
	r.mux.Handle(path+"/", http.StripPrefix(path, handler))

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRouteGroup(r.path, path)
	}
}

// RegisterGroup registers a new router group with a path and a router
func (r *Router) RegisterGroup(path string, router *Router) {
	r.RegisterHandler(path, router.mux)
}

// NewGroup creates a new router group with a path
func (r *Router) NewGroup(
	path string,
	middlewares ...func(next http.Handler) http.Handler,
) *Router {
	// Append the base router middlewares
	middlewares = append(middlewares, r.middlewares...)

	// Create a new group
	newGroup, _ := NewGroup(r, path, middlewares...)
	return newGroup
}
