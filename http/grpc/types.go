package grpc

import (
	"net/http"
	"strings"
	"time"

	gogrpc "github.com/ralvarezdev/go-grpc"
	gogrpcmd "github.com/ralvarezdev/go-grpc/metadata"
	gojwt "github.com/ralvarezdev/go-jwt"
	"google.golang.org/grpc/metadata"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpcookie "github.com/ralvarezdev/go-net/http/cookie"
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

	// Add authorization header if not present to map header names
	if options.MetadataKeysToAuthorizationHeaderNames == nil {
		options.MetadataKeysToAuthorizationHeaderNames = make(map[string]string)
	}
	if _, ok := options.MetadataKeysToAuthorizationHeaderNames[gogrpc.AuthorizationMetadataKey]; !ok {
		options.MetadataKeysToAuthorizationHeaderNames[gogrpc.AuthorizationMetadataKey] = gonethttp.Authorization
	}

	return &DefaultAuthenticationParser{
		options: options,
	}, nil
}

// ParseAuthorizationMetadataAsHeader parses the authorization metadata as a header
//
// Parameters:
//
// - md: The gRPC metadata with authorization information
// - w: http.ResponseWriter
//
// Returns:
//
// - error: error if something goes wrong
func (d DefaultAuthenticationParser) ParseAuthorizationMetadataAsHeader(
	md metadata.MD,
	w http.ResponseWriter,
) error {
	// Iterate over the metadata keys to authorization header names
	for metadataKey, headerName := range d.options.MetadataKeysToAuthorizationHeaderNames {
		// Get the metadata value
		metadataValueSlice, err := gogrpcmd.GetMetadataValue(md, metadataKey)
		if err != nil || len(metadataValueSlice) == 0 {
			continue
		}

		// Get the first value of the metadata
		metadataValue := metadataValueSlice[0]

		// Check if the metadata value is empty, which means that the header must be deleted
		if metadataValue == "" {
			w.Header().Del(headerName)
			continue
		}

		// Check if the authorization is a bearer token
		parts := strings.Split(metadataValue, " ")
		if len(parts) < 2 || parts[0] != gojwt.BearerPrefix {
			continue
		}
		token := parts[1]

		// Set the header
		w.Header().Set(headerName, gojwt.BearerPrefix+" "+token)
	}
	return nil
}

// ParseAuthorizationMetadataAsCookie parses the authorization metadata as a cookie
//
// Parameters:
//
// - md: The gRPC metadata with authorization information
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
		metadataValueSlice, err := gogrpcmd.GetMetadataValue(md, metadataKey)
		if err != nil || len(metadataValueSlice) == 0 {
			continue
		}

		// Get the first value of the metadata
		metadataValue := metadataValueSlice[0]

		// Check if the metadata value is empty, which means that the cookie must be deleted
		if metadataValue == "" {
			gonethttpcookie.DeleteCookie(
				w,
				cookieAttributes,
			)
			continue
		}

		// Check if the authorization is a bearer token
		parts := strings.Split(metadataValue, " ")

		// Check if the authorization is a bearer token
		if len(parts) < 2 || parts[0] != gojwt.BearerPrefix {
			continue
		}

		// Get the token from the metadata
		token := parts[1]

		// Get the expiration time if the function is set
		var expiresAt time.Time
		if d.options.GetExpiresAtFn != nil {
			expiresAt, err = d.options.GetExpiresAtFn(token)
			if err != nil {
				return err
			}
		}

		// Set the cookie
		gonethttpcookie.SetCookie(
			w,
			cookieAttributes,
			token,
			expiresAt,
		)
	}
	return nil
}
