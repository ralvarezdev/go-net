package parser

import (
	gonethttprequest "github.com/ralvarezdev/go-net/http/request"
	gonethttpresponse "github.com/ralvarezdev/go-net/http/response"
)

type (
	// Parser is the interface that handles both the encoder and decoder tasks
	Parser interface {
		gonethttpresponse.Encoder
		gonethttprequest.Decoder
	}
)
