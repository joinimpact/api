package events

import (
	"net/http"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// Post creates a new event.
func Post(eventsService events.Service) http.HandlerFunc {
	type response struct {
		Success bool  `json:"success"`
		EventID int64 `json:"eventId"`
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

		req := events.ModifyEventRequest{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		req.OpportunityID = opportunityID
		req.CreatorID = userID

		id, err := eventsService.CreateEvent(ctx, req)
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
			EventID: id,
		})
	}
}
