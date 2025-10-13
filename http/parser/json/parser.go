package json

import (
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttprequestjson "github.com/ralvarezdev/go-net/http/request/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponsejson "github.com/ralvarezdev/go-net/http/response/json"
)

type (
	// Parser is the parser implementation for JSON encoding and decoding
	Parser struct {
		gonethttprequestjson.Decoder
		gonethttpresponsejson.Encoder
	}
)

// NewParser creates a new JSON parser
//
// Parameters:
//
//   - decoder: the JSON decoder
//   - encoder: the JSON encoder
//
// Returns:
//
//   - *Parser: the created JSON parser
//   - error: the error if any
func NewParser(
	decoder *gonethttprequestjson.Decoder,
	encoder *gonethttpresponsejson.Encoder,
) (*Parser, error) {
	if decoder == nil {
		return nil, gonethttprequest.ErrNilDecoder
	}
	if encoder == nil {
		return nil, gonethttpresponse.ErrNilEncoder
	}

	return &Parser{
		*decoder,
		*encoder,
	}, nil
}
