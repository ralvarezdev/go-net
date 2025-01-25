package factory

import (
	gonethttproute "github.com/ralvarezdev/go-net/http/route"
)

type (
	// ModuleWrapper is the interface for the route module
	ModuleWrapper interface {
		Create(baseRouter gonethttproute.RouterWrapper) error
		CreateSubmodule(submodule ModuleWrapper) error
		CreateService() func() error
		CreateController() func(baseRouter gonethttproute.RouterWrapper) error
		CreateValidator() func() error
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
	// Create the service, validator and controller structs
	if err := m.CreateService()(); err != nil {
		return err
	}
	if err := m.CreateValidator()(); err != nil {
		return err
	}
	if err := m.CreateController()(baseRouter); err != nil {
		return err
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
func (m *Module) CreateService() func() error {
	return func() error {
		// Check if the service is nil
		if m.Service == nil {
			return ErrNilService
		}
		return nil
	}
}

// CreateValidator is a function that creates a new instance of the Validator struct
func (m *Module) CreateValidator() func() error {
	return func() error {
		// Check if the validator is nil
		if m.Validator == nil {
			return ErrNilValidator
		}
		return nil
	}
}

// CreateController is a function that creates a new instance of the Controller struct
func (m *Module) CreateController() func(baseRouter gonethttproute.RouterWrapper) error {
	return func(baseRouter gonethttproute.RouterWrapper) error {
		// Check if the controller is nil
		if m.Controller == nil {
			return ErrNilController
		}
		return nil
	}
}
