package request

import (
	"errors"
)

var (
	ErrCodeNilRequest                 string
	ErrCodeInvalidContentType         string
	ErrCodeInvalidBodyType            string
	ErrCodeUnmarshalRequestBodyFailed string
	ErrCodeSyntaxError                string
	ErrCodeUnmarshalTypeError         string
	ErrCodeUnknownField               string
	ErrCodeEmptyBody                  string
	ErrCodeMaxBodySizeExceeded        string
	ErrCodeNilDecoder                 string
)

const (
	ErrInvalidContentTypeField = "Content-Type"
	ErrMaxBodySizeExceeded     = "json body size exceeds the maximum allowed size, limit is %d bytes"
	ErrSyntaxError             = "json body contains badly-formed JSON at position %d"
	ErrUnknownField            = "json body contains an unknown field %s"
)

var (
	ErrNilRequest          = errors.New("request cannot be nil")
	ErrInvalidContentType  = errors.New("invalid content type, expected application/json")
	ErrInvalidBodyType     = errors.New("invalid body type, expected struct")
	ErrUnexpectedEOF       = errors.New("json body contains badly-formed JSON")
	ErrEmptyBody           = errors.New("json body is empty")
	ErrUnmarshalBodyFailed = errors.New("failed to unmarshal json body")
	ErrNilDecoder          = errors.New("decoder cannot be nil")
)
