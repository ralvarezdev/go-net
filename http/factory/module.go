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
		GetService() ServiceWrapper
		GetValidator() ValidatorWrapper
		GetController() ControllerWrapper
		GetPath() string
		GetLoadFn() func()
		GetSubmodules() *[]ModuleWrapper
	}

	// Module is the struct for the route module
	Module struct {
		path       string
		service    ServiceWrapper
		validator  ValidatorWrapper
		controller ControllerWrapper
		loadFn     func()
		submodules []ModuleWrapper
	}
)

// NewModule is a function that creates a new instance of the Module struct
func NewModule(
	path string,
	service ServiceWrapper,
	validator ValidatorWrapper,
	controller ControllerWrapper,
	loadFn func(),
	submodules ...ModuleWrapper,
) ModuleWrapper {
	return &Module{
		path:       path,
		service:    service,
		validator:  validator,
		controller: controller,
		loadFn:     loadFn,
		submodules: submodules,
	}
}

// Create is a function that creates the router for the controller and its submodules, and loads the module
func (m *Module) Create(
	baseRouter gonethttproute.RouterWrapper,
) error {
	// Create the router for the controller
	if err := m.controller.CreateRouter(baseRouter, m.path); err != nil {
		return err
	}

	// Register the controller routes
	m.controller.RegisterRoutes()

	// Create the submodules controllers router
	router := m.controller.GetRouter()
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

// GetRouter is a function that returns the router
func (m *Module) GetRouter() gonethttproute.RouterWrapper {
	return m.controller.GetRouter()
}

// Handler is a function that returns the handler
func (m *Module) Handler() http.Handler {
	return m.controller.Handler()
}

// GetPath is a function that returns the path
func (m *Module) GetPath() string {
	return m.path
}

// GetService is a function that returns the service
func (m *Module) GetService() ServiceWrapper {
	return m.service
}

// GetValidator is a function that returns the validator
func (m *Module) GetValidator() ValidatorWrapper {
	return m.validator
}

// GetController is a function that returns the controller
func (m *Module) GetController() ControllerWrapper {
	return m.controller
}

// GetLoadFn is a function that returns the load function
func (m *Module) GetLoadFn() func() {
	return m.loadFn
}

// GetSubmodules is a function that returns the submodules
func (m *Module) GetSubmodules() *[]ModuleWrapper {
	return &m.submodules
}
