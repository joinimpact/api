package events

import (
	"net/http"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// ResponsePut adds a user's response to an event.
func ResponsePut(eventsService events.Service) http.HandlerFunc {
	type request struct {
		Response *int `json:"response" validate:"required,max=2"`
	}
	type response struct {
		Success bool `json:"success"`
	}
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

		req := request{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		// Check which value the response matches.
		switch *req.Response {
		case models.EventResponseCanAttend:
			err = eventsService.SetEventResponseCanAttend(ctx, eventID, userID)
		case models.EventResponseCanNotAttend:
			err = eventsService.SetEventResponseCanNotAttend(ctx, eventID, userID)
		default:
			return
		}

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
