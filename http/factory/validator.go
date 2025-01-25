package factory

import (
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/struct/mapper/validator"
)

type (
	// ValidatorWrapper is the interface for the route validator
	ValidatorWrapper interface {
		govalidatormappervalidator.Service
	}
)
