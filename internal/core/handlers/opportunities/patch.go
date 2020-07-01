package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

type patchResponse struct {
	Success bool `json:"success"`
}

// Patch updates an opportunity.
func Patch(opportunitiesService opportunities.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opportunityID, err := idctx.Get(r, "opportunityID")
		if err != nil {
			return
		}

		req := struct {
			Title       string `json:"title" validate:"omitempty,min=4,max=128"`
			Description string `json:"description"`
		}{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		err = opportunitiesService.UpdateOpportunity(ctx, models.Opportunity{
			Model: models.Model{
				ID: opportunityID,
			},
			Title:       req.Title,
			Description: req.Description,
		})
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrOpportunityNotFound, *opportunities.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, patchResponse{
			Success: true,
		})
	}
}
