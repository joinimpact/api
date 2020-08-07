package conversations

import (
	"net/http"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// GetByOrganization gets an organization's conversations.
func GetByOrganization(conversationsService conversations.Service) http.HandlerFunc {
	type response struct {
		Conversations []conversations.ConversationView `json:"conversations"`
		Pages         uint                             `json:"pages"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		organizationID, err := idctx.Get(r, "organizationID")
		if err != nil {
			return
		}

		res, err := conversationsService.GetOrganizationConversations(ctx, organizationID)
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

		resp.OK(w, r, response{
			Conversations: res.Conversations,
			Pages:         res.Pages,
		})
	}
}
