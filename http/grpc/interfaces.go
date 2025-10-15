package grpc

import (
	"context"
	"net/http"
)

type (
	// AuthenticationParser parses the metadata authentication from a gRPC to the header or as a cookie
	AuthenticationParser interface {
		ParseAuthorizationMetadataAsHeader(
			w http.ResponseWriter,
			ctx context.Context,
		) error
		ParseAuthorizationMetadataAsCookie(
			w http.ResponseWriter,
			ctx context.Context,
		) error
	}
)
