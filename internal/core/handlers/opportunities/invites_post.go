package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/apierr"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// InvitesPost creates new invites from user emails.
func InvitesPost(opportunitiesService opportunities.Service) http.HandlerFunc {
	type postInviteEmail struct {
		Email string `json:"email" validate:"email"`
	}
	type request struct {
		Invites []postInviteEmail `json:"invites"`
	}
	type inviteError struct {
		Email          string `json:"email"`
		ErrorReference string `json:"ref"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opportunityID, err := idctx.Get(r, "opportunityID")
		if err != nil {
			return
		}

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		req := request{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		errors := []inviteError{}

		for _, invite := range req.Invites {
			err := opportunitiesService.InviteVolunteer(ctx, userID, opportunityID, invite.Email)
			if err != nil {
				errors = append(errors, inviteError{
					invite.Email,
					apierr.Ref(err),
				})
			}
		}

		if len(errors) > 0 {
			resp.BadRequest(w, r, resp.ErrorRef(32, "errors while sending invites", "opportunities.multiple_invite_errors", errors))
			return
		}

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}
