# go-net

`net/http` routing, middlewares and helpers for Go projects: a router wrapper over `http.ServeMux` with nested modules, JSend-style JSON responses, request decoding and validation, authentication (JWT and gRPC-backed), rate limiting, request size limiting, cookies and request context helpers. Requires Go 1.25.1 (per `go.mod`).

## Installation

```bash
go get github.com/ralvarezdev/go-net
```

Direct dependencies are other `github.com/ralvarezdev` libraries (`go-flags`, `go-grpc`, `go-json`, `go-jwt`, `go-rate-limiter`, `go-reflect`, `go-strings`, `go-validator`) plus `google.golang.org/grpc` and `genproto` RPC types.

## Packages

- **`flags/port`** — `NewFlag(defaultValue)` and `SetFlag` for a `-port` flag.
- **`http`** — `GetClientIP(r)`, `Module` (pattern, middlewares, submodules, `BeforeLoadFn` / `AddHandlersFn` / `AfterLoadFn` hooks, `Create(baseRouter)`), `NewMiddlewares`, `NewSubmodules`.
- **`http/route`** — `RouterWrapper` (`AddHandleFunc`, `AddExactHandleFunc`, `AddEndpointHandler`, `AddRouter`, `ServeStaticFiles`, `Handler`, `Mux`, ...), `NewBaseRouter(mode, handler, logger, middlewares...)`, `NewRouter`, `ChainHandlers`, `GetWildcards`, `SplitPattern`.
- **`http/context`**, **`http/cookie`** — request-scoped body, wildcards and query parameters; `SetCookie`, `SetTimestampCookie`, `GetTimestampCookie`, `DeleteCookie(s)`, `RenovateCookie`.
- **`http/request/...`** — content-type check, decode error handling, JSON and protojson decoders, default requests handler.
- **`http/response/...`** — error/fail constructors (`NewError`, `NewFailFieldError`, `NewFailDataError`, `NewDebugError`), JSend bodies (`StatusSuccess`, `StatusFail`, `StatusError`), gRPC status mapping.
- **`http/handler/...`** — JSON, JSend, protojson and stream handlers (`NewHandler`, `NewStreamHandler`).
- **`http/middleware/auth`** — JWT middleware `NewMiddleware(responsesHandler, validator, options)`; `auth/grpc` delegates authentication to gRPC.
- **`http/middleware/{errorhandler,ratelimiter,sizelimiter,validator}`** — error handling, rate limiting via `go-rate-limiter` (Redis), `Limit(bytesSizeLimit)` using `http.MaxBytesReader`, and request validation (`NewJSONMiddleware`, `NewProtoJSONMiddleware`).
- **`url`** — `IsURL(str)`.

## Usage (sketch)

```go
modeFlag := goflagsmode.NewFlag(goflagsmode.DefaultMode, goflagsmode.AllowedModes)

// handler: a gonethttphandler.Handler (e.g. from http/handler/json/jsend)
router, err := gonethttproute.NewBaseRouter(modeFlag, handler, slog.Default())

router.AddHandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
})
http.ListenAndServe(":8080", router.Handler())
```

Larger services can organise routes with `Module` values (nested submodules, per-module middlewares) and call `Create(baseRouter)`. Check constructor parameters in the source before use.

## Development

```bash
go build ./...
go vet ./...
```

There are no tests.

## License

GNU General Public License v3.0 (see [LICENSE](LICENSE)).
