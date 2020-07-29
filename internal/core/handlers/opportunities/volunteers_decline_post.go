package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// VolunteersDeclinePost accepts a volunteer's request.
func VolunteersDeclinePost(opportunitiesService opportunities.Service) http.HandlerFunc {
	type response struct {
		Success bool `json:"success"`
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

		err = opportunitiesService.DeclineOpportunityMembershipRequest(ctx, opportunityID, volunteerID, userID)
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

		resp.OK(w, r, response{true})
	}
}
