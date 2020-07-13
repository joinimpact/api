package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// PublishPost publishes an opportunity by ID.
func PublishPost(opportunitiesService opportunities.Service) http.HandlerFunc {
	type response struct {
		Success bool `json:"success"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opportunityID, err := idctx.Get(r, "opportunityID")
		if err != nil {
			return
		}

		err = opportunitiesService.PublishOpportunity(ctx, opportunityID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrOpportunityNotPublishable:
				resp.BadRequest(w, r, resp.ErrorInvalidFields(98, "missing or invalid fields", err.(*opportunities.ErrOpportunityNotPublishable).InvalidFields))
			case *opportunities.ErrOpportunityNotFound, *opportunities.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{true})
	}
}
