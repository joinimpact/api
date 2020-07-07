package models

// Organization represents a single volunteering organization.
type Organization struct {
	Model
	Active            bool                       `json:"-"`              // controls whether or not the organization is active
	CreatorID         int64                      `json:"creatorId"`      // the organization's creator's ID (User)
	Name              string                     `json:"name"`           // the organization's name
	Description       string                     `json:"description"`    // a description of the organization
	ProfilePicture    string                     `json:"profilePicture"` // the organization's profile picture/logo
	WebsiteURL        string                     `json:"websiteUrl"`     // the organization's website's URL
	LocationLatitude  float64                    `json:"-"`              // the latitude of the organization's city
	LocationLongitude float64                    `json:"-"`              // the longitude of the organization's city
	ProfileFields     []OrganizationProfileField `json:"profile"`        // fields of the organization's profile
}

// OrganizationRepository represents a repository of organizations.
type OrganizationRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*Organization, error)
	// FindByCreatorID finds multiple entities by the creator's ID.
	FindByCreatorID(creatorID int64) ([]Organization, error)
	// Create creates a new entity.
	Create(organization Organization) error
	// Update updates an entity with the ID in the provided entity.
	Update(organization Organization) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
