package route

import (
	"errors"
)

const (
	ErrNilMiddleware      = "%s: middleware at index %d cannot be nil"
	ErrNilEndpointHandler = "endpoint handler cannot be nil, pattern: %s"
	ErrNilHandlerFunc     = "handler function cannot be nil, pattern: %s"
)

var (
	ErrNilRouter    = errors.New("router cannot be nil")
	ErrEmptyPattern = errors.New("pattern cannot be empty")
	ErrEmptyWildcard  = errors.New("wildcard cannot be empty")
	ErrWildcardNotClosed = errors.New("wildcard not closed")
)
