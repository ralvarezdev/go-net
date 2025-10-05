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
		HandleFunc(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		ExactHandleFunc(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterRoute(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterExactRoute(
			pattern string,
			handler http.HandlerFunc,
			middlewares ...func(next http.Handler) http.Handler,
		)
		RegisterHandler(pattern string, handler http.Handler)
		NewGroup(
			pattern string,
			middlewares ...func(next http.Handler) http.Handler,
		) RouterWrapper
		RegisterGroup(router RouterWrapper)
		Pattern() string
		RelativePath() string
		FullPath() string
		Method() string
		ServeStaticFiles(pattern, path string)
		Logger() *slog.Logger
		Mode() *goflagsmode.Flag
	}
)
