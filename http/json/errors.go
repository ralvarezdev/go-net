package json

import (
	"errors"
	"net/http"

	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

const (
	ErrMaxBodySizeExceeded = "json body size exceeds the maximum allowed size, limit is %d bytes"
	ErrSyntaxError         = "json body contains badly-formed JSON at position %d"
	ErrUnknownField        = "json body contains an unknown field %s"
)

var (
	ErrCodeInvalidContentType         *string
	ErrCodeFailedToReadBody           *string
	ErrCodeNilDestination             *string
	ErrCodeMarshalResponseBodyFailed  *string
	ErrCodeUnmarshalRequestBodyFailed *string
	ErrCodeSyntaxError                *string
	ErrCodeUnmarshalTypeError         *string
	ErrCodeUnknownField               *string
	ErrCodeEmptyBody                  *string
	ErrCodeMaxBodySizeExceeded        *string
)

var (
	ErrNilEncoder          = errors.New("json encoder is nil")
	ErrNilDecoder          = errors.New("json decoder is nil")
	ErrUnmarshalBodyFailed = errors.New("failed to unmarshal json body")
	ErrUnexpectedEOF       = errors.New("json body contains badly-formed JSON")
	ErrEmptyBody           = errors.New("json body is empty")
	ErrInvalidContentType  = gonethttpresponse.NewFailResponseError(
		"Content-Type",
		"invalid content type, expected application/json",
		ErrCodeInvalidContentType,
		http.StatusUnsupportedMediaType,
	)
)
