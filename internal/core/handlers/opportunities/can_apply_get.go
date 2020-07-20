package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// Volunteer status
const (
	VolunteerStatusCanApply       = 0
	VolunteerStatusAlreadyApplied = 1
	VolunteerStatusAlreadyMember  = 2
)

// StatusGet gets a user's status per opportunity.
func StatusGet(opportunitiesService opportunities.Service) http.HandlerFunc {
	type response struct {
		VolunteerStatus int `json:"status"`
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

		_, err = opportunitiesService.GetOpportunityMembership(ctx, opportunityID, userID)
		if err == nil {
			resp.OK(w, r, response{VolunteerStatusAlreadyMember})
			return
		}

		// Check if the user can apply.
		err = opportunitiesService.CanRequestOpportunityMembership(ctx, opportunityID, userID)
		if err != nil {
			resp.OK(w, r, response{VolunteerStatusAlreadyApplied})
			return
		}

		// No request or membership, send code 0 (can apply).
		resp.OK(w, r, response{VolunteerStatusCanApply})
	}
}
