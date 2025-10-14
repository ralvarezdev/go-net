package handler

import (
	"net/http"

	govalidatormappervalidator "github.com/ralvarezdev/go-validator/mapper/validator"
)

type (
	// RequestsHandler interface for handling the requests
	RequestsHandler interface {
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
