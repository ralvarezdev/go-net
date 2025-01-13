package json

import (
	"errors"
)

var (
	ErrNilData                 = errors.New("json data is nil")
	ErrNilEncoder              = errors.New("json encoder is nil")
	ErrNilDecoder              = errors.New("json decoder is nil")
	ErrUnmarshalBodyDataFailed = errors.New("failed to unmarshal json body data")
	ErrFieldInvalidValue       = "field has an invalid value %v, it must be %v"
)
