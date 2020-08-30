package conversations

import (
	"context"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/scopes"
)

// ScopeProviderConversations provides a scope based on the currently logged in user and the conversation they are trying to access.
func ScopeProviderConversations(conversationsService conversations.Service) scopes.ScopeFunction {
	return func(ctx context.Context) scopes.Scope {
		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			return scopes.ScopeAuthenticated
		}

		conversationID, err := idctx.GetFromContext(ctx, "conversationID")
		if err != nil {
			return scopes.ScopeAuthenticated
		}

		// Check that a membership exists for the current user.
		_, err = conversationsService.GetUserConversationMembership(ctx, userID, conversationID)
		if err != nil {
			return scopes.ScopeAuthenticated
		}

		return scopes.ScopeCollaborator
	}
}
