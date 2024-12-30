package route

import (
	"errors"
)

var (
	ErrNilRouter     = errors.New("router cannot be nil")
	ErrNilController = errors.New("route controller cannot be nil")
	ErrNilService    = errors.New("route service cannot be nil")
)
