package models

// OpportunityTag defines an opportunity's area of interest.
type OpportunityTag struct {
	Model
	OpportunityID int64       `json:"-"`
	Opportunity   Opportunity `json:"-"`
	TagID         int64       `json:"-"`
	Tag           Tag
}

// OpportunityTagRepository represents a repository of OpportunityTag.
type OpportunityTagRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OpportunityTag, error)
	// FindByOpportunityID finds entities by Opportunity ID.
	FindByOpportunityID(opportunityID int64) ([]OpportunityTag, error)
	// FindOpportunityTagByID finds a single entity by opportunity ID and tag ID.
	FindOpportunityTagByID(opportunityID int64, tagID int64) (*OpportunityTag, error)
	// Create creates a new entity.
	Create(opportunityTag OpportunityTag) error
	// Update updates an entity with the ID in the provided entity.
	Update(opportunityTag OpportunityTag) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
