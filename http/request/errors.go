package request

import (
	"errors"
)

var (
	ErrNilDecoder = errors.New("decoder is nil")
	ErrNilRequest = errors.New("request cannot be nil")
)
