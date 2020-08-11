package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// VolunteersAcceptPost accepts a volunteer's request.
func VolunteersAcceptPost(opportunitiesService opportunities.Service, conversationsService conversations.Service) http.HandlerFunc {
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

		opportunityID, err := idctx.Get(r, "opportunityID")
		if err != nil {
			return
		}

		volunteerID, err := idctx.Get(r, "userID")
		if err != nil {
			return
		}

		err = opportunitiesService.AcceptOpportunityMembershipRequest(ctx, opportunityID, volunteerID, userID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrRequestNotFound:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		id, err := conversationsService.SendVolunteerRequestAcceptanceMessage(ctx, volunteerID, userID, opportunityID)
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
			Success:   true,
			MessageID: id,
		})
	}
}
