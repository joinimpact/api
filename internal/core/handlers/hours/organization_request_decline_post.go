package hours

import (
	"net/http"

	"github.com/joinimpact/api/internal/hours"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// OrganizationRequestDeclinePost declines a volunteer's hour request.
func OrganizationRequestDeclinePost(hoursService hours.Service) http.HandlerFunc {
	type response struct {
		Success bool `json:"success"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := idctx.Get(r, "userID")
		if err != nil {
			return
		}

		requestID, err := idctx.Get(r, "requestID")
		if err != nil {
			return
		}

		err = hoursService.DeclineRequest(ctx, userID, requestID)
		if err != nil {
			switch err.(type) {
			case *hours.ErrRequestNotFound:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *hours.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
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
