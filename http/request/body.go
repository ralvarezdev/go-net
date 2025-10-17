package request

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

// BodyDecodeErrorHandler handles the error on JSON body decoding
//
// Parameters:
//
//   - err: The error that occurred during decoding
//
// Returns:
//
//   - error: An error if the encoder is nil or if encoding the response fails
func BodyDecodeErrorHandler(
	err error,
) error {
	// Check is there is an UnmarshalTypeError
	var (
		syntaxError        *json.SyntaxError
		maxBytesError      *http.MaxBytesError
		unmarshalTypeError *json.UnmarshalTypeError
	)

	// Check if the error is an UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		// Check which field failed
		fieldName := unmarshalTypeError.Field
		fieldTypeName := unmarshalTypeError.Type.Name()

		// Check if the field name is empty
		if fieldName != "" {
			return gonethttpresponse.NewFailFieldErrorWithCode(
				fieldName,
				fmt.Errorf(
					gonethttpresponse.ErrInvalidFieldValueType,
					fieldTypeName,
				),
				ErrCodeUnmarshalTypeError,
				http.StatusBadRequest,
			)
		}
	}

	// Check if the error is a SyntaxError
	if errors.As(err, &syntaxError) {
		return gonethttpresponse.NewErrorWithCode(
			fmt.Errorf(ErrSyntaxError, syntaxError.Offset),
			ErrCodeSyntaxError,
			http.StatusBadRequest,
		)
	}

	// Check if the error is an ErrUnexpectedEOF
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return gonethttpresponse.NewErrorWithCode(
			ErrUnexpectedEOF,
			ErrCodeSyntaxError,
			http.StatusBadRequest,
		)
	}

	// Check if the error is an unknown field error
	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		// Get the field name
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")

		return gonethttpresponse.NewFailFieldErrorWithCode(
			fieldName,
			fmt.Errorf(ErrUnknownField, fieldName),
			ErrCodeUnknownField,
			http.StatusBadRequest,
		)
	}

	// Check if the error is caused by an empty request body
	if errors.Is(err, io.EOF) {
		return gonethttpresponse.NewErrorWithCode(
			ErrEmptyBody,
			ErrCodeEmptyBody,
			http.StatusBadRequest,
		)
	}

	// Catch the error caused by the request body being too large
	if errors.As(err, &maxBytesError) {
		return gonethttpresponse.NewErrorWithCode(
			fmt.Errorf(ErrMaxBodySizeExceeded, maxBytesError.Limit),
			ErrCodeMaxBodySizeExceeded,
			http.StatusRequestEntityTooLarge,
		)
	}

	return gonethttpresponse.NewDebugErrorWithCode(
		err,
		ErrUnmarshalBodyFailed,
		ErrCodeUnmarshalRequestBodyFailed,
		http.StatusBadRequest,
	)
}
