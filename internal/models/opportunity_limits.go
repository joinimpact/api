package models

// OpportunityLimits represents the limits for an opportunity.
type OpportunityLimits struct {
	Model
	OpportunityID       int64
	VolunteersCapActive bool
	VolunteersCap       int
}

// OpportunityLimitsRepository defines a repository of OpportunityLimits entities.
type OpportunityLimitsRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OpportunityLimits, error)
	// FindByOpportunityID finds an entity by the opportunity ID.
	FindByOpportunityID(opportunityID int64) (*OpportunityLimits, error)
	// Create creates a new entity.
	Create(opportunityLimits OpportunityLimits) error
	// Update updates an entity with the ID in the provided entity.
	Update(opportunityLimits OpportunityLimits) error
	// Save saves all fields in the provided entity.
	Save(opportunityLimits OpportunityLimits) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
