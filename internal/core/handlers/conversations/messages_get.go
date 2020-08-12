package conversations

import (
	"net/http"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// MessagesGet gets messages by conversation ID.
func MessagesGet(conversationsService conversations.Service) http.HandlerFunc {
	type response struct {
		Messages     []conversations.MessageView `json:"messages"`
		Pages        uint                        `json:"pages"`
		TotalResults uint                        `json:"totalResults"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conversationID, err := idctx.Get(r, "conversationID")
		if err != nil {
			return
		}

		res, err := conversationsService.GetConversationMessages(ctx, conversationID)
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
			Messages:     res.Messages,
			Pages:        res.Pages,
			TotalResults: res.TotalResults,
		})
	}
}
