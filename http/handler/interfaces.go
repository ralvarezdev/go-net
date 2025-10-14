package handler

import (
	gonethttprequesthandler "github.com/ralvarezdev/go-net/http/request/handler"
	gonethttpresponsehandler "github.com/ralvarezdev/go-net/http/response/handler"
)

type (
	// Handler is the interface that handles both the requests decoding and responses encoding tasks
	Handler interface {
		gonethttpresponsehandler.ResponsesHandler
		gonethttprequesthandler.RequestsHandler
	}
)
