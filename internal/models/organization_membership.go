package models

// OrganizationPermissions flags
const (
	OrganizationPermissionsMember  = iota
	OrganizationPermissionsOwner   = iota
	OrganizationPermissionsCreator = iota
)

// OrganizationMembership creates a relationship between Organizations and their employees.
type OrganizationMembership struct {
	Model
	Active          bool  `json:"-"`               // controls whether or not the entity is active
	UserID          int64 `json:"-"`               // the ID of the user being granted membership
	OrganizationID  int64 `json:"-"`               // the ID of the organization the user is being granted access to
	PermissionsFlag int   `json:"permissionsFlag"` // a flag which designates permissions the user has
	InviterID       int64 `json:"inviterId"`       // the ID of the user who invited the member to the organization
}

// OrganizationMembershipRepository represents a repository of organization memberships.
type OrganizationMembershipRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OrganizationMembership, error)
	// FindByUserID finds multiple entities by the user ID.
	FindByUserID(userID int64) ([]OrganizationMembership, error)
	// FindByOrganizationID finds multiple entities by the organization ID.
	FindByOrganizationID(organizationID int64) ([]OrganizationMembership, error)
	// FindUserInOrganization finds a user's membership in a specific organization.
	FindUserInOrganization(organizationID, userID int64) (*OrganizationMembership, error)
	// Create creates a new entity.
	Create(organizationMembership OrganizationMembership) error
	// Update updates an entity with the ID in the provided entity.
	Update(organizationMembership OrganizationMembership) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
