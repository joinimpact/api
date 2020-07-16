package events

import (
	"context"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/internal/organizations"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/scopes"
)

// ScopeProviderEvents provides a scope based on a user id and event id.
func ScopeProviderEvents(eventsService events.Service, organizationsService organizations.Service, opportunitiesService opportunities.Service) scopes.ScopeFunction {
	return func(ctx context.Context) scopes.Scope {
		eventID, err := idctx.GetFromContext(ctx, "eventID")
		if err != nil {
			return scopes.NoChange
		}

		event, err := eventsService.GetMinimalEvent(ctx, eventID)
		if err != nil {
			return scopes.NoChange
		}

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			return scopes.NoChange
		}

		opportunityID := event.OpportunityID

		// Get the opportunity so we can check organization membership.
		opportunity, err := opportunitiesService.GetMinimalOpportunity(ctx, opportunityID)
		if err != nil {
			return scopes.NoChange
		}

		// Check organization membership first for manager/admin permissions.
		organizationMembership, err := organizationsService.GetOrganizationMembership(opportunity.OrganizationID, userID)
		if err == nil {
			// If the membership was found, check its type and return.
			switch organizationMembership {
			case models.OrganizationPermissionsCreator:
				return scopes.ScopeOwner
			case models.OrganizationPermissionsOwner:
				return scopes.ScopeAdmin
			case models.OrganizationPermissionsMember:
				return scopes.ScopeManager
			}
		}

		// Get the opportunity membership.
		membership, err := opportunitiesService.GetOpportunityMembership(ctx, opportunityID, userID)
		if err != nil {
			return scopes.NoChange
		}

		// If the membership was found, check its type and return.
		switch membership {
		case models.OpportunityPermissionsMember:
			return scopes.ScopeCollaborator
		}

		return scopes.NoChange
	}
}
