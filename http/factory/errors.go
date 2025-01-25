package factory

import (
	"errors"
)

var (
	ErrNilService    = errors.New("route service cannot be nil")
	ErrNilValidator  = errors.New("route validator cannot be nil")
	ErrNilController = errors.New("route controller cannot be nil")
)
