package handler

import (
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpparser "github.com/ralvarezdev/go-net/http/parser"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponsejsend "github.com/ralvarezdev/go-net/http/response/jsend"
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/mapper/validator"
)

type (
	// DefaultHandler struct
	DefaultHandler struct {
		mode *goflagsmode.Flag
		gonethttpresponse.Decoder
	}
)

// NewHandler creates a new default request handler
//
// Parameters:
//
//   - mode: The flag mode
//   - parser: The HTTP request parser
//
// Returns:
//
//   - *Handler: The default handler
//   - error: The error if any
func NewHandler(
	mode *goflagsmode.Flag,
	parser gonethttpparser.Parser,
) (*Handler, error) {
	// Check if the flag mode or the parser is nil
	if mode == nil {
		return nil, goflagsmode.ErrNilModeFlag
	}
	if parser == nil {
		return nil, gonethttpparser.ErrNilParser
	}

	return &Handler{
		mode,
		parser,
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
func (d Handler) Validate(
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
		d.HandleDebugErrorResponseWithCode(
			w,
			err,
			gonethttp.ErrInternalServerError,
			ErrCodeValidationFailed,
			http.StatusInternalServerError,
		)
		return false
	}

	d.HandleResponse(
		w,
		gonethttpresponsejsend.NewFailResponseWithCode(
			validations,
			ErrCodeValidationFailed,
			http.StatusBadRequest,
		),
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
func (d Handler) DecodeAndValidate(
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
