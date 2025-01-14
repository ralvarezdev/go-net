package json

import (
	"errors"
)

var (
	ErrNilEncoder              = errors.New("json encoder is nil")
	ErrNilDecoder              = errors.New("json decoder is nil")
	ErrUnmarshalBodyDataFailed = errors.New("failed to unmarshal json body data")
	ErrInvalidFieldValueType   = "invalid field value type, expected: %v"
)
