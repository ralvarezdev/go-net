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
	}

	// Module is the struct for the route module
	Module struct {
		path       string
		service    ServiceWrapper
		validator  ValidatorWrapper
		controller ControllerWrapper
	}
)

// NewModule is a function that creates a new instance of the Module struct
func NewModule(
	path string,
	service ServiceWrapper,
	validator ValidatorWrapper,
	controller ControllerWrapper,
) ModuleWrapper {
	return &Module{
		path:       path,
		service:    service,
		validator:  validator,
		controller: controller,
	}
}

// Create is a function that creates a new instance of the Module struct
func (m *Module) Create(
	baseRouter gonethttproute.RouterWrapper,
) error {
	return m.controller.CreateRouter(baseRouter, m.path)
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
