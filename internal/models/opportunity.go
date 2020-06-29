package models

// Opportunity represents an organization's opportunity.
type Opportunity struct {
	Model
	Active         bool         `json:"-"`                    // when false, the opportunity will be treated as if deleted. Useful for suspensions, etc.
	OrganizationID int64        ``                            // the id of the parent organization
	Organization   Organization ``                            //
	CreatorID      int64        ``                            // the id of the user who created the opportunity initially
	Creator        User         `gorm:"foreignkey:CreatorID"` //
	Public         bool         `json:"public"`               // whether or not the opportunity should be shown to volunteers
	Title          string       `json:"title"`                // the title of the opportunity
	Image          string       `json:"image"`                // a url to the opportunity's banner image
	Description    string       `json:"description"`          // a long description of the opportunity and its purpose
}

// OpportunityRepository represents a repository of opportunities.
type OpportunityRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*Opportunity, error)
	// FindByOrganizationID finds multiple entities by the organization ID.
	FindByOrganizationID(organizationID int64) ([]Opportunity, error)
	// FindByCreatorID finds multiple entities by the creator ID.
	FindByCreatorID(creatorID int64) ([]Opportunity, error)
	// Create creates a new entity.
	Create(opportunity Opportunity) error
	// Update updates an entity with the ID in the provided entity.
	Update(opportunity Opportunity) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
