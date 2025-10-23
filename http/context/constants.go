package context

type (
	// ContextKey is the type for context keys
	ContextKey string
)

const (
	// CtxBodyKey is the context key for the body
	CtxBodyKey ContextKey = "body"

	// CtxQueryParametersKey is the context key for the query parameters
	CtxQueryParametersKey ContextKey = "query_parameters"

	// CtxWildcardsKey is the context key for the wildcard
	CtxWildcardsKey ContextKey = "wildcards"
)
