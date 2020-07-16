package opportunities

import (
	"net/http"

	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/internal/users"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// OpportunityVolunteer represents a volunteer in an opportunity.
type OpportunityVolunteer struct {
	models.OpportunityMembership
	users.UserProfile
}

// OpportunityPendingVolunteer represents a pending volunteer in an opportunity.
type OpportunityPendingVolunteer struct {
	models.OpportunityMembershipRequest
	users.UserProfile
}

// OpportunityInvitedVolunteer represents a pending volunteer in an opportunity.
type OpportunityInvitedVolunteer struct {
	EmailOnly bool `json:"emailOnly"`
	models.OpportunityMembershipInvite
	users.UserProfile
}

// VolunteersGet gets all volunteers in a specified opportunity.
func VolunteersGet(opportunitiesService opportunities.Service, usersService users.Service) http.HandlerFunc {
	type response struct {
		Volunteers []OpportunityVolunteer        `json:"volunteers"`
		Pending    []OpportunityPendingVolunteer `json:"pending"`
		Invited    []OpportunityInvitedVolunteer `json:"invited"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opportunityID, err := idctx.Get(r, "opportunityID")
		if err != nil {
			return
		}

		memberships, err := opportunitiesService.GetOpportunityVolunteers(ctx, opportunityID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrOpportunityNotFound, *opportunities.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
		}

		volunteers := []OpportunityVolunteer{}

		for _, membership := range memberships {
			profile, err := usersService.GetMinimalUserProfile(membership.UserID)
			if err != nil {
				switch err.(type) {
				case *users.ErrUserNotFound:
					resp.NotFound(w, r, resp.Error(404, err.Error()))
				case *users.ErrServerError:
					resp.ServerError(w, r, resp.Error(500, err.Error()))
				default:
					resp.ServerError(w, r, resp.UnknownError)
				}
				return
			}

			volunteers = append(volunteers, OpportunityVolunteer{membership, *profile})
		}

		pendingMemberships, err := opportunitiesService.GetOpportunityPendingVolunteers(ctx, opportunityID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrOpportunityNotFound, *opportunities.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
		}

		pendingVolunteers := []OpportunityPendingVolunteer{}

		for _, membership := range pendingMemberships {
			profile, err := usersService.GetMinimalUserProfile(membership.VolunteerID)
			if err != nil {
				switch err.(type) {
				case *users.ErrUserNotFound:
					resp.NotFound(w, r, resp.Error(404, err.Error()))
				case *users.ErrServerError:
					resp.ServerError(w, r, resp.Error(500, err.Error()))
				default:
					resp.ServerError(w, r, resp.UnknownError)
				}
				return
			}

			pendingVolunteers = append(pendingVolunteers, OpportunityPendingVolunteer{membership, *profile})
		}

		invites, err := opportunitiesService.GetOpportunityInvitedVolunteers(ctx, opportunityID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrOpportunityNotFound, *opportunities.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
		}

		invitedVolunteers := []OpportunityInvitedVolunteer{}

		for _, invite := range invites {
			profile, err := usersService.GetMinimalUserProfile(invite.InviteeID)
			if err != nil {
				invitedVolunteers = append(invitedVolunteers, OpportunityInvitedVolunteer{true, invite, users.UserProfile{}})
				continue
			}

			invitedVolunteers = append(invitedVolunteers, OpportunityInvitedVolunteer{false, invite, *profile})
		}

		resp.OK(w, r, response{volunteers, pendingVolunteers, invitedVolunteers})
	}
}
