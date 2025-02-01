package route

import (
	"errors"
)

var (
	ErrNilRouter     = errors.New("router cannot be nil")
	ErrNilMiddleware = "%s: middleware at index %d cannot be nil"
	ErrEmptyPattern  = errors.New("pattern cannot be empty")
)
