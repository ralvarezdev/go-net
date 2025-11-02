package grpc

import (
	"net/http"

	"google.golang.org/grpc/metadata"
)

type (
	// AuthenticationParser parses the metadata authentication from a gRPC to the header or as a cookie
	AuthenticationParser interface {
		ParseAuthorizationMetadataAsHeader(
			md metadata.MD,
			w http.ResponseWriter,
		) error
		ParseAuthorizationMetadataAsCookie(
			md metadata.MD,
			w http.ResponseWriter,
		) error
	}
)
