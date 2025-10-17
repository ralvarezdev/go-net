package protojson

import (
	"errors"
)

var (
	ErrCodeReadBodyFailed           string
	ErrCodeInvalidProtoMessage      string
	ErrCodeUnmarshalProtoJSONFailed string
)

const (
	ErrUnmarshalProtoJSONFailed = "failed to unmarshal proto JSON: %v"
	ErrReadBodyFailed           = "failed to read body: %v"
)

var (
	ErrInvalidProtoMessage = errors.New("invalid proto message")
	ErrNilMarshalOptions   = errors.New("nil marshal options")
)
