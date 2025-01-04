package json

import (
	"errors"
)

var (
	ErrNilJSONData = errors.New("json data is nil")
	ErrNilEncoder  = errors.New("json encoder is nil")
	ErrNilDecoder  = errors.New("json decoder is nil")
)
