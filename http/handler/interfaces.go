package handler

import (
	"net/http"

	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
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

	// RawErrorHandler interface for handling raw errors
	RawErrorHandler interface {
		HandleRawError(
			w http.ResponseWriter,
			err error,
			handleResponseFn func(
				w http.ResponseWriter,
				response gonethttpresponse.Response,
			),
		)
	}

	// ResponsesHandler interface for handling the responses
	ResponsesHandler interface {
		gonethttpresponse.Encoder
		HandleResponse(
			w http.ResponseWriter,
			response gonethttpresponse.Response,
		)
		HandleRawError(
			w http.ResponseWriter,
			err error,
		)
		HandleError(
			w http.ResponseWriter,
			err error,
			httpStatus int,
		)
		HandleErrorWithCode(
			w http.ResponseWriter,
			err error,
			errCode string,
			httpStatus int,
		)
		HandleDebugError(
			w http.ResponseWriter,
			debugErr error,
			err error,
			httpStatus int,
		)
		HandleDebugErrorWithCode(
			w http.ResponseWriter,
			debugErr error,
			err error,
			errCode string,
			httpStatus int,
		)
		HandleFailFieldError(
			w http.ResponseWriter,
			field string,
			err error,
			httpStatus int,
		)
		HandleFailFieldErrorWithCode(
			w http.ResponseWriter,
			field string,
			err error,
			errCode string,
			httpStatus int,
		)
		HandleFailDataError(
			w http.ResponseWriter,
			data interface{},
			httpStatus int,
		)
		HandleFailDataErrorWithCode(
			w http.ResponseWriter,
			data interface{},
			errCode string,
			httpStatus int,
		)
	}

	// Handler is the interface that handles both the requests decoding and responses encoding tasks
	Handler interface {
		ResponsesHandler
		RequestsHandler
	}
)
