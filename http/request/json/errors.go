package json

import (
	"errors"
)

var (
	ErrCodeInvalidContentType         string
	ErrCodeNilDestination             string
	ErrCodeFailedToReadBody           string
	ErrCodeUnmarshalRequestBodyFailed string
	ErrCodeSyntaxError                string
	ErrCodeUnmarshalTypeError         string
	ErrCodeUnknownField               string
	ErrCodeEmptyBody                  string
	ErrCodeMaxBodySizeExceeded        string
)

const (
	ErrMaxBodySizeExceeded = "json body size exceeds the maximum allowed size, limit is %d bytes"
	ErrSyntaxError         = "json body contains badly-formed JSON at position %d"
	ErrUnknownField        = "json body contains an unknown field %s"
)

const (
	ErrInvalidContentTypeField = "Content-Type"
)

var (
	ErrUnexpectedEOF       = errors.New("json body contains badly-formed JSON")
	ErrEmptyBody           = errors.New("json body is empty")
	ErrUnmarshalBodyFailed = errors.New("failed to unmarshal json body")
)

var (
	ErrNilDestination     = errors.New("json destination is nil")
	ErrInvalidContentType = errors.New("invalid content type, expected application/json")
)
