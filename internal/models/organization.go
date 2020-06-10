package models

import "github.com/jinzhu/gorm"

// Organization represents a single volunteering organization.
type Organization struct {
	gorm.Model
	ID             int64                      `gorm:"primary_key" json:"id"` // the organization's unique int64 ID
	Active         bool                       `json:"-"`                     // controls whether or not the organization is active
	CreatorID      int64                      `json:"creatorId"`             // the organization's creator's ID (User)
	Name           string                     `json:"name"`                  // the organization's name
	ProfilePicture string                     `json:"profilePicture"`        // the organization's profile picture/logo
	WebsiteURL     string                     `json:"websiteUrl"`            // the organization's website's URL
	ProfileFields  []OrganizationProfileField `json:"profile"`               // fields of the organization's profile
}
