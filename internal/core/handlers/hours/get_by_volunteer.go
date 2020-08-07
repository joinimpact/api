package hours

import (
	"net/http"

	"github.com/joinimpact/api/internal/hours"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// GetByUser gets a user's hour logs.
func GetByUser(hoursService hours.Service) http.HandlerFunc {
	type response struct {
		VolunteeringHourLogs []models.VolunteeringHourLog `json:"hourLogs"`
		Pages                uint                         `json:"pages"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := idctx.Get(r, "userID")
		if err != nil {
			return
		}

		res, err := hoursService.GetHoursByVolunteer(ctx, userID)
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
			VolunteeringHourLogs: res.VolunteeringHourLogs,
			Pages:                res.Pages,
		})
	}
}
