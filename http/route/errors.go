package route

import (
	"errors"
)

const (
	ErrNilMiddleware = "%s: middleware at index %d cannot be nil"
)

var (
	ErrNilRouter    = errors.New("router cannot be nil")
	ErrEmptyPattern = errors.New("pattern cannot be empty")
)
