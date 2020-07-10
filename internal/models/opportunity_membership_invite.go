package models

// OpportunityMembershipInvite represents an active invite to be a member in an opportunity.
type OpportunityMembershipInvite struct {
	Model
	Accepted      bool   `json:"accepted"`
	InviteeEmail  string `json:"inviteeEmail,omitempty"`
	InviteeID     int64  `json:"inviteeId"`
	Invitee       User   `gorm:"foreignkey:InviteeID"`
	OpportunityID int64  `json:"opportunityId"`
	Opportunity   Opportunity
	InviterID     int64 `json:"inviterId"`
	Inviter       User  `gorm:"foreignkey:InviterID"`
}

// OpportunityMembershipInviteRepository represents the interface for a repository of membership invite entities.
type OpportunityMembershipInviteRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OpportunityMembershipInvite, error)
	// FindByUserID finds multiple entities by the user ID.
	FindByUserID(userID int64) ([]OpportunityMembershipInvite, error)
	// FindByUserEmail finds multiple entities by the user Email.
	FindByUserEmail(userEmail string) ([]OpportunityMembershipInvite, error)
	// FindByOpportunityID finds multiple entities by the opportunity ID.
	FindByOpportunityID(opportunityID int64) ([]OpportunityMembershipInvite, error)
	// Create creates a new entity.
	Create(opportunityMembershipInvite OpportunityMembershipInvite) error
	// Update updates an entity with the ID in the provided entity.
	Update(opportunityMembershipInvite OpportunityMembershipInvite) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
