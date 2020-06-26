package models

// Organization represents a single volunteering organization.
type Organization struct {
	Model
	Active         bool                       `json:"-"`              // controls whether or not the organization is active
	CreatorID      int64                      `json:"creatorId"`      // the organization's creator's ID (User)
	Name           string                     `json:"name"`           // the organization's name
	Description    string                     `json:"description"`    // a description of the organization
	ProfilePicture string                     `json:"profilePicture"` // the organization's profile picture/logo
	WebsiteURL     string                     `json:"websiteUrl"`     // the organization's website's URL
	Location       string                     `json:"location"`
	ProfileFields  []OrganizationProfileField `json:"profile"` // fields of the organization's profile
}
