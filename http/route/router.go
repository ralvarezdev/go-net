package route

import (
	"net/http"
)

type (
	// RouterWrapper is the interface for the routes
	RouterWrapper interface {
		Router() *http.ServeMux
		HandleFunc(string, http.HandlerFunc)
		RegisterRoute(string, http.HandlerFunc)
		RegisterHandler(string, http.Handler)
		RegisterRouteGroup(string, http.Handler)
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

// Router returns the ServeMux
func (r *Router) Router() *http.ServeMux {
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
	r.mux.Handle(path, http.StripPrefix(path, http.StripPrefix(path, handler)))
}

// RegisterRouteGroup registers a new route group with a path and a router
func (r *Router) RegisterRouteGroup(path string, router *Router) {
	r.RegisterHandler(path, router.mux)
}
