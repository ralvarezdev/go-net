package factory

import (
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/struct/mapper/validator"
)

type (
	// Validator is the struct for the route validator
	Validator struct {
		govalidatormappervalidator.Service
	}
)
