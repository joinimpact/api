package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
)

// OrganizationProfileField represents a single field in an organization's profile.
type OrganizationProfileField struct {
	Model
	OrganizationID int64          `json:"organizationId"`  // the ID of the organization the field applies to
	Name           string         `json:"name"`            // the key/name of the profile field
	Value          string         `json:"value,omitempty"` // the value of the profile field (if it is a string)
	ValueInt       int            `json:"value,omitempty"` // the value of the profile field (if it is an int)
	ValueJSON      postgres.Jsonb `json:"value,omitempty"` // the value of the profile field in JSON
	MixedID        int64          `json:"-"`               // represents an ID for a relation to a mixed-type item in the database
}
