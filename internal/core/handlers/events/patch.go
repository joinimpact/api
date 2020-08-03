package events

import (
	"net/http"

	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// Patch updates a single event.
func Patch(eventsService events.Service) http.HandlerFunc {
	type response struct {
		Success bool `json:"success"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		eventID, err := idctx.Get(r, "eventID")
		if err != nil {
			return
		}

		req := events.ModifyEventRequest{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		req.ID = eventID

		req.OpportunityID = 0
		req.CreatorID = 0

		err = eventsService.UpdateEvent(ctx, req)
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

		resp.OK(w, r, response{
			Success: true,
		})
	}
}
