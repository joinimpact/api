package hours

import (
	"net/http"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/hours"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// OrganizationRequestAcceptPost accepts a volunteer's hour request.
func OrganizationRequestAcceptPost(hoursService hours.Service, conversationsService conversations.Service) http.HandlerFunc {
	type response struct {
		Success   bool  `json:"success"`
		MessageID int64 `json:"messageId"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			return
		}

		requestID, err := idctx.Get(r, "requestID")
		if err != nil {
			return
		}

		err = hoursService.AcceptRequest(ctx, userID, requestID)
		if err != nil {
			switch err.(type) {
			case *hours.ErrRequestNotFound:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *hours.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		id, err := conversationsService.SendHoursRequestAcceptedMessage(ctx, userID, requestID)
		if err != nil {
			switch err.(type) {
			case *conversations.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{
			Success:   true,
			MessageID: id,
		})
	}
}
