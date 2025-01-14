package handler

import (
	"errors"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatuserrors "github.com/ralvarezdev/go-net/http/status/errors"
	"net/http"
)

type (
	// Handler interface for handling the responses
	Handler interface {
		Decode(
			w http.ResponseWriter,
			r *http.Request,
			dest interface{},
		) error
		Validate(
			w http.ResponseWriter,
			body interface{},
			validatorFn func(body interface{}) (interface{}, error),
		) bool
		DecodeAndValidate(
			w http.ResponseWriter,
			r *http.Request,
			dest interface{},
			validatorFn func(body interface{}) (interface{}, error),
		) bool
		HandleResponse(
			w http.ResponseWriter,
			response *gonethttpresponse.Response,
		)
		HandleError(
			w http.ResponseWriter,
			err error,
			httpStatus int,
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

// Decode decodes the request body into the destination
func (d *DefaultHandler) Decode(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
) (err error) {
	return d.jsonDecoder.Decode(w, r, dest)
}

// Validate validates the request body
func (d *DefaultHandler) Validate(
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
				gonethttpstatuserrors.InternalServerError,
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

// DecodeAndValidate decodes and validates the request body
func (d *DefaultHandler) DecodeAndValidate(
	w http.ResponseWriter,
	r *http.Request,
	body interface{},
	validatorFn func(body interface{}) (interface{}, error),
) bool {
	// Decode the request body
	if err := d.Decode(w, r, body); err != nil {
		return false
	}

	// Validate the request body
	return d.Validate(w, body, validatorFn)
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
				gonethttpstatuserrors.InternalServerError,
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

// HandleError handles the error response
func (d *DefaultHandler) HandleError(
	w http.ResponseWriter,
	err error,
) {
	// Check if the errors is a request error
	var e gonethttpresponse.RequestError
	if errors.As(err, &e) {
		d.HandleResponse(
			w, gonethttpresponse.NewFailResponseFromRequestError(e),
		)
		return
	}

	d.HandleResponse(
		w, gonethttpresponse.NewDebugErrorResponse(
			gonethttpstatuserrors.InternalServerError,
			err,
			nil, nil, http.StatusInternalServerError,
		),
	)
}
