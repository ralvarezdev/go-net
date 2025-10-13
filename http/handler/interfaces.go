package handler

import (
	"net/http"

	gonethttpparser "github.com/ralvarezdev/go-net/http/parser"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/mapper/validator"
)

type (
	// Handler interface for handling the requests
	Handler interface {
		gonethttpparser.Parser
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
