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
		RegisterRouteGroup(path string, router *Router)
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

// NewRouterGroup creates a new route group
func NewRouterGroup(baseRoute *Router, path string) (*Router, error) {
	// Check if the base route is nil
	if baseRoute == nil {
		return nil, ErrNilRouter
	}

	// Check the base route path
	routerPath := path
	if baseRoute.path != "/" {
		routerPath = baseRoute.path + path
	}

	// Create a new router
	instance := &Router{
		mux:    http.NewServeMux(),
		logger: baseRoute.logger,
		path:   routerPath,
		mode:   baseRoute.mode,
	}

	// Register the route group
	baseRoute.RegisterRouteGroup(path, instance)

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

// RegisterRouteGroup registers a new route group with a path and a router
func (r *Router) RegisterRouteGroup(path string, router *Router) {
	r.RegisterHandler(path, router.mux)
}
