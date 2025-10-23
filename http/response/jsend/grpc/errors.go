package grpc

import (
	"errors"
)

var (
	ErrCodeExpectedStringSliceOnNestedDataMap                 string
	ErrCodeExpectedStringSliceOnDataMap                       string
	ErrCodeExpectedMapOnNestedDataMap                         string
	ErrCodeCtxCanceled                                        string
	ErrCodeCtxDeadlineExceeded                                string
	ErrCodeUnknown                                            string
	ErrCodeBadRequest                                         string
	ErrCodePreconditionFailure                                string
	ErrCodeQuotaFailure                                       string
	ErrCodeRequestInfo                                        string
	ErrCodeHelp                                               string
	ErrCodeResourceInfo                                       string
	ErrCodeLocalizedMessage                                   string
	ErrCodeCCodePrefix                                        string
	ErrCodeExpectedFieldViolationSliceOnDataMap               string
	ErrCodeExpectedPreconditionFailureViolationSliceOnDataMap string
	ErrCodeExpectedQuotaFailureViolationSliceOnDataMap        string
)

var (
	ErrExpectedStringSliceOnNestedDataMap                 = errors.New("expected a string slice on nested data map")
	ErrExpectedStringSliceOnDataMap                       = errors.New("expected a string slice on data map")
	ErrExpectedMapOnNestedDataMap                         = errors.New("expected a map on nested data map")
	ErrExpectedFieldViolationSliceOnDataMap               = errors.New("expected a field violation slice on data map")
	ErrExpectedPreconditionFailureViolationSliceOnDataMap = errors.New(
		"expected a precondition failure violation slice on data map",
	)
	ErrExpectedQuotaFailureViolationSliceOnDataMap = errors.New(
		"expected a quota failure violation slice on data map",
	)
)
