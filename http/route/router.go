package route

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	"net/http"
)

type (
	// RouterWrapper is the interface for the routes
	RouterWrapper interface {
		Handler() *http.ServeMux
		HandleFunc(path string, handler http.HandlerFunc)
		RegisterRoute(path string, handler http.HandlerFunc)
		RegisterHandler(path string, handler http.Handler)
		NewGroup(path string) *Router
		RegisterGroup(path string, router *Router)
	}

	// Router is the route group struct
	Router struct {
		mux    *http.ServeMux
		path   string
		mode   *goflagsmode.Flag
		logger *Logger
	}
)

// NewRouter creates a new router
func NewRouter(path string, mode *goflagsmode.Flag, logger *Logger) *Router {
	// Check if the path is empty
	if path == "" {
		path = "/"
	}

	return &Router{
		mux:    http.NewServeMux(),
		logger: logger,
		path:   path,
		mode:   mode,
	}
}

// NewBaseRouter creates a new base router
func NewBaseRouter(mode *goflagsmode.Flag, logger *Logger) *Router {
	return NewRouter("", mode, logger)
}

// NewGroup creates a new router group
func NewGroup(baseRouter *Router, path string) (*Router, error) {
	// Check if the base router is nil
	if baseRouter == nil {
		return nil, ErrNilRouter
	}

	// Check the base router path
	routerPath := path
	if baseRouter.path != "/" {
		routerPath = baseRouter.path + path
	}

	// Create a new router
	instance := &Router{
		mux:    http.NewServeMux(),
		logger: baseRouter.logger,
		path:   routerPath,
		mode:   baseRouter.mode,
	}

	// Register the group
	baseRouter.RegisterGroup(path, instance)

	return instance, nil
}

// Handler returns the ServeMux
func (r *Router) Handler() *http.ServeMux {
	return r.mux
}

// HandleFunc registers a new route with a path and a handler function
func (r *Router) HandleFunc(path string, handler http.HandlerFunc) {
	// Register the route
	r.mux.HandleFunc(path, handler)

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRoute(r.path, path)
	}
}

// RegisterRoute registers a new route with a path and a handler function
func (r *Router) RegisterRoute(path string, handler http.HandlerFunc) {
	r.HandleFunc(path, handler)
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
func (r *Router) NewGroup(path string) *Router {
	newGroup, _ := NewGroup(r, path)
	return newGroup
}
