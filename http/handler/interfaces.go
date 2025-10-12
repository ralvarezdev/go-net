package handler

import (
	"net/http"

	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/mapper/validator"
)

type (
	// Handler interface for handling the requests
	Handler interface {
		Decode(
			w http.ResponseWriter,
			r *http.Request,
			dest interface{},
		) error
		Validate(
			w http.ResponseWriter,
			body interface{},
			validatorFn govalidatormappervalidator.ValidateFn,
		) bool
		Parse(
			w http.ResponseWriter,
			r *http.Request,
			dest interface{},
			validatorFn govalidatormappervalidator.ValidateFn,
		) bool
		GetParameters(
			r *http.Request,
			keys ...string,
		) map[string]string
		ParseWildcard(
			w http.ResponseWriter,
			r *http.Request,
			wildcardKey string,
			dest interface{},
			toTypeFn ToTypeFn,
		) bool
		HandleResponse(
			w http.ResponseWriter,
			response gonethttpresponse.Response,
		)
		HandleError(
			w http.ResponseWriter,
			err error,
		)
	}
)
