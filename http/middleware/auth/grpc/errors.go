package grpc

import (
	"errors"
)

var (
	ErrCodeInterceptionNotFound string
)

var (
	ErrNilInterceptions     = errors.New("nil interceptions")
	ErrInterceptionNotFound = errors.New("interception not found")
)
