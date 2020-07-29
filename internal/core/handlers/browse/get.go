package browse

import (
	"net/http"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/resp"
)

// Get gets the browse page.
func Get(opportunitiesService opportunities.Service) http.HandlerFunc {
	type response struct {
		Sections []opportunities.Section `json:"sections"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		sections, err := opportunitiesService.GetRecommendations(ctx, userID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrInviteInvalid:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{sections})
	}
}
