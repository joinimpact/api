package models

// OrganizationMembershipInvite represents an active invite to be a member in an organization.
type OrganizationMembershipInvite struct {
	Model
	Accepted       bool   `json:"accepted"`
	InviteeEmail   string `json:"inviteeEmail,omitempty"`
	InviteeID      int64  `json:"inviteeId"`
	Invitee        User   `gorm:"foreignkey:InviteeID"`
	OrganizationID int64  `json:"organizationId"`
	Organization   Organization
	InviterID      int64 `json:"inviterId"`
	Inviter        User  `gorm:"foreignkey:InviterID"`
}

// OrganizationMembershipInviteRepository represents the interface for a repository of membership invite entities.
type OrganizationMembershipInviteRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OrganizationMembershipInvite, error)
	// FindByUserID finds multiple entities by the user ID.
	FindByUserID(userID int64) ([]OrganizationMembershipInvite, error)
	// FindByUserEmail finds multiple entities by the user Email.
	FindByUserEmail(userEmail string) ([]OrganizationMembershipInvite, error)
	// FindByOrganizationID finds multiple entities by the organization ID.
	FindByOrganizationID(organizationID int64) ([]OrganizationMembershipInvite, error)
	// Create creates a new entity.
	Create(organizationMembershipInvite OrganizationMembershipInvite) error
	// Update updates an entity with the ID in the provided entity.
	Update(organizationMembershipInvite OrganizationMembershipInvite) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
