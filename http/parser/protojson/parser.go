package json

import (
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttprequestprotojson "github.com/ralvarezdev/go-net/http/request/protojson"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponseprotojson "github.com/ralvarezdev/go-net/http/response/protojson"
)

type (
	// Parser is the parser implementation for protoJSON encoding and decoding
	Parser struct {
		gonethttprequestprotojson.Decoder
		gonethttpresponseprotojson.Encoder
	}
)

// NewParser creates a new protoJSON parser
//
// Parameters:
//
//   - decoder: the protoJSON decoder
//   - encoder: the protoJSON encoder
//
// Returns:
//
//   - *Parser: the created protoJSON parser
//   - error: the error if any
func NewParser(
	decoder *gonethttprequestprotojson.Decoder,
	encoder *gonethttpresponseprotojson.Encoder,
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
