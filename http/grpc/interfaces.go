package grpc

import (
	"context"
	"net/http"
)

type (
	// AuthenticationParser parses the metadata authentication from a gRPC to the header or as a cookie
	AuthenticationParser interface {
		ParseAuthorizationMetadataAsHeader(
			ctx context.Context,
			w http.ResponseWriter,
		) error
		ParseAuthorizationMetadataAsCookie(
			ctx context.Context,
			w http.ResponseWriter,
		) error
	}
)
