package route

import (
	"errors"
)

var (
	ErrNilGroup      = errors.New("route group cannot be nil")
	ErrNilController = errors.New("route controller cannot be nil")
	ErrNilService    = errors.New("route service cannot be nil")
)
