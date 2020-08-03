package events

import (
	"net/http"

	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// Delete deletes a single event by ID.
func Delete(eventsService events.Service) http.HandlerFunc {
	type response struct {
		Success bool `json:"success"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		eventID, err := idctx.Get(r, "eventID")
		if err != nil {
			return
		}

		err = eventsService.DeleteEvent(ctx, eventID)
		if err != nil {
			switch err.(type) {
			case *events.ErrEventNotFound, *events.ErrResponseNotFound:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *events.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{true})
	}
}
