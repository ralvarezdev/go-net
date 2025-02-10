package http

import (
	"fmt"
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
	"net/http"
)

type (
	// Module is the struct for the route module
	Module struct {
		Pattern          string
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

	// Append the router middlewares to the module middlewares
	if baseRouter.GetMiddlewares() != nil {
		if m.Middlewares == nil {
			m.Middlewares = NewMiddlewares(*baseRouter.GetMiddlewares()...)
		} else {
			// Get the base router middlewares
			moduleMiddlewares := NewMiddlewares(*baseRouter.GetMiddlewares()...)

			// Append the module middlewares to the base router middlewares
			*moduleMiddlewares = append(*moduleMiddlewares, *m.Middlewares...)

			// Set the module middlewares
			m.Middlewares = moduleMiddlewares
		}
	}

	// Set the base route
	if m.Middlewares != nil {
		m.RouterWrapper = baseRouter.NewGroup(m.Pattern, *m.Middlewares...)
	} else {
		m.RouterWrapper = baseRouter.NewGroup(m.Pattern)
	}

	// Create the submodules controllers router
	router := m.GetRouter()
	if m.Submodules != nil {
		for i, submodule := range *m.Submodules {
			if submodule == nil {
				return fmt.Errorf(ErrNilSubmodule, m.Pattern, i)
			}

			if err := submodule.Create(router); err != nil {
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
func (m *Module) GetRouter() gonethttproute.RouterWrapper {
	return m.RouterWrapper
}
