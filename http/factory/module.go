package factory

import (
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
)

type (
	// ModuleWrapper is the interface for the route module
	ModuleWrapper interface {
		Create(baseRouter gonethttproute.RouterWrapper) error
		CreateSubmodule(submodule ModuleWrapper) error
		CreateService() func()
		CreateController() func(baseRouter gonethttproute.RouterWrapper)
		CreateValidator() func()
	}

	// Module is the struct for the route module
	Module struct {
		Service    ServiceWrapper
		Validator  ValidatorWrapper
		Controller ControllerWrapper
	}
)

// Create is a function that creates a new instance of the Module struct
func (m *Module) Create(baseRouter gonethttproute.RouterWrapper) error {
	// Create the service struct
	m.CreateService()()

	// Check if the service is nil
	if m.Service == nil {
		return ErrNilService
	}

	// Create the validator struct
	m.CreateValidator()()

	// Check if the validator is nil
	if m.Validator == nil {
		return ErrNilValidator
	}

	// Create the controller struct
	m.CreateController()(baseRouter)

	// Check if the controller is nil
	if m.Controller == nil {
		return ErrNilController
	}

	// Register the controller routes and groups
	m.Controller.RegisterRoutes()
	m.Controller.RegisterGroups()
	return nil
}

// CreateSubmodule is a function that creates a new submodule of the Module struct
func (m *Module) CreateSubmodule(submodule ModuleWrapper) error {
	return submodule.Create(m.Controller)
}

// CreateService is a function that creates a new instance of the Service struct
func (m *Module) CreateService() func() {
	return func() {}
}

// CreateValidator is a function that creates a new instance of the Validator struct
func (m *Module) CreateValidator() func() {
	return func() {}
}

// CreateController is a function that creates a new instance of the Controller struct
func (m *Module) CreateController() func(baseRouter gonethttproute.RouterWrapper) {
	return func(baseRouter gonethttproute.RouterWrapper) {}
}
