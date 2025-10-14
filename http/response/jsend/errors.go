package jsend

import (
	"errors"
)

var (
	ErrCodeExpectedStringSliceOnNestedDataMap string
	ErrCodeExpectedStringSliceOnDataMap       string
	ErrCodeExpectedMapOnNestedDataMap         string
)

var (
	ErrExpectedStringSliceOnNestedDataMap = errors.New("expected a string slice on nested data map")
	ErrExpectedStringSliceOnDataMap       = errors.New("expected a string slice on data map")
	ErrExpectedMapOnNestedDataMap         = errors.New("expected a map on nested data map")
)
