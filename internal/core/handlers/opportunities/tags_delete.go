package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// TagsDelete deletes a single opportunity tag by ID.
func TagsDelete(opportunitiesService opportunities.Service) http.HandlerFunc {
	type response struct {
		Success bool `json:"success"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opportunityID, err := idctx.Get(r, "opportunityID")
		if err != nil {
			return
		}

		tagID, err := idctx.Get(r, "tagID")
		if err != nil {
			return
		}

		err = opportunitiesService.RemoveOpportunityTag(ctx, opportunityID, tagID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrOpportunityNotFound, *opportunities.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
		}

		resp.OK(w, r, response{true})
	}
}
