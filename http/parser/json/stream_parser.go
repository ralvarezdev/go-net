package json

import (
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttprequestjson "github.com/ralvarezdev/go-net/http/request/json"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
	gonethttpresponsejson "github.com/ralvarezdev/go-net/http/response/json"
)

type (
	// StreamParser is the parser implementation for streaming JSON encoding and decoding
	StreamParser struct {
		gonethttprequestjson.StreamDecoder
		gonethttpresponsejson.StreamEncoder
	}
)

// NewStreamParser creates a new JSON stream parser
//
// Parameters:
//
//   - decoder: the JSON stream decoder
//   - encoder: the JSON stream encoder
//
// Returns:
//
//   - *StreamParser: the created JSON stream parser
//   - error: the error if any
func NewStreamParser(
	decoder *gonethttprequestjson.StreamDecoder,
	encoder *gonethttpresponsejson.StreamEncoder,
) (*StreamParser, error) {
	if decoder == nil {
		return nil, gonethttprequest.ErrNilDecoder
	}
	if encoder == nil {
		return nil, gonethttpresponse.ErrNilEncoder
	}

	return &StreamParser{
		*decoder,
		*encoder,
	}, nil
}
