package scopes

import (
	"context"
	"net/http"
)

// Middleware provides a middleware that injects scope values into the context.
func Middleware(function ScopeFunction) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Run the scope calculation function and add the scope to the context.
			scope := function(ctx)
			if scope == NoChange {
				// If NoChange returned, default to the current scope.
				scope = ScopeFromContext(ctx)
			}

			ctx = context.WithValue(ctx, ScopeKey, scope)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
