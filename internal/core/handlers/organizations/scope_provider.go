package organizations

import (
	"context"

	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/organizations"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/scopes"
)

// ScopeProviderOrganizations provides a scope based on a user id and organization.
func ScopeProviderOrganizations(organizationsService organizations.Service) scopes.ScopeFunction {
	return func(ctx context.Context) scopes.Scope {
		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			return scopes.NoChange
		}

		organizationID, err := idctx.GetFromContext(ctx, "organizationID")
		if err != nil {
			return scopes.NoChange
		}

		membership, err := organizationsService.GetOrganizationMembership(organizationID, userID)
		if err != nil {
			return scopes.NoChange
		}

		switch membership {
		case models.OrganizationPermissionsCreator:
			return scopes.ScopeOwner
		case models.OrganizationPermissionsOwner:
			return scopes.ScopeAdmin
		case models.OrganizationPermissionsMember:
			return scopes.ScopeManager
		}

		return scopes.NoChange
	}
}
