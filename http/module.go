package http

import (
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
	"net/http"
)

type (
	// ModuleWrapper is the interface for the route module
	ModuleWrapper interface {
		Create(baseRouter gonethttproute.RouterWrapper) error
		GetRouter() gonethttproute.RouterWrapper
		GetService() interface{}
		GetValidator() interface{}
		GetController() interface{}
		GetPath() string
		gonethttproute.RouterWrapper
	}

	// Module is the struct for the route module
	Module struct {
		Path             string
		Service          interface{}
		Validator        interface{}
		Controller       interface{}
		LoadFn           func(m *Module)
		RegisterRoutesFn func(m *Module)
		Middlewares      *[]func(next http.Handler) http.Handler
		Submodules       *[]ModuleWrapper
		gonethttproute.RouterWrapper
	}
)

// NewMiddlewares is a function that creates a new middlewares slice
func NewMiddlewares(middlewares ...func(next http.Handler) http.Handler) *[]func(next http.Handler) http.Handler {
	return &middlewares
}

// NewSubmodules is a function that creates a new submodules slice
func NewSubmodules(submodules ...ModuleWrapper) *[]ModuleWrapper {
	return &submodules
}

// Create is a function that creates the router for the controller and its submodules, and loads the module
func (m *Module) Create(
	baseRouter gonethttproute.RouterWrapper,
) error {
	// Check if the base route is nil
	if baseRouter == nil {
		return gonethttproute.ErrNilRouter
	}

	// Set the base route
	if m.Middlewares != nil {
		baseRouter = baseRouter.NewGroup(m.Path, *m.Middlewares...)
	} else {
		m.RouterWrapper = baseRouter.NewGroup(m.Path)
	}

	// Register the routes
	if m.RegisterRoutesFn != nil {
		m.RegisterRoutesFn(m)
	}

	// Create the submodules controllers router
	router := m.GetRouter()
	if m.Submodules != nil {
		for _, submodule := range *m.Submodules {
			if err := submodule.Create(router); err != nil {
				return err
			}
		}
	}

	// Load the module
	if m.LoadFn != nil {
		m.LoadFn(m)
	}
	return nil
}

// GetRouter returns the router
func (m *Module) GetRouter() gonethttproute.RouterWrapper {
	return m.RouterWrapper
}

// GetPath is a function that returns the path
func (m *Module) GetPath() string {
	return m.Path
}

// GetService is a function that returns the service
func (m *Module) GetService() interface{} {
	return m.Service
}

// GetValidator is a function that returns the validator
func (m *Module) GetValidator() interface{} {
	return m.Validator
}

// GetController is a function that returns the controller
func (m *Module) GetController() interface{} {
	return m.Controller
}
