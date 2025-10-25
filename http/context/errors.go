package context

import (
	"errors"
)

var (
	ErrMissingBodyInContext = errors.New("missing body in context")
	ErrInvalidBodyType     = errors.New("invalid body type in context")
)
