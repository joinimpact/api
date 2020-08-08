package hours

import (
	"net/http"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/hours"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// OrganizationRequestsPost requests hours from an organization.
func OrganizationRequestsPost(hoursService hours.Service, conversationsService conversations.Service) http.HandlerFunc {
	type request struct {
		Hours float32 `json:"hours" validate:"min=1,max=100"`
	}
	type response struct {
		MessageID         int64 `json:"messageId"`
		HoursLogRequestID int64 `json:"hoursLogRequestID"`
		Success           bool  `json:"success"`
	}
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

		req := request{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		id, err := hoursService.RequestHours(ctx, userID, organizationID, req.Hours)
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

		messageID, err := conversationsService.SendHoursRequestMessage(ctx, userID, organizationID, id)
		if err != nil {
			switch err.(type) {
			case *conversations.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{
			MessageID:         messageID,
			HoursLogRequestID: id,
			Success:           true,
		})
	}
}
