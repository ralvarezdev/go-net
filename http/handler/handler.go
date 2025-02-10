package handler

import (
	"errors"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatusresponse "github.com/ralvarezdev/go-net/http/status/response"
	"net/http"
)

var (
	ErrCodeValidationFailed      *string
	ErrCodeWildcardParsingFailed *string
	ErrCodeNilResponse           *string
	ErrCodeRequestFatalError     *string
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
		Parse(
			w http.ResponseWriter,
			r *http.Request,
			dest interface{},
			validatorFn func(body interface{}) (interface{}, error),
		) bool
		GetParameters(
			r *http.Request,
			keys ...string,
		) *map[string]string
		ParseWildcard(
			w http.ResponseWriter,
			r *http.Request,
			wildcardKey string,
			dest interface{},
			toTypeFn func(wildcard string, dest interface{}) error,
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
		mode,
		jsonEncoder,
		jsonDecoder,
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

// Parse decodes and validates the request body
func (d *DefaultHandler) Parse(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
	validatorFn func(body interface{}) (interface{}, error),
) bool {
	// Decode the request body
	if err := d.Decode(w, r, dest); err != nil {
		return false
	}

	// Validate the request body
	return d.Validate(w, dest, validatorFn)
}

// GetParameters gets the parameters from the request
func (d *DefaultHandler) GetParameters(
	r *http.Request,
	keys ...string,
) *map[string]string {
	// Check if the request is nil
	if r == nil {
		return nil
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
			continue
		}

		// Add the parameter to the map
		parameters[key] = value
	}
	return &parameters
}

// ParseWildcard parses the wildcard from the request and stores it in the destination
func (d *DefaultHandler) ParseWildcard(
	w http.ResponseWriter,
	r *http.Request,
	wildcardKey string,
	dest interface{},
	toTypeFn func(wildcard string, dest interface{}) error,
) bool {
	// Get the wildcard from the request
	wildcardValue := r.PathValue(wildcardKey)

	// Parse the wildcard value
	if err := toTypeFn(wildcardValue, dest); err != nil {
		// Handle the error
		d.HandleResponse(
			w,
			gonethttpstatusresponse.NewJSendDebugBadRequest(
				err,
				ErrCodeWildcardParsingFailed,
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
	// Check if the errors is a fail body error or a fail request error
	var failResponseErrorTarget *gonethttpresponse.FailResponseError
	if errors.As(err, &failResponseErrorTarget) {
		d.HandleResponse(
			w,
			failResponseErrorTarget.Response(),
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
