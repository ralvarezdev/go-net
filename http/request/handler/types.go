package handler

import (
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gojsondecoder "github.com/ralvarezdev/go-json/decoder"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/mapper/validator"
)

type (
	// DefaultRequestsHandler struct
	DefaultRequestsHandler struct {
		mode             *goflagsmode.Flag
		responsesHandler gonethttphandler.ResponsesHandler
		gonethttprequest.Decoder
	}
)

// NewDefaultRequestsHandler creates a new default request handler
//
// Parameters:
//
//   - mode: The flag mode
//   - decoder: The HTTP request decoder
//   - responsesHandler: The HTTP response handler
//
// Returns:
//
//   - *DefaultRequestsHandler: The default handler
//   - error: The error if any
func NewDefaultRequestsHandler(
	mode *goflagsmode.Flag,
	decoder gonethttprequest.Decoder,
	responsesHandler gonethttphandler.ResponsesHandler,
) (*DefaultRequestsHandler, error) {
	// Check if the flag mode, the decoder or the handler is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	if decoder == nil {
		return nil, gojsondecoder.ErrNilDecoder
	}
	if responsesHandler == nil {
		return nil, gonethttphandler.ErrNilResponsesHandler
	}

	return &DefaultRequestsHandler{
		mode,
		responsesHandler,
		decoder,
	}, nil
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
func (d DefaultRequestsHandler) Validate(
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
		d.responsesHandler.HandleDebugErrorWithCode(
			w,
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeValidationFailed,
			http.StatusInternalServerError,
		)
		return false
	}

	d.responsesHandler.HandleFailDataErrorWithCode(
		w,
		validations,
		ErrCodeValidationFailed,
		http.StatusBadRequest,
	)
	return false
}

// DecodeAndValidate decodes and validates the request body
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
func (d DefaultRequestsHandler) DecodeAndValidate(
	w http.ResponseWriter,
	r *http.Request,
	dest interface{},
	validatorFn govalidatormappervalidator.ValidateFn,
) bool {
	// Decode the request body
	if err := d.DecodeRequest(r, dest); err != nil {
		// Handle the error
		d.responsesHandler.HandleRawError(
			w,
			err,
		)
		return false
	}

	// Validate the request body
	return d.Validate(w, dest, validatorFn)
}
