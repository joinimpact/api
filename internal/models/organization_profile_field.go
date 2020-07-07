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

// OrganizationProfileFieldRepository represents a repository of OrganizationProfileField.
type OrganizationProfileFieldRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*OrganizationProfileField, error)
	// FindByOrganizationID finds entities by OrganizationID.
	FindByOrganizationID(id int64) ([]OrganizationProfileField, error)
	// FinOrganizationFieldByName finds a single entity by Organization ID and field name.
	FindOrganizationFieldByName(id int64, name string) (*OrganizationProfileField, error)
	// Create creates a new entity.
	Create(profileField OrganizationProfileField) error
	// Update updates an entity with the ID in the provided entity.
	Update(profileField OrganizationProfileField) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
