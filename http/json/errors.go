package json

import (
	"errors"
)

var (
	ErrNilJSONData            = errors.New("json data is nil")
	ErrJSONDataMustBeAPointer = errors.New("json data must be a pointer")
	ErrNilJSONEncoder         = errors.New("json encoder is nil")
	ErrNilJSONDecoder         = errors.New("json decoder is nil")
)
