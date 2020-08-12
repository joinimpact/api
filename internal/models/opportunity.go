package models

import "context"

// Opportunity represents an organization's opportunity.
type Opportunity struct {
	Model
	Active                  bool                     `json:"-"`                             // when false, the opportunity will be treated as if deleted. Useful for suspensions, etc.
	OrganizationID          int64                    `json:"opportunityId"`                 // the id of the parent organization
	Organization            Organization             `json:"-"`                             //
	CreatorID               int64                    `json:"creatorId"`                     // the id of the user who created the opportunity initially
	Creator                 User                     `json:"-" gorm:"foreignkey:CreatorID"` //
	Public                  bool                     `json:"public"`                        // whether or not the opportunity should be shown to volunteers
	Title                   string                   `json:"title"`                         // the title of the opportunity
	ProfilePicture          string                   `json:"profilePicture"`                // a url to the opportunity's banner image
	Description             string                   `json:"description"`                   // a long description of the opportunity and its purpose
	OpportunityTags         []OpportunityTag         `json:"-"`
	OpportunityRequirements *OpportunityRequirements `json:"-"`
	OpportunityLimits       *OpportunityLimits       `json:"-"`
}

// OpportunityRepository represents a repository of opportunities.
type OpportunityRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(ctx context.Context, id int64) (*Opportunity, error)
	// FindByIDs finds multiple entities by an array of IDs.
	FindByIDs(ctx context.Context, ids []int64) ([]Opportunity, error)
	// FindByOrganizationID finds multiple entities by the organization ID.
	FindByOrganizationID(ctx context.Context, organizationID int64) ([]Opportunity, error)
	// FindByCreatorID finds multiple entities by the creator ID.
	FindByCreatorID(ctx context.Context, creatorID int64) ([]Opportunity, error)
	// Create creates a new entity.
	Create(ctx context.Context, opportunity Opportunity) error
	// Update updates an entity with the ID in the provided entity.
	Update(ctx context.Context, opportunity Opportunity) error
	// Save saves all fields in the provided entity.
	Save(ctx context.Context, opportunity Opportunity) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(ctx context.Context, id int64) error
}
