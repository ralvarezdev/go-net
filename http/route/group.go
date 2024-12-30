package route

import (
	"net/http"
)

type (
	// GroupWrapper is the interface for the route group
	GroupWrapper interface {
		Group() *http.ServeMux
		HandleFunc(string, http.HandlerFunc)
		RegisterGroup(string, http.Handler)
	}

	// Group is the route group struct
	Group struct {
		mux *http.ServeMux
	}
)

// NewGroup creates a new route group
func NewGroup(mux *http.ServeMux) *Group {
	return &Group{mux: mux}
}

// Group returns the ServeMux
func (g *Group) Group() *http.ServeMux {
	return g.mux
}

// HandleFunc registers a new route with a path and a handler function
func (g *Group) HandleFunc(path string, handler http.HandlerFunc) {
	g.mux.HandleFunc(path, handler)
}

// RegisterGroup registers a new group with a path and a handler function
func (g *Group) RegisterGroup(path string, handler http.Handler) {
	g.mux.Handle(path, http.StripPrefix(path, http.StripPrefix(path, handler)))
}
