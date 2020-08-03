package conversations

import (
	"net/http"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// Get gets a single conversation by ID.
func Get(conversationsService conversations.Service, asOrganization bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conversationID, err := idctx.Get(r, "conversationID")
		if err != nil {
			return
		}

		err = nil
		var conversation *models.Conversation
		if asOrganization {
			conversation, err = conversationsService.GetOrganizationConversation(ctx, conversationID)
		} else {
			conversation, err = conversationsService.GetUserConversation(ctx, conversationID)
		}
		if err != nil {
			switch err.(type) {
			case *conversations.ErrConversationNotFound, *conversations.ErrUserNotFound:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *conversations.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, conversation)
	}
}
