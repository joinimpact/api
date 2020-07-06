package users

import (
	"context"

	"github.com/joinimpact/api/pkg/scopes"
)

// ScopeProviderUsers provides a scope provider function for the users API.
func ScopeProviderUsers() scopes.ScopeFunction {
	return func(ctx context.Context) scopes.Scope {
		reqCtx, ok := ctx.Value(keyUserRequestContext).(*userRequestContext)
		if !ok {
			return scopes.NoChange
		}

		if reqCtx.isSelf {
			return scopes.ScopeOwner
		}

		return scopes.NoChange
	}
}
