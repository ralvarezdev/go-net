package route

import (
	"net/http"

	gonethttpctx "github.com/ralvarezdev/go-net/http/context"
)

// SetCtxWildcardsMiddleware is the middleware to add the wildcards to the context
//
// Parameters:
//
//   - wildcardKeys: The wildcard keys to be added to the context
//
// Returns:
//
//   - func(next http.Handler) http.Handler: The middleware handler
func SetCtxWildcardsMiddleware(
	wildcardKeys []string,
) func(next http.Handler) http.Handler {
	return func(
		next http.Handler,
	) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Add the wildcards to the context
				if wildcardKeys != nil {
					// Create a map to store the wildcards
					wildcards := make(map[string]string)

					for _, key := range wildcardKeys {
						// Get the wildcard from the request
						value := r.PathValue(key)

						// Add the wildcard to the map
						wildcards[key] = value
					}

					// Add the wildcards to the context
					r = gonethttpctx.SetCtxWildcards(r, wildcards)
				}

				// Call the next handler
				next.ServeHTTP(w, r)
			},
		)
	}
}

// SetCtxQueryParametersMiddleware is the middleware to add the query parameters to the context
//
// Parameters:
//
//   - next: The next handler to be executed
//
// Returns:
//
//   - http.Handler: The middleware handler
func SetCtxQueryParametersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Get the URL query
			query := r.URL.Query()
			if query != nil {
				// Convert the query to a map[string][]string
				parameters := map[string][]string(query)

				// Add the parameters to the context
				r = gonethttpctx.SetCtxQueryParameters(r, parameters)
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		},
	)
}
