package events

import (
	"net/http"

	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/users"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// ResponsesGet gets all responses for a single event.
func ResponsesGet(eventsService events.Service, usersService users.Service) http.HandlerFunc {
	type responseObject struct {
		models.EventResponse
		users.UserProfile
	}
	type response struct {
		Responses []responseObject `json:"responses"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		eventID, err := idctx.Get(r, "eventID")
		if err != nil {
			return
		}

		responses, err := eventsService.GetEventResponses(ctx, eventID)
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

		responseObjects := []responseObject{}
		for _, response := range responses {
			profile, err := usersService.GetMinimalUserProfile(response.UserID)
			if err != nil {
				resp.ServerError(w, r, resp.UnknownError)
			}

			responseObjects = append(responseObjects, responseObject{response, *profile})
		}

		resp.OK(w, r, response{responseObjects})
	}
}
