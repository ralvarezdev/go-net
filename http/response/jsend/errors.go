package jsend

import (
	"errors"
)

var (
	ErrCodeExpectedStringSliceOnNestedDataMap string
	ErrCodeExpectedStringSliceOnDataMap       string
	ErrCodeExpectedMapOnNestedDataMap         string
	ErrCodeGRPCCtxCanceled                    string
	ErrCodeGRPCCtxDeadlineExceeded            string
	ErrCodeGRPCUnknown                        string
	ErrCodeGRPCBadRequest                     string
	ErrCodeGRPCPreconditionFailed             string
	ErrCodeGRPCQuotaFailure                   string
	ErrCodeGRPCRequestInfo                    string
	ErrCodeGRPCHelp                           string
	ErrCodeGRPCResourceInfo                   string
	ErrCodeGRPCLocalizedMessage               string
)

var (
	ErrExpectedStringSliceOnNestedDataMap = errors.New("expected a string slice on nested data map")
	ErrExpectedStringSliceOnDataMap       = errors.New("expected a string slice on data map")
	ErrExpectedMapOnNestedDataMap         = errors.New("expected a map on nested data map")
	ErrGRPCPreconditionFailed             = errors.New("grpc precondition failed")
)
