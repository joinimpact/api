package conversations

import (
	"net/http"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// MessagesPost sends a standard message to a single conversation.
func MessagesPost(conversationsService conversations.Service) http.HandlerFunc {
	type messageBody struct {
		Text string `json:"text" validate:"min=1,max=1024"`
	}
	type request struct {
		Body messageBody `json:"body" validate:"dive"`
	}
	type response struct {
		Success   bool  `json:"success"`
		MessageID int64 `json:"messageId"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		conversationID, err := idctx.Get(r, "conversationID")
		if err != nil {
			return
		}

		req := request{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		id, err := conversationsService.SendStandardMessage(ctx, conversationID, userID, req.Body.Text)
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

		resp.OK(w, r, response{true, id})
	}
}
