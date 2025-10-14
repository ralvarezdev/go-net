package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

// Inspired by:
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

// BodyDecodeErrorHandler handles the error on JSON body decoding
//
// Parameters:
//
//   - w: The HTTP response writer
//   - err: The error that occurred during decoding
//   - handler: The handler to use for the response
//
// Returns:
//
//   - error: An error if the encoder is nil or if encoding the response fails
func BodyDecodeErrorHandler(
	w http.ResponseWriter,
	err error,
	handler gonethttphandler.Handler,
) error {
	// Check if the handler is nil
	if handler == nil {
		return gonethttphandler.ErrNilHandler
	}

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
			err = fmt.Errorf(
				gonethttpresponse.ErrInvalidFieldValueType,
				fieldTypeName,
			)
			handler.HandleFieldFailResponseWithCode(
				w,
				fieldName,
				err,
				ErrCodeUnmarshalTypeError,
				http.StatusBadRequest,
			)
			return err
		}
	}

	// Check if the error is a SyntaxError
	if errors.As(err, &syntaxError) {
		err = fmt.Errorf(ErrSyntaxError, syntaxError.Offset)
		handler.HandleErrorResponseWithCode(
			w,
			err,
			ErrCodeSyntaxError,
			http.StatusBadRequest,
		)
		return err
	}

	// Check if the error is an ErrUnexpectedEOF
	if errors.Is(err, io.ErrUnexpectedEOF) {
		handler.HandleErrorResponseWithCode(
			w,
			ErrUnexpectedEOF,
			ErrCodeSyntaxError,
			http.StatusBadRequest,
		)
		return ErrUnexpectedEOF
	}

	// Check if the error is an unknown field error
	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		// Get the field name
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")

		err = fmt.Errorf(ErrUnknownField, fieldName)
		handler.HandleFieldFailResponseWithCode(
			w,
			fieldName,
			err,
			ErrCodeUnknownField,
			http.StatusBadRequest,
		)
		return err
	}

	// Check if the error is caused by an empty request body
	if errors.Is(err, io.EOF) {
		handler.HandleErrorResponseWithCode(
			w,
			ErrEmptyBody,
			ErrCodeEmptyBody,
			http.StatusBadRequest,
		)
		return ErrEmptyBody
	}

	// Catch the error caused by the request body being too large
	if errors.As(err, &maxBytesError) {
		err = fmt.Errorf(ErrMaxBodySizeExceeded, maxBytesError.Limit)
		handler.HandleErrorResponseWithCode(
			w,
			err,
			ErrCodeMaxBodySizeExceeded,
			http.StatusRequestEntityTooLarge,
		)
	}

	handler.HandleDebugErrorResponseWithCode(
		w,
		err,
		ErrUnmarshalBodyFailed,
		ErrCodeUnmarshalRequestBodyFailed,
		http.StatusBadRequest,
	)
	return err
}
