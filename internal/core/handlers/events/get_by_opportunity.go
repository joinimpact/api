package events

import (
	"net/http"

	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// GetByOpportunity gets events in a single opportunity.
func GetByOpportunity(eventsService events.Service) http.HandlerFunc {
	type response struct {
		Events []events.EventView `json:"events"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opportunityID, err := idctx.Get(r, "opportunityID")
		if err != nil {
			return
		}

		res, err := eventsService.GetOpportunityEvents(ctx, opportunityID)
		if err != nil {
			switch err.(type) {
			case *events.ErrEventNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *events.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{res})
	}
}
