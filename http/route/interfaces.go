package route

import (
	"log/slog"
	"net/http"

	goflagsmode "github.com/ralvarezdev/go-flags/mode"
)

type (
	// RouterWrapper is the interface for the routes
	RouterWrapper interface {
		Handler() http.Handler
		Mux() *http.ServeMux
		GetMiddlewares() []func(http.Handler) http.Handler
		AddHandleFunc(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		AddExactHandleFunc(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		AddEndpointHandler(
			pattern string,
			endpointHandler func(w http.ResponseWriter, r *http.Request) error,
			middlewares ...func(next http.Handler) http.Handler,
		)
		AddExactEndpointHandler(
			pattern string,
			endpointHandler func(w http.ResponseWriter, r *http.Request) error,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterHandler(pattern string, handler http.Handler)
		NewRouter(
			pattern string,
			middlewares ...func(next http.Handler) http.Handler,
		) (RouterWrapper, error)
		AddRouter(router RouterWrapper)
		Pattern() string
		RelativePath() string
		FullPath() string
		Method() string
		ServeStaticFiles(pattern, path string)
		Logger() *slog.Logger
		Mode() *goflagsmode.Flag
	}
)
