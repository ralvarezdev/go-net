package handler

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttperrors "github.com/ralvarezdev/go-net/http/status/errors"
	"net/http"
)

type (
	// Handler interface for handling the responses
	Handler interface {
		HandleRequest(
			w http.ResponseWriter,
			r *http.Request,
			body interface{},
		) error
		HandleValidations(
			w http.ResponseWriter,
			body interface{},
			validatorFn func(body interface{}) (interface{}, error),
		) bool
		HandleRequestAndValidations(
			w http.ResponseWriter,
			r *http.Request,
			body interface{},
			validatorFn func(body interface{}) (interface{}, error),
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
	body interface{},
) (err error) {
	return d.jsonDecoder.Decode(w, r, body)
}

// HandleValidations handles the validations
func (d *DefaultHandler) HandleValidations(
	w http.ResponseWriter,
	body interface{},
	validatorFn func(body interface{}) (interface{}, error),
) bool {
	// Validate the request body
	validations, err := validatorFn(body)

	// Check if the error is nil and there are no validations
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
				nil,
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
	body interface{},
	validatorFn func(body interface{}) (interface{}, error),
) bool {
	// Handle the request
	if err := d.HandleRequest(w, r, body); err != nil {
		return false
	}

	// Handle the validations
	return d.HandleValidations(w, body, validatorFn)
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

	// Call the JSON encoder
	_ = d.jsonEncoder.Encode(w, response)
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
