package json

import (
	"encoding/json"
	"errors"
	"fmt"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
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
				gonethttpresponse.NewFailResponseFromRequestError(
					gonethttpresponse.NewFieldError(
						fieldName,
						fmt.Sprintf(
							ErrInvalidFieldValueType,
							fieldTypeName,
						),
						http.StatusBadRequest,
					),
				),
			)
		}
	}

	return encoder.Encode(
		w,
		gonethttpresponse.NewDebugFailResponse(
			ErrUnmarshalBodyDataFailed,
			err,
			nil,
			http.StatusBadRequest,
		),
	)
}
