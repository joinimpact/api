package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// InviteValidatePost validates an invite and returns an opportunity view on success.
func InviteValidatePost(opportunitiesService opportunities.Service) http.HandlerFunc {
	type request struct {
		Key string `json:"key" validate:"min=8,max=128"`
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

		inviteID, err := idctx.Get(r, "inviteID")
		if err != nil {
			return
		}

		req := request{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		res, err := opportunitiesService.GetOpportunityFromInvite(ctx, opportunityID, userID, inviteID, req.Key)
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

		resp.OK(w, r, res)
	}
}
