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
	InviterID     int64  `json:"inviterId"`
	Inviter       User   `gorm:"foreignkey:InviterID"`
	Key           string `json:"-"` // the secret key which is sent to the invitee for authorization
}

// OpportunityMembershipInviteRepository represents the interface for a repository of membership invite entities.
type OpportunityMembershipInviteRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OpportunityMembershipInvite, error)
	// FindByUserID finds multiple entities by the user ID.
	FindByUserID(userID int64) ([]OpportunityMembershipInvite, error)
	// FindByUserEmail finds multiple entities by the user Email.
	FindByUserEmail(userEmail string) ([]OpportunityMembershipInvite, error)
	// FindInOpportunityByUserID finds a membership invite in an opportunity by user ID.
	FindInOpportunityByUserID(opportunityID, userID int64) (*OpportunityMembershipInvite, error)
	// FindInOpportunityByEmail finds a membership invite in an opportunity by user email.
	FindInOpportunityByEmail(opportunityID int64, email string) (*OpportunityMembershipInvite, error)
	// FindByOpportunityID finds multiple entities by the opportunity ID.
	FindByOpportunityID(opportunityID int64) ([]OpportunityMembershipInvite, error)
	// Create creates a new entity.
	Create(opportunityMembershipInvite OpportunityMembershipInvite) error
	// Update updates an entity with the ID in the provided entity.
	Update(opportunityMembershipInvite OpportunityMembershipInvite) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
