package handler

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttperrors "github.com/ralvarezdev/go-net/http/errors"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

type (
	// Handler interface for handling the responses
	Handler interface {
		HandleRequest(
			w http.ResponseWriter,
			r *http.Request,
			data interface{},
		) error
		HandleValidations(
			w http.ResponseWriter,
			r *http.Request,
			validatorFn func() (interface{}, error),
		) bool
		HandleRequestAndValidations(
			w http.ResponseWriter,
			r *http.Request,
			data interface{},
			validatorFn func() (interface{}, error),
		) bool
		HandleResponse(
			w http.ResponseWriter,
			response *gonethttpresponse.Response,
		)
		HandleErrorProneResponse(
			w http.ResponseWriter,
			successResponse *gonethttpresponse.Response,
			errorResponse *gonethttpresponse.Response,
		)
	}

	// DefaultHandler struct
	DefaultHandler struct {
		mode        *goflagsmode.Flag
		jsonEncoder gonethttpjson.Encoder
		jsonDecoder gonethttpjson.Decoder
	}
)

// NewDefaultHandler creates a new default request handler
func NewDefaultHandler(
	mode *goflagsmode.Flag,
	jsonEncoder gonethttpjson.Encoder,
	jsonDecoder gonethttpjson.Decoder,
) (*DefaultHandler, error) {
	// Check if the flag mode, the JSON encoder or the JSON decoder is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	if jsonEncoder == nil {
		return nil, gonethttpjson.ErrNilEncoder
	}
	if jsonDecoder == nil {
		return nil, gonethttpjson.ErrNilDecoder
	}

	return &DefaultHandler{
		mode:        mode,
		jsonEncoder: jsonEncoder,
		jsonDecoder: jsonDecoder,
	}, nil
}

// HandleRequest handles the request
func (d *DefaultHandler) HandleRequest(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) (err error) {
	return d.jsonDecoder.Decode(w, r, data)
}

// HandleValidations handles the validations
func (d *DefaultHandler) HandleValidations(
	w http.ResponseWriter,
	r *http.Request,
	validatorFn func() (interface{}, error),
) bool {
	// Validate the request body
	validations, err := validatorFn()

	// Check if the error is nil are there no validations
	if err == nil && validations == nil {
		return true
	}

	// Check if the error is not nil
	if err != nil {
		d.HandleResponse(
			w,
			gonethttpresponse.NewDebugErrorResponse(
				gonethttperrors.InternalServerError,
				err,
				nil, nil,
				http.StatusInternalServerError,
			),
		)
	} else {
		d.HandleResponse(
			w,
			gonethttpresponse.NewFailResponse(
				validations,
				http.StatusBadRequest,
			),
		)
	}
	return false
}

// HandleRequestAndValidations handles the request and the validations
func (d *DefaultHandler) HandleRequestAndValidations(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
	validatorFn func() (interface{}, error),
) bool {
	// Handle the request
	if err := d.HandleRequest(w, r, data); err != nil {
		return false
	}

	// Handle the validations
	return d.HandleValidations(w, r, validatorFn)
}

// HandleResponse handles the response
func (d *DefaultHandler) HandleResponse(
	w http.ResponseWriter,
	response *gonethttpresponse.Response,
) {
	// Check if the response is nil
	if response == nil {
		d.HandleResponse(
			w,
			gonethttpresponse.NewDebugErrorResponse(
				gonethttperrors.InternalServerError,
				gonethttpresponse.ErrNilResponse,
				nil, nil,
				http.StatusInternalServerError,
			),
		)
		return
	}

	// Check if the response contains the debug response
	if response.DebugResponse != nil && d.mode != nil && d.mode.IsDebug() {
		_ = d.jsonEncoder.Encode(
			w,
			response.DebugResponse,
			response.HTTPStatus,
		)
		return
	}
	_ = d.jsonEncoder.Encode(
		w,
		response.Response,
		response.HTTPStatus,
	)
}

// HandleErrorProneResponse handles the response that may contain an error
func (d *DefaultHandler) HandleErrorProneResponse(
	w http.ResponseWriter,
	successResponse *gonethttpresponse.Response,
	errorResponse *gonethttpresponse.Response,
) {
	// Check if the error response is nil
	if errorResponse != nil {
		d.HandleResponse(w, errorResponse)
		return
	}

	// Handle the success response
	d.HandleResponse(w, successResponse)
}
