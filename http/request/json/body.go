package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

// Inspired by:
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

// NewUnmarshalTypeErrorResponse creates a new response for an UnmarshalTypeError
//
// Parameters:
//
//   - fieldName: The name of the field that caused the error
//   - fieldTypeName: The type name of the field that caused the error
func NewUnmarshalTypeErrorResponse(
	fieldName string,
	fieldTypeName string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		gonethttpresponse.NewFailBodyError(
			fieldName,
			fmt.Sprintf(
				gonethttpresponse.ErrInvalidFieldValueType,
				fieldTypeName,
			),
			ErrCodeUnmarshalTypeError,
		).Body(),
		http.StatusBadRequest,
	)
}

// NewSyntaxErrorResponse creates a new response for a SyntaxError
//
// Parameters:
//
//   - offset: The offset where the error occurred
func NewSyntaxErrorResponse(
	offset int64,
) gonethttpresponse.Response {
	// Create the error
	err := fmt.Errorf(ErrSyntaxError, offset)

	return gonethttpresponse.NewJSendErrorResponse(
		err.Error(),
		ErrCodeSyntaxError,
		http.StatusBadRequest,
	)
}

// NewUnknownFieldErrorResponse creates a new response for an unknown field error
//
// Parameters:
//
//   - fieldName: The name of the unknown field
//
// Returns:
//
//   - gonethttpresponse.Response: The response
func NewUnknownFieldErrorResponse(fieldName string) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		gonethttpresponse.NewFailBodyError(
			fieldName,
			fmt.Sprintf(ErrUnknownField, fieldName),
			ErrCodeUnknownField,
		).Body(),
		http.StatusBadRequest,
	)
}

// NewMaxBodySizeExceededErrorResponse creates a new response for a body size exceeded error
//
// Parameters:
//
//   - limit: The maximum allowed body size
//
// Returns:
//
//   - gonethttpresponse.Response: The response
func NewMaxBodySizeExceededErrorResponse(limit int64) gonethttpresponse.Response {
	// Create the error
	err := fmt.Errorf(ErrMaxBodySizeExceeded, limit)

	return gonethttpresponse.NewJSendErrorResponse(
		err.Error(),
		ErrCodeMaxBodySizeExceeded,
		http.StatusRequestEntityTooLarge,
	)
}

// BodyDecodeErrorHandler handles the error on JSON body decoding
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error that occurred during decoding
//   - encoder: The encoder to use for the response
//
// Returns:
//
//   - error: An error if the encoder is nil or if encoding the response fails
func BodyDecodeErrorHandler(
	w http.ResponseWriter,
	err error,
	encoder gonethttpresponse.Encoder,
) error {
	// Check if the encoder is nil
	if encoder == nil {
		return gonethttpresponse.ErrNilEncoder
	}

	// Check is there is an UnmarshalTypeError
	var syntaxError *json.SyntaxError
	var maxBytesError *http.MaxBytesError
	var unmarshalTypeError *json.UnmarshalTypeError

	// Check if the error is an UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		// Check which field failed
		fieldName := unmarshalTypeError.Field
		fieldTypeName := unmarshalTypeError.Type.Name()

		// Check if the field name is empty
		if fieldName != "" {
			return encoder.Encode(
				w,
				NewUnmarshalTypeErrorResponse(fieldName, fieldTypeName),
			)
		}
	}

	// Check if the error is a SyntaxError
	if errors.As(err, &syntaxError) {
		return encoder.Encode(
			w,
			NewSyntaxErrorResponse(syntaxError.Offset),
		)
	}

	// Check if the error is an ErrUnexpectedEOF
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return encoder.Encode(
			w,
			gonethttpresponse.NewJSendErrorResponse(
				ErrUnexpectedEOF.Error(),
				ErrCodeSyntaxError,
				http.StatusBadRequest,
			),
		)
	}

	// Check if the error is an unknown field error
	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		// Get the field name
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")

		return encoder.Encode(
			w,
			NewUnknownFieldErrorResponse(fieldName),
		)
	}

	// Check if the error is caused by an empty request body
	if errors.Is(err, io.EOF) {
		return encoder.Encode(
			w,
			gonethttpresponse.NewJSendErrorResponse(
				ErrEmptyBody.Error(),
				ErrCodeEmptyBody,
				http.StatusBadRequest,
			),
		)
	}

	// Catch the error caused by the request body being too large
	if errors.As(err, &maxBytesError) {
		return encoder.Encode(
			w,
			NewMaxBodySizeExceededErrorResponse(maxBytesError.Limit),
		)
	}

	return encoder.Encode(
		w,
		gonethttpresponse.NewJSendErrorDebugResponse(
			ErrUnmarshalBodyFailed.Error(),
			err.Error(),
			ErrCodeUnmarshalRequestBodyFailed,
			http.StatusBadRequest,
		),
	)
}
