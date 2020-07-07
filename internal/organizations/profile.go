package organizations

import (
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/pkg/location"
)

// OrganizationProfile represents an organization's profile.
type OrganizationProfile struct {
	ID             int64                             `json:"id"`
	CreatorID      int64                             `json:"creatorId" scope:"collaborator"` // the organization's creator's ID (User)
	Name           string                            `json:"name"`                           // the organization's name
	Description    string                            `json:"description"`                    // a description of the organization
	ProfilePicture string                            `json:"profilePicture,omitempty"`       // the url for the organization's profile picture
	WebsiteURL     string                            `json:"websiteUrl"`                     // the organization's website's URL
	Tags           []models.Tag                      `json:"tags,omitempty"`                 // the model's tags
	Location       *location.Location                `json:"location,omitempty"`             // the location of the organization
	ProfileFields  []models.OrganizationProfileField `json:"profile,omitempty"`              // fields of the organization's profile
}
