package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
)

// Privacy constants
const (
	PrivacyPublic            = iota
	PrivacyOrganizationsOnly = iota
)

// UserProfileField represents a single field in a user's profile.
type UserProfileField struct {
	Model
	UserID    int64          `json:"userId"`
	Name      string         `json:"name"`            // the key/name of the profile field
	Value     string         `json:"value,omitempty"` // the value of the profile field (if it is a string)
	ValueInt  int            `json:"value,omitempty"` // the value of the profile field (if it is an int)
	ValueJSON postgres.Jsonb `json:"value,omitempty"` // the value of the profile field in JSON
	MixedID   int64          `json:"-"`               // represents an ID for a relation to a mixed-type item in the database.
	Privacy   int            `json:"privacy"`
}

// UserProfileFieldRepository represents a repository of UserProfileField.
type UserProfileFieldRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*UserProfileField, error)
	// FindByUserID finds entities by UserID.
	FindByUserID(id int64) ([]UserProfileField, error)
	// FindUserFieldByID finds a single entity by UserID and field name.
	FindUserFieldByID(id int64, name string) (*UserProfileField, error)
	// Create creates a new entity.
	Create(profileField UserProfileField) error
	// Update updates an entity with the ID in the provided entity.
	Update(profileField UserProfileField) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
