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

	// ResponsesHandler interface for handling the responses
	ResponsesHandler interface {
		gonethttpresponse.Encoder
		HandleResponse(
			w http.ResponseWriter,
			response gonethttpresponse.Response,
		)
		HandleError(
			w http.ResponseWriter,
			err error,
		)
		HandleErrorResponse(
			w http.ResponseWriter,
			err error,
			httpStatus int,
		)
		HandleErrorResponseWithCode(
			w http.ResponseWriter,
			err error,
			errCode string,
			httpStatus int,
		)
		HandleDebugErrorResponse(
			w http.ResponseWriter,
			debugErr error,
			err error,
			httpStatus int,
		)
		HandleDebugErrorResponseWithCode(
			w http.ResponseWriter,
			debugErr error,
			err error,
			errCode string,
			httpStatus int,
		)
		HandleFieldFailResponse(
			w http.ResponseWriter,
			field string,
			err error,
			httpStatus int,
		)
		HandleFieldFailResponseWithCode(
			w http.ResponseWriter,
			field string,
			err error,
			errCode string,
			httpStatus int,
		)
		HandleFailResponse(
			w http.ResponseWriter,
			data interface{},
			httpStatus int,
		)
		HandleFailResponseWithCode(
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
