package models

// OrganizationTag defines a single user's area of interest.
type OrganizationTag struct {
	Model
	OrganizationID int64        `json:"-"`
	Organization   Organization `json:"-"`
	TagID          int64        `json:"-"`
	Tag            Tag
}

// OrganizationTagRepository represents a repository of OrganizationTag.
type OrganizationTagRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OrganizationTag, error)
	// FindByOrganizationID finds entities by Organization ID.
	FindByOrganizationID(organizationID int64) ([]OrganizationTag, error)
	// FindOrganizationTagByID finds a single entity by OrganizationID and tag ID.
	FindOrganizationTagByID(organizationID int64, tagID int64) (*OrganizationTag, error)
	// Create creates a new entity.
	Create(userTag OrganizationTag) error
	// Update updates an entity with the ID in the provided entity.
	Update(userTag OrganizationTag) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
