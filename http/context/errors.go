package context

import (
	"errors"
)

var (
	ErrMissingBodyInContext = errors.New("missing body in context")
)
