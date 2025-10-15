package grpc

import (
	"context"
	"net/http"
	"time"

	gogrpcmd "github.com/ralvarezdev/go-grpc/metadata"
	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpcookie "github.com/ralvarezdev/go-net/http/cookie"
	gonethttphandler "github.com/ralvarezdev/go-net/http/handler"
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
		AuthorizationCookieAttributes          *gonethttpcookie.Attributes
		GetExpiresAtFn                         GetExpiresAtFn
	}

	// GetExpiresAtFn is a function that returns the expiration time for a token
	GetExpiresAtFn func(metadataValue string) (time.Time, error)
)

// NewDefaultAuthenticationParser parses the metadata authentication from a gRPC to the header or as a cookie
//
// Parameters:
//
// - handler: the HTTP handler to handle errors
// - options: the options for the DefaultAuthenticationParser
//
// Returns:
//
// - ParseAuthorizationMetadataAsHeader: parses the authorization metadata as a header
// - error: error if something goes wrong
func NewDefaultAuthenticationParser(
	handler gonethttphandler.Handler,
	options *Options,
) (
	*DefaultAuthenticationParser,
	error,
) {
	// Check if the handler is nil
	if handler == nil {
		return nil, gonethttphandler.ErrNilHandler
	}

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
// - w: http.ResponseWriter
// - ctx: context.Context
func (d DefaultAuthenticationParser) ParseAuthorizationMetadataAsHeader(
	w http.ResponseWriter,
	ctx context.Context,
) {
	// Get the metadata from the context
	md, err := gogrpcmd.GetCtxMetadata(ctx)
	if err != nil {
		d.handler.HandleErrorResponse(
			w,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Get the authorization metadata from the context
	authorization, err := gogrpcmd.GetMetadataAuthorizationToken(md)
	if err != nil {

	}

	// Set the authorization header
	w.Header().Set(gonethttp.Authorization, authorization)

	// Iterate over the metadata keys to authorization header names
	for metadataKey, headerName := range d.options.MetadataKeysToAuthorizationHeaderNames {
		// Get the metadata value
		if metadataValueSlice, ok := md[metadataKey]; ok {
			// Check if the metadata value is empty
			if metadataValueSlice == nil || len(metadataValueSlice) == 0 {
				continue
			}

			// Get the first value of the metadata
			metadataValue := metadataValueSlice[0]

			// Set the header
			w.Header().Set(headerName, metadataValue)
		}
	}
}

// ParseAuthorizationMetadataAsCookie parses the authorization metadata as a cookie
//
// Parameters:
//
// - w: http.ResponseWriter
// - ctx: context.Context
func (d DefaultAuthenticationParser) ParseAuthorizationMetadataAsCookie(
	w http.ResponseWriter,
	ctx context.Context,
) error {
	// Get the metadata from the context
	md, err := gogrpcmd.GetCtxMetadata(ctx)
	if err != nil {
		return err
	}

	// Get the authorization metadata from the context
	authorization, err := gogrpcmd.GetMetadataAuthorizationToken(md)
	if err != nil {
		return err
	}

	// Check if the GetExpiresAtFn is nil
	var expiresAt time.Time
	if d.options.GetExpiresAtFn != nil {
		expiresAt, err = d.options.GetExpiresAtFn(authorization)
		if err != nil {
			return err
		}
	}

	// Set the authorization cookie
	gonethttpcookie.SetCookie(
		w,
		d.options.AuthorizationCookieAttributes,
		authorization,
		expiresAt,
	)

	// Iterate over the metadata keys to cookies attributes
	for metadataKey, cookieAttributes := range d.options.MetadataKeysToCookiesAttributes {
		// Get the metadata value
		if metadataValueSlice, ok := md[metadataKey]; ok {
			// Check if the metadata value is empty
			if metadataValueSlice == nil || len(metadataValueSlice) == 0 {
				continue
			}

			// Get the first value of the metadata
			metadataValue := metadataValueSlice[0]

			// Get the expiration time if the function is set
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
	}
	return nil
}
