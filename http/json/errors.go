package json

import (
	"errors"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	"net/http"
)

var (
	ErrCodeInvalidContentType *string
)

var (
	ErrNilEncoder          = errors.New("json encoder is nil")
	ErrNilDecoder          = errors.New("json decoder is nil")
	ErrUnmarshalBodyFailed = errors.New("failed to unmarshal json body")
	ErrInvalidContentType  = gonethttpresponse.NewHeaderError(
		"Content-Type",
		"invalid content type, expected application/json",
		ErrCodeInvalidContentType,
		http.StatusUnsupportedMediaType,
	)
	ErrMaxBodySizeExceeded = "json body size exceeds the maximum allowed size, limit is %d bytes"
	ErrSyntaxError         = "json body contains badly-formed JSON at position %d"
	ErrUnexpectedEOF       = errors.New("json body contains badly-formed JSON")
	ErrEmptyBody           = errors.New("json body is empty")
	ErrUnknownField        = "json body contains an unknown field %s"
)
