package http

import (
	"fmt"
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
	"net/http"
)

type (
	// Module is the struct for the route module
	Module struct {
		Path             string
		Service          interface{}
		Controller       interface{}
		BeforeLoadFn     func(m *Module)
		RegisterRoutesFn func(m *Module)
		AfterLoadFn      func(m *Module)
		Middlewares      *[]func(next http.Handler) http.Handler
		Submodules       *[]*Module
		gonethttproute.RouterWrapper
	}
)

// NewMiddlewares is a function that creates a new middlewares slice
func NewMiddlewares(middlewares ...func(next http.Handler) http.Handler) *[]func(next http.Handler) http.Handler {
	return &middlewares
}

// NewSubmodules is a function that creates a new submodules slice
func NewSubmodules(submodules ...*Module) *[]*Module {
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

	// Run the before load function
	if m.BeforeLoadFn != nil {
		m.BeforeLoadFn(m)
	}

	// Set the base route
	if m.Middlewares != nil {
		m.RouterWrapper = baseRouter.NewGroup(m.Path, *m.Middlewares...)
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
		for i, submodule := range *m.Submodules {
			if submodule == nil {
				return fmt.Errorf(ErrNilSubmodule, m.Path, i)
			}

			if err := submodule.Create(router); err != nil {
				return err
			}
		}
	}

	// Run the after load function
	if m.AfterLoadFn != nil {
		m.AfterLoadFn(m)
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

// GetController is a function that returns the controller
func (m *Module) GetController() interface{} {
	return m.Controller
}
