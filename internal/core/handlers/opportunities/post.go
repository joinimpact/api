package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

type postResponse struct {
	Success       bool  `json:"success"`
	OpportunityID int64 `json:"opportunityId"`
}

// Post creates a new opportunity.
func Post(opportunitiesService opportunities.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		organizationID, err := idctx.Get(r, "organizationID")
		if err != nil {
			return
		}

		req := struct {
			Title        string                      `json:"title" validate:"omitempty,min=4,max=128"`
			Description  string                      `json:"description"`
			Requirements *opportunities.Requirements `json:"requirements"`
			Limits       *opportunities.Limits       `json:"limits"`
		}{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		id, err := opportunitiesService.CreateOpportunity(ctx, opportunities.OpportunityView{
			CreatorID:      userID,
			OrganizationID: organizationID,
			Title:          req.Title,
			Description:    req.Description,
			Requirements:   req.Requirements,
			Limits:         req.Limits,
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

		resp.OK(w, r, postResponse{
			Success:       true,
			OpportunityID: id,
		})
	}
}
