package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// GetByVolunteer gets all of a volunteer's currently enrolled
// opportunities.
func GetByVolunteer(opportunitiesService opportunities.Service) http.HandlerFunc {
	type response struct {
		Opportunities []opportunities.OpportunityView `json:"opportunities"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := idctx.Get(r, "userID")
		if err != nil {
			return
		}

		res, err := opportunitiesService.GetVolunteerOpportunities(ctx, userID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrOpportunityNotFound, *opportunities.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{res})
	}
}
