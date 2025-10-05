package handler

import (
	"errors"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttpjson "github.com/ralvarezdev/go-net/http/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpstatusresponse "github.com/ralvarezdev/go-net/http/status/response"
	govalidatorstructmappervalidator "github.com/ralvarezdev/go-validator/struct/mapper/validator"
)

type (
	// DefaultHandler struct
	DefaultHandler struct {
		mode        *goflagsmode.Flag
		jsonEncoder gonethttpjson.Encoder
		jsonDecoder gonethttpjson.Decoder
	}
)

// NewDefaultHandler creates a new default request handler
//
// Parameters:
//
//   - mode: The flag mode
//   - jsonEncoder: The JSON encoder
//   - jsonDecoder: The JSON decoder
//
// Returns:
//
//   - *DefaultHandler: The default handler
//   - error: The error if any
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
//
// Parameters:
//
//   - w: The HTTP response writer
//   - r: The HTTP request
//   - dest: The destination to decode the request body into
//
// Returns:
//
//   - error: The error if any
func (d DefaultHandler) Decode(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
) error {
	return d.jsonDecoder.Decode(w, r, dest)
}

// Validate validates the request body
//
// Parameters:
//
//   - w: The HTTP response writer
//   - body: The request body to validate
//   - validatorFn: The validator function
//
// Returns:
//
//   - bool: True if the request body is valid, false otherwise
func (d DefaultHandler) Validate(
	w http.ResponseWriter,
	body interface{},
	validatorFn govalidatorstructmappervalidator.ValidateFn,
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
		return false
	}

	d.HandleResponse(
		w,
		gonethttpresponse.NewJSendFailResponse(
			validations,
			ErrCodeValidationFailed,
			http.StatusBadRequest,
		),
	)
	return false
}

// Parse decodes and validates the request body
//
// Parameters:
//
//   - w: The HTTP response writer
//   - r: The HTTP request
//   - dest: The destination to decode the request body into
//   - validatorFn: The validator function
//
// Returns:
//
//   - bool: True if the request body is valid, false otherwise
func (d DefaultHandler) Parse(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
	validatorFn govalidatorstructmappervalidator.ValidateFn,
) bool {
	// Decode the request body
	if err := d.Decode(w, r, dest); err != nil {
		return false
	}

	// Validate the request body
	return d.Validate(w, dest, validatorFn)
}

// GetParameters gets the parameters from the request
//
// Parameters:
//
//   - r: The HTTP request
//   - keys: The keys of the parameters to get
//
// Returns:
//
//   - map[string]string: The parameters map
func (d DefaultHandler) GetParameters(
	r *http.Request,
	keys ...string,
) map[string]string {
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
	return parameters
}

// ParseWildcard parses the wildcard from the request and stores it in the destination
//
// Parameters:
//
//   - w: The HTTP response writer
//   - r: The HTTP request
//   - wildcardKey: The key of the wildcard to parse
//   - dest: The destination to store the parsed wildcard
//   - toTypeFn: The function to convert the wildcard value to the desired type
//
// Returns:
//
//   - bool: True if the wildcard was parsed successfully, false otherwise
func (d DefaultHandler) ParseWildcard(
	w http.ResponseWriter,
	r *http.Request,
	wildcardKey string,
	dest interface{},
	toTypeFn ToTypeFn,
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
//
// Parameters:
//
//   - w: The HTTP response writer
//   - response: The response to handle
func (d DefaultHandler) HandleResponse(
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
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error to handle
func (d DefaultHandler) HandleError(
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
