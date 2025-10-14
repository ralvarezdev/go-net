package handler

import (
	"net/http"

	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// Handler interface for handling the responses
	Handler interface {
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
	}
)
