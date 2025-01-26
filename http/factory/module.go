package factory

import (
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
	"net/http"
)

type (
	// ModuleWrapper is the interface for the route module
	ModuleWrapper interface {
		Create(baseRouter gonethttproute.RouterWrapper) error
		GetRouter() gonethttproute.RouterWrapper
		Handler() http.Handler
		GetService() interface{}
		GetValidator() interface{}
		GetController() interface{}
		GetPath() string
		SetLoadFn(loadFn func())
		GetLoadFn() func()
		SetRegisterRoutesFn(registerRoutesFn func())
		GetSubmodules() *[]ModuleWrapper
		gonethttproute.RouterWrapper
	}

	// Module is the struct for the route module
	Module struct {
		path             string
		service          interface{}
		validator        interface{}
		controller       interface{}
		loadFn           func()
		registerRoutesFn func()
		submodules       []ModuleWrapper
		gonethttproute.RouterWrapper
	}
)

// NewModule is a function that creates a new instance of the Module struct
func NewModule(
	path string,
	service,
	validator,
	controller interface{},
	submodules ...ModuleWrapper,
) ModuleWrapper {
	return &Module{
		path:       path,
		service:    service,
		validator:  validator,
		controller: controller,
		submodules: submodules,
	}
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
	m.RouterWrapper = baseRouter.NewGroup(m.path)

	// Register the routes
	if m.registerRoutesFn != nil {
		m.registerRoutesFn()
	}

	// Create the submodules controllers router
	router := m.GetRouter()
	for _, submodule := range m.submodules {
		if err := submodule.Create(router); err != nil {
			return err
		}
	}

	// Load the module
	if m.loadFn != nil {
		m.loadFn()
	}
	return nil
}

// GetRouter returns the router
func (m *Module) GetRouter() gonethttproute.RouterWrapper {
	return m.RouterWrapper
}

// Handler is a function that returns the handler
func (m *Module) Handler() http.Handler {
	return m.Handler()
}

// GetPath is a function that returns the path
func (m *Module) GetPath() string {
	return m.path
}

// GetService is a function that returns the service
func (m *Module) GetService() interface{} {
	return m.service
}

// GetValidator is a function that returns the validator
func (m *Module) GetValidator() interface{} {
	return m.validator
}

// GetController is a function that returns the controller
func (m *Module) GetController() interface{} {
	return m.controller
}

// SetLoadFn is a function that sets the load function
func (m *Module) SetLoadFn(loadFn func()) {
	m.loadFn = loadFn
}

// GetLoadFn is a function that returns the load function
func (m *Module) GetLoadFn() func() {
	return m.loadFn
}

// SetRegisterRoutesFn is a function that sets the register routes function
func (m *Module) SetRegisterRoutesFn(registerRoutesFn func()) {
	m.registerRoutesFn = registerRoutesFn
}

// GetSubmodules is a function that returns the submodules
func (m *Module) GetSubmodules() *[]ModuleWrapper {
	return &m.submodules
}
