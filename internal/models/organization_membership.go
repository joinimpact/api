package models

import "github.com/jinzhu/gorm"

// OrganizationMembership creates a relationship between Organizations and their employees.
type OrganizationMembership struct {
	gorm.Model
	ID              int64 `gorm:"primary_key" json:"id"` // the membership's unique int64 ID
	Active          bool  `json:"-"`                     // controls whether or not the entity is active
	UserID          int64 `json:"userId"`                // the ID of the user being granted membership
	OrganizationID  int64 `json:"organizationId"`        // the ID of the organization the user is being granted access to
	PermissionsFlag int   `json:"permissionsFlag"`       // a flag which designates permissions the user has
	InviterID       int64 `json:"inviterId"`             // the ID of the user who invited the member to the organization
}
