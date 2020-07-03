package scopes

import (
	"context"
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

// ScopeFromContext gets a Scope from a context or defaults to the
// ScopeUnauthenticated scope if one is not present.
func ScopeFromContext(ctx context.Context) Scope {
	scope, ok := ctx.Value(ScopeKey).(Scope)
	if !ok {
		return ScopeUnauthenticated
	}

	return scope
}

// MarshalFromContext uses a context to marshal data.
func MarshalFromContext(ctx context.Context, input interface{}) interface{} {
	scope := ScopeFromContext(ctx)
	return Marshal(scope, input)
}
