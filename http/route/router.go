package route

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	"net/http"
)

type (
	// RouterWrapper is the interface for the routes
	RouterWrapper interface {
		Handler() *http.ServeMux
		ChainHandlers(handlers ...http.HandlerFunc) http.HandlerFunc
		HandleFunc(path string, handlers ...http.HandlerFunc)
		RegisterRoute(path string, handlers ...http.HandlerFunc)
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

// ChainHandlers chains the handlers functions
func (r *Router) ChainHandlers(handlers ...http.HandlerFunc) http.HandlerFunc {
	// Check if the handlers are empty
	if len(handlers) == 0 {
		return nil
	}

	// Set the handlers
	var firstHandler, modifiedHandler, nextHandler http.HandlerFunc
	for i, h := range handlers {
		// Check if it is not the last handler
		if i < len(handlers)-1 {
			nextHandler = handlers[i+1]
		} else {
			nextHandler = nil
		}

		// Set the handler
		modifiedHandler = func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			// Call the handler
			h(writer, request)

			// Check if there is a next handler
			if nextHandler != nil {
				nextHandler(writer, request)
			}
		}

		// Check if it is the first handler
		if i == 0 {
			firstHandler = modifiedHandler
		}
	}
	return firstHandler
}

// HandleFunc registers a new route with a path and the handlers functions
// It's not needed to call the next handler function in each handler
func (r *Router) HandleFunc(path string, handlers ...http.HandlerFunc) {
	// Chain the handlers
	firstHandler := r.ChainHandlers(handlers...)

	// Register the route
	r.mux.HandleFunc(path, firstHandler)

	if r.logger != nil && r.mode != nil && !r.mode.IsProd() {
		r.logger.RegisterRoute(r.path, path)
	}
}

// RegisterRoute registers a new route with a path and the handlers functions
// It's not needed to call the next handler function in each handler
func (r *Router) RegisterRoute(path string, handlers ...http.HandlerFunc) {
	r.HandleFunc(path, handlers...)
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
