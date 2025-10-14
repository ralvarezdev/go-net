package handler

import (
	"net/http"

	govalidatormappervalidator "github.com/ralvarezdev/go-validator/mapper/validator"
)

type (
	// Handler interface for handling the requests
	Handler interface {
		Validate(
			w http.ResponseWriter,
			body interface{},
			validatorFn govalidatormappervalidator.ValidateFn,
		) bool
		DecodeAndValidate(
			w http.ResponseWriter,
			r *http.Request,
			dest interface{},
			validatorFn govalidatormappervalidator.ValidateFn,
		) bool
	}
)
