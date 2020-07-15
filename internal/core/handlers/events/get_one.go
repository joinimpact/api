package events

import (
	"net/http"

	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// GetOne gets a single event by ID.
func GetOne(eventsService events.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		eventID, err := idctx.Get(r, "eventID")
		if err != nil {
			return
		}

		event, err := eventsService.GetEvent(ctx, eventID)
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

		resp.OK(w, r, event)
	}
}
