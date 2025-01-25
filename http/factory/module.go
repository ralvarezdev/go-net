package factory

import (
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
)

type (
	// ModuleWrapper is the interface for the route module
	ModuleWrapper interface {
		Create(baseRouter gonethttproute.RouterWrapper, path string) error
	}

	// Module is the struct for the route module
	Module struct {
		Service    ServiceWrapper
		Validator  ValidatorWrapper
		Controller ControllerWrapper
	}
)

// NewModule is a function that creates a new instance of the Module struct
func NewModule(
	service ServiceWrapper,
	validator ValidatorWrapper,
	controller ControllerWrapper,
) ModuleWrapper {
	return &Module{
		Service:    service,
		Validator:  validator,
		Controller: controller,
	}
}

// Create is a function that creates a new instance of the Module struct
func (m *Module) Create(
	baseRouter gonethttproute.RouterWrapper,
	path string,
) error {
	return m.Controller.CreateRouter(baseRouter, path)

}
