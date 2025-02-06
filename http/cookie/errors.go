package cookie

import (
	"errors"
)

var (
	ErrNilRequest    = errors.New("nil request")
	ErrNilAttributes = errors.New("nil cookie attributes")
)
