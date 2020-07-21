package events

import (
	"net/http"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// ResponseGet gets a user's response to an event.
func ResponseGet(eventsService events.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		eventID, err := idctx.Get(r, "eventID")
		if err != nil {
			return
		}

		response, err := eventsService.GetUserEventResponse(ctx, userID, eventID)
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

		resp.OK(w, r, response)
	}
}
