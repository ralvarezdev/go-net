package grpc

import (
	"context"
	"net/http"
	"time"

	gogrpcmd "github.com/ralvarezdev/go-grpc/metadata"
	"google.golang.org/grpc/metadata"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpcookie "github.com/ralvarezdev/go-net/http/cookie"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// DefaultAuthenticationParser is the default implementation of AuthenticationParser
	DefaultAuthenticationParser struct {
		options *Options
	}

	// Options is the options for the DefaultAuthenticationParser
	Options struct {
		MetadataKeysToAuthorizationHeaderNames map[string]string
		MetadataKeysToCookiesAttributes        map[string]*gonethttpcookie.Attributes
		GetExpiresAtFn                         GetExpiresAtFn
	}

	// GetExpiresAtFn is a function that returns the expiration time for a token
	GetExpiresAtFn func(metadataValue string) (time.Time, error)
)

// NewDefaultAuthenticationParser parses the metadata authentication from a gRPC to the header or as a cookie
//
// Parameters:
//
// - options: the options for the DefaultAuthenticationParser
//
// Returns:
//
// - ParseAuthorizationMetadataAsHeader: parses the authorization metadata as a header
// - error: error if something goes wrong
func NewDefaultAuthenticationParser(
	options *Options,
) (
	*DefaultAuthenticationParser,
	error,
) {
	// Check if the options are nil
	if options == nil {
		return nil, ErrNilOptions
	}

	return &DefaultAuthenticationParser{
		options: options,
	}, nil
}

// ParseAuthorizationMetadataAsHeader parses the authorization metadata as a header
//
// Parameters:
//
// - ctx: context.Context
// - w: http.ResponseWriter
//
// Returns:
//
// - error: error if something goes wrong
func (d DefaultAuthenticationParser) ParseAuthorizationMetadataAsHeader(
	ctx context.Context,
	w http.ResponseWriter,
) error {
	// Get the metadata from the context
	md, err := gogrpcmd.GetIncomingCtxMetadata(ctx)
	if err != nil {
		return gonethttpresponse.NewDebugError(
			err,
			gonethttp.ErrInternalServerError,
			http.StatusInternalServerError,
		)
	}

	// Get the authorization metadata from the context
	authorization, err := gogrpcmd.GetMetadataAuthorizationToken(md)
	if err != nil {
		return err
	}

	// Set the authorization header
	w.Header().Set(gonethttp.Authorization, authorization)

	// Iterate over the metadata keys to authorization header names
	for metadataKey, headerName := range d.options.MetadataKeysToAuthorizationHeaderNames {
		// Get the metadata value
		if metadataValueSlice, ok := md[metadataKey]; ok {
			// Check if the metadata value is empty
			if len(metadataValueSlice) == 0 {
				continue
			}

			// Get the first value of the metadata
			metadataValue := metadataValueSlice[0]

			// Set the header
			w.Header().Set(headerName, metadataValue)
		}
	}
	return nil
}

// ParseAuthorizationMetadataAsCookie parses the authorization metadata as a cookie
//
// Parameters:
//
// - ctx: context.Context
// - w: http.ResponseWriter
//
// Returns:
//
// - error: error if something goes wrong
func (d DefaultAuthenticationParser) ParseAuthorizationMetadataAsCookie(
	md metadata.MD,
	w http.ResponseWriter,
) error {
	// Iterate over the metadata keys to cookies attributes
	for metadataKey, cookieAttributes := range d.options.MetadataKeysToCookiesAttributes {
		// Get the metadata value
		metadataValueSlice, ok := md[metadataKey]
		if !ok || len(metadataValueSlice) == 0 {
			continue
		}

		// Get the first value of the metadata
		metadataValue := metadataValueSlice[0]

		// Get the expiration time if the function is set
		var (
			expiresAt time.Time
			err       error
		)
		if d.options.GetExpiresAtFn != nil {
			expiresAt, err = d.options.GetExpiresAtFn(metadataValue)
			if err != nil {
				return err
			}
		}

		// Set the cookie
		gonethttpcookie.SetCookie(
			w,
			cookieAttributes,
			metadataValue,
			expiresAt,
		)
	}
	return nil
}
