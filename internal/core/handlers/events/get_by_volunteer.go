package events

import (
	"net/http"

	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// GetByVolunteer gets all of a volunteers upcoming events.
func GetByVolunteer(eventsService events.Service) http.HandlerFunc {
	type response struct {
		Events []events.EventView `json:"events"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := idctx.Get(r, "userID")
		if err != nil {
			return
		}

		res, err := eventsService.GetUserEvents(ctx, userID)
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
