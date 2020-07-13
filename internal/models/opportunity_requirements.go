package models

// OpportunityRequirements represents the requirements for an opportunity.
type OpportunityRequirements struct {
	Model
	OpportunityID       int64
	AgeLimitActive      bool
	AgeLimitFrom        int
	AgeLimitTo          int
	ExpectedHoursActive bool
	ExpectedHours       int
}

// OpportunityRequirementsRepository defines a repository of OpportunityRequirements entities.
type OpportunityRequirementsRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OpportunityRequirements, error)
	// FindByOpportunityID finds an entity by the opportunity ID.
	FindByOpportunityID(opportunityID int64) (*OpportunityRequirements, error)
	// Create creates a new entity.
	Create(opportunityRequirements OpportunityRequirements) error
	// Update updates an entity with the ID in the provided entity.
	Update(opportunityRequirements OpportunityRequirements) error
	// Save saves all fields in the provided entity.
	Save(opportunityRequirements OpportunityRequirements) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
