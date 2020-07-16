package opportunities

import (
	"context"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/internal/organizations"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/scopes"
)

// ScopeProviderOpportunities provides a scope based on a user id and organization.
func ScopeProviderOpportunities(organizationsService organizations.Service, opportunitiesService opportunities.Service) scopes.ScopeFunction {
	return func(ctx context.Context) scopes.Scope {
		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			return scopes.NoChange
		}

		opportunityID, err := idctx.GetFromContext(ctx, "opportunityID")
		if err != nil {
			return scopes.NoChange
		}

		// Get the opportunity so we can check organization membership.
		opportunity, err := opportunitiesService.GetMinimalOpportunity(ctx, opportunityID)
		if err != nil {
			return scopes.NoChange
		}

		// If no opportunity membership is found,
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
