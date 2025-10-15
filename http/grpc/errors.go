package grpc

import (
	"errors"
)

var (
	ErrNilContext = errors.New("context cannot be nil")
	ErrNilOptions = errors.New("options cannot be nil")
)
