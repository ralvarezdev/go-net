package route

import (
	"net/http"
)

// Group is the route group struct
type Group struct {
	mux *http.ServeMux
}

// NewGroup creates a new group
func NewGroup(mux *http.ServeMux) *Group {
	return &Group{mux: mux}
}

// Mux returns the ServeMux
func (g *Group) Mux() *http.ServeMux {
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
