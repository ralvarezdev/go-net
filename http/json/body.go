package json

import (
	"encoding/json"
	"errors"
	"fmt"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"io"
	"net/http"
	"strings"
)

// Inspired by:
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

var (
	ErrCodeUnmarshalRequestBodyFailed *string
	ErrCodeSyntaxError                *string
	ErrCodeUnmarshalTypeError         *string
	ErrCodeUnknownField               *string
	ErrCodeEmptyBody                  *string
	ErrCodeMaxBodySizeExceeded        *string
)

// NewUnmarshalTypeErrorResponse creates a new response for an UnmarshalTypeError
func NewUnmarshalTypeErrorResponse(
	fieldName string,
	fieldTypeName string,
) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		gonethttpresponse.NewBodyFromFailBodyError(
			gonethttpresponse.NewFailBodyError(
				fieldName,
				fmt.Sprintf(
					gonethttpresponse.ErrInvalidFieldValueType,
					fieldTypeName,
				),
				ErrCodeUnmarshalTypeError,
			),
		),
		http.StatusBadRequest,
	)
}

// NewSyntaxErrorResponse creates a new response for a SyntaxError
func NewSyntaxErrorResponse(
	debugErr error,
	offset int64,
) gonethttpresponse.Response {
	// Create the error
	err := fmt.Errorf(ErrSyntaxError, offset)

	return gonethttpresponse.NewJSendErrorResponse(
		nil,
		err.Error(),
		ErrCodeSyntaxError,
		http.StatusBadRequest,
	)
}

// NewUnknownFieldErrorResponse creates a new response for an unknown field error
func NewUnknownFieldErrorResponse(fieldName string) gonethttpresponse.Response {
	return gonethttpresponse.NewResponse(
		gonethttpresponse.NewBodyFromFailBodyError(
			gonethttpresponse.NewFailBodyError(
				fieldName,
				fmt.Sprintf(ErrUnknownField, fieldName),
				ErrCodeUnknownField,
			),
		),
		http.StatusBadRequest,
	)
}

// NewMaxBodySizeExceededErrorResponse creates a new response for a body size exceeded error
func NewMaxBodySizeExceededErrorResponse(limit int64) gonethttpresponse.Response {
	// Create the error
	err := fmt.Errorf(ErrMaxBodySizeExceeded, limit)

	return gonethttpresponse.NewJSendErrorResponse(
		nil,
		err.Error(),
		ErrCodeMaxBodySizeExceeded,
		http.StatusRequestEntityTooLarge,
	)
}

// BodyDecodeErrorHandler handles the error on JSON body decoding
func BodyDecodeErrorHandler(
	w http.ResponseWriter,
	err error,
	encoder Encoder,
) error {
	// Check if the encoder is nil
	if encoder == nil {
		return ErrNilEncoder
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
			NewSyntaxErrorResponse(err, syntaxError.Offset),
		)
	}

	// Check if the error is an ErrUnexpectedEOF
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return encoder.Encode(
			w,
			gonethttpresponse.NewJSendErrorResponse(
				nil,
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
				nil,
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
			nil,
			ErrUnmarshalBodyFailed.Error(),
			err.Error(),
			ErrCodeUnmarshalRequestBodyFailed,
			http.StatusBadRequest,
		),
	)
}
