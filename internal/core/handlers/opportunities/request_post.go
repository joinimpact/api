package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// RequestPost creates an opportunity request for a user on the specified opportunity.
func RequestPost(opportunitiesService opportunities.Service, conversationsService conversations.Service) http.HandlerFunc {
	type request struct {
		Message string `json:"message" validate:"omitempty,min=12,max=128"`
	}
	type response struct {
		Success        bool  `json:"success"`
		ConversationID int64 `json:"conversationId"`
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

		req := request{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		opportunity, err := opportunitiesService.GetOpportunity(ctx, opportunityID)
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

		// Create the opportunity membership request.
		requestID, err := opportunitiesService.RequestOpportunityMembership(ctx, opportunityID, userID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrOpportunityNotFound, *opportunities.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *opportunities.ErrMembershipAlreadyRequested:
				resp.BadRequest(w, r, resp.Error(400, err.Error()))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		if req.Message == "" {
			req.Message = "Hi! I'd like to join your opportunity. Here is my profile:"
		}

		// Create the conversation.
		conversationID, err := conversationsService.CreateOpportunityMembershipRequestConversation(ctx, opportunity.OrganizationID, opportunityID, requestID, userID, req.Message)
		if err != nil {
			switch err.(type) {
			case *conversations.ErrConversationNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *conversations.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{true, conversationID})
	}
}
