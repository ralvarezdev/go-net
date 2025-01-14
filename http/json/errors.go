package json

import (
	"errors"
)

var (
	ErrNilEncoder          = errors.New("json encoder is nil")
	ErrNilDecoder          = errors.New("json decoder is nil")
	ErrUnmarshalBodyFailed = errors.New("failed to unmarshal json body")
)
