package models

// OpportunityMembershipRequest represents an active request to be a member in an opportunity.
type OpportunityMembershipRequest struct {
	Model
	Accepted      bool        `json:"accepted"`
	VolunteerID   int64       `json:"volunteerID"`
	Volunteer     User        `json:"-" gorm:"foreignkey:VolunteerID"`
	OpportunityID int64       `json:"opportunityId"`
	Opportunity   Opportunity `json:"-"`
}

// OpportunityMembershipRequestRepository represents the interface for a repository of membership request entities.
type OpportunityMembershipRequestRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OpportunityMembershipRequest, error)
	// FindByVolunteerID finds multiple entities by the Volunteer ID.
	FindByVolunteerID(volunteerID int64) ([]OpportunityMembershipRequest, error)
	// FindByOpportunityID finds multiple entities by the opportunity ID.
	FindByOpportunityID(opportunityID int64) ([]OpportunityMembershipRequest, error)
	// FindInOpportunityByVolunteerID finds a single entity by opportunity and volunteer ID.
	FindInOpportunityByVolunteerID(opportunityID, volunteerID int64) (*OpportunityMembershipRequest, error)
	// Create creates a new entity.
	Create(opportunityMembershipRequest OpportunityMembershipRequest) error
	// Update updates an entity with the ID in the provided entity.
	Update(opportunityMembershipRequest OpportunityMembershipRequest) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
