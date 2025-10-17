package request

import (
	"errors"
)

var (
	ErrCodeNilRequest         string
	ErrCodeInvalidContentType string
	ErrCodeInvalidBodyType    string
)

const (
	ErrInvalidContentTypeField = "Content-Type"
)

var (
	ErrNilDecoder         = errors.New("decoder is nil")
	ErrNilRequest         = errors.New("request cannot be nil")
	ErrInvalidContentType = errors.New("invalid content type, expected application/json")
	ErrInvalidBodyType    = errors.New("invalid body type, expected struct")
)
