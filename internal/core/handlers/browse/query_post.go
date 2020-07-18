package browse

import (
	"net/http"

	"github.com/joinimpact/api/internal/opportunities"
	opportunitiesSearch "github.com/joinimpact/api/internal/search/stores/opportunities"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// QueryPost queries opportunities.
func QueryPost(opportunitiesService opportunities.Service) http.HandlerFunc {
	type request struct {
		opportunitiesSearch.Query
	}
	type response struct {
		Opportunities []opportunities.OpportunityView `json:"opportunities"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := request{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		res, err := opportunitiesService.Search(ctx, req.Query)
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

		resp.OK(w, r, response{res})
	}
}
