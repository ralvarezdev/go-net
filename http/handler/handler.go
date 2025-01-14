package handler

import (
	"errors"
	"fmt"
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatusresponse "github.com/ralvarezdev/go-net/http/status/response"
	"net/http"
)

var (
	ErrCodeParameterNotFound      *string
	ErrCodeParameterParsingFailed *string
	ErrCodeValidationFailed       *string
	ErrCodeNilResponse            *string
	ErrCodeRequestFatalError      *string
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
			validatorFn func(body interface{}) (interface{}, error),
		) bool
		DecodeAndValidate(
			w http.ResponseWriter,
			r *http.Request,
			dest interface{},
			validatorFn func(body interface{}) (interface{}, error),
		) bool
		GetParameters(
			w http.ResponseWriter,
			r *http.Request,
			keys ...string,
		) (*map[string]string, bool)
		ParseParameter(
			w http.ResponseWriter,
			parameter string,
			dest interface{},
			toTypeFn func(parameter string, dest interface{}) error,
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
			gonethttpstatusresponse.NewJSendDebugInternalServerError(
				err,
				ErrCodeValidationFailed,
			),
		)
	} else {
		d.HandleResponse(
			w,
			gonethttpresponse.NewJSendFailResponse(
				validations,
				ErrCodeValidationFailed,
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

// GetParameters gets the parameters from the request
func (d *DefaultHandler) GetParameters(
	w http.ResponseWriter,
	r *http.Request,
	keys ...string,
) (*map[string]string, bool) {
	// Check if the request is nil
	if r == nil {
		return nil, false
	}

	// Initialize the parameters map
	parameters := make(map[string]string)

	// Get the URL query
	query := r.URL.Query()

	// Iterate over the keys
	for _, key := range keys {
		// Get the parameter value
		value := query.Get(key)

		// Check if the value was not found
		if value == "" {
			// Handle the error
			d.HandleResponse(
				w,
				gonethttpstatusresponse.NewJSendDebugInternalServerError(
					fmt.Errorf(ErrParameterNotFound, key),
					ErrCodeParameterNotFound,
				),
			)
			return nil, false
		}

		// Add the parameter to the map
		parameters[key] = value
	}
	return &parameters, true
}

// ParseParameter parses the parameter
func (d *DefaultHandler) ParseParameter(
	w http.ResponseWriter,
	parameter string,
	dest interface{},
	toTypeFn func(parameter string, dest interface{}) error,
) bool {
	// Parse the parameter
	if err := toTypeFn(parameter, dest); err != nil {
		// Handle the error
		d.HandleResponse(
			w,
			gonethttpstatusresponse.NewJSendDebugInternalServerError(
				err,
				ErrCodeParameterParsingFailed,
			),
		)
		return false
	}
	return true
}

// HandleResponse handles the response
func (d *DefaultHandler) HandleResponse(
	w http.ResponseWriter,
	response gonethttpresponse.Response,
) {
	// Check if the response is nil
	if response == nil {
		d.HandleResponse(
			w,
			gonethttpstatusresponse.NewJSendDebugInternalServerError(
				gonethttpresponse.ErrNilResponse,
				ErrCodeNilResponse,
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
			w, gonethttpresponse.NewResponseFromRequestError(e),
		)
		return
	}

	d.HandleResponse(
		w,
		gonethttpstatusresponse.NewJSendDebugInternalServerError(
			err,
			ErrCodeRequestFatalError,
		),
	)
}
