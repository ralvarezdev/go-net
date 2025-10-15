package handler

import (
	"errors"
)

var (
	ErrNilHandler          = errors.New("requests and responses handler cannot be nil")
	ErrNilRequestsHandler  = errors.New("requests handler cannot be nil")
	ErrNilRawErrorHandler  = errors.New("raw error handler cannot be nil")
	ErrNilResponsesHandler = errors.New("responses handler cannot be nil")
)
