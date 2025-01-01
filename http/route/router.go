package route

import (
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
		mux *http.ServeMux
	}
)

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

// NewRouterGroup creates a new route group
func NewRouterGroup(baseRoute *Router, path string) (*Router, error) {
	// Check if the base route is nil
	if baseRoute == nil {
		return nil, ErrNilRouter
	}

	// Create a new router
	instance := &Router{mux: http.NewServeMux()}

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
	r.mux.HandleFunc(path, handler)
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

	r.mux.Handle(path+"/", http.StripPrefix(path, handler))
}

// RegisterRouteGroup registers a new route group with a path and a router
func (r *Router) RegisterRouteGroup(path string, router *Router) {
	r.RegisterHandler(path, router.mux)
}
