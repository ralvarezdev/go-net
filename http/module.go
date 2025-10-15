package http

import (
	"fmt"
	"net/http"

	gonethttproute "github.com/ralvarezdev/go-net/http/route"
)

type (
	// Module is the struct for the route module
	Module struct {
		Pattern          string
		BeforeLoadFn     func(m *Module)
		RegisterRoutesFn func(m *Module)
		AfterLoadFn      func(m *Module)
		Middlewares      []func(next http.Handler) http.Handler
		Submodules       []*Module
		gonethttproute.RouterWrapper
	}
)

// NewMiddlewares is a function that creates a new middlewares slice
//
// Usage: NewMiddlewares(middleware1, middleware2, ...)
//
// Parameters:
//
//   - middlewares: variadic list of middleware functions
//
// Returns:
//
//   - []func(next http.ResponsesHandler) http.ResponsesHandler: Slice of middleware functions
func NewMiddlewares(middlewares ...func(next http.Handler) http.Handler) []func(next http.Handler) http.Handler {
	return middlewares
}

// NewSubmodules is a function that creates a new submodules slice
//
// Usage: NewSubmodules(submodule1, submodule2, ...)
//
// Parameters:
//
//   - submodules: variadic list of submodule pointers
//
// Returns:
//
//   - []*Module: Slice of submodule pointers
func NewSubmodules(submodules ...*Module) []*Module {
	return submodules
}

// Create is a function that creates the router for the controller and its submodules, and loads the module
//
// Parameters:
//
//   - baseRouter: The base router to create the module's router group
//
// Returns:
//
//   - error: The error if any
func (m *Module) Create(
	baseRouter gonethttproute.RouterWrapper,
) error {
	if m == nil {
		return ErrNilModule
	}

	// Check if the base route is nil
	if baseRouter == nil {
		return gonethttproute.ErrNilRouter
	}

	// Run the before load function
	if m.BeforeLoadFn != nil {
		m.BeforeLoadFn(m)
	}

	// Set the base route
	var err error
	if m.Middlewares != nil {
		m.RouterWrapper, err = baseRouter.NewRouter(m.Pattern, m.Middlewares...)
	} else {
		m.RouterWrapper, err = baseRouter.NewRouter(m.Pattern)
	}
	if err != nil {
		return err
	}

	// Create the submodules controllers router
	router := m.GetRouter()
	if m.Submodules != nil {
		for i, submodule := range m.Submodules {
			if submodule == nil {
				return fmt.Errorf(ErrNilSubmodule, m.Pattern, i)
			}

			if err = submodule.Create(router); err != nil {
				return err
			}
		}
	}

	// Register the routes
	if m.RegisterRoutesFn != nil {
		m.RegisterRoutesFn(m)
	}

	// Run the after load function
	if m.AfterLoadFn != nil {
		m.AfterLoadFn(m)
	}
	return nil
}

// GetRouter returns the router
//
// Returns:
//
//   - gonethttproute.RouterWrapper: The router
func (m *Module) GetRouter() gonethttproute.RouterWrapper {
	if m == nil {
		return nil
	}
	return m.RouterWrapper
}
