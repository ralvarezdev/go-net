package json

import (
	"encoding/json"
	"errors"
	"fmt"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

var (
	ErrCodeUnmarshalRequestBodyFailed *string
	ErrCodeRequestBodyFieldError      *string
)

// bodyDecodeErrorHandler handles the error on JSON body decoding
func bodyDecodeErrorHandler(
	w http.ResponseWriter,
	err error,
	encoder Encoder,
) error {
	// Check if the encoder is nil
	if encoder == nil {
		return ErrNilEncoder
	}

	// Check is there is an UnmarshalTypeError
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		// Check which field failed
		fieldName := unmarshalTypeError.Field
		fieldTypeName := unmarshalTypeError.Type.Name()

		// Check if the field name is empty
		if fieldName != "" {
			return encoder.Encode(
				w,
				gonethttpresponse.NewResponseFromRequestError(
					gonethttpresponse.NewFieldError(
						fieldName,
						fmt.Sprintf(
							ErrInvalidFieldValueType,
							fieldTypeName,
						),
						http.StatusBadRequest,
						ErrCodeRequestBodyFieldError,
					),
				),
			)
		}
	}

	return encoder.Encode(
		w,
		gonethttpresponse.NewJSendErrorResponse(
			ErrUnmarshalBodyFailed,
			err,
			nil,
			ErrCodeUnmarshalRequestBodyFailed,
			http.StatusBadRequest,
		),
	)
}
