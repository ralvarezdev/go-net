package handler

import (
	"errors"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/mapper/validator"
)

type (
	// DefaultHandler struct
	DefaultHandler struct {
		mode    *goflagsmode.Flag
		encoder gonethttpresponse.Encoder
		decoder gonethttprequest.Decoder
	}
)

// NewDefaultHandler creates a new default request handler
//
// Parameters:
//
//   - mode: The flag mode
//   - encoder: The encoder
//   - decoder: The decoder
//
// Returns:
//
//   - *DefaultHandler: The default handler
//   - error: The error if any
func NewDefaultHandler(
	mode *goflagsmode.Flag,
	encoder gonethttpresponse.Encoder,
	decoder gonethttprequest.Decoder,
) (*DefaultHandler, error) {
	// Check if the flag mode, the JSON encoder or the JSON decoder is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	if encoder == nil {
		return nil, gonethttpresponse.ErrNilEncoder
	}
	if decoder == nil {
		return nil, gonethttprequest.ErrNilDecoder
	}

	return &DefaultHandler{
		mode,
		encoder,
		decoder,
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
	return d.decoder.Decode(w, r, dest)
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
	validatorFn govalidatormappervalidator.ValidateFn,
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
			gonethttpresponse.NewJSendDebugInternalServerError(
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
	validatorFn govalidatormappervalidator.ValidateFn,
) bool {
	// Decode the request body
	if err := d.Decode(w, r, dest); err != nil {
		return false
	}

	// Validate the request body
	return d.Validate(w, dest, validatorFn)
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
			gonethttpresponse.NewJSendDebugInternalServerError(
				gonethttpresponse.ErrNilResponse,
				ErrCodeNilResponse,
			),
		)
		return
	}

	// Call the encoder
	_ = d.encoder.Encode(w, response)
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
		gonethttpresponse.NewJSendDebugInternalServerError(
			err,
			ErrCodeRequestFatalError,
		),
	)
}
