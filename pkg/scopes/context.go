package scopes

import (
	"context"
	"net/http"
)

// ContextKey represents a key used to access scopes from contexts.
type ContextKey int

// ScopeKey represents the key used for the scope context.
const (
	ScopeKey ContextKey = iota
)

// ScopeFunction is a type that represents a function that takes a context
// and returns a Scope.
type ScopeFunction func(ctx context.Context) Scope

// scopeFromContext gets a Scope from a context or defaults to the
// ScopeUnauthenticated scope if one is not present.
func scopeFromContext(ctx context.Context) Scope {
	scope, ok := ctx.Value(ScopeKey).(Scope)
	if !ok {
		return ScopeUnauthenticated
	}

	return scope
}

// Middleware provides a middleware that injects scope values into the context.
func Middleware(function ScopeFunction) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Run the scope calculation function and add the scope to the context.
			scope := function(ctx)
			ctx = context.WithValue(ctx, ScopeKey, scope)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MarshalFromContext uses a context to marshal data.
func MarshalFromContext(ctx context.Context, input interface{}) interface{} {
	scope := scopeFromContext(ctx)
	return Marshal(scope, input)
}
