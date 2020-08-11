package conversations

import (
	"net/http"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// GetByUser gets a user's conversations.
func GetByUser(conversationsService conversations.Service) http.HandlerFunc {
	type response struct {
		Conversations []conversations.ConversationView `json:"conversations"`
		Pages         uint                             `json:"pages"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := idctx.Get(r, "userID")
		if err != nil {
			return
		}

		res, err := conversationsService.GetUserConversations(ctx, userID)
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
